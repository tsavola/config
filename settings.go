// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Setting documents a settable configuration path.
type Setting struct {
	Path string
	Type reflect.Type
}

func (s Setting) String() string {
	return s.Path
}

// Settings lists the settable configuration paths.
func Settings(config interface{}) (settings []Setting) {
	settings = enumerate(settings, "", reflect.ValueOf(config))
	return
}

func enumerate(list []Setting, prefix string, node reflect.Value) []Setting {
	if node.Type().Kind() == reflect.Ptr {
		node = reflect.Indirect(node)
	}

	for i := 0; i < node.Type().NumField(); i++ {
		field := node.Type().Field(i)

		path := prefix
		if !field.Anonymous {
			if len(path) > 0 {
				path += "."
			}
			path += strings.ToLower(field.Name)
		}

		t := field.Type
		kind := t.Kind()

		if kind == reflect.Ptr {
			list = enumerate(list, path, node.Field(i))
		} else {
			switch kind {
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
				list = append(list, Setting{path, t})

			case reflect.Struct:
				list = enumerate(list, path, node.Field(i))
			}
		}
	}

	return list
}

// PrintSettings of the given configuration.  Writer defaults to the default
// flag set's output.
func PrintSettings(w io.Writer, config interface{}) {
	if w == nil {
		w = flag.CommandLine.Output()
	}

	for _, s := range Settings(config) {
		fmt.Fprintf(w, "  %s %s\n", s.Path, s.Type)
	}
}

// FlagUsage creates a function which may be used as flag.Usage.  It includes
// the default usage and the configuration settings.
func FlagUsage(config interface{}) func() {
	stdUsage := flag.Usage

	return func() {
		stdUsage()
		fmt.Fprintf(flag.CommandLine.Output(), "\nConfiguration settings:\n")
		PrintSettings(nil, config)
	}
}
