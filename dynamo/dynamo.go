package dynamo

import (
	"github.com/alessiosavi/GoStatOgame/datastructure/players"
	"github.com/alessiosavi/GoStatOgame/datastructure/score"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"time"
)

// AWS limit for put/get batch operation
const batchGetItemSize = 100

func retrieveDataFromDynamo(p []score.Player, tableName string, svc *dynamodb.DynamoDB) ([]players.PlayerData, error) {
	// Retrieve the Item from dynamo
	var i = 0
	var data []players.PlayerData
	var result *dynamodb.BatchGetItemOutput
	var err error
	for ; i < len(p)-batchGetItemSize; i += batchGetItemSize {
		getInput := dynamodb.BatchGetItemInput{RequestItems: map[string]*dynamodb.KeysAndAttributes{tableName: {}}}
		for j := i; j < i+batchGetItemSize; j++ {
			key := map[string]*dynamodb.AttributeValue{
				"ID": &dynamodb.AttributeValue{S: aws.String(p[j].ID)}}
			getInput.RequestItems[tableName].Keys = append(getInput.RequestItems[tableName].Keys, key)
		}
		if result, err = svc.BatchGetItem(&getInput); err != nil {
			return nil, err
		}
		var tmp []players.PlayerData
		if err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[tableName], &tmp); err != nil {
			return nil, err
		}
		data = append(data, tmp...)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	{
		getInput := dynamodb.BatchGetItemInput{RequestItems: map[string]*dynamodb.KeysAndAttributes{tableName: {}}}
		for ; i < len(p); i++ {
			key := map[string]*dynamodb.AttributeValue{
				"ID": &dynamodb.AttributeValue{S: aws.String(p[i].ID)}}
			getInput.RequestItems[tableName].Keys = append(getInput.RequestItems[tableName].Keys, key)
		}
		if result, err = svc.BatchGetItem(&getInput); err != nil {
			return nil, err
		}
		var tmp []players.PlayerData
		if err = dynamodbattribute.UnmarshalListOfMaps(result.Responses[tableName], &tmp); err != nil {
			return nil, err
		}
		data = append(data, tmp...)
	}
	return data, nil
}
