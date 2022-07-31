package main

import "fmt"
import "go-utils/src/fs"

func main(){
	fmt.Println("hello world")
	fmt.Println(fs.GetParentAbsDir("."))
}