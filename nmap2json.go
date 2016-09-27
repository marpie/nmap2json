// Copyright 2016 Markus Pieton (marpie). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nmap2json converts a NMap XML file to a JSON file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	nmap "github.com/lair-framework/go-nmap"
)

// usage modifies the default usage message to include
// the positional arguments
func usage() {
	fmt.Printf("Usage: %s [OPTIONS] [FILE] ([FILE] [FILE]...)\n", os.Args[0])
	flag.PrintDefaults()
}

// parseScan uses the go-nmap library to parse NMap XML scans.
func parseScan(filename string) (*nmap.NmapRun, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return nmap.Parse(b)
}

func main() {
	flag.Usage = usage
	outdir := flag.String("outdir", ".", "Output Directory.")
	pretty := flag.Bool("prettify", false, "Prettify JSON Output.")
	flag.Parse()

	for _, filename := range flag.Args() {
		ext := filepath.Ext(filename)
		outpath := filepath.Join(*outdir, filename[0:len(filename)-len(ext)]+".json")

		scan, err := parseScan(filename)
		if err != nil {
			fmt.Println("[E]", err)
		}

		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(scan, "", "\t")
		} else {
			b, err = json.Marshal(scan)
		}
		if err != nil {
			fmt.Println("[E]", err)
		}
		// Write to JSON file
		if err := ioutil.WriteFile(outpath, b, 0644); err != nil {
			fmt.Println("[E]", err)
		}
	}
}
