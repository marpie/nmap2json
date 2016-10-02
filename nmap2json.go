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

	nmap "github.com/marpie/go-nmap"
	s2m "github.com/marpie/struct2elasticMapping"
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

// writeMapping writes a ElasticSearch-Mapping of NmapRun to
// the specified file.
func writeMapping(filename string) error {
	name, mapping, err := s2m.Analyze(nmap.NmapRun{}, "json")
	if err != nil {
		return err
	}
	data, err := s2m.MappingAsJson(name, mapping)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func main() {
	flag.Usage = usage
	outdir := flag.String("outdir", "./", "Output Directory.")
	pretty := flag.Bool("prettify", true, "Prettify JSON Output.")
	mapping := flag.String("write-mapping", "", "Write an ElasticSearch-Mapping to the specified file.")
	flag.Parse()

	if len(*mapping) > 0 {
		fmt.Println("[*] Writing Mapping to file: " + *mapping)
		if err := writeMapping(*mapping); err != nil {
			PrintError(err)
			os.Exit(1)
		}
	}

	for _, filename := range flag.Args() {
		ext := filepath.Ext(filename)
		basename := filepath.Base(filename)
		outpath := filepath.Join(*outdir, basename[0:len(basename)-len(ext)]+".json")

		fmt.Println("[*] Parsing: " + filename)
		scan, err := parseScan(filename)
		if err != nil {
			PrintError(err)
			continue
		}

		if scan.Scanner != "nmap" {
			ErrorOut("Not a NMap XML file!\n")
			continue
		}

		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(scan, "", "\t")
		} else {
			b, err = json.Marshal(scan)
		}
		if err != nil {
			PrintError(err)
			continue
		}

		// Write to JSON file
		fmt.Println("[+] Writing file to: " + outpath)
		if err := ioutil.WriteFile(outpath, b, 0644); err != nil {
			PrintError(err)
		}
	}
}
