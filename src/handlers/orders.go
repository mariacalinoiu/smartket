package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"datasources"
	"repositories"
)

func HandleOrders(w http.ResponseWriter, r *http.Request, db datasources.DBClient, logger *log.Logger) {
	var response []byte
	var status int
	var err error

	switch r.Method {
	case http.MethodGet:
		response, status, err = getOrders(db)
	case http.MethodPost, http.MethodPut:
		response, status, err = insertOrder(r, db)
	case http.MethodDelete:
		status, err = deleteOrder(r, db)
	default:
		status = http.StatusBadRequest
		err = errors.New("wrong method type for /orders route")
	}

	if err != nil {
		logger.Printf("Status: %d %s", status, http.StatusText(status))
		http.Error(w, err.Error(), status)

		return
	}

	_, err = w.Write(response)
	if err != nil {
		status = http.StatusInternalServerError
		logger.Printf("Status: %d %s", status, http.StatusText(status))
		http.Error(w, err.Error(), status)

		return
	}

	status = http.StatusOK
	logger.Printf("Status: %d %s", status, http.StatusText(status))
}

func getOrders(db datasources.DBClient) ([]byte, int, error) {
	orders, err := db.GetOrders()
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not get orders")
	}

	response, err := json.Marshal(orders)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal orders response json")
	}

	return response, http.StatusOK, nil
}

func extractOrderParams(r *http.Request) (repositories.Order, error) {
	var unmarshalledOrder repositories.Order

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return repositories.Order{}, err
	}

	err = json.Unmarshal(body, &unmarshalledOrder)
	if err != nil {
		return repositories.Order{}, err
	}

	return unmarshalledOrder, nil
}

func insertOrder(r *http.Request, db datasources.DBClient) ([]byte, int, error) {
	var orderID int

	order, err := extractOrderParams(r)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Order information sent on request body does not match required format")
	}

	if r.Method == http.MethodPost {
		orderID, err = db.InsertOrder(order)
	} else {
		err = db.EditOrder(order)
		orderID = order.ID
	}
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not save Order")
	}

	response, err := json.Marshal(fmt.Sprintf("orderID:%d", orderID))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal orderID response json")
	}

	return response, http.StatusOK, nil
}

func deleteOrder(r *http.Request, db datasources.DBClient) (int, error) {
	params, ok := r.URL.Query()["orderID"]

	if !ok || len(params[0]) < 1 {
		return http.StatusBadRequest, errors.New("mandatory parameter 'orderID' not found")
	}

	orderID, err := strconv.Atoi(params[0])
	if err != nil {
		return http.StatusBadRequest, errors.New("could not convert parameter 'orderID' to integer")
	}
	err = db.DeleteOrder(orderID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("could not delete Order")
	}

	return http.StatusOK, nil
}
