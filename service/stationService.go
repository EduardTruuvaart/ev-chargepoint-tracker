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
	ApiKey            string
	getStatus         func(stationID string) *model.Station
	search            func(location model.Location) []model.Station
	getDetails        func(stationID string) *model.Station
	fulfillAllDetails func(stations []model.Station) []model.Station
}

func NewStationService(apiKey string) *StationService {
	return &StationService{ApiKey: apiKey}
}

func (service *StationService) GetStatus(stationID string) *model.Station {
	client := &http.Client{}

	requestURI := fmt.Sprintf("https://api.zap-map.com/v5/chargepoints/location/%s/status", stationID)

	req, err := http.NewRequest("GET", requestURI, nil)
	req.Header.Add("X-Api-Key", service.ApiKey)

	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body Error")
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	var chargePointsData = result["resources"].(map[string]interface{})["chargepoint_location_status"].(map[string]interface{})["data"]
	var devices = chargePointsData.(map[string]interface{})["devices"]
	devicesArr := devices.([]interface{})
	var statusHistory = devicesArr[0].((map[string]interface{}))["status_history"].([]interface{})
	var currentStatus = statusHistory[0].(map[string]interface{})["description"].(string)
	var stationIDFloat = chargePointsData.(map[string]interface{})["id"].(float64)
	stationIDStr := strconv.Itoa((int(stationIDFloat)))

	station := model.NewStation(stationIDStr)
	station.Status = currentStatus
	return station
}

func (service *StationService) Search(location model.Location) []model.Station {
	fmt.Println("Key: ", service.ApiKey)
	client := &http.Client{}

	requestURI := fmt.Sprintf("https://api.zap-map.com/v5/chargepoints/locations/search?lat=%v&long=%v&radius=2&unit=KM&connector-types=&networks=&payments=&location-types=&access=2&ev-models=", location.Latitude, location.Longitude)

	req, err := http.NewRequest("GET", requestURI, nil)
	req.Header.Add("X-Api-Key", service.ApiKey)

	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body Error")
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	var chargePointLocationsData = result["resources"].(map[string]interface{})["search_chargepoint_locations"].(map[string]interface{})["data"]
	dataArr := chargePointLocationsData.([]interface{})

	stations := []model.Station{}
	for _, element := range dataArr {
		stationIDFloat := element.(map[string]interface{})["id"].(float64)
		station := model.NewStation(strconv.Itoa((int(stationIDFloat))))
		stations = append(stations, *station)
	}

	return stations
}

func (service *StationService) GetDetails(stationID string) *model.Station {
	client := &http.Client{}

	requestURI := fmt.Sprintf("https://api.zap-map.com/v5/chargepoints/locations/placecards?id=%s", stationID)

	req, err := http.NewRequest("GET", requestURI, nil)
	req.Header.Add("X-Api-Key", service.ApiKey)

	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Get Error")
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Body Error")
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	var data = result["resources"].(map[string]interface{})["chargepoint_locations_placecards"].(map[string]interface{})["data"].([]interface{})[0]
	addressInfo := data.(map[string]interface{})["address_info"]
	name := data.(map[string]interface{})["name"]
	networkName := data.(map[string]interface{})["primary_network"].(map[string]interface{})["name"]
	formattedAddress := addressInfo.(map[string]interface{})["formatted"]
	postCode := addressInfo.(map[string]interface{})["postcode"]

	station := model.NewStation(stationID)
	station.Name = name.(string)
	station.NetworkName = networkName.(string)
	station.FormattedAddress = formattedAddress.(string)
	station.PostCode = postCode.(string)

	return station
}

func (service *StationService) FulfillAllDetails(stations []model.Station) []model.Station {
	fulfilledStations := []model.Station{}
	for _, element := range stations {
		stationDetails := service.GetDetails(element.ID)
		stationStatus := service.GetStatus(element.ID)

		station := model.NewStation(element.ID)
		station.FormattedAddress = stationDetails.FormattedAddress
		station.Name = stationDetails.Name
		station.NetworkName = stationDetails.NetworkName
		station.PostCode = stationDetails.PostCode
		station.Status = stationStatus.Status
		fulfilledStations = append(fulfilledStations, *station)
	}

	return fulfilledStations
}
