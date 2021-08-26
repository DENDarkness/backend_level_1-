package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	fmt.Println(genExp())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	nick := make([]byte, 256)

	conn.Read(nick)

	go clientWriter(conn, ch)

	//who := conn.RemoteAddr().String()
	who := string(nick)
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func genExp() (string, float64) {
	e := []string{"+", "-", "*", "/"}

	rand.Seed(time.Now().UnixNano())
	one := rand.Intn(10)
	two := rand.Intn(10)
	exp := e[rand.Intn(4)]
	switch exp {
	case "+":
		return fmt.Sprintf("%d + %d = ?", one, two), float64(one + two)
	case "-":
		return fmt.Sprintf("%d - %d = ?", one, two), float64(one - two)
	case "*":
		return fmt.Sprintf("%d * %d = ?", one, two), float64(one * two)
	case "/":
		return fmt.Sprintf("%d / %d = ?", one, two), float64(one / two)
	}
	return "", 0
}
