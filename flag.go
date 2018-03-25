// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
)

// FileReader makes a ``dynamic value'' which reads files into the
// configuration as its receives filenames.
func FileReader(config interface{}) flag.Value {
	return fileReader{config}
}

type fileReader struct {
	config interface{}
}

func (fr fileReader) Set(filename string) error {
	return ReadFile(filename, fr.config)
}

func (fileReader) String() string {
	return ""
}

// Assigner makes a ``dynamic value'' which sets fields in the configuration as
// it receives assignment expressions.
func Assigner(config interface{}) flag.Value {
	return assigner{config}
}

type assigner struct {
	config interface{}
}

func (a assigner) Set(expr string) error {
	return Assign(a.config, expr)
}

func (assigner) String() string {
	return ""
}
