package dal

import (
	"encoding/json"
	"errors"
	"hot-coffee/config"
	"hot-coffee/logging"
	"hot-coffee/models"
	"hot-coffee/utils"
	"os"
)

type OrderRepository interface {
	ReadItems() ([]models.Order, error)
	SaveItems([]models.Order) error

	ReadClosedOrders() ([]models.Order, error)
}

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (o *OrderService) ReadItems() ([]models.Order, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Reading orders from file")

	data, err := os.ReadFile(config.OrdersFile)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Warn("Orders file does not exist, returning empty list")
			return []models.Order{}, nil
		}
		logging.Error("Failed to read orders file", err)
		return nil, err
	}

	logging.Info("Raw JSON data:", string(data))

	var orders []models.Order
	if err := json.Unmarshal(data, &orders); err == nil {
		logging.Info("Successfully unmarshalled orders")
		return orders, nil
	}

	var singleOrder models.Order
	if err := json.Unmarshal(data, &singleOrder); err == nil {
		logging.Info("Successfully unmarshalled single order")
		return []models.Order{singleOrder}, nil
	}

	logging.Error("Invalid JSON structure", err)
	return nil, errors.New("invalid JSON structure")
}

func (o *OrderService) SaveItems(orders []models.Order) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Saving orders to file")

	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		logging.Error("Failed to marshal orders", err)
		return err
	}

	file, err := os.OpenFile(config.OrdersFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		logging.Error("Failed to open orders file", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		logging.Error("Failed to write to orders file", err)
		return err
	}

	logging.Info("Successfully saved orders to file")
	return nil
}

func (o *OrderService) ReadClosedOrders() ([]models.Order, error) {
	data, err := os.ReadFile(config.OrdersFile)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Info("Orders file not found, returning empty orders")
			return []models.Order{}, nil
		}
		logging.Error("Failed to read orders file", err)
		return nil, err
	}

	var orders []models.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		logging.Error("Failed to parse orders file", err)
		return nil, err
	}

	// Filter orders that are closed
	var closedOrders []models.Order
	for _, order := range orders {
		if order.Status == "closed" {
			closedOrders = append(closedOrders, order)
		}
	}

	return closedOrders, nil
}
