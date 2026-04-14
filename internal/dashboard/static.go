package dashboard

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public/*
var assetsFS embed.FS

func Index() ([]byte, error) {
	return assetsFS.ReadFile("public/index.html")
}

func AssetsHandler() (http.Handler, error) {
	sub, err := fs.Sub(assetsFS, "public")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}
