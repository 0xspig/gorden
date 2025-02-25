package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"gorden.tsmckee.com/server"
)

func main() {
	initFlag := flag.Bool("init", false, "initalize new Gorden site in CWD")
	serveFlag := flag.Bool("serve", false, "Serve local gorden")
	draftFlag := flag.Bool("drafts", false, "enable draft rendering")
	addrFlag := flag.String("addr", "localhost:3000", "local server address (defaults to localhost:3000)")
	dirFlag := flag.String("dir", "", "optional base dir (default is CWD)")
	flag.Parse()
	fmt.Printf("init: %t\nserve:%t\ndraft:%t\n", *initFlag, *serveFlag, *draftFlag)
	if *initFlag {
		GordenInit(*dirFlag)
		return
	}
	app := server.Application{}
	if *serveFlag {
		app.Init(*dirFlag, *draftFlag)
		app.Start(*addrFlag)
	}

}

func GordenInit(dir string) {
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Initializing Gorden Site in %s", dir)
	f, err := os.Open(dir)
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
