package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type APIResponse struct {
	Poster string
}

func main() {
	var title string
	var apiKey string
	flag.StringVar(&title, "title", "", "Movie title that you wanna download poster")
	flag.StringVar(&apiKey, "key", "", "omdbapi API Key")
	flag.Parse()

	url := fmt.Sprintf("https://www.omdbapi.com/?t=%v&apikey=%v", title, apiKey)

	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response APIResponse
	json.Unmarshal(body, &response)

	reqPoster, _ := http.NewRequest("GET", response.Poster, nil)
	resPoster, err := http.DefaultClient.Do(reqPoster)
	if err != nil {
		log.Fatalln(err)
	}
	defer resPoster.Body.Close()
	// open a file for writing
	date := time.Now()
	dateFormated := date.Format("2006-Jan-02")
	fileName := fmt.Sprintf("./%v-%v.jpg", title, dateFormated)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resPoster.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")

}
