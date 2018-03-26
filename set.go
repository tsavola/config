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

// SetFromString sets a field of the configuration object.  The value
// representation is parsed according to the type of the field.
//
// Valid boolean representations are "true", "false", "yes", "no", "y", "n",
// "on", and "off".
//
// If the representation of a string list is the empty string, the field will
// be an empty list.  If the representation starts with "[", it is assumed to
// be a JSON-encoded string array.  Otherwise the length will be one, and repr
// will be the single item.
func SetFromString(config interface{}, path string, repr string) (err error) {
	defer func() {
		err = asError(recover())
	}()

	MustSetFromString(config, path, repr)
	return
}

// MustSetFromString sets a field of the configuration object.  The value
// representation is parsed according to the type of the field.  Panic if the
// field doesn't exist or parsing fails.
//
// See SetFromString for parsing rules.
func MustSetFromString(config interface{}, path string, repr string) {
	node := lookup(config, path)

	switch node.Kind() {
	case reflect.Bool:
		setBoolFromString(node, repr)

	case reflect.Int:
		setIntFromString(node, repr, intBitSize)

	case reflect.Int8:
		setIntFromString(node, repr, 8)

	case reflect.Int16:
		setIntFromString(node, repr, 16)

	case reflect.Int32:
		setIntFromString(node, repr, 32)

	case reflect.Int64:
		if node.Type() == durationType {
			d, err := time.ParseDuration(repr)
			if err != nil {
				panic(err)
			}
			node.SetInt(int64(d))
		} else {
			setIntFromString(node, repr, 64)
		}

	case reflect.Uint:
		setUintFromString(node, repr, intBitSize)

	case reflect.Uint8:
		setUintFromString(node, repr, 8)

	case reflect.Uint16:
		setUintFromString(node, repr, 16)

	case reflect.Uint32:
		setUintFromString(node, repr, 32)

	case reflect.Uint64:
		setUintFromString(node, repr, 64)

	case reflect.Float32:
		setFloatFromString(node, repr, 32)

	case reflect.Float64:
		setFloatFromString(node, repr, 64)

	case reflect.String:
		node.SetString(repr)

	case reflect.Slice:
		if node.Type().Elem().Kind() == reflect.String {
			setSliceFromString(node, repr)
			break
		}
		fallthrough
	default:
		panic(fmt.Errorf("unsupported field type: %s", node.Type()))
	}
}

func setBoolFromString(node reflect.Value, repr string) {
	switch strings.ToLower(repr) {
	case "false", "no", "n", "off":
		node.SetBool(false)

	case "true", "yes", "y", "on":
		node.SetBool(true)

	default:
		panic(fmt.Errorf("invalid boolean string: %q", repr))
	}
}

func setIntFromString(node reflect.Value, repr string, bitSize int) {
	i, err := strconv.ParseInt(repr, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetInt(i)
}

func setUintFromString(node reflect.Value, repr string, bitSize int) {
	i, err := strconv.ParseUint(repr, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetUint(i)
}

func setFloatFromString(node reflect.Value, repr string, bitSize int) {
	f, err := strconv.ParseFloat(repr, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetFloat(f)
}

func setSliceFromString(node reflect.Value, repr string) {
	var slice []string

	switch {
	case repr == "":
		// ok

	case strings.HasPrefix(repr, "["):
		if err := json.Unmarshal([]byte(repr), &slice); err != nil {
			panic(err)
		}

	default:
		slice = []string{repr}
	}

	node.Set(reflect.ValueOf(slice))
}

// Assign a value to a field of the configuration object.  The field's path and
// string representation are parsed from an expression of the form "path=repr".
//
// See SetFromString for parsing rules.
func Assign(config interface{}, expr string) (err error) {
	defer func() {
		err = asError(recover())
	}()

	MustAssign(config, expr)
	return
}

// Assign a value to a field of the configuration object.  The field's path and
// string representation are parsed from an expression of the form "path=repr".
// Panic if the field doesn't exist or parsing fails.
//
// See SetFromString for parsing rules.
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
