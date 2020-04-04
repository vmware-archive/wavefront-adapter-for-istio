package main

import "fmt"

import (
	"log"
	"net"
	"os"
)

func main() {

	// Harcoded credentials
	username := "admin"
	password := "f62e5bcda4fae4f82370da0c6f20697b8f8447ef"
	fmt.Println("Doing something with: ", username, password)

	// SampleCodeG102 code snippet for network binding
	l, err := net.Listen("tcp", "0.0.0.0:2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// Changing file permissions of important files
	os.Chmod("/etc/passwd", 0777)
	os.Chmod("~/.bashrc", 0777)

}
