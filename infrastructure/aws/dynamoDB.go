package aws

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	database "readmodels/internal/db"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog/log"
)

type DynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamodbClient(config aws.Config) *DynamoDBClient {
	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(config),
	}
}

// Clean deletes all tables from DynamoDB.
// This method is destructive and will remove all tables and data.
// Errors are logged but not returned.
func (dc *DynamoDBClient) Clean() {
	ctx := context.TODO()
	log.Info().Msg("Starting DynamoDB clean process")

	var nextToken *string
	tablesToDelete := []string{}

	// List all tables
	for {
		listTablesInput := &dynamodb.ListTablesInput{
			ExclusiveStartTableName: nextToken,
		}

		response, err := dc.client.ListTables(ctx, listTablesInput)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Failed to list DynamoDB tables")
			return
		}

		tablesToDelete = append(tablesToDelete, response.TableNames...)

		// Check if we need to continue pagination
		if response.LastEvaluatedTableName == nil {
			break
		}
		nextToken = response.LastEvaluatedTableName
	}

	if len(tablesToDelete) == 0 {
		log.Info().Msg("No tables found to delete")
		return
	}

	log.Info().Msgf("Found %d tables to delete", len(tablesToDelete))

	// Delete all tables
	failedCount := 0
	for _, tableName := range tablesToDelete {
		log.Info().Msgf("Deleting table: %s", tableName)

		_, err := dc.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})

		if err != nil {
			log.Error().Stack().Err(err).Msgf("Failed to delete table %s", tableName)
			failedCount++
			continue
		}

		// Wait for table deletion to complete
		waiter := dynamodb.NewTableNotExistsWaiter(dc.client)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}, 5*time.Minute)

		if err != nil {
			log.Warn().Err(err).Msgf("Wait for table deletion failed for %s", tableName)
			continue
		}

		log.Info().Msgf("Table %s successfully deleted", tableName)
	}

	if failedCount > 0 {
		log.Error().Msgf("Failed to delete %d tables", failedCount)
	} else {
		log.Info().Msgf("Successfully cleaned up %d DynamoDB tables", len(tablesToDelete))
	}
}

func (dc *DynamoDBClient) TableExists(tableName string) bool {
	exists := true
	_, err := dc.client.DescribeTable(
		context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
		} else {
			log.Error().Stack().Err(err).Msgf("Couldn't determine existence of table %v", tableName)
		}
		exists = false
	}
	return exists
}

func (dc *DynamoDBClient) CreateTable(tableName string, keys *[]database.TableAttributes, ctx context.Context) error {
	keySchemas, attributeDefinitions, err := mapTableKeys(keys)
	if err != nil {
		return err
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: *attributeDefinitions,
		KeySchema:            *keySchemas,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}

	table, err := dc.client.CreateTable(ctx, input)

	var tableDesc *types.TableDescription
	if err != nil {
		log.Fatal().Stack().Err(err).Msgf("Couldn't create table %v", tableName)
		return err
	} else {
		waiter := dynamodb.NewTableExistsWaiter(dc.client)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Warn().Err(err).Msgf("Wait for table exists failed")
		}
		tableDesc = table.TableDescription
	}

	log.Info().Msgf("Created table: %s", *tableDesc.TableName)
	return nil
}

func (dc *DynamoDBClient) CreateIndexesOnTable(tableName, indexName string, indexes *[]database.TableAttributes, ctx context.Context) error {
	keySchemas, attributeDefinitions, err := mapTableKeys(indexes)
	if err != nil {
		return err
	}

	gsi := types.GlobalSecondaryIndexUpdate{
		Create: &types.CreateGlobalSecondaryIndexAction{
			IndexName: aws.String(indexName),
			KeySchema: *keySchemas,
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	}

	input := &dynamodb.UpdateTableInput{
		TableName:                   aws.String(tableName),
		AttributeDefinitions:        *attributeDefinitions,
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{gsi},
	}

	_, err = dc.client.UpdateTable(ctx, input)

	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to update table")
	}

	log.Info().Msgf("GSI %s created on table %s", indexName, tableName)
	return nil
}

func (dc *DynamoDBClient) InsertData(tableName string, attributes any) error {
	item, err := attributevalue.MarshalMap(attributes)
	if err != nil {
		return err
	}

	_, err = dc.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't put item to table %s", tableName)
		return err
	}

	return nil
}

func (dc *DynamoDBClient) GetData(tableName string, key any, result any) error {
	k, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't map %v key to AttributeValues", key)
	}

	response, err := dc.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: k, TableName: aws.String(tableName),
	})
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't get info about %s", tableName)
		return err
	}
	if response.Item == nil {
		err = database.NewNotFoundError(tableName, key)
		log.Error().Stack().Err(err).Msg("Item was not found")
		return err
	}

	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't unmarshal response")
		return err
	}

	return nil
}

func (dc *DynamoDBClient) GetMultipleData(tableName string, keys []any, results any) error {
	if len(keys) == 0 {
		return nil
	}

	err := validateIsPointerToSlice(results)
	if err != nil {
		return err
	}

	keyItems := make([]map[string]types.AttributeValue, 0, len(keys))

	keysAndAttributes := &types.KeysAndAttributes{
		Keys: keyItems,
	}

	for _, key := range keys {
		k, err := attributevalue.MarshalMap(key)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Couldn't map key %v to AttributeValues", key)
			return err
		}
		keysAndAttributes.Keys = append(keysAndAttributes.Keys, k)
	}

	requestItems := map[string]types.KeysAndAttributes{
		tableName: *keysAndAttributes,
	}

	response, err := dc.client.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: requestItems,
	})

	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't get batch info from %s", tableName)
		return err
	}

	responseItems, ok := response.Responses[tableName]
	if !ok || len(responseItems) == 0 {
		err = database.NewNotFoundError(tableName, keys)
		log.Error().Stack().Err(err).Msg("No items were found")
		return err
	}
	resultsVal := reflect.ValueOf(results)

	sliceType := resultsVal.Elem().Type()
	elemType := sliceType.Elem()
	newSlice := reflect.MakeSlice(sliceType, 0, len(responseItems))

	// Process each item in the response
	for _, item := range responseItems {
		// Create a new element of the slice's element type
		newElem := reflect.New(elemType).Interface()

		// Unmarshal the item into the new element
		err = attributevalue.UnmarshalMap(item, &newElem)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Couldn't unmarshal response item")
			// We don't know which key this corresponds to, so we can't add to errorMap
			continue
		}

		// Add the element to the new slice
		newSlice = reflect.Append(newSlice, reflect.ValueOf(newElem).Elem())
	}

	resultsVal.Elem().Set(newSlice)

	return nil
}

func (dc *DynamoDBClient) RemoveMultipleData(tableName string, keys []any) error {
	writeRequests := make([]types.WriteRequest, len(keys))
	for i, key := range keys {
		k, err := attributevalue.MarshalMap(key)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Couldn't map %v key to AttributeValues", key)
		}
		writeRequests[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: k,
			},
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeRequests,
		},
	}

	_, err := dc.client.BatchWriteItem(context.TODO(), input)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to batch delete items %v from table %s", keys, tableName)
		return err
	}

	return nil
}

func (dc *DynamoDBClient) UpdateData(tableName string, key any, updateAttributes map[string]any) error {
	k, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't map %v key to AttributeValues", key)
		return err
	}

	updateExp := "set "
	expAttrNames := make(map[string]string)
	expAttrValues := make(map[string]types.AttributeValue)

	i := 0
	for attrName, attrValue := range updateAttributes {
		placeholder := ":val" + string(rune(97+i)) // :vala, :valb, :valc, etc.
		nameHolder := "#n" + string(rune(97+i))    // #na, #nb, #nc, etc.

		// Engadir o atributo á expresión
		if i > 0 {
			updateExp += ", "
		}
		updateExp += nameHolder + " = " + placeholder

		// Engadir o nome e valor do atributo ás expresións
		expAttrNames[nameHolder] = attrName

		// Convertir o valor a tipo DynamoDB
		av, err := attributevalue.Marshal(attrValue)
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Couldn't map %v key to AttributeValues", attrValue)
			return err
		}
		expAttrValues[placeholder] = av

		i++
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       k,
		UpdateExpression:          aws.String(updateExp),
		ExpressionAttributeNames:  expAttrNames,
		ExpressionAttributeValues: expAttrValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	// Executar a operación
	result, err := dc.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't update the element in the table %s", tableName)
		return err
	}

	log.Info().Msgf("Element correctly updated: %v", result.Attributes)
	return nil
}

func (dc *DynamoDBClient) IncrementCounter(tableName string, key any, counterFieldName string, incrementValue int) error {
	k, err := attributevalue.MarshalMap(key)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Non se puido converter a clave %v a AttributeValues", key)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(tableName),
		Key:              k,
		UpdateExpression: aws.String(fmt.Sprintf("set #field = #field + :val")),
		ExpressionAttributeNames: map[string]string{
			"#field": counterFieldName,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberN{Value: strconv.Itoa(incrementValue)},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	_, err = dc.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Non se puido actualizar o contador na táboa %s", tableName)
		return err
	}

	return nil
}

func (dc *DynamoDBClient) GetPostsByIndexUser(username string, lastPostId, lastPostCreatedAt string, limit int) ([]*database.PostMetadata, string, string, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("PostMetadata"),
		IndexName:              aws.String("UserIndex"),
		KeyConditionExpression: aws.String("#user = :user"),
		ExpressionAttributeNames: map[string]string{
			"#user": "Username",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user": &types.AttributeValueMemberS{Value: username},
		},
		Limit: aws.Int32(int32(limit)),
	}

	if lastPostId != "" {
		input.ExclusiveStartKey = map[string]types.AttributeValue{
			"Username":  &types.AttributeValueMemberS{Value: username},
			"PostId":    &types.AttributeValueMemberS{Value: lastPostId},
			"CreatedAt": &types.AttributeValueMemberS{Value: lastPostCreatedAt},
		}
	}

	response, err := dc.client.Query(context.TODO(), input)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't get info about Posts")
		return nil, "", "", err
	}

	var results []*database.PostMetadata
	for _, item := range response.Items {

		var result database.PostMetadata
		err = attributevalue.UnmarshalMap(item, &result)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Couldn't unmarshal response")
			return nil, "", "", err
		}

		results = append(results, &result)
	}

	lastPostId = ""
	lastPostCreatedAt = ""
	if response.LastEvaluatedKey != nil {
		if val, ok := response.LastEvaluatedKey["PostId"]; ok {
			if postId, ok := val.(*types.AttributeValueMemberS); ok {
				lastPostId = postId.Value
			}
		}
		if val, ok := response.LastEvaluatedKey["CreatedAt"]; ok {
			if postCreatedAt, ok := val.(*types.AttributeValueMemberS); ok {
				lastPostCreatedAt = postCreatedAt.Value
			}
		}
	}

	return results, lastPostId, lastPostCreatedAt, nil
}

func mapTableKeys(keys *[]database.TableAttributes) (*[]types.KeySchemaElement, *[]types.AttributeDefinition, error) {
	var keySchemas []types.KeySchemaElement
	var attributeDefinitions []types.AttributeDefinition
	isPartitionKey := true

	for _, key := range *keys {
		keySchema, attributeDefinition, err := mapTableKey(key, isPartitionKey)
		if err != nil {
			return nil, nil, err
		}
		keySchemas = append(keySchemas, *keySchema)
		attributeDefinitions = append(attributeDefinitions, *attributeDefinition)
		isPartitionKey = false
	}

	return &keySchemas, &attributeDefinitions, nil
}

func mapTableKey(key database.TableAttributes, isPartitionKey bool) (*types.KeySchemaElement, *types.AttributeDefinition, error) {
	attributeType, err := mapAttributeType(key.AttributeType)
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
			AttributeName: aws.String(key.Name),
			KeyType:       keyType,
		},
		&types.AttributeDefinition{
			AttributeName: aws.String(key.Name),
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

func validateIsPointerToSlice(results any) error {
	resultsVal := reflect.ValueOf(results)

	if resultsVal.Kind() != reflect.Ptr || resultsVal.Elem().Kind() != reflect.Slice {
		err := database.NewInvalidSlicePointerError(string(resultsVal.Kind()))
		log.Error().Stack().Err(err).Msg("Invalid results parameter")
		return err
	}

	return nil
}

func makeResultSlice(results any, lenght int) reflect.Value {
	resultsVal := reflect.ValueOf(results)
	sliceType := resultsVal.Elem().Type()
	return reflect.MakeSlice(sliceType, 0, lenght)
}

func makeResultElement(results any) reflect.Type {
	resultsVal := reflect.ValueOf(results)
	sliceType := resultsVal.Elem().Type()
	return sliceType.Elem()
}
