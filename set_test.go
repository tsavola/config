// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"reflect"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := new(testConfig)
	c.Foo.Key2 = 67890

	if err := Set(c, "foo.key1", true); err != nil {
		t.Error(err)
	}
	if !c.Foo.Key1 {
		t.Fail()
	}

	if err := Set(c, "foo.key2", true); err == nil {
		t.Fail()
	}
	if c.Foo.Key2 != 67890 {
		t.Fail()
	}

	if err := Set(c, "foo.key2b", 10); err == nil {
		t.Fail()
	}
	if c.Foo.Key2b != 0 {
		t.Fail()
	}

	if err := Set(c, "foo.key3", 10); err == nil {
		t.Fail()
	}
	if c.Foo.Key3 != 0 {
		t.Fail()
	}

	if err := Set(c, "foo.key3", int32(10)); err != nil {
		t.Error(err)
	}
	if c.Foo.Key3 != 10 {
		t.Fail()
	}

	if err := Set(c, "foo.key4", 10); err == nil {
		t.Fail()
	}
	if c.Foo.Key4 != 0 {
		t.Fail()
	}

	if err := Set(c, "foo.key4", int64(10)); err != nil {
		t.Error(err)
	}
	if c.Foo.Key4 != 10 {
		t.Fail()
	}

	if err := Set(c, "foo.key9", "Hello, World"); err == nil {
		t.Fail()
	}
	if c.Foo.Key9 != 0 {
		t.Fail()
	}

	if err := Set(c, "foo.key10", "Hello, World"); err != nil {
		t.Error(err)
	}
	if c.Foo.Key10 != "Hello, World" {
		t.Fail()
	}

	if err := Set(c, "foo.key11", []string{"Hello", "World"}); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c.Foo.Key11, []string{"Hello", "World"}) {
		t.Fail()
	}

	if err := Set(c, "baz.interval", int64(time.Second)); err == nil {
		t.Fail()
	}
	if c.Baz.Interval != 0 {
		t.Fail()
	}

	if err := Set(c, "baz.interval", time.Second); err != nil {
		t.Error(err)
	}
	if c.Baz.Interval != time.Second {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	c := new(testConfig)
	c.Bar = 67890

	if x, err := Get(c, "bar"); err != nil {
		t.Error(err)
	} else if x.(int) != 67890 {
		t.Fail()
	}
}
