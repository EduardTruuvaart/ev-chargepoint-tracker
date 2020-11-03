package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StationRepository struct {
	findByID func(ID int64) (*model.Station, error)
}

func (*StationRepository) FindByID(ID int64) (*model.Station, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "station"
	stationID := strconv.FormatInt(ID, 10)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"stationId": {
				S: aws.String(stationID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := svc.GetItem(params)

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		msg := "Could not find station"
		return nil, errors.New(msg)
	}

	station := model.Station{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &station)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println("Found item:")
	fmt.Println("ID:  ", station.ID)
	fmt.Println("Status: ", station.Status)

	return &station, nil
}
