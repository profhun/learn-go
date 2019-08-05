package main

import (
	"testing"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"math"
)

var dbMockResult string = `{"total_rows":3,"bookmark":"random-str","rows":[{"id":"aa915feea5ecf2f87d0c7bca672da2d5","order":[1.4142135381698608,0],"fields":{"lat":53.630389,"lon":9.988228,"name":"Hamburg"}},{"id":"aa915feea5ecf2f87d0c7bca672db3ae","order":[1.4142135381698608,1],"fields":{"lat":51.432447,"lon":12.241633,"name":"Leipzig Halle"}},{"id":"aa915feea5ecf2f87d0c7bca672d298a","order":[1.4142135381698608,1],"fields":{"lat":51.394783,"lon":4.960194,"name":"Weelde"}}]}`

func geoDbMock(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, dbMockResult)
}

func TestCreateDbLink(t *testing.T) {
	os.Setenv("ENV", "development")

	r := float64(2)
	lon := float64(10)
	lat := float64(5)

	lon_from := math.Mod(lon - r/2/100, 50)
	lon_to := math.Mod(lon + r/2/100, 50)
	lat_from := math.Mod(lat - r/2/100, 205)
	lat_to := math.Mod(lat + r/2/100, 205)

	query := url.QueryEscape(fmt.Sprintf("lon:[%f TO %f] AND lat:[%f TO %f]", lon_from, lon_to, lat_from, lat_to))

	prodLink := "https://mikerhodes.cloudant.com/airportdb/_design/view1/_search/geo?q=" + query
	link := CreateDbLink(float64(10), float64(5), float64(2))
	if link != prodLink {
		t.Errorf("Url was incorrect, got: %s, want: %s.", link, prodLink)
	}

	os.Setenv("ENV", "test")

	testLink := "http://localhost:3000/geo?q=" + query
	link = CreateDbLink(float64(10), float64(5), float64(2))

	if link != testLink {
		t.Errorf("Url was incorrect, got: %s, want: %s.", link, testLink)
	}
}

func TestGetResultsFromDb(t *testing.T) {
	srv := &http.Server{Addr: ":3000"}
	http.HandleFunc("/geo", geoDbMock)

	go func() {
		srv.ListenAndServe()
	}()

	os.Setenv("ENV", "test")

	link := CreateDbLink(float64(10), float64(5), float64(2))

	resultJson := GetResultsFromDb(link)

	if resultJson != dbMockResult {
		t.Errorf("DB Json result was incorrect, got: %s, want: %s.", resultJson, dbMockResult)
	}

	srv.Shutdown(context.TODO())

}
