package conf

import (
	"os"
	"testing"
)

type TestConfig struct {
	StringField string   `conf:"STRING_VAR,default_value"`
	IntField    int      `conf:"INT_VAR,10"`
	BoolField   bool     `conf:"BOOL_VAR,false"`
	FloatField  float64  `conf:"FLOAT_VAR,1.23"`
	SliceField  []string `conf:"SLICE_VAR,a;b;c"`
}

func TestLoad(t *testing.T) {
	os.Setenv("STRING_VAR", "test_value")
	os.Setenv("INT_VAR", "42")
	os.Setenv("BOOL_VAR", "true")
	os.Setenv("FLOAT_VAR", "3.14")
	os.Setenv("SLICE_VAR", "x;y;z")

	var cfg TestConfig
	err := Load(&cfg)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.StringField != "test_value" {
		t.Errorf("expected StringField to be 'test_value', got %s", cfg.StringField)
	}
	if cfg.IntField != 42 {
		t.Errorf("expected IntField to be 42, got %d", cfg.IntField)
	}
	if cfg.BoolField != true {
		t.Errorf("expected BoolField to be true, got %t", cfg.BoolField)
	}
	if cfg.FloatField != 3.14 {
		t.Errorf("expected FloatField to be 3.14, got %f", cfg.FloatField)
	}
	if len(cfg.SliceField) != 3 || cfg.SliceField[0] != "x" || cfg.SliceField[1] != "y" || cfg.SliceField[2] != "z" {
		t.Errorf("expected SliceField to be [x y z], got %v", cfg.SliceField)
	}

	// Cleanup
	os.Unsetenv("STRING_VAR")
	os.Unsetenv("INT_VAR")
	os.Unsetenv("BOOL_VAR")
	os.Unsetenv("FLOAT_VAR")
	os.Unsetenv("SLICE_VAR")
}
