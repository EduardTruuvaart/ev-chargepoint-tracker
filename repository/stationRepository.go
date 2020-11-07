package repository

import (
	"fmt"

	"github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StationRepository struct {
	findByID func(ID string) (*model.Station, error)
	save     func(station *model.Station)
}

func (*StationRepository) FindByID(ID string) (*model.Station, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "station"

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := svc.GetItem(params)

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		fmt.Println("Could not find station")
		return nil, nil
	}

	station := model.Station{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &station)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal record, %v", err))
	}

	return &station, nil
}

func (*StationRepository) Save(station *model.Station) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := "station"

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(station.Status),
			},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("set #station_status = :status"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(station.ID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#station_status": aws.String("status"),
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to update record, %v", err))
	}
}
