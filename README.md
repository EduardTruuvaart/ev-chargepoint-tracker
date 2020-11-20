# EV Chargepoint Tracker
App is tracking availability of EV charge points

## Running using Docker container
```
docker build --tag tracker:latest .
docker run -it --rm -e APIKEY=<api_key> -e AWS_ACCESS_KEY_ID=<aws_key> -e AWS_SECRET_ACCESS_KEY=<aws_secret> -e AWS_DEFAULT_REGION=<aws_region> -e STATIONID=<station_id> tracker
```
