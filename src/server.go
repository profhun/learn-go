package main

import (
    "fmt"
    "net/http"
	"net/url"
	"log"
	"io/ioutil"
)

func getResultsFromDb(lat, lon float32) string {
	lon_from := 0.000
	lon_to := 50.0000
	lat_from := 0.0000
	lat_to := 205.0000
	query := url.QueryEscape(fmt.Sprintf("lon:[%f TO %f] AND lat:[%f TO %f]", lon_from, lon_to, lat_from, lat_to))
	link := "https://mikerhodes.cloudant.com/airportdb/_design/view1/_search/geo?q=" + query
	resp, err := http.Get(link)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func getList(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, getResultsFromDb(0, 10))
}

func main() {
	http.HandleFunc("/getlist", getList)

	http.ListenAndServe(":8080", nil)
}