package main

import (
	"fmt"
	"github/adedaryorh/pooler_Remmitance_Application/api"
	"io/ioutil"
	"os"
)

func main() {

	fmt.Println("WELCOME TO POOLER BANK")

	// Open the banner.txt file
	file, err := os.Open("banner.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the content of the file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(string(content))

	server := api.NewServer(".")
	server.Start(3000)
}
