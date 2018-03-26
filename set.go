// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var intBitSize = int(unsafe.Sizeof(int(0)) * 8)
var durationType = reflect.TypeOf(time.Second)

// Set a field of the configuration object.  The value must have the same type
// as the field.
func Set(config interface{}, path string, value interface{}) (err error) {
	defer func() {
		err = asError(recover())
	}()

	MustSet(config, path, value)
	return
}

// MustSet a field of the configuration object.  The value must have the same
// type as the field.  Panic if the field doesn't exist or the types don't
// match.
func MustSet(config interface{}, path string, value interface{}) {
	lookup(config, path).Set(reflect.ValueOf(value))
}

// SetFromString sets a field of the configuration object.  The value is parsed
// according to the type of the field.
func SetFromString(config interface{}, path string, value string) (err error) {
	defer func() {
		err = asError(recover())
	}()

	MustSetFromString(config, path, value)
	return
}

// MustSetFromString sets a field of the configuration object.  The value is
// parsed according to the type of the field.  Panic if the field doesn't exist
// or parsing fails.
func MustSetFromString(config interface{}, path string, value string) {
	node := lookup(config, path)

	switch node.Kind() {
	case reflect.Bool:
		setBoolFromString(node, value)

	case reflect.Int:
		setIntFromString(node, value, intBitSize)

	case reflect.Int8:
		setIntFromString(node, value, 8)

	case reflect.Int16:
		setIntFromString(node, value, 16)

	case reflect.Int32:
		setIntFromString(node, value, 32)

	case reflect.Int64:
		if node.Type() == durationType {
			d, err := time.ParseDuration(value)
			if err != nil {
				panic(err)
			}
			node.SetInt(int64(d))
		} else {
			setIntFromString(node, value, 64)
		}

	case reflect.Uint:
		setUintFromString(node, value, intBitSize)

	case reflect.Uint8:
		setUintFromString(node, value, 8)

	case reflect.Uint16:
		setUintFromString(node, value, 16)

	case reflect.Uint32:
		setUintFromString(node, value, 32)

	case reflect.Uint64:
		setUintFromString(node, value, 64)

	case reflect.Float32:
		setFloatFromString(node, value, 32)

	case reflect.Float64:
		setFloatFromString(node, value, 64)

	case reflect.String:
		node.SetString(value)

	case reflect.Slice:
		if node.Type().Elem().Kind() == reflect.String {
			setSliceFromString(node, value)
			break
		}
		fallthrough
	default:
		panic(fmt.Errorf("unsupported field type: %s", node.Type()))
	}
}

func setBoolFromString(node reflect.Value, value string) {
	switch strings.ToLower(value) {
	case "false", "no", "n", "off":
		node.SetBool(false)

	case "true", "yes", "y", "on":
		node.SetBool(true)

	default:
		panic(fmt.Errorf("invalid boolean string: %q", value))
	}
}

func setIntFromString(node reflect.Value, value string, bitSize int) {
	i, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetInt(i)
}

func setUintFromString(node reflect.Value, value string, bitSize int) {
	i, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetUint(i)
}

func setFloatFromString(node reflect.Value, value string, bitSize int) {
	f, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetFloat(f)
}

func setSliceFromString(node reflect.Value, value string) {
	var slice []string

	switch {
	case value == "":
		// ok

	case strings.HasPrefix(value, `[`):
		if err := json.Unmarshal([]byte(value), &slice); err != nil {
			panic(err)
		}

	default:
		slice = []string{value}
	}

	node.Set(reflect.ValueOf(slice))
}

// Assign a value to a field of the configuration object.  The field's path and
// value are parsed from an expression of the form "path=value".
func Assign(config interface{}, expr string) (err error) {
	defer func() {
		err = asError(recover())
	}()

	MustAssign(config, expr)
	return
}

// Assign a value to a field of the configuration object.  The field's path and
// value are parsed from an expression of the form "path=value".  Panic if the
// field doesn't exist or parsing fails.
func MustAssign(config interface{}, expr string) {
	tokens := strings.SplitN(expr, "=", 2)
	if len(tokens) != 2 {
		panic(fmt.Errorf("invalid assignment expression: %q", expr))
	}

	MustSetFromString(config, strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1]))
}

// Get the value of a field of the configuration object.
func Get(config interface{}, path string) (value interface{}, err error) {
	defer func() {
		err = asError(recover())
	}()

	value = lookup(config, path).Interface()
	return
}

func lookup(config interface{}, path string) (node reflect.Value) {
	node = reflect.ValueOf(config)

	for _, nodeName := range strings.Split(path, ".") {
		if node.Kind() == reflect.Ptr {
			node = node.Elem()
		}

		field, ok := node.Type().FieldByNameFunc(func(fieldName string) bool {
			return strings.ToLower(fieldName) == nodeName
		})
		if !ok {
			panic(fmt.Errorf("unknown config key: %q", path))
		}

		node = node.FieldByIndex(field.Index)
	}

	return
}

func asError(x interface{}) (err error) {
	if x != nil {
		err, _ = x.(error)
		if err == nil {
			err = fmt.Errorf("%v", x)
		}
	}
	return
}
