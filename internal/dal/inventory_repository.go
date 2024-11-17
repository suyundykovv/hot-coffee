package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/config"
	"hot-coffee/logging"
	"hot-coffee/models"
	"net/http"
	"os"

	"hot-coffee/utils" // Import your utils package
)

type InventoryRepository interface {
	ReadItem() ([]models.InventoryItem, error)
	SaveItem([]models.InventoryItem) error
}

type InventoryItemService struct {
	models.InventoryItem
}

func (i *InventoryItemService) ReadItem() ([]models.InventoryItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Reading inventory items", "file", config.InventoryFile)
	var inventoryItems []models.InventoryItem

	data, err := os.ReadFile(config.InventoryFile)
	if err != nil {
		logging.Error("Failed to read inventory file", fmt.Errorf("Error"), config.InventoryFile, "error", err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, "Failed to read inventory file", err)
	}

	err = json.Unmarshal(data, &inventoryItems)
	if err != nil {
		logging.Error("Failed to unmarshal inventory data", fmt.Errorf("Error"), config.InventoryFile, "error", err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, "Failed to unmarshal inventory data", err)
	}

	logging.Info("Successfully read inventory items", "count", len(inventoryItems))
	return inventoryItems, nil
}

func (i *InventoryItemService) SaveItem(inventoryItems []models.InventoryItem) error {
	// Use defer to catch panics in the method
	defer utils.CatchCriticalPoint()

	logging.Info("Saving inventory items", "file", config.InventoryFile, "count", len(inventoryItems))

	data, err := json.MarshalIndent(inventoryItems, "", "  ")
	if err != nil {
		logging.Error("Failed to marshal inventory items", fmt.Errorf("Error"), config.InventoryFile, "error", err.Error())
		return models.NewErrorResponse(http.StatusInternalServerError, "Failed to marshal inventory items", err)
	}

	file, err := os.OpenFile(config.InventoryFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		logging.Error("Failed to open inventory file for writing", fmt.Errorf("Error"), config.InventoryFile, "error", err.Error())
		return models.NewErrorResponse(http.StatusInternalServerError, "Failed to open inventory file for writing", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		logging.Error("Failed to write inventory data to file", fmt.Errorf("Error"), config.InventoryFile, "error", err.Error())
		return models.NewErrorResponse(http.StatusInternalServerError, "Failed to write inventory data to file", err)
	}

	logging.Info("Successfully saved inventory items", "file", config.InventoryFile, "count", len(inventoryItems))
	return nil
}
