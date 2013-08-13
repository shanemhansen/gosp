package main

import "os"
import "io"
import "bufio"

func main() {
	data := make([]byte, 1)
	LITERAL_OUTPUT := 0
	CODE := 1
	WAITINGFORBEGINCODE := 2
	WAITINGFORENDCODE := 3

	state := LITERAL_OUTPUT
	reader := bufio.NewReader(os.Stdin)
	os.Stdout.Write([]byte("package main\nfunc main() {\nprint(`"))
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
				os.Stdout.Write([]byte{char})
			}
		case CODE:
			if char == '%' {
				state = WAITINGFORENDCODE
				counter = 0
			} else if char == '=' && counter == 1 {
				expressionFlag = 1
				os.Stdout.Write([]byte("print("))
			} else {
				os.Stdout.Write([]byte{char})
			}
		case WAITINGFORBEGINCODE:
			if char == '%' {
				state = CODE
				counter = 0
				os.Stdout.Write([]byte("`)\n"))
			} else {
				os.Stdout.Write([]byte{'<'})
				state = LITERAL_OUTPUT
				counter = 0
			}
		case WAITINGFORENDCODE:
			if char == '>' {
				state = LITERAL_OUTPUT
				counter = 0
				if expressionFlag == 1 {
					expressionFlag = 0
					os.Stdout.Write([]byte(")"))
				}
				os.Stdout.Write([]byte("\nprint(`"))
			} else {
				os.Stdout.Write([]byte{'%'})
				state = CODE
				counter = 0
			}
		}
	}
	os.Stdout.Write([]byte("`)\n}\n"))

}
