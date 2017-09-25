// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"

	"github.com/CanonicalLtd/jujushell/config"
)

var readTests = []struct {
	about          string
	content        []byte
	expectedConfig *config.Config
	expectedError  string
}{{
	about: "valid config",
	content: mustMarshalYAML(map[string]interface{}{
		"juju-addrs": []string{"1.2.3.4", "4.3.2.1"},
		"log-level":  "debug",
		"port":       8047,
	}),
	expectedConfig: &config.Config{
		JujuAddrs: []string{"1.2.3.4", "4.3.2.1"},
		LogLevel:  zapcore.DebugLevel,
		Port:      8047,
	},
}, {
	about:         "unreadable config",
	content:       []byte("not a yaml"),
	expectedError: `cannot parse ".*": yaml: unmarshal errors:\n.*`,
}, {
	about:         "invalid config",
	expectedError: `invalid configuration at ".*": missing fields juju-addrs, port`,
}}

func TestRead(t *testing.T) {
	for _, test := range readTests {
		t.Run(test.about, func(t *testing.T) {
			c := qt.New(t)

			// Create the config file.
			f, err := ioutil.TempFile("", "config")
			c.Assert(err, qt.Equals, nil)
			defer f.Close()
			defer os.Remove(f.Name())
			_, err = f.Write(test.content)
			c.Assert(err, qt.Equals, nil)

			// Read the config file.
			conf, err := config.Read(f.Name())
			if test.expectedError != "" {
				c.Assert(err, qt.ErrorMatches, test.expectedError)
				c.Assert(conf, qt.IsNil)
				return
			}
			c.Assert(err, qt.Equals, nil)
			c.Assert(conf, qt.DeepEquals, test.expectedConfig)
		})
	}
}

func mustMarshalYAML(v interface{}) []byte {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}