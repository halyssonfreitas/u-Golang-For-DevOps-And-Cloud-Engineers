// buid : go build MySlowReader
// execute 1 :
//   - ./MySlowReader
package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type MySlowReader struct {
	Contents string
	Position int
}

func (m *MySlowReader) Read(p []byte) (n int, err error) {
	if m.Position >= len(m.Contents) {
		return 0, io.EOF
	}
	p[0] = m.Contents[m.Position]
	m.Position = m.Position + 1
	return 1, nil
}

func main() {

	mySlowReaderIntance := &MySlowReader{
		Contents: "hello world"}

	out, err := io.ReadAll(mySlowReaderIntance)

	if err != nil && strings.Compare("EOF", err.Error()) != 0 {
		log.Println("Aqui")
		log.Fatal(err)

	}

	fmt.Printf("output: %s\n", out)
}
