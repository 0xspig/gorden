# Gorden - Golang Digital Garden
Gorden is a digital garden written in Golang.
It's basically a lightweight SSG shipped with a node graph backend.

# Install
Creating a digital garden with Gorden involves two steps.
First - generating a site.
Second - running the backend.

## Site Generation
To create a new site navigate to an empty folder and run:

``` gorden -init "Site Name" ```

This will populate the current working directory with the base files.

## Running the backend
The backend can be installed as a systemd service on your server.
By default the server runs on localhost:3000.
It should be reverse proxied by your http server.

# Layout of a Goden site
Goden relies on markdown files and go templates to generate pages.
Every markdown file generates regular node.
Frontmatter tags generate Tag nodes.
Folders in the content directory generate Category nodes.

External links in markdown will generate external link nodes.
Internal links can be specified with the format `{Link Text}[NodeID]`.

The NodeID is always the name of the markdown file in the case of regular nodes (for example `{Test Post}[test-post.md]`).

For tags or categories it is simply the name (for example `{Test Category}[tests]`).

Note that each post, category, and node must have a unique specifier. That is, a category named "test" precludes usage of the keyword "test" as a tag in the front matter.

Nodes are NOT specified by folder/category. That is, "pictures/test.md" and "music/test.md" will both generate the same NodeID and such collisions should be avoided.



## content
The content directory contains blog content.
Content can be organized by two specifiers: Categories or Tags.

Each folder within the content folder will be rendered as a Category node.
Categories are good ways to separate large groups of content.

Tags are specified within the frontmatter of each markdown file using YAML.

Categories being delimited by folders means that each post can belong to only one category.
There is no limit, however, on the number of tags that can be attributed to a post.

## src
This folder will contain any source files which need to be built as npm dependent javascript or scss.

## static
This folder contains all static files.
This should include images and any files generated from the src dir.

## templates
This folder contains template files that will be used to render your final site.

