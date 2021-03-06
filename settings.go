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
	Path    string
	Type    reflect.Type
	Default string
}

func (s Setting) String() string {
	return s.Path
}

// Settings lists the settable configuration paths.
func Settings(config interface{}) []Setting {
	return enumerate(nil, "", reflect.ValueOf(config))
}

func enumerate(list []Setting, prefix string, node reflect.Value) []Setting {
	if node.Type().Kind() == reflect.Ptr {
		if node.IsNil() {
			return list
		}
		node = node.Elem()
	}

	for i := 0; i < node.Type().NumField(); i++ {
		value := node.Field(i)
		if !value.CanInterface() {
			continue
		}

		field := node.Type().Field(i)

		path := prefix
		if !field.Anonymous {
			if len(path) > 0 {
				path += "."
			}
			path += strings.ToLower(field.Name)
		}

		kind := field.Type.Kind()

		if kind == reflect.Ptr {
			list = enumerate(list, path, value)
		} else {
			switch kind {
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
				s := Setting{
					Path: path,
					Type: value.Type(),
				}
				if x := value.Interface(); x != reflect.Zero(value.Type()).Interface() {
					s.Default = fmt.Sprint(x)
				}
				list = append(list, s)

			case reflect.Slice:
				switch value.Type().Elem().Kind() {
				case reflect.String:
					s := Setting{
						Path: path,
						Type: value.Type(),
					}
					if repr := fmt.Sprintf("%q", value.Interface()); len(repr) > 2 {
						s.Default = repr
					}
					list = append(list, s)
				}

			case reflect.Struct:
				list = enumerate(list, path, value)
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
		if s.Default == "" {
			fmt.Fprintf(w, "  %s %s\n", s.Path, s.Type)
		} else {
			fmt.Fprintf(w, "  %s %s (%s)\n", s.Path, s.Type, s.Default)
		}
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
