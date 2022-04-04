package main

import (
	"mime"
	"net/http"
	"strings"
)

func filename(resp *http.Response) string {
	array := strings.Split(resp.Request.RequestURI, "/")
	name := array[len(array)-1]
	getDispos := resp.Header.Get("content-disposition")
	if getDispos != "" {
		_, params, err := mime.ParseMediaType(getDispos)
		if err != nil {
			return name
		}
		filename, ok := params["filename"]
		if ok {
			return filename
		}
	}
	return name
}

func urlJoin(base string, pathes ...string) string {
	for _, path := range pathes {
		if strings.HasSuffix(base, "/") {
			base = base + path
		} else {
			base = base + "/" + path
		}
	}
	return base
}
