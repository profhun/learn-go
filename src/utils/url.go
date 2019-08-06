package utils

import (
    "fmt"
    "net/http"
	"strconv"
	"strings"
	"net/url"
	"os"
)

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


func CheckParams(w http.ResponseWriter, req *http.Request) (bool, float64, float64, float64, []string) {
	var isCorrect bool = true
	var err error
	var lonParam float64
	var latParam float64
	var rParam float64
	var errors []string

	if req.URL.Query()["lon"] == nil {
		errors = append(errors, "Please set query string ?lon=float64")
		isCorrect = false
	} else if lonParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lon"], ""), 64); err != nil {
		errors = append(errors, "param lon is not a float")
		isCorrect = false
	}

	if req.URL.Query()["lat"] == nil {
		errors = append(errors, "Please set query string ?lon=float64")
		isCorrect = false
	} else if latParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["lat"], ""), 64); err != nil {
		errors = append(errors, "param lat is not a float")
		isCorrect = false
	}

	if req.URL.Query()["r"] == nil {
		errors = append(errors, "Please set query string ?r=float64 as range")
		isCorrect = false
	} else if rParam, err = strconv.ParseFloat(strings.Join(req.URL.Query()["r"], ""), 64); err != nil {
		errors = append(errors, "param r is not a float")
		isCorrect = false
	}

	if (lonParam < -180 || lonParam > 180) {
		errors = append(errors, "param lon must be between -180 and 180")
		isCorrect = false
	}

	if (latParam < -90 || latParam > 90) {
		errors = append(errors, "param lon must be between -180 and 180")
		isCorrect = false
	}

	if rParam <= 0 {
		errors = append(errors, "range must be greater then zero")
		isCorrect = false
	}

	return isCorrect, lonParam, latParam, rParam, errors
}