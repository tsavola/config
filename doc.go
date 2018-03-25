// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

Package config is an ergonomic configuration parsing toolkit.  The schema is
declared using a struct type, and values can be read from YAML files or set via
command-line flags.

A pointer to a preallocated configuration object of a user-defined struct type
must be passed to all functions.  The type can have an arbitrary number of
nested structs (either embedded or through an initialized pointer).  Only
exported fields can be used.  The object can be initialized with default
values.

The field names are spelled in lower case in YAML files and on the
command-line.  The accessor functions and flag values use dotted paths to
identify the field, such as "audio.samplerate".

Supported field types are bool, int, int8, int16, int32, int64, uint, uint8,
uint16, uint32, uint64, float32, float64, string, and time.Duration.

The Get method is provided for completeness; the intended way to access
configuration values is through direct struct field access.

Short example:

	c := &myConfig{}

	flag.Usage = config.FlagUsage(c)
	flag.Var(config.FileReader(c), "f", "read config from YAML files")
	flag.Var(config.Assigner(c), "c", "set config keys (path.to.key=value)")
	flag.Parse()

Longer example:

	package main

	import (
		"flag"
		"fmt"
		"log"

		"github.com/tsavola/config"
	)

	type myConfig struct {
		Comment string

		Size struct {
			Width  uint32
			Height uint32
		}

		Audio struct {
			Enabled    bool
			SampleRate int
		}
	}

	func main() {
		c := new(myConfig)
		c.Size.Width = 640
		c.Size.Height = 480
		c.Audio.SampleRate = 44100

		if err := config.ReadFileIfExists("defaults.yaml", c); err != nil {
			log.Print(err)
		}

		if x, _ := config.Get(c, "audio.samplerate"); x.(int) <= 0 {
			config.MustSet(c, "audio.enabled", false)
		}

		dump := flag.Bool("dump", false, "create defaults.yaml")
		flag.Var(config.FileReader(c), "f", "read config from YAML files")
		flag.Var(config.Assigner(c), "c", "set config keys (path.to.key=value)")
		flag.Usage = config.FlagUsage(c)
		flag.Parse()

		if *dump {
			if err := config.WriteFile("defaults.yaml", c); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("Comment is %q\n", c.Comment)
		fmt.Printf("Size is %dx%d\n", c.Size.Width, c.Size.Height)
		if c.Audio.Enabled {
			fmt.Printf("Sample rate is %d\n", c.Audio.SampleRate)
		}
	}

Example usage output:

	$ example -help
	Usage of example:
	  -c value
	    	set config keys (path.to.key=value)
	  -dump
	    	create defaults.yaml
	  -f value
	    	read config from YAML files

	Configuration settings:
	  comment string
	  size.width uint32
	  size.height uint32
	  audio.enabled bool
	  audio.samplerate int

*/
package config
