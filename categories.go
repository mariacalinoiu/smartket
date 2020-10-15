package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type category struct {
	ID           int    `json:"ID"`
	Name         string `json:"name"`
	DepartmentId int    `json:"departmentID"`
}

func handleCategories(w http.ResponseWriter, r *http.Request, db DBClient, logger *log.Logger) {
	var response []byte
	var status int
	var err error

	switch r.Method {
	case http.MethodGet:
		response, status, err = getCategories(r, db)
	default:
		status = http.StatusBadRequest
		err = errors.New("wrong method type for /categories route")
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

func getCategories(r *http.Request, db DBClient) ([]byte, int, error) {
	params, ok := r.URL.Query()["departmentID"]

	if !ok || len(params[0]) < 1 {
		return nil, http.StatusBadRequest, errors.New("mandatory parameter 'departmentID' not found")
	}

	categoryId, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("could not convert parameter 'departmentID' to integer")
	}
	categories, err := db.getCategoriesByDepartmentID(categoryId)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not get categories in department")
	}

	response, err := json.Marshal(categories)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal categories response json")
	}

	return response, http.StatusOK, nil
}
