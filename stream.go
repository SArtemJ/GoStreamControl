package main

import "github.com/SArtemJ/GoStreamControlAPI/libstream"

func main() {
	app := libstream.NewApplication()
	app.Init()
	app.Run()
}