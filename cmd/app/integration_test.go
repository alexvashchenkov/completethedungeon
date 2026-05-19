package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	cfgAdapter "impulse/internal/adapters/input/config"
)

func TestReadmeExamplePipeline(t *testing.T) {
	tmpDir := t.TempDir()

	cfgContent := `{"Floors": 2, "Monsters": 2, "OpenAt": "14:05:00", "Duration": 2}`
	cfgPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	events := `[14:00:00] 1 1
[14:00:00] 2 1
[14:10:00] 2 2
[14:10:00] 3 2
[14:11:00] 2 5
[14:12:00] 3 3
[14:14:00] 2 3
[14:27:00] 2 11 60
[14:29:00] 2 11 50
[14:40:00] 1 2
[14:41:00] 1 3
[14:44:00] 1 11 50
[14:45:00] 1 3
[14:48:00] 1 4
[14:48:00] 1 6
[14:49:00] 1 11 25
[14:49:02] 1 10 80
[14:50:00] 1 11 65
[14:59:00] 1 7
[15:04:00] 1 8
`
	eventsPath := filepath.Join(tmpDir, "events")
	if err := os.WriteFile(eventsPath, []byte(events), 0o644); err != nil {
		t.Fatalf("write events: %v", err)
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		t.Fatalf("open cfg: %v", err)
	}
	defer f.Close()
	cfgSrc := cfgAdapter.NewJsonConfigSource(f)
	cfg, err := cfgSrc.Load()
	if err != nil {
		t.Fatalf("load cfg: %v", err)
	}

	ef, err := os.Open(eventsPath)
	if err != nil {
		t.Fatalf("open events: %v", err)
	}
	defer ef.Close()

	var outBuf bytes.Buffer

	eng, source, sink := SetupPipeline(&cfg, ef, &outBuf)

	for {
		ev, err := source.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("source error: %v", err)
		}

		outs, err := eng.Process(ev)
		if err != nil {
			t.Fatalf("process err: %v", err)
		}

		if err := sink.WriteMany(outs); err != nil {
			t.Fatalf("sink write many: %v", err)
		}
	}

	final := eng.Report()

	expected := "Final report:\n[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35\n[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0\n[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100"

	if final != expected {
		t.Fatalf("unexpected final report:\nexpected:\n%s\n\nactual:\n%s", expected, final)
	}
}
