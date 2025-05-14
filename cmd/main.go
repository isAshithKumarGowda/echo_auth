package main

import "flag"

func main() {
	addr := flag.String("addr", ":8080", "Address at which the server will run")
	flag.Parse()

	db := NewDatabase()
	db.CheckStatus()
	defer db.Close()

	Start(db.db, addr)
}
