package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mariacalinoiu/smartket/src/datasources"
)

func HandleProducts(w http.ResponseWriter, r *http.Request, db datasources.DBClient, logger *log.Logger) {
	var response []byte
	var status int
	var err error

	switch r.Method {
	case http.MethodGet:
		response, status, err = getProducts(r, db, logger)
	default:
		status = http.StatusBadRequest
		err = errors.New("wrong method type for /products route")
	}

	if err != nil {
		logger.Printf("Error: %s; Status: %d %s", err.Error(), status, http.StatusText(status))
		http.Error(w, err.Error(), status)

		return
	}

	_, err = w.Write(response)
	if err != nil {
		status = http.StatusInternalServerError
		logger.Printf("Error: %s; Status: %d %s", err.Error(), status, http.StatusText(status))
		http.Error(w, err.Error(), status)

		return
	}

	status = http.StatusOK
	logger.Printf("Status: %d %s", status, http.StatusText(status))
}

func getProducts(r *http.Request, db datasources.DBClient, logger *log.Logger) ([]byte, int, error) {
	params, ok := r.URL.Query()["categoryID"]

	if !ok || len(params[0]) < 1 {
		return nil, http.StatusBadRequest, errors.New("mandatory parameter 'categoryID' not found")
	}

	categoryId, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("could not convert parameter 'categoryID' to integer")
	}
	products, err := db.GetProductsByCategoryID(categoryId)
	if err != nil {
		logger.Printf("Internal error: %s", err.Error())
		return nil, http.StatusInternalServerError, errors.New("could not get products in Category")
	}

	response, err := json.Marshal(products)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal products response json")
	}

	return response, http.StatusOK, nil
}
