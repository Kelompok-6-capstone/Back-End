package helper

import "fmt"

func ParseOrderID(orderID string) (int, error) {
	var consultationID int
	_, err := fmt.Sscanf(orderID, "consultation-%d", &consultationID)
	if err != nil {
		return 0, fmt.Errorf("invalid order ID format")
	}
	return consultationID, nil
}
