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

var menuitem service.MenuService

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	// Log incoming request method and URL
	logging.Info("Received request", "method", r.Method, "url", r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	menuitemRepo := &dal.MenuItemService{}
	menuitem = service.NewMenuService(menuitemRepo)
	item, itemId, _ := splitPath(r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		handleGetMenu(w, item, itemId)
	case http.MethodPost:
		handlePostMenu(w, r)
	case http.MethodPut:
		handlePutMenu(w, r, itemId)
	case http.MethodDelete:
		handleDeleteMenu(w, itemId)
	default:
		logging.Warn("Invalid HTTP method", "method", r.Method)
		writeJSONError(w, http.StatusMethodNotAllowed, "Invalid HTTP method")
	}
}

func splitPath(path string) (string, string, error) {
	defer utils.CatchCriticalPoint()

	// Split the path into parts
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid path: %s", path)
	}

	// First part should be "menu", second part should be itemId
	item := parts[1]   // "menu"
	itemId := parts[2] // menu item ID

	return item, itemId, nil
}

func handleGetMenu(w http.ResponseWriter, item string, itemId string) {
	defer utils.CatchCriticalPoint()

	// Log the GET request
	logging.Info("Handling GET request", "item", item, "itemId", itemId)

	if itemId == "" {
		menuItems, err := menuitem.FetchAllMenuItems()
		if err != nil {
			logging.Error("Failed to fetch all menu items", err)
			writeJSONError(w, http.StatusInternalServerError, "Failed to fetch all menu items")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(menuItems)
	} else {
		menuItem, err := menuitem.FindMenuItemByID(itemId)
		if err != nil {
			logging.Error("Failed to fetch menu item by ID", err, "itemId", itemId)
			writeJSONError(w, http.StatusNotFound, "Menu item not found")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(menuItem)
	}
}

func handlePostMenu(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	// Log the POST request
	logging.Info("Handling POST request")

	var newItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		logging.Error("Failed to decode request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	// Log the parsed menu item
	logging.Info("Parsed Menu Item", "item", newItem)

	// Fetch the existing menu items
	menuItems, err := menuitem.FetchAllMenuItems()
	if err != nil {
		logging.Error("Failed to fetch all menu items", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch menu items")
		return
	}

	// Append the new item
	menuItems = append(menuItems, newItem)

	// Save the updated list of menu items
	if err := menuitem.CreateMenuItem(newItem); err != nil {
		logging.Error("Failed to create menu item", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create menu item")
		return
	}

	// Respond with the newly created item
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
	logging.Info("Successfully created menu item", "item", newItem)
}

func handlePutMenu(w http.ResponseWriter, r *http.Request, itemId string) {
	defer utils.CatchCriticalPoint()

	// Log the PUT request
	logging.Info("Handling PUT request", "itemId", itemId)

	var updatedItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		logging.Error("Failed to decode request body", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	// Log the updated menu item
	logging.Info("Parsed updated menu item", "itemId", itemId, "updatedItem", updatedItem)

	// Check if the item exists
	if err := menuitem.UpdateMenuItemByID(itemId, updatedItem); err != nil {
		logging.Error("Failed to update menu item", err, "itemId", itemId)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update menu item")
		return
	}

	// Respond with the updated item
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
	logging.Info("Successfully updated menu item", "itemId", itemId, "updatedItem", updatedItem)
}

func handleDeleteMenu(w http.ResponseWriter, itemId string) {
	defer utils.CatchCriticalPoint()

	// Log the DELETE request
	logging.Info("Handling DELETE request", "itemId", itemId)

	if err := menuitem.DeleteMenuItemByID(itemId); err != nil {
		logging.Error("Failed to delete menu item", err, "itemId", itemId)
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete menu item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logging.Info("Successfully deleted menu item", "itemId", itemId)
}

// writeJSONError writes a structured JSON error response.
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errorResponse)
}
