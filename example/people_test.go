package example
import "fmt"
import "io"

var _ fmt.Stringer
type Template func(io.Writer)
func PeopleTest(family Person) (func(io.Writer)) {
return func(output io.Writer) {
    output.Write([]byte(``))
 for _, child := range family.Everyone() {
output.Write([]byte(`
`))
fmt.Fprintf(output, "%v", child.Name )
output.Write([]byte(` is `))
fmt.Fprintf(output, "%v",child.Age)
output.Write([]byte(`
`))
 } 
output.Write([]byte(`
`))
}
}
