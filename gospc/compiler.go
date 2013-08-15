package main

import (
    "os"
    "io"
    "fmt"
    "flag"
    "bytes"
    "path"
    "regexp"
)

var packageName string
var printedMeta = false
func main() {
    ext := ".gosp"
    flag.StringVar(&packageName,"package", "template", "The package name to store your template under")
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
func compile(reader io.Reader, out io.Writer, funcName string) {
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
var _ fmt.Stringer
`, packageName)
    if (!printedMeta) {
        printedMeta = true
        out.Write([]byte(`type Template func(io.Writer, Template)
`))
    }
    fmt.Fprintf(out,
        `func %s(output io.Writer, content Template) {
    output.Write([]byte(%s`,
        funcName,
        "`",
    )
	counter := 0
	expressionFlag :=0
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
	out.Write([]byte("`))\n}\n"))

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
