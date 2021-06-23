package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "https://gmail.googleapis.com/gmail/v1/users/i.joni40@gmail.com/messages/1787936fa2a028c0"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer ya29.a0AfH6SMDG9bWbDm7rGSRl-SF5IPZCFVco8UXoS6Q0VUFwRJ0tHrbi_oB-2Ovp4rykoPbX1ByjOigBDpPGWTsIrrIbF4rhQ8icLwvh92QdGTEv99UXKzqR2sXDlEWEkRKFCJwNVlFrLasoeetLKqcwlYiwX84tXg")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Print(res)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var gmData interface{}
	json.Unmarshal(body, &gmData) //extracting the json file
	// fmt.Println(gmData)

	// taking required values
	HistoryId := gmData.(map[string]interface{})["historyId"]
	snippet := gmData.(map[string]interface{})["snippet"]
	// payload := gmData.(map[string]interface{})["payload"]
	labelIds := gmData.(map[string]interface{})["labelIds"]
	// headers := payload.(map[string]interface{})["headers"]
	fmt.Println(HistoryId, "\n", snippet, "\n", labelIds)
}
