// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
	"testing"
)

func TestSetter(t *testing.T) {
	c := new(testConfig)

	s := flag.NewFlagSet("test", flag.PanicOnError)
	s.Var(Setter(c), "c", "path.to.key=value")

	if err := s.Parse([]string{
		"-c", "foo.key1=true",
		"-c", "foo.key2=-10",
		"-c", "foo.key3=-11",
		"-c", "foo.key4=-100000000000000",
		"-c", "foo.key5=10",
		"-c", "foo.key6=11",
		"-c", "foo.key7=100000000000000",
		"-c", "foo.key8=1.5",
		"-c", "foo.key9=1.0000000000005",
		"-c", "foo.key10=hello, world",
		"-c", "bar=12345",
		"-c", "baz.quux.key_a=true",
		"-c", "baz.quux.key_b=true",
	}); err != nil {
		t.Fatal(err)
	}

	testConfigValues(t, c)
}
