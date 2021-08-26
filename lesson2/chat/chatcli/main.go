package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	nick := ""
	fmt.Print("Введите nickname: ")
	fmt.Scan(&nick)

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "%s", nick)

	defer conn.Close()

	go func() {
		io.Copy(os.Stdout, conn)
	}()

	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
