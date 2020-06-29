package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

//APIOutput Struct helps to get data from external APIs
type APIOutput struct {
	Name      string      `json:"name"`
	Character []Character `json:"character"`
}

//Character Struct to retrieve characters from external APIs
type Character struct {
	Name     string `json:"name"`
	MaxPower int64  `json:"max_power"`
}

//TrueCharacter Struct helps the common map data structure
type TrueCharacter struct {
	Name      string
	MaxPower  int64
	CurrPower int64
	Frequency int64
}

//ReturnResponse is the structure for the output to be sent to server
type ReturnResponse struct {
	Character    string `json:"character"`
	CurrentPower int64  `json:"currentpowder"`
}

var response map[string]*TrueCharacter
var backup map[string]*TrueCharacter
var wg2 sync.WaitGroup
var flag bool
var flag2 bool = true

//APICall does all the processing of external API data
func APICall(url string, wg *sync.WaitGroup) {

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error")
	}

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error")
	}

	var APIoutput APIOutput

	json.Unmarshal(html, &APIoutput)

	for char := range APIoutput.Character {

		if len(response) > 15 {
			var name string
			var min int64
			for key, val := range response {
				if min > val.MaxPower {
					min = val.MaxPower
					name = key
				}
			}

			backup[name] = response[name]

			delete(response, name)
		}

		var Truechar TrueCharacter
		Truechar.Name = APIoutput.Character[char].Name
		Truechar.MaxPower = APIoutput.Character[char].MaxPower
		Truechar.CurrPower = 0
		//Truechar.Frequency = 0
		response[APIoutput.Character[char].Name] = &Truechar

	}

	wg.Done()

}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

//ChangePower changes power of the character every 10 seconds
func ChangePower() {

	if flag == true {
		wg2.Add(1)
	}

	for item := range response {

		rand.Seed(time.Now().Unix())
		response[item].CurrPower = int64(rand.Intn(int(response[item].MaxPower) - 0))

	}

	wg2.Done()

}

//APIOuter is a shell which calls the encapsulated APICall and helps in connectivity with main
func APIOuter(w http.ResponseWriter, r *http.Request) {

	flag = false

	character := r.FormValue("char")
	fmt.Println(character)
	var wg sync.WaitGroup

	if flag2 == true {
		response = make(map[string]*TrueCharacter, 15)

		urlAvenger := "https://www.mocky.io/v2/5ecfd5dc3200006200e3d64b/"
		//urlAnti := "http://www.mocky.io/v2/5ecfd630320000f1aee3d64d"
		urlMutant := "http://www.mocky.io/v2/5ecfd6473200009dc1e3d64e"

		wg.Add(1)
		go APICall(urlAvenger, &wg)
		wg.Add(1)
		//go APICall(urlAnti)
		//wg.Add(1)
		go APICall(urlMutant, &wg)

		flag2 = false

		wg.Wait()
	}

	_, found := response[character]

	if found == false {

		var min int64
		var char string
		for name, val := range response {

			if min >= val.Frequency {
				min = val.Frequency
				char = name
			}
		}

		backup[char] = response[char]
		delete(response, char)

		response[character] = backup[character]
	}

	response[character].Frequency++

	go doEvery(10*time.Second, ChangePower)
	wg2.Add(1)
	wg2.Wait()
	flag = true

	fmt.Println("Response")
	fmt.Println(*response[character])
	//fmt.Println(len(response))

	Returnresponse := &ReturnResponse{
		Character:    character,
		CurrentPower: response[character].CurrPower}

	result, errFor := json.Marshal(Returnresponse)

	if errFor != nil {
		fmt.Println(errFor)
	}

	w.Write(result)
}
