// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Read YAML into the object pointed to by target.
func Read(r io.Reader, target interface{}) error {
	return yaml.NewDecoder(r).Decode(target)
}

// Read a YAML file into the object pointed to by target.
func ReadFile(filename string, target interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	return Read(f, target)
}

// Read a YAML file into the object pointed to by target.  No error is returned
// if the file doesn't exist,
func ReadFileIfExists(filename string, target interface{}) (err error) {
	err = ReadFile(filename, target)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

// Write the user-defined object as YAML.
func Write(w io.Writer, source interface{}) error {
	return yaml.NewEncoder(w).Encode(source)
}

// Write the user-defined object to a YAML file.
func WriteFile(filename string, source interface{}) (err error) {
	data, err := yaml.Marshal(source)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}
