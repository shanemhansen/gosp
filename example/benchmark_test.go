package example

import "testing"
import "text/template"
import "io/ioutil"

type Person struct {
    Name string
    Age int
    Metadata []byte
    Childen []*Person
}
func (this *Person) Everyone() (people []Person) {
    queue := make([]*Person, 1)
    queue[0] = this
    for {
        if len(queue) == 0 {
            return
        }
        person := queue[len(queue)-1]
        people = append(people, *person)
        queue = queue[:len(queue)-1]
        for _, kid := range person.Childen {
            queue = append(queue, kid)
        }
    }
    return
}

func BenchmarkRendering(b *testing.B) {
    family := Person{
        Name : "Parent",
        Age: 28,
        Metadata: []byte("some binary data"),
        Childen: []*Person{
            &Person{Name:"child1", Age: 5},
            &Person{Name:"child2", Age: 6},
            &Person{Name:"child3", Age: 7},
        },
    }
    tmpl := template.Must(template.New("test").Funcs(map[string]interface{}{"People": func (t Person)[]Person { return t.Everyone() }}).Parse(`
{{range $i, $child := .|People}}
{{$child}}
{{end}}
`))    
    
    for i:=0; i< b.N; i++ {
        err := tmpl.Execute(ioutil.Discard, family)
        if err != nil { panic(err) }
    }
}
func BenchmarkRenderingGOSP(b *testing.B) {
    family := Person{
        Name : "Parent",
        Age: 28,
        Metadata: []byte("some binary data"),
        Childen: []*Person{
            &Person{Name:"child1", Age: 5},
            &Person{Name:"child2", Age: 6},
            &Person{Name:"child3", Age: 7},
        },
    }
    for i:=0; i< b.N; i++ {
        PeopleTest(family)(ioutil.Discard)
    }
}
