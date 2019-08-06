package utils

import (
    "net/http"
	"log"
	"io/ioutil"
)

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
