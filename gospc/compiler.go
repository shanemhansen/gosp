package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
    "bufio"
	"path"
	"regexp"
    "strings"
)

var packageName string
var printedMeta = false
type KeyValue [2]string
type Directive struct {
    Imports []string
    Params []KeyValue
}


func main() {
	ext := ".gosp"
	flag.StringVar(&packageName, "package", "template", "The package name to store your template under")
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
		output, err := os.Create(path.Join(path.Dir(fname), ofname))
		if err != nil {
			panic(err)
		}
		compile(input, output, funcName)
	}
}
func compile(in io.Reader, out io.Writer, funcName string) {
    reader  := bufio.NewReader(in)
    directive := Directive{}
    //process all directives at beginning of file
    for {
        //check for line beginnig with '@'
        peekaboo, err := reader.Peek(1)
        if err != nil {
            panic(err)
        }
        if peekaboo[0] != '@' {
            break
        }
        //found line, read it in and process it.
        line ,err := reader.ReadString('\n')
        if err != nil {
            panic(err)
        }
        if strings.HasPrefix(line, "@import") {
            directive.Imports = append(directive.Imports, line[1:])
            continue
        }
        if strings.HasPrefix(line, "@(") {
            line = line[2:len(line)-2]
            for _, parameter := range strings.Split(line, ",") {
                parameter := strings.Trim(parameter, " ")
                kv := strings.Split(parameter, " ")
                param := KeyValue{kv[0], kv[1]}
                directive.Params = append(directive.Params, param)
            }
            continue
        }
    }
	data := make([]byte, 1)

	LITERAL_OUTPUT := 0
	CODE := 1
	WAITINGFORBEGINCODE := 2
	WAITINGFORENDCODE := 3
    
	state := LITERAL_OUTPUT
	fmt.Fprintf(out,
		`package %s
import "fmt"
import "io"
`, packageName)
    for _, pkg := range directive.Imports {
        fmt.Fprintf(out, "%s\n", pkg)
    }
	if !printedMeta {
		printedMeta = true
		out.Write([]byte(`
var _ fmt.Stringer
type Template func(io.Writer)
`))
	}
    var params string
    for _, param := range directive.Params {
        params += "," + param[0] + " " + param[1]
    }
	fmt.Fprintf(out,
		`func %s(content Template%s) (func(io.Writer)) {
return func(output io.Writer) {
    output.Write([]byte(%s`,
		funcName,
        params,
		"`",
	)
	counter := 0
	expressionFlag := 0
OUTER:
	for {
		n, err := reader.Read(data)
		switch err {
		case io.EOF:
			break OUTER
		case nil:
			break
		default:
			panic(err)
		}
		if n != len(data) {
			panic("short read")
		}
		char := data[0]
		counter++
		switch state {
		case LITERAL_OUTPUT:
			if char == '<' {
				state = WAITINGFORBEGINCODE
				counter = 0
			} else {
				out.Write([]byte{char})
			}
		case CODE:
			if char == '%' {
				state = WAITINGFORENDCODE
				counter = 0
			} else if char == '=' && counter == 1 {
				expressionFlag = 1
				out.Write([]byte(`fmt.Fprintf(output, "%v",`))
			} else {
				out.Write([]byte{char})
			}
		case WAITINGFORBEGINCODE:
			if char == '%' {
				state = CODE
				counter = 0
				out.Write([]byte("`))\n"))
			} else {
				out.Write([]byte{'<', char})
				state = LITERAL_OUTPUT
				counter = 0
			}
		case WAITINGFORENDCODE:
			if char == '>' {
				state = LITERAL_OUTPUT
				counter = 0
				if expressionFlag == 1 {
					expressionFlag = 0
					out.Write([]byte(")"))
				}
				out.Write([]byte("\noutput.Write([]byte(`"))
			} else {
				out.Write([]byte{'%'})
				state = CODE
				counter = 0
			}
		}
	}
	out.Write([]byte("`))\n}\n}\n"))

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
