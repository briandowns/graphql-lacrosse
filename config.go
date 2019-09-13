package main

// db
type db struct {
	user string
	pass string
	host string
	port string
}

// config
type config struct {
	db *db
}
