// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const default_sleeptime = 1 * time.Second

var runtimes = map[string]time.Duration{
	"nehe03": 5 * time.Second,
}

func runExample(t *testing.T, path string, files []string) {
	println(strings.Repeat("=", 80))

	get := exec.Command("go", "get", "-v", "./"+path)
	get.Stdout = os.Stdout
	get.Stderr = os.Stderr
	err := get.Run()
	if err != nil {
		panic(err)
	}

	bin_name := filepath.Join("bin", path)
	bld := exec.Command("go", "build", "-v", "-o", bin_name, "./"+path)
	bld.Stdout = os.Stdout
	bld.Stderr = os.Stderr
	err = bld.Run()
	if err != nil {
		panic(err)
	}

	println(strings.Repeat("-", 80))

	cmd := exec.Command(bin_name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	done := false

	go func() {
		sleeptime, ok := runtimes[path]
		if !ok {
			sleeptime = default_sleeptime
		}
		time.Sleep(sleeptime)
		done = true
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}()

	err = cmd.Run()

	// If the done flag is true, then we made it through five seconds of runtime
	if !done && err != nil {
		//panic(err)
		t.Fatal("Process died unexpectedly: ", err)
	}
}

// Runs all examples in turn using "go build" and "dirname/dirname".
// They run in an arbitrary order.
func TestExamples(t *testing.T) {
	example_files := map[string][]string{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		dir := filepath.Dir(path)
		if dir == "." {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			example_files[dir] = append(example_files[dir], path)
		}
		return err
	})
	for k := range example_files {
		runExample(t, k, example_files[k])
	}
}