package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/config"
	"hot-coffee/logging"
	"hot-coffee/models"
	"hot-coffee/utils"
	"os"
)

type MenuRepository interface {
	ReadItems() ([]models.MenuItem, error)
	SaveItems([]models.MenuItem) error
}

type MenuItemService struct{}

// ReadItems reads the menu items from the file, logs the process, and returns the items.
func (m *MenuItemService) ReadItems() ([]models.MenuItem, error) {
	defer utils.CatchCriticalPoint()
	logging.Info("Reading menu items from file", "file", config.MenuFile)

	data, err := os.ReadFile(config.MenuFile)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Info("Menu file not found, returning empty menu items", "file", config.MenuFile)
			return []models.MenuItem{}, nil
		}

		logging.Error("Failed to read menu file", fmt.Errorf("Error"), config.MenuFile, "error", err.Error())
		return nil, err
	}

	var items []models.MenuItem
	if err := json.Unmarshal(data, &items); err == nil {
		logging.Info("Successfully read menu items", "count", len(items))
		return items, nil
	}

	var singleItem models.MenuItem
	if err := json.Unmarshal(data, &singleItem); err == nil {
		logging.Info("Successfully read single menu item", "item", singleItem.Name)
		return []models.MenuItem{singleItem}, nil
	}

	logging.Error("Invalid JSON structure in menu file", fmt.Errorf("Error"), config.MenuFile)
	return nil, errors.New("invalid JSON structure")
}

func (m *MenuItemService) SaveItems(items []models.MenuItem) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Saving menu items to file", "file", config.MenuFile, "count", len(items))

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		logging.Error("Failed to marshal menu items", fmt.Errorf("Error"), err.Error())
		return err
	}

	file, err := os.OpenFile(config.MenuFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		logging.Error("Failed to open menu file for writing", fmt.Errorf("Error"), config.MenuFile, "error", err.Error())
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		logging.Error("Failed to write data to menu file", fmt.Errorf("Error"), config.MenuFile, "error", err.Error())
		return err
	}

	logging.Info("Successfully saved menu items", "file", config.MenuFile, "count", len(items))
	return nil
}
