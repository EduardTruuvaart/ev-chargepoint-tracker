package stations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/EduardTruuvaart/ev-chargepoint-tracker/service/geo"
)

type StationService struct {
	ApiKey            string
	getStatus         func(stationID string) *model.Station
	search            func(location model.Location) []model.Station
	getDetails        func(stationID string) *model.Station
	fulfillAllDetails func(currentLocation model.Location, stations []model.Station) []model.Station
	getAllDetails     func(stationID string) *model.Station
}

func NewStationService(apiKey string) *StationService {
	return &StationService{ApiKey: apiKey}
}

func (service *StationService) GetStatus(stationID string) []model.Device {
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

	devicesObjArr := []model.Device{}
	for _, element := range devicesArr {
		deviceIDFloat := element.(map[string]interface{})["id"].(float64)
		deviceStatus := element.(map[string]interface{})["status"].(map[string]interface{})["description"]
		deviceConnectorsArr := element.(map[string]interface{})["connectors"].([]interface{})
		statusHistoryArr := element.(map[string]interface{})["status_history"].([]interface{})

		device := model.NewDevice(strconv.Itoa((int(deviceIDFloat))))
		device.Status = deviceStatus.(string)

		for _, connectorElem := range deviceConnectorsArr {
			connectorIDStr := connectorElem.(map[string]interface{})["id"].(string)
			connectorName := connectorElem.(map[string]interface{})["name"]
			connectorSpeed := connectorElem.(map[string]interface{})["speed_group_summary"].(map[string]interface{})["speed_group_name"]
			connectorStatus := connectorElem.(map[string]interface{})["status"].(map[string]interface{})["description"]

			connector := model.NewConnector(connectorIDStr)
			connector.Name = connectorName.(string)
			connector.Status = connectorStatus.(string)
			connector.Speed = connectorSpeed.(string)
			device.Connectors = append(device.Connectors, *connector)
		}

		for _, statusHistoryElem := range statusHistoryArr {
			dateElem := statusHistoryElem.(map[string]interface{})["date"]
			historyDescriptionStr := statusHistoryElem.(map[string]interface{})["description"].(string)
			historyDateStr := dateElem.(map[string]interface{})["value"].(string)
			historyDateTitle := dateElem.(map[string]interface{})["title"].(string)
			historyDate, _ := time.Parse(time.RFC3339, historyDateStr)

			status := model.Status{
				Description: historyDescriptionStr,
				DateTitle:   historyDateTitle,
				Date:        historyDate,
			}
			device.StatusHistory = append(device.StatusHistory, status)
		}

		devicesObjArr = append(devicesObjArr, *device)
	}

	return devicesObjArr
}

func (service *StationService) Search(location model.Location) []model.Station {
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

	station.Location = model.Location{
		Latitude:  addressInfo.(map[string]interface{})["latitude"].(float64),
		Longitude: addressInfo.(map[string]interface{})["longitude"].(float64),
	}

	return station
}

func (service *StationService) FulfillAllDetails(currentLocation model.Location, stations []model.Station) []model.Station {
	var geoService geo.GeoService

	fulfilledStations := []model.Station{}
	for _, element := range stations {
		stationDetails := service.GetDetails(element.ID)
		devices := service.GetStatus(element.ID)
		station := model.NewStation(element.ID)
		station.FormattedAddress = stationDetails.FormattedAddress
		station.Name = stationDetails.Name
		station.NetworkName = stationDetails.NetworkName
		station.PostCode = stationDetails.PostCode
		station.Location = stationDetails.Location
		station.DistanceInKm = geoService.CalculateDistanceInKm(currentLocation, stationDetails.Location)
		station.Devices = devices
		fulfilledStations = append(fulfilledStations, *station)
	}
	service.SortByDistance(fulfilledStations)

	return fulfilledStations
}

func (service *StationService) GetAllDetails(stationID string) *model.Station {
	stationDetails := service.GetDetails(stationID)
	devices := service.GetStatus(stationID)
	stationDetails.Devices = devices

	return stationDetails
}

func (service *StationService) SortByDistance(stations []model.Station) {
	sort.Slice(stations, func(i, j int) bool {
		return stations[i].DistanceInKm < stations[j].DistanceInKm
	})
}
