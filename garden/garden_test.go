package garden

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	g := CreateGarden()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	home = filepath.Join(home, "prg/tsm-go/")
	os.Chdir(home)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(wd)
	entries, err := os.ReadDir(wd)
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Println(e)
	}
	// TODO make content dir in config or something to search files in
	// for now Im just going to hack in static
	g.PopulateGardenFromDir("ui/content")
	g.ParseAllConnections()
	g.ConnectNodes("peepee.md", "firstFile.md")

	g.GenAssets()

	fmt.Printf("distance between home and blog: %d\n", g.shortestPath("home.md", "blog"))
	fmt.Printf("distance between hips and honey: %d\n", g.shortestPath("hips.md", "honey-hollow.md"))
	center := g.findCenter()
	fmt.Printf("Center of graph: %s\n %d", center, len(center))

	fmt.Printf("peepee\n")
}
