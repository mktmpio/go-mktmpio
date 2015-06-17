package mktmpio

import (
	"testing"
)

func TestConfigLoading(t *testing.T) {
	err, cfg := LoadConfig()
	if err != nil {
		t.Fail()
	}
	if len(cfg.Token) < 10 {
		t.Fail()
	}
	t.Logf("client: %v\n", cfg)
}
