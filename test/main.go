package main

import "time"

func main() {
	println("test sleep 5 second")
	time.Sleep(5 * time.Second)
	println("test exit")
}
