package main

import (
	"flag"
	"os"
   	"path"
    "bytes"
	"regexp"
    "github.com/shanemhansen/gosp"
)



func main() {
	ext := ".gosp"
	flag.Parse()
	args := flag.Args()
	for _, fname := range args {
		input, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		funcName := camelCase(path.Base(fname[:len(fname)-len(ext)]))
		ofname := path.Base(fname)
		ofname = ofname[:len(ofname)-2]
        packageName := path.Base(path.Dir(fname))
		output, err := os.Create(path.Join(path.Dir(fname), ofname))
		if err != nil {
			panic(err)
		}
		gosp.Compile(input, output, funcName, packageName)
	}
}
func camelCase(src string) string {
	camelingRegex := regexp.MustCompile("[0-9A-Za-z]+")
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		chunks[idx] = bytes.Title(val)
	}
	return string(bytes.Join(chunks, nil))
}
