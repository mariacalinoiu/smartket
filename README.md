Server documentation
------------------

/departments
    
    method:         GET
    parameters:     -
    returns:        a JSON of departments
    example URL:    http://localhost:8081/departments


/categories
    
    method:         GET
    parameters:     departmentID int
    returns:        a JSON of categories in the given departmentID
    example URL:    http://localhost:8081/categories


/products
    
    method:         GET
    parameters:     categoryID int
    returns:        a JSON of products in the given categoryID
    example URL:    http://localhost:8081/products


/orders
    
    method:         GET
    parameters:     -
    returns:        a JSON of orders
    example URL:    http://localhost:8080/orders
    

    method:         POST
    body:           an order, along with ordered product IDs
    returns:        the corresponding orderID
    example URL:    http://localhost:8080/orders
    

    method:         PUT
    body:           an order
    returns:        the corresponding orderID
    example URL:    http://localhost:8080/orders
    

    method:         DELETE
    parameters:     orderID int
    returns:        -
    example URL:    http://localhost:8080/orders
    
------------------
 
Running the server
------------------

Run the server with: `cd src && go run .`
