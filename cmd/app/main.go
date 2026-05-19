package main

import (
	"fmt"
	"impulse/internal/adapters/input/config"
	"impulse/internal/adapters/input/events"
	"impulse/internal/adapters/output"
	engine2 "impulse/internal/application/engine"
	"impulse/internal/application/processing"
	"impulse/internal/application/processing/handlers"
	"impulse/internal/application/report"
	"impulse/internal/domain/models"
	events2 "impulse/internal/domain/models/events"
	"impulse/internal/infrastructure/storage"
	"io"
	"log"
	"os"
)

func LoadConfig() (models.Config, error) {
	cfgFile, err := os.Open("./config.json")
	if err != nil {
		return models.Config{}, err
	}
	defer func() {
		if cerr := cfgFile.Close(); cerr != nil {
			log.Printf("unable to close config file: %v\n", cerr)
		}
	}()

	configSource := config.NewJsonConfigSource(cfgFile)

	cfg, err := configSource.Load()
	if err != nil {
		return models.Config{}, err
	}

	return cfg, nil
}

func LoadProcessor(cfg *models.Config) *processing.Processor {
	if cfg == nil {
		log.Fatal("nil config")
	}

	processor := processing.NewProcessor()
	processor.RegisterHandler(events2.UserEnteredDungeonEventID, handlers.NewEnteredDungeonHandler(cfg))
	processor.RegisterHandler(events2.UserKilledMonsterEventID, handlers.NewKillMonsterHandler(cfg))
	processor.RegisterHandler(events2.UserWentToNextFloorEventID, handlers.NewNextFloorHandler(cfg))
	processor.RegisterHandler(events2.UserWentToPreviousFloorEventID, handlers.NewPreviousFloorHandler(cfg))
	processor.RegisterHandler(events2.UserEnteredBossFloorEventID, handlers.NewEnteredBossFloorHandler(cfg))
	processor.RegisterHandler(events2.UserKilledBossEventID, handlers.NewKillBossHandler())
	processor.RegisterHandler(events2.UserLeftDungeonEventID, handlers.NewLeftDungeonHandler())
	processor.RegisterHandler(events2.UserCantContinueDueToReasonEventID, handlers.NewCannotContinueHandler())
	processor.RegisterHandler(events2.UserRestoredHealthEventID, handlers.NewRestoredHealthHandler())
	processor.RegisterHandler(events2.UserReceivedDamageEventID, handlers.NewDamageHandler())

	return processor
}

func SetupPipeline(cfg *models.Config, eventsReader io.Reader, outWriter io.Writer) (*engine2.Engine, events.EventSource, output.EventSink) {
	if cfg == nil {
		log.Fatal("config is nil")
	}

	sessionRepo := storage.NewInMemorySessionStore()
	processor := LoadProcessor(cfg)
	reportBuilder := report.NewReportBuilder(cfg)
	reportFormatter := report.NewStringFormatter()

	source := events.NewFileEventSource(eventsReader)
	sink := output.NewStdoutEventSink(outWriter)
	engine := engine2.NewEngine(sessionRepo, processor, reportBuilder, reportFormatter)

	return engine, source, sink
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./events")
	if err != nil {
		panic(err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("unable to close events file: %v\n", cerr)
		}
	}()

	engine, source, sink := SetupPipeline(&cfg, file, os.Stdout)

	for {

		event, err := source.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("unable to read event: %v\n", err)
			continue
		}

		out, err := engine.Process(event)
		if err != nil {
			log.Printf("unable to process event: %v\n", err)
			continue
		}

		err = sink.WriteMany(out)
		if err != nil {
			log.Printf("%v\n", err)
		}
	}

	finalReport := engine.Report()
	fmt.Println(finalReport)
}
