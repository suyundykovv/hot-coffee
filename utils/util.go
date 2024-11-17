package utils

import (
	"errors"
	"fmt"
	"hot-coffee/models"
	"log"
	"math/rand"
	"time"
)

func CatchCriticalPoint() {
	if r := recover(); r != nil {
		log.Printf("Recovered from error")
	}
}

// validateUpdatedOrder ensures that the updated order is valid
func ValidateUpdatedOrder(order models.Order) error {
	if order.ID == "" {
		return errors.New("order ID cannot be empty")
	}
	if order.CustomerName == "" {
		return errors.New("customer name cannot be empty")
	}
	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	for _, item := range order.Items {
		if item.ProductID == "" {
			return errors.New("item product ID cannot be empty")
		}
		if item.Quantity <= 0 {
			return errors.New("item quantity must be greater than zero")
		}
	}

	validStatuses := []string{"open", "closed"}
	if !contains(validStatuses, order.Status) {
		return fmt.Errorf("invalid order status: %s", order.Status)
	}

	if order.CreatedAt == "" {
		return errors.New("order creation date cannot be empty")
	}

	return nil
}

// Helper function to check if a status is valid
func contains(validStatuses []string, status string) bool {
	for _, validStatus := range validStatuses {
		if validStatus == status {
			return true
		}
	}
	return false
}

func ValidateUpdatedMenuItem(item models.MenuItem) error {
	if item.ID == "" {
		return errors.New("menu item ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("menu item name cannot be empty")
	}
	if item.Price <= 0 {
		return errors.New("menu item price must be greater than zero")
	}
	if item.Ingredients == nil {
		return errors.New("menu item ingredient cannot be empty")
	}
	if item.Description == "" {
		return errors.New("menu item ingredient cannot be empty")
	}
	return nil
}

// validateUpdatedInventoryItem ensures that the updated inventory item is valid
func ValidateUpdatedInventoryItem(item models.InventoryItem) error {
	if item.IngredientID == "" {
		return errors.New("ingredient ID cannot be empty")
	}
	if item.Name == "" {
		return errors.New("ingredient name cannot be empty")
	}
	if item.Quantity < 0 {
		return errors.New("ingredient quantity cannot be negative")
	}
	if item.Unit == "" {
		return errors.New("ingredient unit cannot be empty")
	}

	return nil
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min, max int) int {
	// Initialize the random number generator with the current time as the seed
	rand.Seed(time.Now().UnixNano())

	// Return a random integer between min and max (inclusive)
	return rand.Intn(max-min+1) + min
}
