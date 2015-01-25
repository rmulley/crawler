package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
) //import

type PostData_t struct {
	Urls []string `json:"urls"`
} //PostData_t

func main() {
	var (
		r = mux.NewRouter()
	) //var

	// Routes for HTTP requests
	r.HandleFunc("/status/{id}", StatusHandler)
	r.HandleFunc("/", CrawlHandler)

	// Build server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	} //http.Server

	// Start server and listen for HTTP requests
	log.Fatalln(server.ListenAndServe())
} //main

func CrawlHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		err  error
		data PostData_t
	) //var

	if err = json.NewDecoder(req.Body).Decode(&data); err != nil {
		log.Println(err)
	} //if

	for _, urlSupplied := range data.Urls {
		go func(urlFromPost string) {
			var (
				err      error
				content  string
				imgs     []string
				urlToGet *url.URL
			) //var

			// Parse URL
			if urlToGet, err = url.Parse(urlFromPost); err != nil {
				log.Fatalln(err)
			} //if

			// Retrieve content of URL
			if content, err = getUrlContent(urlToGet.String()); err != nil {
				log.Fatalln(err)
			} //if

			// Clean up HTML entities
			content = html.UnescapeString(content)

			// Retrieve image URLs
			if imgs, err = parseImages(urlToGet, content); err != nil {
				log.Fatalln(err)
			} //if

			for _, img := range imgs {
				log.Println(img)
			} //for
		}(urlSupplied) //goroutine
	} //for
} //CrawlHandler

func StatusHandler(resp http.ResponseWriter, req *http.Request) {
	log.Println("STATUS!")
} //StatusHandler

func getUrlContent(urlToGet string) (string, error) {
	var (
		err     error
		content []byte
		resp    *http.Response
	) //var

	// GET content of URL
	if resp, err = http.Get(urlToGet); err != nil {
		return "", err
	} //if
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != 200 {
		return "", err
	} //if

	// Read the body of the HTTP response
	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} //if

	return string(content), err
} //getUrlContent

func parseImages(urlToGet *url.URL, content string) ([]string, error) {
	var (
		err        error
		imgs       []string
		matches    [][]string
		findImages = regexp.MustCompile("<img.*?src=\"(.*?)\"")
	) //var

	// Retrieve all image URLs from string
	matches = findImages.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var imgUrl *url.URL

		// Parse the image URL
		if imgUrl, err = url.Parse(val[1]); err != nil {
			return imgs, err
		} //if

		// If the URL is absolute, add it to the slice
		// If the URL is relative, build an absolute URL
		if imgUrl.IsAbs() {
			imgs = append(imgs, imgUrl.String())
		} else {
			imgs = append(imgs, urlToGet.Scheme+"://"+urlToGet.Host+imgUrl.String())
		} //else
	} //for

	return imgs, err
} //parseImages
