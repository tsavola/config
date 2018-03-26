// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"reflect"
	"testing"
	"time"
)

type testConfig struct {
	Foo struct {
		Key1  bool
		Key2  int
		Key2b int8
		Key3a int16
		Key3  int32
		Key4  int64
		Key5  uint
		Key5b uint8
		Key6a uint16
		Key6  uint32
		Key7  uint64
		Key8  float32
		Key9  float64
		Key10 string
		Key11 []string
	}

	Bar int

	Baz struct {
		Quux     testConfigQuux
		Interval time.Duration
	}
}

type testConfigQuux struct {
	Key_a string
	Key_b bool
}

var testConfigYAML = `foo:
  key1: true
  key2: -10
  key2b: -128
  key3a: -32768
  key3: -11
  key4: -100000000000000
  key5: 10
  key5b: 255
  key6a: 65535
  key6: 11
  key7: 100000000000000
  key8: 1.5
  key9: 1.0000000000005
  key10: hello, world
  key11:
  - hello
  - world
bar: 12345
baz:
  quux:
    key_a: "true"
    key_b: true
  interval: 10h9m8.007006005s
`

func testConfigValues(t *testing.T, c *testConfig) {
	if !c.Foo.Key1 {
		t.Fail()
	}
	if c.Foo.Key2 != -10 {
		t.Fail()
	}
	if c.Foo.Key2b != -128 {
		t.Fail()
	}
	if c.Foo.Key3a != -32768 {
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
	if c.Foo.Key5b != 255 {
		t.Fail()
	}
	if c.Foo.Key6a != 65535 {
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
	if !reflect.DeepEqual(c.Foo.Key11, []string{"hello", "world"}) {
		t.Fail()
	}
	if c.Bar != 12345 {
		t.Fail()
	}
	if c.Baz.Quux.Key_a != "true" {
		t.Fail()
	}
	if !c.Baz.Quux.Key_b {
		t.Fail()
	}
	if c.Baz.Interval != 10*time.Hour+9*time.Minute+8*time.Second+7*time.Millisecond+6*time.Microsecond+5*time.Nanosecond {
		t.Fail()
	}
}
