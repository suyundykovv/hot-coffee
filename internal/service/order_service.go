package service

import (
	"errors"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
	"time"

	"hot-coffee/logging" // Import the logging package
)

type OrderService interface {
	CreateOrder(order models.Order) error
	FetchAllOrders() ([]models.Order, error)
	FindOrderByID(id string) (models.Order, error)
	UpdateOrderByID(id string, updatedOrder models.Order) error
	DeleteOrderByID(id string) error
	CloseOrder(orderID string) error
	TotalSalesCount() (map[string]int, error)
}

type orderService struct {
	orderRepo dal.OrderRepository
}

func NewOrderService(orderRepo dal.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) CreateOrder(order models.Order) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to create order", "customerName", order.CustomerName)

	// Read the current orders
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read orders", err)
		return err
	}

	// If there are existing orders, generate the next ID based on the last order ID
	if len(orders) > 0 {
		// Get the last order ID (assuming the order ID is numeric or can be parsed as such)
		lastOrder := orders[len(orders)-1]
		// Assuming the order ID is numeric (e.g., "order123" => 123), adjust based on your ID format.
		var lastOrderID int
		_, err := fmt.Sscanf(lastOrder.ID, "order%d", &lastOrderID)
		if err != nil {
			logging.Error("Failed to parse last order ID", err)
			return err
		}

		// Increment the last order ID
		order.ID = fmt.Sprintf("order%d", lastOrderID+1)
	} else {
		// If no orders exist, start with the first order (order1)
		order.ID = "order1"
	}

	// Set the default status if not provided
	if order.Status == "" {
		order.Status = "open" // Set default status to "open"
	}

	// Set the created_at field to the current time minus 1 hour in Nur-Sultan (Asia/Almaty time zone)
	if order.CreatedAt == "" {
		loc, err := time.LoadLocation("Asia/Almaty") // Load the Asia/Almaty time zone (for Nur-Sultan)
		if err != nil {
			logging.Error("Failed to load time zone", err)
			return err
		}
		// Get current time in Nur-Sultan time zone and subtract one hour
		order.CreatedAt = time.Now().In(loc).Add(-time.Hour).Format(time.RFC3339)
	}

	// Append the new order to the existing list
	orders = append(orders, order)

	// Save the updated orders list back to the file
	if err := s.orderRepo.SaveItems(orders); err != nil {
		logging.Error("Failed to save new order", err)
		return err
	}

	logging.Info("Successfully created order", "orderID", order.ID)
	return nil
}

func (s *orderService) FetchAllOrders() ([]models.Order, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching all orders")

	// Fetch all orders
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch orders", err)
		return nil, err
	}

	logging.Info("Fetched all orders", "count", len(orders))
	return orders, nil
}

func (s *orderService) FindOrderByID(id string) (models.Order, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching order by ID", "orderID", id)

	// Fetch all orders
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch orders", err)
		return models.Order{}, err
	}

	// Search for the order by ID
	for _, order := range orders {
		if order.ID == id {
			logging.Info("Found order by ID", "orderID", id)
			return order, nil
		}
	}

	logging.Warn("Order not found", "orderID", id)
	return models.Order{}, errors.New("order not found")
}

func (s *orderService) UpdateOrderByID(id string, updatedOrder models.Order) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to update order", "orderID", id)

	// Check if the updated order contains valid data
	if err := utils.ValidateUpdatedOrder(updatedOrder); err != nil {
		logging.Warn("Invalid updated order data", "orderID", id, "error", err)
		return err
	}

	// Read current orders from the repository
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read orders", err)
		return err
	}

	// Iterate over the orders to find the order by ID
	for i, order := range orders {
		if order.ID == id {
			// Check if the updated status is "closed" before performing the update
			// If the order is already closed, prevent further modifications
			if order.Status == "closed" {
				logging.Warn("Order is already closed and cannot be modified", "orderID", id)
				return errors.New("order is already closed and cannot be modified")
			}

			// Apply the updates if status is valid
			orders[i] = updatedOrder

			// Save the updated list back to the repository
			if err := s.orderRepo.SaveItems(orders); err != nil {
				logging.Error("Failed to save updated order", err)
				return err
			}

			logging.Info("Successfully updated order", "orderID", id)
			return nil
		}
	}

	logging.Warn("Order not found for update", "orderID", id)
	return errors.New("order not found")
}

func (s *orderService) DeleteOrderByID(id string) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to delete order", "orderID", id)

	// Read current orders
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read orders", err)
		return err
	}
	for _, order := range orders {
		if order.ID == id {
			if order.Status == "closed" {
				logging.Warn("Order is already closed and cannot be deleted", "orderID", id)
				return errors.New("order is already closed and cannot be deleted")
			}
		}
	}
	// Create a new list excluding the order to be deleted
	var newOrders []models.Order
	for _, order := range orders {
		if order.ID != id {
			newOrders = append(newOrders, order)
		}
	}
	if len(newOrders) == len(orders) {
		logging.Warn("Orders item not found for deletion", "orderID", id)
		return errors.New("order not found")
	}
	// Save the updated list
	if err := s.orderRepo.SaveItems(newOrders); err != nil {
		logging.Error("Failed to save updated orders after deletion", err)
		return err
	}

	logging.Info("Successfully deleted order", "orderID", id)
	return nil
}

func (s *orderService) CloseOrder(orderID string) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to close order", "orderID", orderID)

	// Retrieve current orders
	orders, err := s.orderRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read orders", err)
		return err
	}

	// Find the order by ID
	var orderToUpdate *models.Order
	for i, order := range orders {
		if order.ID == orderID {
			orderToUpdate = &orders[i]
			break
		}
	}

	if orderToUpdate == nil {
		logging.Warn("Order not found for closing", "orderID", orderID)
		return errors.New("order not found")
	}

	// Check if the status is already closed
	if orderToUpdate.Status == "closed" {
		logging.Warn("Order is already closed", "orderID", orderID)
		return errors.New("order is already closed")
	}

	// Create instances of MenuItemService and InventoryItemService
	menuService := dal.MenuItemService{}
	inventoryService := dal.InventoryItemService{}

	// Read menu items to verify ordered products
	menuItems, err := menuService.ReadItems()
	if err != nil {
		logging.Error("Failed to read menu items", err)
		return err
	}

	// Map the menu items for fast lookup by product ID
	menuItemMap := make(map[string]models.MenuItem)
	for _, menuItem := range menuItems {
		menuItemMap[menuItem.ID] = menuItem
	}

	// Read the inventory to update quantities based on the order
	inventoryItems, err := inventoryService.ReadItem()
	if err != nil {
		logging.Error("Failed to read inventory items", err)
		return err
	}

	// Map inventory items for fast lookup by ingredient ID
	inventoryMap := make(map[string]*models.InventoryItem)
	for i := range inventoryItems {
		inventoryMap[inventoryItems[i].IngredientID] = &inventoryItems[i]
	}

	// Check each order item against the menu and update inventory
	for _, orderItem := range orderToUpdate.Items {
		menuItem, exists := menuItemMap[orderItem.ProductID]
		if !exists {
			logging.Warn("Product not found in menu", "productID", orderItem.ProductID)
			return errors.New("product not found in menu: " + orderItem.ProductID)
		}

		// Deduct inventory for the ingredients of the ordered menu item
		for _, ingredient := range menuItem.Ingredients {
			inventoryItem, found := inventoryMap[ingredient.IngredientID]
			if !found {
				logging.Warn("Ingredient not found in inventory", "ingredientID", ingredient.IngredientID)
				return errors.New("ingredient not found in inventory: " + ingredient.IngredientID)
			}

			requiredQty := ingredient.Quantity * float64(orderItem.Quantity)
			if inventoryItem.Quantity < requiredQty {
				logging.Warn("Insufficient inventory for ingredient", "ingredientID", ingredient.IngredientID)
				return errors.New("insufficient inventory for ingredient: " + ingredient.IngredientID)
			}

			// Deduct the required quantity
			inventoryItem.Quantity -= requiredQty
		}
	}

	// Save the updated inventory
	if err := inventoryService.SaveItem(inventoryItems); err != nil {
		logging.Error("Failed to save updated inventory", err)
		return err
	}

	// Update the order status to closed
	orderToUpdate.Status = "closed"
	orderToUpdate.CreatedAt = time.Now().Format(time.RFC3339)

	// Save the updated order list
	if err := s.orderRepo.SaveItems(orders); err != nil {
		logging.Error("Failed to save updated orders after closing", err)
		return err
	}

	logging.Info("Successfully closed order", "orderID", orderID)
	return nil
}

func (s *orderService) TotalSalesCount() (map[string]int, error) {
	orders, err := s.orderRepo.ReadClosedOrders()
	if err != nil {
		logging.Error("Failed to read closed orders", err)
		return nil, err
	}

	// Create a map to hold the sales count for each menu item
	salesCount := make(map[string]int)

	// Loop through orders and sum the quantities for each menu item
	for _, order := range orders {
		for _, item := range order.Items {
			// Add the quantity of each item to the sales count
			salesCount[item.ProductID] += item.Quantity
		}
	}

	logging.Info("Total sales count calculated", "salesCount", salesCount)
	return salesCount, nil
}
