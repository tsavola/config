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
	return yaml.NewEncoder(w).Encode(config)
}

// Write the configuration to a YAML file.
func WriteFile(filename string, config interface{}) (err error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}
