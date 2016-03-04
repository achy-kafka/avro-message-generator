package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
)

var (
	schema = flag.String("schema", "", "Avro schema to generate a message for, optionally pass this in as the 1st arg")
)

func main() {
	checkArgs()
	test := GenerateMessage(*schema)
	fmt.Println(test)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func checkArgs() {
	flag.Parse()

	if *schema == "" {
		if !terminal.IsTerminal(0) {
			bytes, err := ioutil.ReadAll(os.Stdin)
			must(err)
			*schema = string(bytes)
		}
		if *schema == "" {
			printUsageErrorAndExit("You must provide a schema to generate messages for")
		}
	}
}

func printUsageErrorAndExit(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", fmt.Sprintf(format, values...))
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Available command line options:")
	flag.PrintDefaults()
	os.Exit(64)
}
