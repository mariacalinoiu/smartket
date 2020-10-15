package main

import (
	"database/sql"
	"fmt"
	"time"
)

type DBClient struct {
	db *sql.DB
}

func GetClient(user string, password string, dbName string) DBClient {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@/%s", user, password, dbName),
	)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3000)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	return DBClient{db: db}
}

func (client DBClient) getProductsByCategoryID(categoryID int) ([]product, error) {
	var products []product
	var id int
	var name string
	var imageURL string
	var description string
	var price float32

	rows, err := client.db.Query(
		"SELECT ID, name, imageURL, description, price FROM Products WHERE categoryID = ?",
		categoryID,
	)
	if err != nil {
		return products, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &imageURL, &description, &price)
		if err != nil {
			return products, err
		}

		products = append(
			products,
			product{
				ID:          id,
				Name:        name,
				ImageURL:    imageURL,
				Description: description,
				Price:       price,
				CategoryID:  categoryID,
			},
		)
	}

	err = rows.Err()
	if err != nil {
		return products, err
	}

	return products, nil
}

func (client DBClient) getCategoriesByDepartmentID(departmentID int) ([]category, error) {
	var categories []category
	var id int
	var name string

	rows, err := client.db.Query(
		"SELECT ID, name FROM Categories WHERE departmentID = ?",
		departmentID,
	)
	if err != nil {
		return categories, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return categories, err
		}

		categories = append(
			categories,
			category{
				ID:           id,
				Name:         name,
				DepartmentId: departmentID,
			},
		)
	}

	err = rows.Err()
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (client DBClient) getDepartments() ([]department, error) {
	var departments []department
	var id int
	var name string

	rows, err := client.db.Query(
		"SELECT ID, name FROM Departments",
	)
	if err != nil {
		return departments, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return departments, err
		}

		departments = append(
			departments,
			department{
				ID:   id,
				Name: name,
			},
		)
	}

	err = rows.Err()
	if err != nil {
		return departments, err
	}

	return departments, nil
}
