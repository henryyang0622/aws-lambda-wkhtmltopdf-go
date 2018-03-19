package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	ff, err := ioutil.ReadFile("base64.txt")
	if err != nil {
		fmt.Println("base64.txt not find")
		log.Fatalln(err)
	}
	var sss = string(ff)
	decodeBytes, err := base64.StdEncoding.DecodeString(sss)
	if err != nil {
		fmt.Println("base64.txt not base 64 file")
		log.Fatalln(err)
	}
	//fmt.Println(string(decodeBytes))

	tofile, _ := os.Create("base64.pdf") // 建立字檔
	tofile.WriteString(string(decodeBytes))
	fmt.Println("see base64.pdf")
	
}
