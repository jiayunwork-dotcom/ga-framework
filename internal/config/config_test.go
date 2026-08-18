package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefaultConfig_Valid(t *testing.T) {
	c := DefaultConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestValidate_InvalidPopSize(t *testing.T) {
	c := DefaultConfig()
	c.PopSize = 0
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "pop_size") {
		t.Fatalf("expected pop_size error, got: %v", err)
	}
}

func TestValidate_InvalidMutateRate(t *testing.T) {
	c := DefaultConfig()
	c.MutateRate = 1.5
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "mutate_rate") {
		t.Fatalf("expected mutate_rate error, got: %v", err)
	}
}

func TestValidate_UnknownSelection(t *testing.T) {
	c := DefaultConfig()
	c.Selection = "magic"
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "selection") {
		t.Fatalf("expected selection error, got: %v", err)
	}
}

func TestReadWriteJSON_Roundtrip(t *testing.T) {
	original := DefaultConfig()
	original.Seed = 42
	var buf bytes.Buffer
	err := original.WriteJSON(&buf)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	loaded, err := ReadJSON(&buf)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if loaded.PopSize != original.PopSize || loaded.Seed != original.Seed {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", original, loaded)
	}
}

func TestMerge(t *testing.T) {
	base := DefaultConfig()
	overlay := GAConfig{PopSize: 200, Selection: "rank"}
	merged := Merge(base, overlay)
	if merged.PopSize != 200 {
		t.Fatalf("expected PopSize 200, got %d", merged.PopSize)
	}
	if merged.Selection != "rank" {
		t.Fatalf("expected rank selection, got %s", merged.Selection)
	}
	if merged.Genes != base.Genes {
		t.Fatalf("genes should keep base value")
	}
}

func TestSummary(t *testing.T) {
	c := DefaultConfig()
	s := c.Summary()
	if !strings.Contains(s, "pop=50") {
		t.Fatalf("summary missing pop: %s", s)
	}
	if !strings.Contains(s, "sel=tournament") {
		t.Fatalf("summary missing selection: %s", s)
	}
}

func TestValidate_InvalidBounds(t *testing.T) {
	c := DefaultConfig()
	c.BoundsLo = 10
	c.BoundsHi = 5
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "bounds") {
		t.Fatalf("expected bounds error, got: %v", err)
	}
}

func TestValidate_InvalidElite(t *testing.T) {
	c := DefaultConfig()
	c.Elite = 100 // >= PopSize(50)
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "elite") {
		t.Fatalf("expected elite error, got: %v", err)
	}
}
