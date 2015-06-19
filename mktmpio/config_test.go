package mktmpio

import (
	"testing"
)

func TestConfigLoading(t *testing.T) {
	err, cfg := LoadConfig()
	if err != nil {
		t.Error("LoadConfig returned an error:", err)
	}
	if len(cfg.Token) < 10 {
		t.Error("config token too short:", cfg)
	}
}
