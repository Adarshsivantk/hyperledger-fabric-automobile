
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CarContract contract for managing CRUD operations for cars.
type CarContract struct {
	contractapi.Contract
}

// Car represents a car asset stored on the ledger.
type Car struct {
	AssetType         string `json:"AssetType"`
	CarId             string `json:"CarId"`
	Color             string `json:"Color"`
	DateOfManufacture string `json:"DateOfManufacture"`
	OwnedBy           string `json:"OwnedBy"`
	Make              string `json:"Make"`
	Model             string `json:"Model"`
	Status            string `json:"Status"`
}

// PaginatedQueryResult represents paginated car query results.
type PaginatedQueryResult struct {
	Records             []*Car `json:"records"`
	FetchedRecordsCount int32  `json:"fetchedRecordsCount"`
	Bookmark            string `json:"bookmark"`
}

// HistoryQueryResult represents a historical car record.
type HistoryQueryResult struct {
	Record    *Car   `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

// EventData represents the payload of a chaincode event.
type EventData struct {
	Type  string
	Model string
}

// CarExists returns true when an asset with the given ID exists.
func (c *CarContract) CarExists(
	ctx contractapi.TransactionContextInterface,
	carID string,
) (bool, error) {
	data, err := ctx.GetStub().GetState(carID)
	if err != nil {
		return false, fmt.Errorf(
			"failed to read from world state: %v",
			err,
		)
	}

	return data != nil, nil
}

// CreateCar creates a new car and emits a chaincode event.
func (c *CarContract) CreateCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
	make string,
	model string,
	color string,
	manufacturerName string,
	dateOfManufacture string,
) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf(
			"user under following MSPID: %v can't perform this action",
			clientOrgID,
		)
	}

	exists, err := c.CarExists(ctx, carID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state: %s", err)
	}

	if exists {
		return "", fmt.Errorf("the asset %s already exists", carID)
	}

	car := Car{
		AssetType:         "car",
		CarId:             carID,
		Color:             color,
		DateOfManufacture: dateOfManufacture,
		Make:              make,
		Model:             model,
		OwnedBy:           manufacturerName,
		Status:            "In Factory",
	}

	carBytes, err := json.Marshal(car)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(carID, carBytes)
	if err != nil {
		return "", err
	}

	// Create event data for the newly created car.
	addCarEventData := EventData{
		Type:  "Car creation",
		Model: model,
	}

	// Convert event data to JSON.
	eventDataByte, err := json.Marshal(addCarEventData)
	if err != nil {
		return "", err
	}

	// Emit the chaincode event.
	err = ctx.GetStub().SetEvent("CreateCar", eventDataByte)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("successfully added car %v", carID), nil
}

// ReadCar retrieves a car from the world state.
func (c *CarContract) ReadCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
) (*Car, error) {
	exists, err := c.CarExists(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state: %s", err)
	}

	if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", carID)
	}

	carBytes, err := ctx.GetStub().GetState(carID)
	if err != nil {
		return nil, err
	}

	car := new(Car)

	err = json.Unmarshal(carBytes, car)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal world state data to type Car: %v",
			err,
		)
	}

	return car, nil
}

// DeleteCar removes a car from the world state.
func (c *CarContract) DeleteCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf(
			"user under following MSPID: %v cannot perform this action",
			clientOrgID,
		)
	}

	exists, err := c.CarExists(ctx, carID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state: %s", err)
	}

	if !exists {
		return "", fmt.Errorf("the asset %s does not exist", carID)
	}

	err = ctx.GetStub().DelState(carID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Car with ID %v is deleted from the world state.",
		carID,
	), nil
}

// GetMatchingOrders retrieves orders matching a car's details.
func (c *CarContract) GetMatchingOrders(
	ctx contractapi.TransactionContextInterface,
	carID string,
) ([]*Order, error) {
	exists, err := c.CarExists(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state: %s", err)
	}

	if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", carID)
	}

	car, err := c.ReadCar(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("error reading car: %v", err)
	}

	queryString := fmt.Sprintf(
		`{"selector":{"assetType":"Order","make":"%s","model":"%s","color":"%s"}}`,
		car.Make,
		car.Model,
		car.Color,
	)

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(
		getCollectionName(),
		queryString,
	)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return orderResultIteratorFunction(resultsIterator)
}

// MatchOrder matches a car with a compatible order.
func (c *CarContract) MatchOrder(
	ctx contractapi.TransactionContextInterface,
	carID string,
	orderID string,
) (string, error) {
	order, err := ReadPrivateState(ctx, orderID)
	if err != nil {
		return "", err
	}

	car, err := c.ReadCar(ctx, carID)
	if err != nil {
		return "", err
	}

	if car.Make != order.Make ||
		car.Color != order.Color ||
		car.Model != order.Model {
		return "", fmt.Errorf("order is not matching")
	}

	car.OwnedBy = order.DealerName
	car.Status = "assigned to a dealer"

	carBytes, err := json.Marshal(car)
	if err != nil {
		return "", err
	}

	collectionName := getCollectionName()

	err = ctx.GetStub().DelPrivateData(collectionName, orderID)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(carID, carBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Deleted order %v and assigned %v to %v",
		orderID,
		car.CarId,
		order.DealerName,
	), nil
}

// RegisterCar registers a car to its buyer.
func (c *CarContract) RegisterCar(
	ctx contractapi.TransactionContextInterface,
	carID string,
	ownerName string,
	registrationNumber string,
) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID != "Org3MSP" {
		return "", fmt.Errorf(
			"user under following MSPID: %v cannot perform this action",
			clientOrgID,
		)
	}

	exists, err := c.CarExists(ctx, carID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state: %s", err)
	}

	if !exists {
		return "", fmt.Errorf("car %v does not exist", carID)
	}

	car, err := c.ReadCar(ctx, carID)
	if err != nil {
		return "", err
	}

	car.Status = fmt.Sprintf(
		"Registered to %v with plate number %v",
		ownerName,
		registrationNumber,
	)
	car.OwnedBy = ownerName

	carBytes, err := json.Marshal(car)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(carID, carBytes)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Car %v successfully registered to %v",
		carID,
		ownerName,
	), nil
}

// GetAllCars retrieves all assets with AssetType "car".
func (c *CarContract) GetAllCars(
	ctx contractapi.TransactionContextInterface,
) ([]*Car, error) {
	queryString := `{"selector":{"AssetType":"car"},"sort":[{"CarId":"desc"}]}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return carResultIteratorFunction(resultsIterator)
}

// GetCarsByRange retrieves cars between a start key and an end key.
func (c *CarContract) GetCarsByRange(
	ctx contractapi.TransactionContextInterface,
	startKey string,
	endKey string,
) ([]*Car, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(
		startKey,
		endKey,
	)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return carResultIteratorFunction(resultsIterator)
}

// GetCarsWithPagination retrieves cars using pagination.
func (c *CarContract) GetCarsWithPagination(
	ctx contractapi.TransactionContextInterface,
	pageSize int32,
	bookmark string,
) (*PaginatedQueryResult, error) {
	queryString := `{"selector":{"AssetType":"car"}}`

	resultsIterator, responseMetadata, err :=
		ctx.GetStub().GetQueryResultWithPagination(
			queryString,
			pageSize,
			bookmark,
		)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	cars, err := carResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             cars,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetCarHistory returns the history of a car since its creation.
func (c *CarContract) GetCarHistory(
	ctx contractapi.TransactionContextInterface,
	carID string,
) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(carID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car

		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &car)
			if err != nil {
				return nil, err
			}
		} else {
			car = Car{
				CarId: carID,
			}
		}

		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &car,
			IsDelete:  response.IsDelete,
		}

		records = append(records, &record)
	}

	return records, nil
}

// carResultIteratorFunction converts query results into Car objects.
func carResultIteratorFunction(
	resultsIterator shim.StateQueryIteratorInterface,
) ([]*Car, error) {
	var cars []*Car

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car

		err = json.Unmarshal(queryResult.Value, &car)
		if err != nil {
			return nil, err
		}

		cars = append(cars, &car)
	}

	return cars, nil
}