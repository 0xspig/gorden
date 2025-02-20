# Gorden - Golang Digital Garden
Gorden is a digital garden written in Golang.
It's basically a lightweight SSG shipped with a node graph backend.

# Install
Creating a digital garden with Gorden involves two steps.
First - generating a site.
Second - running the backend.

## Site Generation
To create a new site navigate to an empty folder and run:

``` gorden init "Site Name" ```

This will populate the current working directory with the base files.

## Running the backend
The backend can be installed as a systemd service on your server.
By default the server runs on localhost:3000.
It should be reverse proxied by your http server.

# Layout of a Goden site
Goden relies on markdown files and go templates to generate pages.