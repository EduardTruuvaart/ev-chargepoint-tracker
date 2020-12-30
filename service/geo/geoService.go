package geo

import (
	"math"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
)

type GeoService struct {
	calculateDistanceInKm func(locationFrom model.Location, locationTo model.Location) float64
}

func (*GeoService) CalculateDistanceInKm(locationFrom model.Location, locationTo model.Location) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * locationFrom.Latitude / 180)
	radlat2 := float64(PI * locationTo.Latitude / 180)

	theta := float64(locationFrom.Longitude - locationTo.Longitude)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344

	return dist
}
