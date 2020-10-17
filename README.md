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
    example URL:    http://localhost:8081/categories?departmentID=1


/products
    
    method:         GET
    parameters:     categoryID int
    returns:        a JSON of products in the given categoryID
    example URL:    http://localhost:8081/products?categoryID=1


/orders
    
    method:         GET
    parameters:     orderID int (optional)
    returns:        a JSON of orders
    example URL:    http://localhost:8081/orders
                    http://localhost:8081/orders?orderID=1
    

    method:         POST
    body:           an order, along with ordered product details
    returns:        the corresponding orderID
    example URL:    http://localhost:8081/orders
    

    method:         PUT
    body:           an order
    returns:        the corresponding orderID
    example URL:    http://localhost:8081/orders
    

    method:         DELETE
    parameters:     orderID int
    returns:        -
    example URL:    http://localhost:8081/orders?orderID=1
    
------------------
 
Running the server
------------------

Reach project folder:  `cd ~/environment/smartket/`

Build the server with: `go build src/server.go`

Start running the server with `./server`
