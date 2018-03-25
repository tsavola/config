// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
)

// Loader makes a flag.Value which reads configuration files.
func Loader(target interface{}) flag.Value {
	return loader{target}
}

type loader struct {
	config interface{}
}

func (l loader) Set(filename string) error {
	return ReadFile(filename, l.config)
}

func (loader) String() string {
	return ""
}

// Setter makes a flag.Value which applies configuration expressions.
func Setter(target interface{}) flag.Value {
	return setter{target}
}

type setter struct {
	config interface{}
}

func (s setter) Set(expr string) error {
	return Apply(s.config, expr)
}

func (setter) String() string {
	return ""
}
