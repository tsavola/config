// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
)

// FileReader makes a flag.Value which reads configuration files.
func FileReader(target interface{}) flag.Value {
	return fileReader{target}
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

// Assigner makes a flag.Value which applies configuration expressions.
func Assigner(target interface{}) flag.Value {
	return assigner{target}
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
