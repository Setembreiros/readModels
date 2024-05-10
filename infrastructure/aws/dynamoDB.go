package aws

import (
	"context"
	"errors"
	"log"
	"time"

	database "readmodels/infrastructure/db"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	client   *dynamodb.Client
}

func NewDynamodbClient(config aws.Config, infoLog, errorLog *log.Logger) *DynamoDBClient {
	return &DynamoDBClient{
		infoLog:  infoLog,
		errorLog: errorLog,
		client:   dynamodb.NewFromConfig(config),
	}
}

func (dc DynamoDBClient) TableExists(tableName string) bool {
	exists := true
	_, err := dc.client.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
		} else {
			dc.errorLog.Printf("Couldn't determine existence of table %v. error: %v\n", tableName, err)
		}
		exists = false
	}
	return exists
}

func (dc DynamoDBClient) CreateTable(tableName string, attributes []database.TableAttributes, ctx context.Context) error {
	keySchemas, attributeDefinitions, err := mapTableAttributes(attributes)
	if err != nil {
		return err
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: attributeDefinitions,
		KeySchema:            keySchemas,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}

	table, err := dc.client.CreateTable(ctx, input)

	var tableDesc *types.TableDescription
	if err != nil {
		dc.errorLog.Fatalf("Couldn't create table %v. error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(dc.client)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. error: %v\n", err)
		}
		tableDesc = table.TableDescription
	}

	dc.infoLog.Printf("Created table: %s\n", *tableDesc.TableName)
	return nil
}

func mapTableAttributes(attributes []database.TableAttributes) ([]types.KeySchemaElement, []types.AttributeDefinition, error) {
	var keySchemas []types.KeySchemaElement
	var attributeDefinitions []types.AttributeDefinition
	isPartitionKey := true

	for _, attribute := range attributes {
		keySchema, attributeDefinition, err := mapTableAttribute(attribute, isPartitionKey)
		if err != nil {
			return nil, nil, err
		}
		keySchemas = append(keySchemas, *keySchema)
		attributeDefinitions = append(attributeDefinitions, *attributeDefinition)
		isPartitionKey = false
	}

	return keySchemas, attributeDefinitions, nil
}

func mapTableAttribute(attribute database.TableAttributes, isPartitionKey bool) (*types.KeySchemaElement, *types.AttributeDefinition, error) {
	attributeType, err := mapAttributeType(attribute.AttributeType)
	if err != nil {
		return nil, nil, err
	}

	var keyType types.KeyType
	if isPartitionKey {
		keyType = types.KeyTypeHash
	} else {
		keyType = types.KeyTypeRange
	}

	return &types.KeySchemaElement{
			AttributeName: aws.String(attribute.Name),
			KeyType:       keyType,
		},
		&types.AttributeDefinition{
			AttributeName: aws.String(attribute.Name),
			AttributeType: attributeType,
		}, nil
}

func mapAttributeType(attributeType string) (types.ScalarAttributeType, error) {
	switch attributeType {
	case "number":
		return types.ScalarAttributeTypeN, nil
	case "string":
		return types.ScalarAttributeTypeS, nil
	case "binary":
		return types.ScalarAttributeTypeB, nil
	default:
		return "", errors.New("attribute type " + attributeType + " doesn't exist")
	}
}