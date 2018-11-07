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
	return yaml.NewEncoder(w).Encode(sanitizeContainer(nil, reflect.ValueOf(config).Elem()))
}

// Write the configuration to a YAML file.
func WriteFile(filename string, config interface{}) (err error) {
	data, err := yaml.Marshal(sanitizeContainer(nil, reflect.ValueOf(config).Elem()))
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}

func sanitizeContainer(sane yaml.MapSlice, node reflect.Value) yaml.MapSlice {
	switch node.Kind() {
	case reflect.Map:
		return sanitizeMap(sane, node)

	case reflect.Struct:
		return sanitizeStruct(sane, node)

	default:
		panic("must be a struct or a map")
	}
}

func sanitizeMap(sane yaml.MapSlice, node reflect.Value) yaml.MapSlice {
	for _, key := range reflectMapKeyStrings(node) {
		value := node.MapIndex(reflect.ValueOf(key))

		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		if x := sanitizeValue(sane, value, false); x != nil {
			sane = append(sane, yaml.MapItem{
				Key:   key,
				Value: x,
			})
		}
	}

	return sane
}

func sanitizeStruct(sane yaml.MapSlice, node reflect.Value) yaml.MapSlice {
	for i := 0; i < node.Type().NumField(); i++ {
		value := node.Field(i)
		if !value.CanInterface() {
			continue
		}

		field := node.Type().Field(i)

		if x := sanitizeValue(sane, value, field.Anonymous); x != nil {
			sane = append(sane, yaml.MapItem{
				Key:   strings.ToLower(field.Name),
				Value: x,
			})
		}
	}

	return sane
}

func sanitizeValue(sane yaml.MapSlice, value reflect.Value, anonymous bool) (x interface{}) {
	switch value.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		x = value.Interface()

	case reflect.Slice:
		switch value.Type().Elem().Kind() {
		case reflect.String:
			x = value.Interface()
		}

	case reflect.Map:
		if s := sanitizeMap(nil, value); len(s) > 0 {
			x = s
		}

	case reflect.Ptr:
		if value.IsNil() {
			break
		}
		if value.Type().Elem().Kind() != reflect.Struct {
			break
		}
		value = value.Elem()
		fallthrough

	case reflect.Struct:
		if anonymous {
			sane = sanitizeStruct(sane, value)
		} else if s := sanitizeStruct(nil, value); len(s) > 0 {
			x = s
		}
	}

	return
}
