// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestSettings(t *testing.T) {
	c := new(testConfig)
	c.Bar = 12345
	c.Baz.Embed2.TestConfigEmbed = new(TestConfigEmbed)

	if ss := Settings(c); !reflect.DeepEqual(ss, []Setting{
		{"foo.key1", reflect.TypeOf(false), ""},
		{"foo.key2", reflect.TypeOf(0), ""},
		{"foo.key2b", reflect.TypeOf(int8(0)), ""},
		{"foo.key3a", reflect.TypeOf(int16(0)), ""},
		{"foo.key3", reflect.TypeOf(int32(0)), ""},
		{"foo.key4", reflect.TypeOf(int64(0)), ""},
		{"foo.key5", reflect.TypeOf(uint(0)), ""},
		{"foo.key5b", reflect.TypeOf(uint8(0)), ""},
		{"foo.key6a", reflect.TypeOf(uint16(0)), ""},
		{"foo.key6", reflect.TypeOf(uint32(0)), ""},
		{"foo.key7", reflect.TypeOf(uint64(0)), ""},
		{"foo.key8", reflect.TypeOf(float32(0)), ""},
		{"foo.key9", reflect.TypeOf(0.0), ""},
		{"foo.key10", reflect.TypeOf(""), ""},
		{"foo.key11", reflect.TypeOf([]string{}), ""},
		{"bar", reflect.TypeOf(0), "12345"},
		{"baz.quux.key_a", reflect.TypeOf(""), ""},
		{"baz.quux.key_b", reflect.TypeOf(false), ""},
		{"baz.interval", reflect.TypeOf(time.Duration(0)), ""},
		{"baz.embedded", reflect.TypeOf(false), ""},
		{"baz.embed1.embedded", reflect.TypeOf(false), ""},
		{"baz.embed2.embedded", reflect.TypeOf(false), ""},
	}) {
		t.Errorf("%#v", ss)
	}

	b := new(bytes.Buffer)
	PrintSettings(b, c)
	t.Logf("\n%s", b)
}
