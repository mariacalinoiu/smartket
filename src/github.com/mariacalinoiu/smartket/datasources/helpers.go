package datasources

import (
	"time"

	"github.com/mariacalinoiu/smartket/repositories"
)

func GetOrderID(orderID int) repositories.OrderIDResponse {
	return repositories.OrderIDResponse{OrderID: orderID}
}

func ParseTimestamp(timestamp int) string {
	tm := time.Unix(int64(timestamp), 0)
	layout := "2006-01-02 15:04:05"

	return tm.Format(layout)
}
