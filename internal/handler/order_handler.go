package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/service"
	"hot-coffee/logging"
	"hot-coffee/models"
	"hot-coffee/utils"
	"net/http"
	"strings"
)

var orderService service.OrderService

func InitializeOrderService(svc service.OrderService) {
	orderService = svc
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	logging.Info("Received request", "method", r.Method, "url", r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	if orderService == nil {
		orderRepo := &dal.OrderService{}
		orderService = service.NewOrderService(orderRepo)
	}

	item, itemId, _ := splitPath(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		handleGetOrder(w, item, itemId)
	case http.MethodPost:
		if strings.HasSuffix(r.URL.Path, "/close") {
			CloseOrderHandler(w, r)
		} else {
			handlePostOrder(w, r)
		}
	case http.MethodPut:
		handlePutOrder(w, r, itemId)
	case http.MethodDelete:
		handleDeleteOrder(w, itemId)
	default:
		logging.Warn("Invalid HTTP method", "method", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "Invalid HTTP method")
	}
}

func handleGetOrder(w http.ResponseWriter, item string, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling GET request", "item", item, "itemId", itemId)

	if itemId == "" {
		orders, err := orderService.FetchAllOrders()
		if err != nil {
			logging.Error("Failed to fetch all orders", err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to fetch all orders")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(orders)
	} else {
		order, err := orderService.FindOrderByID(itemId)
		if err != nil {
			logging.Error("Failed to fetch order by ID", err, "itemId", itemId)
			writeJSONError(w, http.StatusNotFound, "Order not found")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(order)
	}
}

func handlePostOrder(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling POST request")

	var newOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		logging.Error("Failed to decode request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	logging.Info("Parsed order", "order", newOrder)

	if err := orderService.CreateOrder(newOrder); err != nil {
		logging.Error("Failed to create order", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
	logging.Info("Successfully created order", "order", newOrder)
}

func handlePutOrder(w http.ResponseWriter, r *http.Request, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling PUT request", "itemId", itemId)

	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		logging.Error("Failed to decode request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	logging.Info("Parsed updated order", "itemId", itemId, "updatedOrder", updatedOrder)

	if err := orderService.UpdateOrderByID(itemId, updatedOrder); err != nil {
		logging.Error("Failed to update order", err, "itemId", itemId)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedOrder)
	logging.Info("Successfully updated order", "itemId", itemId, "updatedOrder", updatedOrder)
}

func handleDeleteOrder(w http.ResponseWriter, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling DELETE request", "itemId", itemId)

	if err := orderService.DeleteOrderByID(itemId); err != nil {
		logging.Error("Failed to delete order", err, "itemId", itemId)
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logging.Info("Successfully deleted order", "itemId", itemId)
}

func CloseOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling CLOSE order request")

	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		logging.Error("Invalid order ID in URL", fmt.Errorf("Error"), "url", r.URL.Path)
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID in URL")
		return
	}
	orderID := pathParts[2]

	err := orderService.CloseOrder(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			logging.Error("Order not found", err, "orderID", orderID)
			writeJSONError(w, http.StatusNotFound, "Order not found")
		} else if err.Error() == "order is already closed" {
			logging.Error("Order is already closed", err, "orderID", orderID)
			writeJSONError(w, http.StatusBadRequest, "Order is already closed")
		} else {
			logging.Error("Failed to close order", err, "orderID", orderID)
			writeJSONError(w, http.StatusInternalServerError, "Failed to close order")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order closed successfully"})
	logging.Info("Successfully closed order", "orderID", orderID)
}
