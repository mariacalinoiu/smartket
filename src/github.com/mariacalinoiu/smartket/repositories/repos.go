package repositories

const DefaultOrderStatus = "in asteptare"

type (
	Department struct {
		ID   int    `json:"ID"`
		Name string `json:"name"`
	}

	Category struct {
		ID           int    `json:"ID"`
		Name         string `json:"name"`
		DepartmentId int    `json:"departmentID"`
	}

	Order struct {
		ID                 int              `json:"ID"`
		FirstName          string           `json:"firstName"`
		LastName           string           `json:"lastName"`
		Email              string           `json:"email"`
		PhoneNumber        string           `json:"phoneNumber"`
		City               string           `json:"city"`
		Address            string           `json:"address"`
		VoucherCode        string           `json:"voucherCode"`
		DiscountPercentage int              `json:"discountPercentage"`
		PaymentMethod      string           `json:"paymentMethod"`
		Status             string           `json:"status"`
		Timestamp          int              `json:"timestamp"`
		Value              float32          `json:"value"`
		ProductsOrdered    []OrderedProduct `json:"products"`
	}

	OrderedProduct struct {
		ProductID int     `json:"productID"`
		OrderID   int     `json:"orderID"`
		Quantity  int     `json:"quantity"`
		Product   Product `json:"productDetails"`
	}

	Product struct {
		ID          int     `json:"ID"`
		Name        string  `json:"name"`
		ImageURL    string  `json:"imageURL"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
		CategoryID  int     `json:"categoryID"`
	}
)
