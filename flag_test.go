// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestAssigner(t *testing.T) {
	c := new(testConfig)
	c.Baz.Embed2.TestConfigEmbed = new(TestConfigEmbed)

	s := flag.NewFlagSet("test", flag.PanicOnError)
	s.Var(Assigner(c), "c", "path.to.key=value")

	t.Run("Ok", func(t *testing.T) {
		if err := s.Parse([]string{
			"-c", "foo.key1=true",
			"-c", "foo.key2=-10",
			"-c", "foo.key2b=-128",
			"-c", "foo.key3a=-32768",
			"-c", "foo.key3=-11",
			"-c", "foo.key4=-100000000000000",
			"-c", "foo.key5=10",
			"-c", "foo.key5b=255",
			"-c", "foo.key6a=65535",
			"-c", "foo.key6=11",
			"-c", "foo.key7=100000000000000",
			"-c", "foo.key8=1.5",
			"-c", "foo.key9=1.0000000000005",
			"-c", "foo.key10=hello, world",
			"-c", `foo.key11=["hello", "world"]`,
			"-c", "bar=12345",
			"-c", "baz.quux.key_a=true",
			"-c", "baz.quux.key_b=yes",
			"-c", "baz.interval=10h9m8s7ms6µs5ns",
		}); err != nil {
			t.Fatal(err)
		}

		testConfigValues(t, c)
	})

	s.SetOutput(ioutil.Discard)

	for i := 1; i < 10; i++ {
		t.Run("WrongType", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fail()
				}
			}()
			s.Parse([]string{
				"-c", fmt.Sprintf("foo.key%d=this is a string", i),
			})
			t.Fail()
		})
	}

	for _, path := range []string{
		"nonexistent",
		"non.existent",
		"foo.nonexistent",
		"nonexistent.key1",
		"baz.quux.nonexistent",
		"baz.quux.0",
		"0",
	} {
		t.Run("UnknownPath", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fail()
				}
			}()
			s.Parse([]string{
				"-c", fmt.Sprintf("%s=1", path),
			})
			t.Fail()
		})
	}

	for _, spec := range []struct {
		num int
		arg string
	}{
		{0, ""},
		{0, `[]`},
		{1, "one"},
		{1, `["one"]`},
		{1, `[""]`},
		{2, `["two","two"]`},
		{2, `["",""]`},
		{3, ` [ "th", "r", "ee" ] `},
	} {
		t.Run("StringSlice", func(t *testing.T) {
			if err := s.Parse([]string{
				"-c", fmt.Sprintf("foo.key11=%s", spec.arg),
			}); err != nil {
				t.Fatal(err)
			}
			if len(c.Foo.Key11) != spec.num {
				t.Fail()
			}
		})
	}
}
