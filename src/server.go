package main

import (
    "fmt"
    "net/http"
	"net/url"
	"log"
	"io/ioutil"
	"math"
	"strings"
	"strconv"
	"os"
	"encoding/json"
)

type JsonRowItemFields struct {
	lat float64
	lon float64
	name string
}

type JsonRowItem struct {
	id string
	// order []float64
	fields JsonRowItemFields
}

type JsonRoot struct {
	total_rows int16
	bookmark string
	rows []JsonRowItem
}

func CreateDbLink(lon, lat, r float64) string {
	lon_from := math.Mod(lon - r/2/100, 50)
	lon_to := math.Mod(lon + r/2/100, 50)
	lat_from := math.Mod(lat - r/2/100, 205)
	lat_to := math.Mod(lat + r/2/100, 205)

	query := url.QueryEscape(fmt.Sprintf("lon:[%f TO %f] AND lat:[%f TO %f]", lon_from, lon_to, lat_from, lat_to))

	var link string
	if os.Getenv("ENV") == "test" {
		link = "http://localhost:3000/geo?q="
	} else {
		link = "https://mikerhodes.cloudant.com/airportdb/_design/view1/_search/geo?q="
	}

	return link + query
}

func GetResultsFromDb(link string) string {

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

func FilterList(list []JsonRowItem, lon, lat, r float64) []JsonRowItem {

	filteredList := make([]JsonRowItem, 0)
	for _, v := range list {

        if distance(lon, lat, v.fields.lon, v.fields.lat) <= r {
            filteredList = append(filteredList, v)
        }
    }

    return filteredList
}

func checkParams(w http.ResponseWriter, req *http.Request) bool {
	var missingParams bool = false
	if req.URL.Query()["lon"] == nil {
		fmt.Fprintf(w, "Please set query string ?lon=float64")
		missingParams = true
	}

	if req.URL.Query()["lat"] == nil {
		fmt.Fprintf(w, "Please set query string ?lon=float64")
		missingParams = true
	}

	if req.URL.Query()["r"] == nil {
		fmt.Fprintf(w, "Please set query string ?r=float64 as range")
		missingParams = true
	}

	return missingParams;
}

func GetList(w http.ResponseWriter, req *http.Request) {

	if checkParams(w, req) == false {

		var lonParam float64
		var latParam float64
		var rParam float64
		var err error
		if lonParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lon"], ""), 64); err != nil {
			log.Fatalln(err)
		}
		if latParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lat"], ""), 64); err != nil {
			log.Fatalln(err)
		}
		if rParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["r"], ""), 64); err != nil {
			log.Fatalln(err)
		}

		link := CreateDbLink(lonParam, latParam, rParam)
		jsonData := (GetResultsFromDb(link))

		fmt.Println(jsonData)

		var message JsonRoot
		err = json.Unmarshal([]byte(jsonData), &message)

		if err != nil {
			log.Fatalln(err)
			fmt.Println(err)
		}

		fmt.Println(message)

		if len(message.rows) > 0 {
			fmt.Println(len(message.rows))
			fmt.Println(lonParam)
			fmt.Println(latParam)
			fmt.Println(rParam)


			list := FilterList(message.rows, lonParam, latParam, rParam)
			listRowsCount := len(list)

			for i := 0; i < listRowsCount; i++ {
				fmt.Fprintf(w, "name: %s lon: %f lat: %f", list[i].fields.name, list[i].fields.lat, list[i].fields.lon)
			}
		} else {
			fmt.Fprintf(w, "No Airport found in the selected area")
		}
	}

}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	return dist * 1.609344
}

func main() {
	http.HandleFunc("/getlist", GetList)
	http.ListenAndServe(":8080", nil)
}