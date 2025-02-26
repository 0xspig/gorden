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
	fmt.Printf("init: %t\nserve:%t\ndraft:%t\nname:%s\ndir:%s\n", *initFlag, *serveFlag, *draftFlag, flag.Arg(0), *dirFlag)
	if *initFlag {
		GordenInit(*dirFlag, flag.Arg(0))
		return
	}
	app := server.Application{}
	if *serveFlag {
		app.Init(*dirFlag, *draftFlag)
		app.Start(*addrFlag)
	}

}

func GordenInit(dir string, siteName string) {
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		os.Chdir(dir)
	}
	fmt.Printf("Initializing Gorden Site in %s\n", dir)
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != io.EOF {
		fmt.Println("FAILED: Directory must be empty to initialize new Gorden site")
		return
	} else {
		fmt.Println("SUCCESS: Gorden site initialized")
	}

	os.Mkdir("src", 0755)
	os.Mkdir("static", 0755)
	os.Mkdir("static/gen", 0755)
	os.Mkdir("templates", 0755)
	os.Mkdir("content", 0755)

	templates := os.DirFS("/usr/local/share/gorden/templates")
	os.CopyFS("templates", templates)

	defaults := os.DirFS("/usr/local/share/gorden/def/")
	os.CopyFS(".", defaults)

}
