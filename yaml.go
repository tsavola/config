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

func Read(target interface{}, r io.Reader) error {
	return yaml.NewDecoder(r).Decode(target)
}

func ReadFile(target interface{}, filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	return Read(target, f)
}

func ReadFileIfExists(target interface{}, filename string) (err error) {
	err = ReadFile(target, filename)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func Write(w io.Writer, source interface{}) error {
	return yaml.NewEncoder(w).Encode(source)
}

func WriteFile(filename string, source interface{}) (err error) {
	data, err := yaml.Marshal(source)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}
