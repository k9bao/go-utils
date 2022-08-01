package main

import "fmt"
import "go-utils/src/fs"
import "go-utils/src/hash"
import "go-utils/src/httputil"

func main(){
	fmt.Println("hello world")
	fmt.Println(fs.GetParentAbsDir("."))
	fmt.Println(hash.MD5("12345"))
	fmt.Println(httputil.GetDefaultHeader())
}