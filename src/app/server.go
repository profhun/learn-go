package main

import (
    "fmt"
    "net/http"
	"log"
	"encoding/json"
	"../utils"
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

func FilterList(list []JsonRowItem, lon, lat, r float64) []JsonRowItem {

	filteredList := make([]JsonRowItem, 0)
	for _, v := range list {

		distance := utils.Distance(lon, lat, v.Fields.Lon, v.Fields.Lat)
        if distance <= r {
			v.Fields.Distance = distance
            filteredList = append(filteredList, v)
        }
    }

    return filteredList
}

func GetList(w http.ResponseWriter, req *http.Request) {

	isCorrect, lonParam, latParam, rParam := utils.CheckParams(w, req)

	if isCorrect == true {

		var while bool = true;
		var bookmark string = "";

		for while == true {

			link := utils.CreateDbLink(lonParam, latParam, rParam, bookmark)
			jsonData := (utils.GetResultsFromDb(link))

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

func main() {
	http.HandleFunc("/getlist", GetList)
	http.ListenAndServe(":8080", nil)
}