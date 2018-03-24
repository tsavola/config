// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"testing"
)

type testConfig struct {
	Foo struct {
		Key1  bool
		Key2  int
		Key3  int32
		Key4  int64
		Key5  uint
		Key6  uint32
		Key7  uint64
		Key8  float32
		Key9  float64
		Key10 string
	}

	Bar int

	Baz struct {
		Quux testConfigQuux
	}
}

type testConfigQuux struct {
	Key_a string
	Key_b bool
}

var testConfigYAML = `foo:
  key1: true
  key2: -10
  key3: -11
  key4: -100000000000000
  key5: 10
  key6: 11
  key7: 100000000000000
  key8: 1.5
  key9: 1.0000000000005
  key10: hello, world
bar: 12345
baz:
  quux:
    key_a: "true"
    key_b: true
`

func testConfigValues(t *testing.T, c *testConfig) {
	if c.Foo.Key1 != true {
		t.Fail()
	}
	if c.Foo.Key2 != -10 {
		t.Fail()
	}
	if c.Foo.Key3 != -11 {
		t.Fail()
	}
	if c.Foo.Key4 != -100000000000000 {
		t.Fail()
	}
	if c.Foo.Key5 != 10 {
		t.Fail()
	}
	if c.Foo.Key6 != 11 {
		t.Fail()
	}
	if c.Foo.Key7 != 100000000000000 {
		t.Fail()
	}
	if c.Foo.Key8 != 1.5 {
		t.Fail()
	}
	if c.Foo.Key9 != 1.0000000000005 {
		t.Fail()
	}
	if c.Foo.Key10 != "hello, world" {
		t.Fail()
	}
	if c.Bar != 12345 {
		t.Fail()
	}
	if c.Baz.Quux.Key_a != "true" {
		t.Fail()
	}
	if c.Baz.Quux.Key_b != true {
		t.Fail()
	}
}
