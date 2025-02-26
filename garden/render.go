package garden

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gorilla/feeds"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func (garden *Garden) NodeToHTML(nodeID string) []byte {
	node := garden.Masterlist[nodeID]
	if node == nil {
		return []byte("<h1>File Not Found</h1>")
	}
	switch node.Data_type {
	case CONTENT_TYPE_MARKDOWN:
		return garden.mdToHTML(node)
	case CONTENT_TYPE_HTML:
		return []byte("<h1>HTML unsupported</h1>")
	case CONTENT_TYPE_TAG:
		return garden.tagToHtml(node)
	case CONTENT_TYPE_CATEGORY:
		return garden.catToHtml(node)
	default:
		return []byte("<h1>File Not Found</h1>")
	}
}

func (garden *Garden) NodeLinksToHTML(nodeID string) []byte {
	ts := garden.Templates["links_template"]
	var buf bytes.Buffer
	err := ts.ExecuteTemplate(&buf, "links", garden.Masterlist[nodeID])
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (garden *Garden) mdToRSS(node *Node) []byte {
	source, err := os.ReadFile(node.Data_source)
	if err != nil {
		panic(err)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	internal_regex, err := regexp.Compile(`\{([^\}]*)\}\(([^\)]*)\)`)
	if err != nil {
		panic(err)
	}

	data := internal_regex.ReplaceAll(buf.Bytes(), []byte(`<a class="internal-link" href="/$2">$1</a>`))

	return data
}

func (garden *Garden) mdToHTML(node *Node) []byte {
	source, err := os.ReadFile(node.Data_source)
	if err != nil {
		panic(err)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	internal_regex, err := regexp.Compile(`\{([^\}]*)\}\(([^\)]*)\)`)
	if err != nil {
		panic(err)
	}

	data := internal_regex.ReplaceAll(buf.Bytes(), []byte(`<a class="internal-link" href="/$2" onmouseover="highlightNode('$2')" onmouseout="highlightNode('')" onClick="return targetNode('$2')">$1</a>`))

	data = append([]byte(`{{define "content"}}`), data...)
	data = append(data, []byte(`{{end}}`)...)

	switch node.Metadata.Class {
	default:
		ts, err := template.ParseFiles("templates/single.html")
		if err != nil {
			panic(err)
		}
		fmt.Println("post template")

		ts, err = ts.Parse(string(data))
		if err != nil {
			panic(err)
		}
		var template_buf bytes.Buffer
		ts.Execute(&template_buf, node)

		return template_buf.Bytes()
	case "home":
		fmt.Println("home template")
		ts, err := template.ParseFiles("templates/index.html")
		if err != nil {
			panic(err)
		}
		var template_buf bytes.Buffer
		ts.Execute(&template_buf, node.ParentGarden)

		return template_buf.Bytes()
	}

}

func (garden *Garden) tagToHtml(node *Node) []byte {
	ts := garden.Templates["list_template"]
	var buf bytes.Buffer
	ts.Execute(&buf, node)

	return buf.Bytes()
}

func (garden *Garden) catToHtml(node *Node) []byte {
	ts := garden.Templates["list_template"]
	var buf bytes.Buffer
	ts.Execute(&buf, node)

	return buf.Bytes()

}

// Generate staic assets
/*
	Currently this function is completely bloated.
	It's doing 100% of the template work and its all hardcoded.
	Ideally we should have one function that generates default templates and a second function that generates other templates.
	First, what are the necessary defaults?
	baseof - base template... obivously
	single - default template for posts

*/
func (garden *Garden) GenAssets() {
	// cache node data
	json_data, err := garden.genJSONData()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("static/gen/graph-data.json", json_data, 0644)
	if err != nil {
		panic(err)
	}

	me := &feeds.Author{Name: garden.Info.Author, Email: garden.Info.Email}

	// create Atom feed
	feed := feeds.Feed{
		Title:       garden.Info.Title,
		Link:        &feeds.Link{Href: garden.Info.Link},
		Description: garden.Info.Description,
		Author:      me,
		Created:     time.Now(),
	}

	for _, post := range garden.Masterlist {
		if post.Data_type != CONTENT_TYPE_MARKDOWN {
			continue
		}
		date, err := time.Parse(time.RFC3339, post.Metadata.Date)
		if err != nil {
			fmt.Println("error parsing date in post: " + post.ID)
			date = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
		}
		item := &feeds.Item{
			Title:   post.Metadata.Title,
			Link:    &feeds.Link{Href: garden.Info.Link + post.ID},
			Author:  me,
			Updated: date,
			Content: string(garden.mdToRSS(post)),
		}
		feed.Add(item)
	}

	xml, err := os.Create("static/index.xml")
	if err != nil {
		panic(err)
	}
	feed.WriteAtom(xml)

	// parse templates
	partials, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		panic(err)
	}
	base_files := []string{
		"./templates/baseof.html",
	}
	base_files = append(base_files, partials...)
	base_template, err := template.ParseFiles(base_files...)
	if err != nil {
		panic(err)
	}
	garden.Templates["base_template"] = base_template

	list_template, err := template.ParseFiles("templates/list.html")
	if err != nil {
		panic(err)
	}
	garden.Templates["list_template"] = list_template
	//category_template, err := template.ParseFiles("templates/cat.template.html", "templates/footer.template.html")
	//if err != nil {
	//	panic(err)
	//}
	//garden.Templates["category_template"] = category_template

	//tag_template, err := template.ParseFiles("templates/tag.template.html", "templates/footer.template.html")
	//if err != nil {
	//	panic(err)
	//}
	//garden.Templates["tag_template"] = tag_template

	links_template, err := template.ParseFiles("templates/partials/links.html")
	if err != nil {
		panic(err)
	}
	garden.Templates["links_template"] = links_template

	post_template, err := template.ParseFiles(partials...)
	post_template.ParseFiles("./templates/single.html")
	if err != nil {
		panic(err)
	}
	garden.Templates["post_template"] = post_template

	home_template, err := template.ParseFiles(base_files...)
	if err != nil {
		panic(err)
	}
	home_template, err = home_template.ParseFiles("./templates/index.html")
	if err != nil {
		panic(err)
	}
	garden.Templates["home_template"] = home_template

}
