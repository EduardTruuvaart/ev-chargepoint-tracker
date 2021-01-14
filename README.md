# EV Chargepoint Telegram Bot
Telegram messanger chat bot that shows EV charging stations around your location

## Deploying Telegram handler to AWS Lambda
```
1. cd cmd/bot/
2. GOOS=linux go build main.go
3. zip function.zip main
4. deploy function.zip to AWS Lambda
```
