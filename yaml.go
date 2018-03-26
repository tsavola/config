// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// Read YAML into the configuration.
func Read(r io.Reader, config interface{}) error {
	return yaml.NewDecoder(r).Decode(config)
}

// Read a YAML file into the configuration.
func ReadFile(filename string, config interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	return Read(f, config)
}

// Read a YAML file into the configuration.  No error is returned if the file
// doesn't exist.
func ReadFileIfExists(filename string, config interface{}) (err error) {
	err = ReadFile(filename, config)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

// Write the configuration as YAML.
func Write(w io.Writer, config interface{}) error {
	return yaml.NewEncoder(w).Encode(sanitize(reflect.ValueOf(config).Elem()))
}

// Write the configuration to a YAML file.
func WriteFile(filename string, config interface{}) (err error) {
	data, err := yaml.Marshal(sanitize(reflect.ValueOf(config).Elem()))
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}

func sanitize(struc reflect.Value) (sane yaml.MapSlice) {
	for i := 0; i < struc.Type().NumField(); i++ {
		value := struc.Field(i)
		if !value.CanInterface() {
			continue
		}

		var (
			field = struc.Type().Field(i)
			kind  = field.Type.Kind()
			x     interface{}
		)

		switch kind {
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
			x = value.Interface()

		case reflect.Slice:
			switch value.Type().Elem().Kind() {
			case reflect.String:
				x = value.Interface()
			}

		case reflect.Struct:
			if s := sanitize(value); len(s) > 0 {
				x = s
			}

		case reflect.Ptr:
			if !value.IsNil() {
				switch value.Type().Elem().Kind() {
				case reflect.Struct:
					if s := sanitize(value.Elem()); len(s) > 0 {
						x = s
					}
				}
			}
		}

		if x != nil {
			sane = append(sane, yaml.MapItem{
				Key:   strings.ToLower(field.Name),
				Value: x,
			})
		}
	}

	return
}
