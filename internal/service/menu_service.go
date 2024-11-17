package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"

	"hot-coffee/logging" // Import the logging package
)

type MenuService interface {
	CreateMenuItem(item models.MenuItem) error
	FetchAllMenuItems() ([]models.MenuItem, error)
	FindMenuItemByID(id string) (models.MenuItem, error)
	UpdateMenuItemByID(id string, item models.MenuItem) error
	DeleteMenuItemByID(id string) error
	GetPopularMenuItems() ([]models.MenuItem, error)
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(menuRepo dal.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (s *menuService) CreateMenuItem(item models.MenuItem) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to create menu item", "itemID", item.ID)

	// Read current menu items
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read menu items", err)
		return err
	}

	// Check if the item already exists
	for _, existingItem := range items {
		if existingItem.ID == item.ID {
			logging.Warn("Menu item with this ID already exists", "itemID", item.ID)
			return errors.New("menu item with this ID already exists")
		}
	}
	if err := utils.ValidateUpdatedMenuItem(item); err != nil {
		logging.Warn("Invalid updated menu item data", "error", err)
		return err
	}
	// Add the new item
	items = append(items, item)

	// Save the updated menu items
	err = s.menuRepo.SaveItems(items)
	if err != nil {
		logging.Error("Failed to save new menu item", err)
		return err
	}

	// Log success
	logging.Info("Successfully created menu item", "itemID", item.ID)
	return nil
}

func (s *menuService) FetchAllMenuItems() ([]models.MenuItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching all menu items")

	// Fetch all menu items
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch menu items", err)
		return nil, err
	}

	// Log success
	logging.Info("Fetched all menu items", "count", len(items))
	return items, nil
}

func (s *menuService) FindMenuItemByID(id string) (models.MenuItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching menu item by ID", "itemID", id)

	// Fetch all menu items
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch menu items", err)
		return models.MenuItem{}, err
	}

	// Search for the item by ID
	for _, item := range items {
		if item.ID == id {
			logging.Info("Found menu item", "itemID", id)
			return item, nil
		}
	}

	logging.Warn("Menu item not found", "itemID", id)
	return models.MenuItem{}, errors.New("menu item not found")
}

func (s *menuService) UpdateMenuItemByID(id string, updatedItem models.MenuItem) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to update menu item", "itemID", id)

	if err := utils.ValidateUpdatedMenuItem(updatedItem); err != nil {
		logging.Warn("Invalid updated menu item data", "itemID", id, "error", err)
		return err
	}

	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read menu items", err)
		return err
	}

	for i, item := range items {
		if item.ID == id {
			items[i] = updatedItem
			err := s.menuRepo.SaveItems(items)
			if err != nil {
				logging.Error("Failed to save updated menu items", err)
				return err
			}

			// Log success
			logging.Info("Successfully updated menu item", "itemID", id)
			return nil
		}
	}

	logging.Warn("Menu item not found for update", "itemID", id)
	return errors.New("menu item not found")
}

func (s *menuService) DeleteMenuItemByID(id string) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to delete menu item", "itemID", id)

	// Read current menu items
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to read menu items", err)
		return err
	}

	// Create a new list excluding the item to be deleted
	var updatedItems []models.MenuItem
	for _, item := range items {
		if item.ID != id {
			updatedItems = append(updatedItems, item)
		}
	}
	if len(updatedItems) == len(items) {
		logging.Warn("Menu item not found for deletion", "MenuID", id)
		return errors.New("Menu item not found")
	}
	// Save the updated list
	err = s.menuRepo.SaveItems(updatedItems)
	if err != nil {
		logging.Error("Failed to save updated menu items after deletion", err)
		return err
	}

	// Log success
	logging.Info("Successfully deleted menu item", "itemID", id)
	return nil
}

func (s *menuService) GetPopularMenuItems() ([]models.MenuItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching popular menu items")

	// Fetch and return all menu items (assuming popular items are part of the entire list)
	items, err := s.menuRepo.ReadItems()
	if err != nil {
		logging.Error("Failed to fetch popular menu items", err)
		return nil, err
	}

	// Log success
	logging.Info("Fetched popular menu items", "count", len(items))
	return items, nil
}
