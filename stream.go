package main

import "github.com/SArtemJ/GoStreamControlAPI/libstream"

func main() {
	app := libstream.NewApplication()
	app.Configure()
	app.Run()
}