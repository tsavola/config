// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var intBitSize = int(unsafe.Sizeof(int(0)) * 8)

func Set(target interface{}, path string, value interface{}) (err error) {
	defer func() {
		err = settingError(path, recover())
	}()

	SetPanic(target, path, value)
	return
}

func SetPanic(target interface{}, path string, value interface{}) {
	lookup(target, path).Set(reflect.ValueOf(value))
}

func SetString(target interface{}, path string, value string) (err error) {
	defer func() {
		err = settingError(path, recover())
	}()

	SetStringPanic(target, path, value)
	return
}

func SetStringPanic(target interface{}, path string, value string) {
	node := lookup(target, path)
	kind := node.Type().Kind()

	switch kind {
	case reflect.Bool:
		setStringBool(node, value)

	case reflect.Int:
		setStringInt(node, value, intBitSize)

	case reflect.Int32:
		setStringInt(node, value, 32)

	case reflect.Int64:
		setStringInt(node, value, 64)

	case reflect.Uint:
		setStringUint(node, value, intBitSize)

	case reflect.Uint32:
		setStringUint(node, value, 32)

	case reflect.Uint64:
		setStringUint(node, value, 64)

	case reflect.Float32:
		setStringFloat(node, value, 32)

	case reflect.Float64:
		setStringFloat(node, value, 64)

	case reflect.String:
		node.SetString(value)

	default:
		panic(fmt.Errorf("unsupported field type: %s", kind))
	}
}

func setStringBool(node reflect.Value, value string) {
	switch strings.ToLower(value) {
	case "false", "no", "n", "off":
		node.SetBool(false)

	case "true", "yes", "y", "on":
		node.SetBool(true)

	default:
		panic(fmt.Errorf("invalid boolean string: %q", value))
	}
}

func setStringInt(node reflect.Value, value string, bitSize int) {
	i, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetInt(i)
}

func setStringUint(node reflect.Value, value string, bitSize int) {
	i, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetUint(i)
}

func setStringFloat(node reflect.Value, value string, bitSize int) {
	f, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		panic(err)
	}
	node.SetFloat(f)
}

func SetExpr(target interface{}, expr string) (err error) {
	defer func() {
		err = settingError(expr, recover())
	}()

	SetExprPanic(target, expr)
	return
}

func SetExprPanic(target interface{}, expr string) {
	tokens := strings.SplitN(expr, "=", 2)
	if len(tokens) != 2 {
		panic(fmt.Errorf("invalid expression: %q", expr))
	}

	SetStringPanic(target, tokens[0], tokens[1])
}

func Get(target interface{}, path string) (value interface{}, err error) {
	defer func() {
		err = settingError(path, recover())
	}()

	value = GetPanic(target, path)
	return
}

func GetPanic(target interface{}, path string) interface{} {
	return lookup(target, path).Interface()
}

func lookup(target interface{}, path string) (node reflect.Value) {
	node = reflect.ValueOf(target)

	for _, nodeName := range strings.Split(path, ".") {
		if node.Type().Kind() == reflect.Ptr {
			node = reflect.Indirect(node)
		}

		node = node.FieldByNameFunc(func(fieldName string) bool {
			return strings.ToLower(fieldName) == nodeName
		})
	}

	return
}

func settingError(s string, x interface{}) (err error) {
	if x != nil {
		err = fmt.Errorf("config: %q: %v", s, x)
	}
	return
}
