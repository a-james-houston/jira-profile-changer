package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	configFile         = "env.json"
	sequentialOrdering = "SEQUENTIAL"
	randomOrdering     = "RANDOM"
	urlApiString       = "%s/rest/api/2/user/avatar?username=%s"
)

type Config struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	JiraBaseURL string `json:"base_url"`
	AvatarIDs   []int  `json:"avatar_ids"`
	UsageOrder  string `json:"usage_order"`
	DataFile    string `json:"data_file"`
}

func main() {
	// Open config file
	jsonFile, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	idToUse := getNextImageID(config)
	ok, err := setAvatar(config, idToUse)
	if ok {
		fmt.Printf("SUCCESS: %d\n", idToUse)
		updateDataFile(config, idToUse)
	} else {
		panic(err)
	}
}

func updateDataFile(config Config, id int) {
	err := os.WriteFile(config.DataFile, []byte(strconv.Itoa(id)), 0666)
	if err != nil {
		panic(err)
	}
}

func setAvatar(config Config, idToUse int) (bool, error) {
	client := &http.Client{}
	url := fmt.Sprintf(urlApiString, config.JiraBaseURL, config.Username)
	postBody := fmt.Sprintf("{\"id\":%s}", strconv.Itoa(idToUse))
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(postBody))
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 204 {
		return false, nil
	}
	return true, nil
}

func getNextImageID(config Config) int {
	switch config.UsageOrder {
	case sequentialOrdering:
		id := getLastImageID(config)
		if id == 0 {
			return config.AvatarIDs[0]
		}
		index := getIntIndex(id, config.AvatarIDs)
		if index == len(config.AvatarIDs)-1 {
			return config.AvatarIDs[0]
		}
		return config.AvatarIDs[index+1]
	case randomOrdering:
		return config.AvatarIDs[rand.Intn(len(config.AvatarIDs))]
	default:
		// default to random is ordering is not a valid selection
		return config.AvatarIDs[rand.Intn(len(config.AvatarIDs))]
	}
}

// gets the last image id that was set as the user profile
// returns zero and creates data file if the data file does not exist
func getLastImageID(config Config) int {
	dataFile, err := ioutil.ReadFile(config.DataFile)
	if os.IsNotExist(err) {
		os.Create(config.DataFile)
		return 0
	} else if err != nil {
		panic(err)
	}
	val, err := strconv.Atoi(string(dataFile))
	if err != nil {
		panic(err)
	}
	return val
}

// returns the index of an element in an int slice, returns -1 if element is not in slice
func getIntIndex(elem int, slice []int) int {
	for i, val := range slice {
		if val == elem {
			return i
		}
	}
	return -1
}
