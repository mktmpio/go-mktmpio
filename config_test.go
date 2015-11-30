// Copyright 2015 Datajin Technologies, Inc. All rights reserved.
// Use of this source code is governed by an Artistic-2
// license that can be found in the LICENSE file.

package mktmpio

import (
	"os"
	"testing"
)

func testEquivalent(t *testing.T, a *Config, b *Config) {
	if a.Token != b.Token {
		t.Error("Tokens to not match", a.Token, b.Token)
	}
	if a.URL != b.URL {
		t.Error("URLs to not match", a.URL, b.URL)
	}
}

func TestConfigLoading(t *testing.T) {
	if err := os.Setenv("MKTMPIO_TOKEN", "1234-5678-90abcdef"); err != nil {
		t.Error("Could not set env var in test", err)
	}
	cfg := LoadConfig()
	if cfg.err != nil {
		t.Error("LoadConfig returned an error:", cfg.err)
	}
	if len(cfg.Token) < 10 {
		t.Error("config token too short:", cfg)
	}
}

func TestConfigAppy(t *testing.T) {
	a := new(Config)
	b := a.Apply(&Config{})
	if a == b {
		t.Error("new instance should be separate")
	}
	testEquivalent(t, a, b)
	a.Token = "ATOK"
	b.Token = "BTOK"
	testEquivalent(t, b, a.Apply(b))
	a.URL = "AURL"
	testEquivalent(t, &Config{Token: b.Token, URL: a.URL}, a.Apply(b))
	testEquivalent(t, a, b.Apply(a))
}

func TestConfigFile(t *testing.T) {
	c := FileConfig("example.mktmpio.yml")
	if c.URL != "" {
		t.Error("example config file unexpectedly set the URL")
	}
	if c.Token != "01234567890abcdefghijkl" {
		t.Error("token does not match expected example config file")
	}
	if c.err != nil {
		t.Error("failed to load example config file")
	}
	c = FileConfig("file-that-does-not-exist")
	if c.err == nil {
		t.Error("config should report an error if the file it is from doesn't exist")
	}
}

func TestConfigSave(t *testing.T) {
	tmp := "temp_mktmpio.test.yml"
	c := new(Config)
	c.Save(tmp)
	from := FileConfig(tmp)
	if from.err != nil {
		t.Error("failed to load freshly created file", from.err)
	}
	testEquivalent(t, c, from)
	c.Token = "New Token"
	c.Save(tmp)
	from = FileConfig(tmp)
	if from.err != nil {
		t.Error("failed to load freshly created file", from.err)
	}
	testEquivalent(t, c, from)
}

func TestConfigStringer(t *testing.T) {
	c := new(Config)
	if c.String() != "token: \"\"\n" {
		t.Errorf("expected '%s' to be '%s'", c, "token: \"\"\n")
	}
	c.Token = "NEW TOKEN"
	if c.String() != "token: NEW TOKEN\n" {
		t.Errorf("expected '%s' to be '%s'", c, "token: NEW TOKEN\n")
	}
	c.URL = "http://url/here"
	if c.String() != "token: NEW TOKEN\nurl: http://url/here\n" {
		t.Errorf("expected '%s' to be '%s'", c, "token: NEW TOKEN\nurl: http://url/here\n")
	}
}
