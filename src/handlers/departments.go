package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mariacalinoiu/smartket/src/datasources"
)

func HandleDepartments(w http.ResponseWriter, r *http.Request, db datasources.DBClient, logger *log.Logger) {
	var response []byte
	var status int
	var err error

	switch r.Method {
	case http.MethodGet:
		response, status, err = getDepartments(db)
	default:
		status = http.StatusBadRequest
		err = errors.New("wrong method type for /departments route")
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

func getDepartments(db datasources.DBClient) ([]byte, int, error) {
	departments, err := db.GetDepartments()
	if err != nil {
		logger.Printf("Internal error: %s", err.Error())
		return nil, http.StatusInternalServerError, errors.New("could not get departments")
	}

	response, err := json.Marshal(departments)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("could not marshal departments response json")
	}

	return response, http.StatusOK, nil
}
