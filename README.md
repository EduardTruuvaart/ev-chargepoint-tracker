# EV Chargepoint Tracker
Tracking availability of electric vehicle charge points

## Running using Docker container
```
docker build --tag tracker:latest .
docker run -it --rm -e APIKEY='mykey' tracker <stationID>
```
