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
	Lat float64
	Lon float64
	Name string
	Distance float64
}

type JsonRowItem struct {
	Id string
	Order []float64
	Fields JsonRowItemFields
}

type JsonRoot struct {
	Total_rows int
	Bookmark string
	Rows []JsonRowItem
}

func CreateDbLink(lon, lat, r float64, bookmark string) string {
	// one means nearly 101km
	lon_from := lon - r/100
	if lon_from < -180 {
		lon_from = -180
	}
	lon_to := lon + r/100
	if lon_to > 180 {
		lon_to = 180
	}
	lat_from := lat - r/100
	if lat_from < -90 {
		lat_from = -90
	}
	lat_to := lat + r/100
	if lat_to > 90 {
		lat_to = 90
	}

	query := url.QueryEscape(fmt.Sprintf(`lon:[%f TO %f] AND lat:[%f TO %f]`, lon_from, lon_to, lat_from, lat_to))

	var link string
	if os.Getenv("ENV") == "test" {
		link = "http://localhost:3000/geo?q="
	} else {
		link = "https://mikerhodes.cloudant.com/airportdb/_design/view1/_search/geo?q="
	}

	pagination := fmt.Sprintf(`&limit=10&sort="<distance,lon,lat,%f,%f,km>"`, lon, lat)

	if bookmark != "" {
		pagination += fmt.Sprintf(`&bookmark=%s`, bookmark)
	}

	return link + query + pagination
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

		distance := distance(lon, lat, v.Fields.Lon, v.Fields.Lat)
        if distance <= r {
			v.Fields.Distance = distance
            filteredList = append(filteredList, v)
        }
    }

    return filteredList
}

func checkParams(w http.ResponseWriter, req *http.Request) (bool, float64, float64, float64) {
	var isCorrect bool = true
	var err error

	if req.URL.Query()["lon"] == nil {
		fmt.Fprintf(w, "Please set query string ?lon=float64")
		isCorrect = false
	}

	if req.URL.Query()["lat"] == nil {
		fmt.Fprintf(w, "Please set query string ?lon=float64")
		isCorrect = false
	}

	if req.URL.Query()["r"] == nil {
		fmt.Fprintf(w, "Please set query string ?r=float64 as range")
		isCorrect = false
	}

	var lonParam float64
	if lonParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lon"], ""), 64); err != nil {
		fmt.Fprintf(w, "param lon is not a float")
		isCorrect = false
	}
	var latParam float64
	if latParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lat"], ""), 64); err != nil {
		fmt.Fprintf(w, "param lat is not a float")
		isCorrect = false
	}
	var rParam float64
	if rParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["r"], ""), 64); err != nil {
		fmt.Fprintf(w, "param r is not a float")
		isCorrect = false
	}

	if (lonParam < -180 || lonParam > 180) {
		fmt.Fprintf(w, "param lon must be between -180 and 180")
		isCorrect = false
	}

	if (latParam < -90 || latParam > 90) {
		fmt.Fprintf(w, "param lon must be between -180 and 180")
		isCorrect = false
	}

	if rParam <= 0 {
		fmt.Fprintf(w, "range must be greater then zero")
		isCorrect = false
	}

	return isCorrect, lonParam, latParam, rParam
}

func GetList(w http.ResponseWriter, req *http.Request) {

	isCorrect, lonParam, latParam, rParam := checkParams(w, req)

	if isCorrect == true {

		var while bool = true;
		var bookmark string = "";

		for while == true {

			link := CreateDbLink(lonParam, latParam, rParam, bookmark)
			jsonData := (GetResultsFromDb(link))

			var message JsonRoot
			err := json.Unmarshal([]byte(jsonData), &message)

			if err != nil {
				log.Fatalln(err)
			}

			bookmark = message.Bookmark

			if len(message.Rows) > 0 {

				list := FilterList(message.Rows, lonParam, latParam, rParam)
				listRowsCount := len(list)

				for i := 0; i < listRowsCount; i++ {
					fmt.Fprintf(w, "name: %s; distance: %fKm lon: %f lat: %f \r\n", list[i].Fields.Name, list[i].Fields.Distance, list[i].Fields.Lat, list[i].Fields.Lon)
				}
			} else {
				if (message.Total_rows == 0) {
					fmt.Fprintf(w, "No Airport found in the selected area")
				} else {
					// pagination is ended
				}
				while = false;
			}
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