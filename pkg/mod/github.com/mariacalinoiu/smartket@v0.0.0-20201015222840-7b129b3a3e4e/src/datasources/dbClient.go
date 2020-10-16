package datasources

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"repositories"
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

func (client DBClient) GetProductsByCategoryID(categoryID int) ([]repositories.Product, error) {
	var (
		products    []repositories.Product
		id          int
		name        string
		imageURL    string
		description string
		price       float32
	)

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
			repositories.Product{
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

func (client DBClient) GetCategoriesByDepartmentID(departmentID int) ([]repositories.Category, error) {
	var (
		categories []repositories.Category
		id         int
		name       string
	)

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
			repositories.Category{
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

func (client DBClient) GetDepartments() ([]repositories.Department, error) {
	var (
		departments []repositories.Department
		id          int
		name        string
	)

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
			repositories.Department{
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

func (client DBClient) InsertOrder(order repositories.Order) (int, error) {
	isVoucherValid := client.isVoucherValid(order.VoucherCode)
	if !isVoucherValid {
		return 0, errors.New("the voucher code provided is invalid")
	}

	stmt, err := client.db.Prepare("INSERT INTO Orders(firstName, lastName, email, phoneNumber, city, address, voucherCode, paymentMethod, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(
		order.FirstName,
		order.LastName,
		order.Email,
		order.PhoneNumber,
		order.City,
		order.Address,
		order.VoucherCode,
		order.PaymentMethod,
		order.Status,
	)
	if err != nil {
		return 0, err
	}
	orderID, err := res.LastInsertId()

	stmt, err = client.db.Prepare("INSERT INTO ProductOrders(orderID, productID, quantity) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}

	for _, product := range order.ProductsOrdered {
		_, err = stmt.Exec(
			orderID,
			product.ProductID,
			product.Quantity,
		)
		if err != nil {
			return 0, err
		}
	}

	return int(orderID), nil
}

func (client DBClient) EditOrder(order repositories.Order) error {
	isVoucherValid := client.isVoucherValid(order.VoucherCode)
	if !isVoucherValid {
		return errors.New("the voucher code provided is invalid")
	}

	stmt, err := client.db.Prepare("REPLACE INTO Orders(ID, firstName, lastName, email, phoneNumber, city, address, voucherCode, paymentMethod, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		order.ID,
		order.FirstName,
		order.LastName,
		order.Email,
		order.PhoneNumber,
		order.City,
		order.Address,
		order.VoucherCode,
		order.PaymentMethod,
		order.Status,
	)

	return err
}

func (client DBClient) DeleteOrder(orderID int) error {
	_, err := client.db.Exec(
		"DELETE FROM ProductOrders WHERE orderID = ?",
		orderID,
	)
	if err != nil {
		return err
	}

	_, err = client.db.Exec(
		"DELETE FROM Orders WHERE ID = ?",
		orderID,
	)

	return err
}

func (client DBClient) GetOrders(orderIDProvided ...int) ([]repositories.Order, error) {
	var (
		orderRows *sql.Rows
		err       error

		orders             []repositories.Order
		orderID            int
		firstName          string
		lastName           string
		email              string
		phoneNumber        string
		city               string
		address            string
		voucherCode        string
		paymentMethod      string
		status             string
		discountPercentage int
	)

	if len(orderIDProvided) == 1 {
		orderRows, err = client.db.Query(
			`
				SELECT o.ID, o.firstName, o.lastName, o.email, o.phoneNumber, o.city, o.address, o.voucherCode, o.paymentMethod, o.status, v.discountPercentage 
				FROM Orders o, Vouchers v
				WHERE o.voucherCode = v.code AND ID = ?
			`,
			orderIDProvided[0],
		)
	} else {
		orderRows, err = client.db.Query(
			`
				SELECT o.ID, o.firstName, o.lastName, o.email, o.phoneNumber, o.city, o.address, o.voucherCode, o.paymentMethod, o.status, v.discountPercentage 
				FROM Orders o, Vouchers v
				WHERE o.voucherCode = v.code
			`,
		)
	}
	if err != nil {
		return orders, err
	}

	defer orderRows.Close()
	for orderRows.Next() {
		err := orderRows.Scan(&orderID, &firstName, &lastName, &email, &phoneNumber, &city, &address, &voucherCode, &paymentMethod, &status, &discountPercentage)
		if err != nil {
			return orders, err
		}
		products, err := client.getOrderedProducts(orderID)
		if err != nil {
			return orders, err
		}

		orders = append(
			orders,
			repositories.Order{
				ID:                 orderID,
				FirstName:          firstName,
				LastName:           lastName,
				Email:              email,
				PhoneNumber:        phoneNumber,
				City:               city,
				Address:            address,
				VoucherCode:        voucherCode,
				DiscountPercentage: discountPercentage,
				PaymentMethod:      paymentMethod,
				Status:             status,
				ProductsOrdered:    products,
			},
		)
	}

	err = orderRows.Err()
	if err != nil {
		return orders, err
	}

	return orders, nil
}

func (client DBClient) getOrderedProducts(orderID int) ([]repositories.OrderedProduct, error) {
	var (
		products  []repositories.OrderedProduct
		productID int
		quantity  int
	)
	productOrderRows, err := client.db.Query(
		"SELECT productID, quantity FROM ProductOrders WHERE orderID = ?",
		orderID,
	)
	if err != nil {
		return products, err
	}

	for productOrderRows.Next() {
		err := productOrderRows.Scan(&productID, &quantity)
		if err != nil {
			return products, err
		}

		products = append(
			products,
			repositories.OrderedProduct{
				ProductID: productID,
				OrderID:   orderID,
				Quantity:  quantity,
			},
		)

		productOrderRows.Close()
	}

	err = productOrderRows.Err()
	if err != nil {
		return products, err
	}

	return products, nil
}

func (client DBClient) isVoucherValid(voucherCode string) bool {
	rows, err := client.db.Query(
		"SELECT discountPercentage FROM Vouchers WHERE code = ?",
		voucherCode,
	)
	if err != nil {
		return false
	}

	defer rows.Close()
	for rows.Next() {
		return true
	}

	return false
}
