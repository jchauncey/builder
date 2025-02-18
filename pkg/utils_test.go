package pkg

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	dtime "github.com/deis/deis/pkg/time"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	// we don't have to do anything here, since the buffer is just some data in memory
	return nil
}

func stringInSlice(list []string, s string) bool {
	for _, li := range list {
		if li == s {
			return true
		}
	}
	return false
}

func TestParseConfigGood(t *testing.T) {
	// mock the controller response
	resp := bytes.NewBufferString(`{"owner": "test",
			"app": "example-go",
			"values": {"FOO": "bar", "CAR": 1234},
			"memory": {},
			"cpu": {},
			"tags": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"}`)

	config, err := ParseConfig(resp.Bytes())

	if err != nil {
		t.Error(err)
	}

	if config.Values["FOO"] != "bar" {
		t.Errorf("expected FOO='bar', got FOO='%v'", config.Values["FOO"])
	}

	if car, ok := config.Values["CAR"].(float64); ok {
		if car != 1234 {
			t.Errorf("expected CAR=1234, got CAR=%d", config.Values["CAR"])
		}
	} else {
		t.Error("expected CAR to be of type float64")
	}
}

func TestGetDefaultTypeGood(t *testing.T) {
	goodData := [][]byte{[]byte(`default_process_types:
  web: while true; do echo hello; sleep 1; done`),
		[]byte(`foo: bar
default_process_types:
  web: while true; do echo hello; sleep 1; done`),
		[]byte(``)}

	for _, data := range goodData {
		defaultType, err := GetDefaultType(data)
		if err != nil {
			t.Error(err)
		}
		if defaultType != `{"web":"while true; do echo hello; sleep 1; done"}` && string(data) != "" {
			t.Errorf("incorrect default type, got %s", defaultType)
		}
		if string(data) == "" && defaultType != "{}" {
			t.Errorf("incorrect default type, got %s", defaultType)
		}
	}
}

func TestParseControllerConfigGood(t *testing.T) {
	// mock controller config response
	resp := []byte(`{"owner": "test",
		"app": "example-go",
		"values": {"FOO": "bar", "CAR": "star"},
		"memory": {},
		"cpu": {},
		"tags": {},
		"created": "2014-01-01T00:00:00UTC",
		"updated": "2014-01-01T00:00:00UTC",
		"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
	}`)

	config, err := ParseControllerConfig(resp)

	if err != nil {
		t.Errorf("expected to pass, got '%v'", err)
	}

	if len(config) != 2 {
		t.Errorf("expected 2, got %d", len(config))
	}

	if !stringInSlice(config, " -e CAR=\"star\"") {
		t.Error("expected ' -e CAR=\"star\"' in slice")
	}
}

func TestTimeSerialize(t *testing.T) {
	time, err := json.Marshal(&dtime.Time{Time: time.Now().UTC()})

	if err != nil {
		t.Errorf("expected to be able to serialize time as json, got '%v'", err)
	}

	if !strings.Contains(string(time), "UTC") {
		t.Errorf("could not find 'UTC' in datetime, got '%s'", string(time))
	}
}
