package main

import (
	"fmt"
	"strings"
)

func main() {
	//builderUsage()
	readerUsage()
}
func builderUsage() {
	var builder1 strings.Builder
	builder1.WriteString("hello world")
	fmt.Printf("the first output(%d):\n%q\n", builder1.Len(), builder1.String())
	fmt.Println()
	//Grow
	builder1.Grow(10)
	fmt.Printf("the second output(%d): \n\"%s\"\n", builder1.Len(), builder1.String())

	builder1.Reset()
	fmt.Printf("the third output(%d):\n%q\n", builder1.Len(), builder1.String())
	fmt.Println()
}

func readerUsage() {
	reader1 := strings.NewReader("new reader")
	fmt.Printf("the size ofi reader: %d \n", reader1.Size())
	//Len方法返回未读取到的byte数
	fmt.Printf("the reading index in reader: %d \n", reader1.Size()-int64(reader1.Len()))
	fmt.Println()

	buf1 := make([]byte, 3)
	n,_ := reader1.Read(buf1)
	fmt.Printf("%d bytes were read. \n", n)
	fmt.Printf("has read bytes: %s\n", string(buf1))
	fmt.Printf("the reading index in reader: %d \n", reader1.Size() - int64(reader1.Len()))
}























