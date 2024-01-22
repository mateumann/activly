package main

import (
	"activly/server"

	_ "github.com/lib/pq"
)

func main() {
	server.Start()
}
