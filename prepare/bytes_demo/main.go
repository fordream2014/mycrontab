package main

import (
	"bytes"
	"fmt"
)

func main() {
	var contents string
	buffer1 := bytes.NewBufferString(contents)
	fmt.Printf("the length of new buffer with contents %q: %d\n", contents, buffer1.Len())
	fmt.Printf("the capacity of new buffer with contents %q: %d\n", contents, buffer1.Cap())
	fmt.Println()
}

func demo01() {
	var buffer1 bytes.Buffer
	contents := "hello world"
	buffer1.WriteString(contents)
	fmt.Printf("the length of buffer: %d \n", buffer1.Len())
	fmt.Printf("the capacity of buffer: %d \n", buffer1.Cap())
	fmt.Println()

	p1 := make([]byte, 7)
	n,_ := buffer1.Read(p1)
	fmt.Printf("%d bytes were read. \n", n)
	fmt.Printf("the length of buffer: %d\n", buffer1.Len())
	fmt.Printf("the capacity of buffer: %d\n", buffer1.Cap())
}
