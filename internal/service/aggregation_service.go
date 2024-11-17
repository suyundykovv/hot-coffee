package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/logging"
	"hot-coffee/models"
	"hot-coffee/utils"
)

type ReportService interface {
	GetMostPopularItem() (models.MenuItem, error)
	GetDailyItem() (models.MenuItem, error)
	TotalSalesAmount() (float64, error)
}

type reportService struct {
	menuRepo        dal.MenuRepository
	aggregationRepo dal.AggregationRepository
	orderRepo       dal.OrderRepository
}

func NewReportService(menuRepo dal.MenuRepository, aggregationRepo dal.AggregationRepository, orderRepo dal.OrderRepository) ReportService {
	return &reportService{
		menuRepo:        menuRepo,
		aggregationRepo: aggregationRepo,
		orderRepo:       orderRepo,
	}
}

// TotalSalesAmount calculates the total sales amount from all closed orders.
func (s *reportService) TotalSalesAmount() (float64, error) {
	defer utils.CatchCriticalPoint()

	orders, err := s.orderRepo.ReadClosedOrders()
	if err != nil {
		logging.Error("Failed to read closed orders", err)
		return 0, err
	}

	// Get all menu items for price lookup
	menuItems, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch menu items", err)
		return 0, err
	}

	// Create a map to quickly look up menu items by their ID
	menuItemMap := make(map[string]models.MenuItem)
	for _, item := range menuItems {
		menuItemMap[item.ID] = item
	}

	// Calculate the total sales amount
	var totalSalesAmount float64
	for _, order := range orders {
		for _, orderItem := range order.Items {
			menuItem, exists := menuItemMap[orderItem.ProductID]
			if !exists {
				logging.Warn("Menu item not found for order item", "menuItemID", orderItem.ProductID)
				continue
			}

			// Calculate the sales amount for this item (quantity * price)
			totalSalesAmount += float64(orderItem.Quantity) * menuItem.Price
		}
	}

	logging.Info("Total sales amount calculated", "totalSalesAmount", totalSalesAmount)
	return totalSalesAmount, nil
}

// GetMostPopularItem calculates the most frequently ordered product by product ID
func (s *reportService) GetMostPopularItem() (models.MenuItem, error) {
	defer utils.CatchCriticalPoint()

	// Fetch all closed orders
	orders, err := s.orderRepo.ReadClosedOrders()
	if err != nil {
		logging.Error("Failed to read closed orders", err)
		return models.MenuItem{}, err
	}

	// Map to count frequency of each ProductID
	productFrequency := make(map[string]int)

	// Loop through all orders and count the frequency of each product
	for _, order := range orders {
		for _, orderItem := range order.Items {
			productFrequency[orderItem.ProductID]++ // Increment frequency of each ProductID
		}
	}

	// Find the ProductID with the highest frequency
	var mostPopularProductID string
	var highestFrequency int
	for productID, frequency := range productFrequency {
		if frequency > highestFrequency {
			mostPopularProductID = productID
			highestFrequency = frequency
		}
	}

	// If we have a most popular product, fetch its details from the menu
	if mostPopularProductID != "" {
		// Fetch all menu items
		menuItems, err := s.menuRepo.ReadItems()
		if err != nil {
			logging.Error("Failed to fetch menu items", err)
			return models.MenuItem{}, err
		}

		// Search for the MenuItem corresponding to the most popular ProductID
		for _, menuItem := range menuItems {
			if menuItem.ID == mostPopularProductID {
				logging.Info("Most popular item found", "productID", mostPopularProductID, "frequency", highestFrequency)
				return menuItem, nil
			}
		}
	}

	// If no popular product found, return error
	logging.Info("No popular item found", "productID", mostPopularProductID)
	return models.MenuItem{}, nil
}

// GetDailyItem selects a random menu item from the available items.
func (s *reportService) GetDailyItem() (models.MenuItem, error) {
	defer utils.CatchCriticalPoint()

	// Fetch all menu items
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch menu items", err)
		return models.MenuItem{}, err
	}

	if len(items) == 0 {
		logging.Warn("No menu items available")
		return models.MenuItem{}, errors.New("no menu items available")
	}

	// Select a random menu item
	dailyItem := items[utils.RandomInt(0, len(items)-1)]

	logging.Info("Daily item selected", "itemID", dailyItem.ID)

	return dailyItem, nil
}
