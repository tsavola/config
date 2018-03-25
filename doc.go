// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Package config is an ergonomic configuration parsing toolkit.  The schema is
declared using a struct type, and values can be read from YAML files or set via
command-line flags.

Example:

	package main

	import (
		"flag"
		"fmt"
		"log"

		"github.com/tsavola/config"
	)

	type Config struct {
		Comment string

		Size struct {
			Width  uint32
			Height uint32
		}

		Sampling   bool
		SampleRate int
	}

	func main() {
		c := new(Config)
		c.SampleRate = 44100
		c.Size.Width = 640
		c.Size.Height = 480

		if err := config.ReadFileIfExists("defaults.yaml"); err != nil {
			log.Print(err)
		}

		flag.Var(config.Loader(c), "f", "read config from YAML files")
		flag.Var(config.Setter(c), "c", "set config keys (path.to.key=value)")
		flag.Parse()

		fmt.Printf("Comment is %q", c.Comment)
		fmt.Printf("Size is %dx%d", c.Size.Width, c.Size.Height)
		if c.Sampling {
			fmt.Printf("Sampling rate is %d", c.SampleRate)
		}
	}

*/
package config
