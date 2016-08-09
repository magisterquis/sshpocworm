package main

/*
 * binary.go
 * Read in the binary
 * By J. Stuart McMurray
 * Created 20160808
 * Last Modified 20160808
 */

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/kardianos/osext"
)

var BINARY []byte

/* ReadBinary reads this binary into its own memory. */
func ReadBinary() {
	/* Path to binary */
	path, err := osext.Executable()
	if nil != err {
		log.Fatalf("Unable to get binary path: %v", err)
	}

	/* Slurp it */
	BINARY, err = ioutil.ReadFile(path)
	if nil != err {
		log.Fatalf("Unable to read binary: %v", err)
	}

	log.Printf("Read binary from %v", path)

	/* Remove it */
	if err := os.Remove(path); nil != err {
		log.Printf("Unable to remove %v: %v", path, err)
	} else {
		log.Printf("Removed %v", path)
	}
}
