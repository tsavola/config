// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	c := new(testConfig)
	c.Bar = 67890

	if err := Read(strings.NewReader(testConfigYAML), c); err != nil {
		t.Fatal(err)
	}

	testConfigValues(t, c)
}

func TestReadFileIfExists(t *testing.T) {
	if err := ReadFileIfExists("/nonexistent", nil); err != nil {
		t.Error(err)
	}

	if ReadFileIfExists("/dev/zero", nil) == nil {
		t.Fail()
	}
}

func TestWrite(t *testing.T) {
	c := new(testConfig)

	if err := Read(strings.NewReader(testConfigYAML), c); err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)

	if err := Write(b, c); err != nil {
		t.Fatal(err)
	}

	if s := b.String(); s != testConfigYAML {
		t.Error(s)
	}
}
