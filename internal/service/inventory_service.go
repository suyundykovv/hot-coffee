package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
	"hot-coffee/utils"
	"strings"

	"hot-coffee/logging" // Import the logging package
)

type InventoryService interface {
	AddInventoryItem(item models.InventoryItem) error
	GetAllInventoryItems() ([]models.InventoryItem, error)
	GetInventoryItemByID(id string) (models.InventoryItem, error)
	UpdateInventoryItem(id string, item models.InventoryItem) error
	DeleteInventoryItem(id string) error
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *inventoryService) AddInventoryItem(item models.InventoryItem) error {
	defer utils.CatchCriticalPoint()

	// Log adding inventory item
	logging.Info("Attempting to add inventory item", "ingredientID", item.IngredientID)

	// Read existing items
	items, err := s.inventoryRepo.ReadItem()
	if err != nil {
		logging.Error("Failed to read inventory items", err)
		return err
	}

	// Check for duplicate ingredientID
	for _, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			logging.Warn("Inventory item with IngredientID already exists", "ingredientID", item.IngredientID)
			return errors.New("inventory item with this IngredientID already exists")
		}
	}
	if err := utils.ValidateUpdatedInventoryItem(item); err != nil {
		logging.Warn("Invalid create inventory item data", "error", err)
		return err
	}
	// Add the new item to the list
	items = append(items, item)

	// Save updated items
	err = s.inventoryRepo.SaveItem(items)
	if err != nil {
		logging.Error("Failed to save updated inventory", err)
		return err
	}

	// Log success
	logging.Info("Successfully added inventory item", "ingredientID", item.IngredientID)
	return nil
}

func (s *inventoryService) GetAllInventoryItems() ([]models.InventoryItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching all inventory items")

	// Fetch and return all inventory items
	items, err := s.inventoryRepo.ReadItem()
	if err != nil {
		logging.Error("Failed to fetch inventory items", err)
		return nil, err
	}

	logging.Info("Fetched all inventory items", "count", len(items))
	return items, nil
}

func (s *inventoryService) GetInventoryItemByID(id string) (models.InventoryItem, error) {
	defer utils.CatchCriticalPoint()

	logging.Info("Fetching inventory item by ID", "ingredientID", id)

	// Fetch all inventory items
	items, err := s.inventoryRepo.ReadItem()
	if err != nil {
		logging.Error("Failed to fetch inventory items", err)
		return models.InventoryItem{}, err
	}

	// Search for the item by ID
	for _, item := range items {
		if item.IngredientID == id {
			logging.Info("Found inventory item", "ingredientID", id)
			return item, nil
		}
	}

	logging.Warn("Inventory item not found", "ingredientID", id)
	return models.InventoryItem{}, errors.New("inventory item not found")
}

func (s *inventoryService) UpdateInventoryItem(id string, updatedItem models.InventoryItem) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to update inventory item", "ingredientID", id)

	// Validate the updated inventory item before proceeding
	if err := utils.ValidateUpdatedInventoryItem(updatedItem); err != nil {
		logging.Warn("Invalid updated inventory item data", "ingredientID", id, "error", err)
		return err
	}

	// Fetch all items from the inventory
	items, err := s.inventoryRepo.ReadItem()
	if err != nil {
		logging.Error("Failed to fetch inventory items", err)
		return err
	}

	// Look for the item to update
	for i, item := range items {
		if item.IngredientID == id {
			// Update the item with the new data
			items[i] = updatedItem
			err := s.inventoryRepo.SaveItem(items)
			if err != nil {
				logging.Error("Failed to save updated inventory", err)
				return err
			}

			// Log success
			logging.Info("Successfully updated inventory item", "ingredientID", id)
			return nil
		}
	}

	logging.Warn("Inventory item not found for update", "ingredientID", id)
	return errors.New("inventory item not found")
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	defer utils.CatchCriticalPoint()

	logging.Info("Attempting to delete inventory item", "ingredientID", id)

	// Fetch all items
	items, err := s.inventoryRepo.ReadItem()
	if err != nil {
		logging.Error("Failed to fetch inventory items", err)
		return err
	}

	if len(items) == 0 {
		logging.Warn("Inventory is empty, cannot delete item", "ingredientID", id)
		return errors.New("inventory is empty")
	}

	// Create a new list excluding the item to be deleted
	var updatedItems []models.InventoryItem
	for _, item := range items {
		if strings.EqualFold(item.IngredientID, id) {
			continue
		}
		updatedItems = append(updatedItems, item)
	}

	if len(updatedItems) == len(items) {
		logging.Warn("Inventory item not found for deletion", "ingredientID", id)
		return errors.New("inventory item not found")
	}

	// Save the updated list
	err = s.inventoryRepo.SaveItem(updatedItems)
	if err != nil {
		logging.Error("Failed to save updated inventory after deletion", err)
		return err
	}

	logging.Info("Successfully deleted inventory item", "ingredientID", id)
	return nil
}
