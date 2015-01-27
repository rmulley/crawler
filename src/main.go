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
	"runtime"
	"strconv"
	"time"
) //import

type PostData_t struct {
	Urls []string `json:"urls"`
} //PostData_t

type JobResult_t struct {
	ID      uint          `json:"id"`
	Results []UrlResult_t `json:"results"`
} //JobResult_t

type UrlResult_t struct {
	Url    string   `json:"url"`
	Images []string `json:"images"`
} //UrlResult_t

type Status_t struct {
	ID         uint `json:"id"`
	Completed  uint `json:"completed"`
	InProgress uint `json:"inProgress"`
} //Status_t

var (
	jobCtr   uint                  = 1
	results  map[uint]*JobResult_t = make(map[uint]*JobResult_t)
	statuses map[uint]*Status_t    = make(map[uint]*Status_t)
) //var

func main() {
	var (
		r = mux.NewRouter()
	) //var

	// Use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Routes for HTTP requests
	r.HandleFunc("/status/{id}", statusHandler)
	r.HandleFunc("/result/{id}", resultHandler)
	r.HandleFunc("/", crawlHandler)

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

func crawlHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		err  error
		data PostData_t
	) //var

	if err = json.NewDecoder(req.Body).Decode(&data); err != nil {
		log.Println(err)
	} //if

	for _, urlToCrawl := range data.Urls {
		var (
			ndx uint = 0
		) //var

		// Create new status entry for Job
		statuses[jobCtr] = &Status_t{
			ID: jobCtr,
		} //Status_t

		// Create new result entry for Job
		results[jobCtr] = &JobResult_t{
			ID:      jobCtr,
			Results: make([]UrlResult_t, 0),
		} //Result_t

		//parseLinks(urlToGet, content)

		results[jobCtr].Results = append(results[jobCtr].Results, UrlResult_t{
			Url: urlToCrawl,
		}) //append

		//ndx++

		go crawlUrl(ndx, urlToCrawl, statuses[jobCtr], results[jobCtr])
		jobCtr++

		// REMOVE THIS NOW
		ndx++
	} //for
} //crawlHandler

func crawlUrl(ndx uint, urlToCrawl string, status *Status_t, result *JobResult_t) {
	var (
		err      error
		content  string
		imgs     []string
		urlToGet *url.URL
	) //var

	// Update status and result of current job
	status.InProgress++

	// Parse URL
	if urlToGet, err = url.Parse(urlToCrawl); err != nil {
		log.Println(err)
		return
	} //if

	// Retrieve content of URL
	if content, err = getUrlContent(urlToGet.String()); err != nil {
		log.Println(err)
		return
	} //if

	// Clean up HTML entities
	content = html.UnescapeString(content)

	// Retrieve image URLs
	if imgs, err = parseImages(urlToGet, content); err != nil {
		log.Println(err)
		return
	} //if

	for _, img := range imgs {
		result.Results[ndx].Images = append(result.Results[ndx].Images, img)
	} //for

	status.InProgress--
	status.Completed++
} //crawlUrl

func resultHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		err error
		id  int
	) //var

	vars := mux.Vars(req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		log.Println(err)
		return
	} //if

	// Set response headers
	resp.Header().Set("Accept", "application/json")
	resp.Header().Set("Content-Type", "application/json")

	// Encode and write the result
	if err = json.NewEncoder(resp).Encode(results[uint(id)]); err != nil {
		log.Println(err)
	} //if

	return
} //resultHandler

func statusHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		status Status_t
	) //var

	status.Completed = 1
	status.InProgress = 2

	// Set response headers
	resp.Header().Set("Accept", "application/json")
	resp.Header().Set("Content-Type", "application/json")

	// Encode and write the result
	if err = json.NewEncoder(resp).Encode(&status); err != nil {
		log.Println(err)
	} //if

	return
} //statusHandler

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

func parseLinks(urlToGet *url.URL, content string) ([]string, error) {
	var (
		err       error
		links     []string
		matches   [][]string
		findLinks = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	) //var

	// Retrieve all anchor tag URLs from string
	matches = findLinks.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var linkUrl *url.URL

		// Parse the anchr tag URL
		if linkUrl, err = url.Parse(val[1]); err != nil {
			return links, err
		} //if

		// If the URL is absolute, add it to the slice
		// If the URL is relative, build an absolute URL
		if linkUrl.IsAbs() {
			links = append(links, linkUrl.String())
		} else {
			links = append(links, urlToGet.Scheme+"://"+urlToGet.Host+linkUrl.String())
		} //else
	} //for

	return links, err
} //parseImages
