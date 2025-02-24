package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"gorden.tsmckee.com/server"
)

func main() {
	initFlag := flag.Bool("init", true, "initalize new Gorden site in CWD")
	if *initFlag {
		GordenInit()
		return
	}

}

func GordenServe() {
	app := server.Application{}
	app.Start()
}

func GordenInit() {
	fmt.Println("Initializing Gorden Site")
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(cwd)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != io.EOF {
		fmt.Println("Directory must be empty to initialize new Gorden site")
		return
	}

	os.Mkdir("src", 0755)
	os.Mkdir("static", 0755)
	os.Mkdir("templates", 0755)
}
