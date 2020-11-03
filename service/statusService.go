package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
)

type StationService struct {
	getStatus func(stationID string, apiKey string) *model.Station
}

func (service *StationService) GetStatus(stationID string, apiKey string) *model.Station {
	fmt.Println("Key: ", apiKey)
	fmt.Println("Selected station ID: ", stationID)
	client := &http.Client{}

	requestURI := fmt.Sprintf("https://api.zap-map.com/v5/chargepoints/location/%s/status", stationID)

	req, err := http.NewRequest("GET", requestURI, nil)
	req.Header.Add("X-Api-Key", apiKey)

	if err != nil {
		fmt.Println("Get Error")
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Get Error")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body Error")
	}

	bodyStr := string(body)

	var result map[string]interface{}
	json.Unmarshal([]byte(bodyStr), &result)

	var chargePointsData interface{} = result["resources"].(map[string]interface{})["chargepoint_location_status"].(map[string]interface{})["data"]
	var devices interface{} = chargePointsData.(map[string]interface{})["devices"]
	devicesArr := devices.([]interface{})
	var statusHistory []interface{} = devicesArr[0].((map[string]interface{}))["status_history"].([]interface{})
	var currentStatus string = statusHistory[0].(map[string]interface{})["description"].(string)
	var stationIDFloat float64 = chargePointsData.(map[string]interface{})["id"].(float64)
	stationIDStr := strconv.Itoa((int(stationIDFloat)))

	return model.NewStation(stationIDStr, currentStatus)
}
