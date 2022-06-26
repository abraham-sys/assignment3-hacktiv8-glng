package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
)

const PORT string = ":8000"

type Data struct {
	Status       map[string]uint `json:"status"`
	NatureStatus string
}

var (
	natureData          Data
	pageAlreadyReloaded bool
)

func main() {
	http.HandleFunc("/siaga", func(w http.ResponseWriter, r *http.Request) {

		jsonNaturesData, err := os.Open("./data/water_wind_status.json")

		defer jsonNaturesData.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		byteValue, err := ioutil.ReadAll(jsonNaturesData)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal([]byte(byteValue), &natureData)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "GET" {

			html, err := template.ParseFiles("./public/windwater_status.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if pageAlreadyReloaded {
				natureData.Status["water"] = uint(rand.Intn(100))
				natureData.Status["wind"] = uint(rand.Intn(100))
				writeJsonFile(natureData)
			} else {
				pageAlreadyReloaded = true
			}

			natureData = getStatus(natureData)

			fmt.Println(natureData, "here")

			html.Execute(w, natureData)
			return

		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	})

	http.ListenAndServe(PORT, nil)

}

func getStatus(natureData Data) Data {
	if natureData.Status["water"] < 5 {
		natureData.NatureStatus = "aman"
	} else if natureData.Status["water"] <= 8 {
		natureData.NatureStatus = "siaga"
	} else {
		natureData.NatureStatus = "bahaya"
	}

	if natureData.Status["wind"] < 6 {
		natureData.NatureStatus = "aman"
	} else if natureData.Status["wind"] <= 15 {
		natureData.NatureStatus = "siaga"
	} else {
		natureData.NatureStatus = "bahaya"
	}
	return natureData
}

func writeJsonFile(data Data) {
	var newdata = map[string]map[string]uint{
		"status": data.Status,
	}

	file, _ := json.MarshalIndent(newdata, "", " ")

	_ = ioutil.WriteFile("./data/water_wind_status.json", file, 0644)
}
