package goexmars

import (
	"strings"
	"testing"
)

func TestFightConfigValidateValid(t *testing.T) {
	if err := DefaultConfig.Validate(); err != nil {
		t.Fatalf("expected DefaultConfig to validate, got %v", err)
	}
}

func TestFightConfigValidateInvalid(t *testing.T) {
	cfg := DefaultConfig
	cfg.CoreSize = 0
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "CoreSize") {
		t.Fatalf("expected CoreSize validation error, got %v", err)
	}

	cfg = DefaultConfig
	cfg.MinSep = cfg.CoreSize + 1
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "MinSep") {
		t.Fatalf("expected MinSep validation error, got %v", err)
	}
}
