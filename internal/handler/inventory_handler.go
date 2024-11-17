package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/service"
	"hot-coffee/logging"
	"hot-coffee/models"
	"hot-coffee/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

var Inventory service.InventoryService

// InventoryHandler handles different HTTP methods for the inventory endpoint.
func InventoryHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	logging.Info("Received request", "method", r.Method, "url", r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	inventoryRepo := &dal.InventoryItemService{}
	Inventory = service.NewInventoryService(inventoryRepo)
	item, itemId, _ := splitPath(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		handleGetInventory(w, item, itemId)
	case http.MethodPost:
		handlePostInventory(w, r)
	case http.MethodPut:
		handlePutInventory(w, r, itemId)
	case http.MethodDelete:
		handleDeleteInventory(w, itemId)
	default:
		logging.Warn("Invalid HTTP method", "method", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "Invalid HTTP method")
	}
}

// handleGetInventory handles the GET request for fetching inventory items.
func handleGetInventory(w http.ResponseWriter, item string, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling GET request", "item", item, "itemId", itemId)

	if itemId == "" {
		inventoryItems, err := Inventory.GetAllInventoryItems()
		if err != nil {
			logging.Error("Failed to fetch all inventory items", err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to fetch all inventory items")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inventoryItems)
	} else {
		inventoryItem, err := Inventory.GetInventoryItemByID(itemId)
		if err != nil {
			logging.Error("Failed to fetch inventory item by ID", err, "itemId", itemId)
			writeJSONError(w, http.StatusNotFound, "Inventory item not found")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(inventoryItem)
	}
}

// handlePostInventory handles the POST request for adding a new inventory item.
func handlePostInventory(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling POST request")

	var newItem models.InventoryItem
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logging.Error("Failed to read request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	logging.Info("Received request body", "body", string(body))

	if err := json.Unmarshal(body, &newItem); err != nil {
		logging.Error("Invalid JSON input", err, "body", string(body))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	if err := validateInventoryItem(newItem); err != nil {
		logging.Error("Validation failed", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	logging.Info("Parsed Inventory Item", "item", newItem)

	if err := Inventory.AddInventoryItem(newItem); err != nil {
		logging.Error("Failed to add inventory item", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to add inventory item")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
	logging.Info("Successfully added inventory item", "item", newItem)
}

// handlePutInventory handles the PUT request for updating an inventory item.
func handlePutInventory(w http.ResponseWriter, r *http.Request, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling PUT request", "itemId", itemId)

	var updatedItem models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		logging.Error("Failed to decode request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	// Check if item exists
	_, err := Inventory.GetInventoryItemByID(itemId)
	if err != nil {
		logging.Error("Item not found", err, "itemId", itemId)
		writeJSONError(w, http.StatusNotFound, "Inventory item not found")
		return
	}

	// Validate the input
	if err := validateInventoryItem(updatedItem); err != nil {
		logging.Error("Validation failed", err)
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := Inventory.UpdateInventoryItem(itemId, updatedItem); err != nil {
		logging.Error("Failed to update inventory item", err, "itemId", itemId)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update inventory item")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
	logging.Info("Successfully updated inventory item", "itemId", itemId, "updatedItem", updatedItem)
}

// handleDeleteInventory handles the DELETE request for removing an inventory item.
func handleDeleteInventory(w http.ResponseWriter, itemId string) {
	defer utils.CatchCriticalPoint()

	logging.Info("Handling DELETE request", "itemId", itemId)

	err := Inventory.DeleteInventoryItem(itemId)
	if err != nil {
		if err.Error() == "inventory item not found" {
			logging.Error("Inventory item not found", err, "itemId", itemId)
			writeJSONError(w, http.StatusNotFound, "Inventory item not found")
		} else {
			logging.Error("Failed to delete inventory item", err, "itemId", itemId)
			writeJSONError(w, http.StatusInternalServerError, "Failed to delete inventory item")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logging.Info("Successfully deleted inventory item", "itemId", itemId)
}

// validateInventoryItem validates the fields of an inventory item.
func validateInventoryItem(item models.InventoryItem) error {
	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("inventory item name is required")
	}
	if item.Quantity < 0 {
		return fmt.Errorf("inventory item quantity cannot be negative")
	}

	return nil
}
