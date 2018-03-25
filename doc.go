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

		Audio struct {
			Enabled    bool
			SampleRate int
		}
	}

	func main() {
		c := new(Config)
		c.Size.Width = 640
		c.Size.Height = 480
		c.Audio.SampleRate = 44100

		if err := config.ReadFileIfExists(c, "defaults.yaml"); err != nil {
			log.Print(err)
		}

		if config.GetPanic(c, "audio.samplerate").(int) <= 0 {
			config.Set(c, "audio.enabled", false)
		}

		dump := flag.Bool("dump", false, "create defaults.yaml")
		flag.Var(config.Loader(c), "f", "read config from YAML files")
		flag.Var(config.Setter(c), "c", "set config keys (path.to.key=value)")
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

*/
package config
