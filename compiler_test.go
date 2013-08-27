package gosp

import "os"
import "bytes"
import "testing"

func TestTemplate(t *testing.T) {
	tests := []string{
		"",
		"Hello",
		`Hello <% output.Write([]byte("world")) %>`,
		`<title><% for i := 0;i<10;i++ { %>
          yo dawg <%= i%>
         <% } %>`,
	}
	for _, test := range tests {
		var actual bytes.Buffer
		tmpl := bytes.NewBuffer([]byte(test))
		Compile(tmpl, &actual, "TESTING", "TESTING")
	}
}
func ExampleHello() {
	input := bytes.NewBuffer([]byte(`<% println("hello")%>`))
	Compile(input, os.Stdout, "main", "main")
	/* Output:package main
	import "fmt"
	import "io"
	func main() (func(io.Writer)) {
	return func(output io.Writer) {
	    output.Write([]byte(``))
	 println("hello")
	output.Write([]byte(``))
	}
	}
	*/
}
