package config

import "fmt"

var (
	AggregationFile string
	Port            string
	StorageDir      string
	InventoryFile   string
	MenuFile        string
	OrdersFile      string
	LogFile         string
	RestrictedDirs  = []string{"flags", "handlers", "models", "servers", "storage", "utils", "../"}
)

// Default content for inventory.json
func DefaultInventory() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"ingredient_id": "espresso_shot",
			"name":          "Espresso Shot",
			"quantity":      372,
			"unit":          "shots",
		},
		{
			"ingredient_id": "milk",
			"name":          "Milk",
			"quantity":      3400,
			"unit":          "ml",
		},
		{
			"ingredient_id": "flour",
			"name":          "Flour",
			"quantity":      9400,
			"unit":          "g",
		},
		{
			"ingredient_id": "blueberries",
			"name":          "Blueberries",
			"quantity":      1800,
			"unit":          "g",
		},
		{
			"ingredient_id": "sugar",
			"name":          "Sugar",
			"quantity":      4750,
			"unit":          "g",
		},
	}
}

// Default content for menu_item.json
func DefaultMenuItems() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"product_id":  "latte",
			"name":        "Caffe Latte",
			"description": "Espresso with steamed milk",
			"price":       3.5,
			"ingredients": []map[string]interface{}{
				{"ingredient_id": "espresso_shot", "quantity": 1},
				{"ingredient_id": "milk", "quantity": 200},
			},
		},
		{
			"product_id":  "muffin",
			"name":        "Blueberry Muffin",
			"description": "Freshly baked muffin with blueberries",
			"price":       2.0,
			"ingredients": []map[string]interface{}{
				{"ingredient_id": "flour", "quantity": 100},
				{"ingredient_id": "blueberries", "quantity": 20},
				{"ingredient_id": "sugar", "quantity": 30},
			},
		},
		{
			"product_id":  "espresso",
			"name":        "Espresso",
			"description": "Strong and bold coffee",
			"price":       2.5,
			"ingredients": []map[string]interface{}{
				{"ingredient_id": "espresso_shot", "quantity": 10},
			},
		},
		{
			"product_id":  "shit",
			"name":        "Blueberry Muffin",
			"description": "Freshly baked muffin with blueberries",
			"price":       2.0,
			"ingredients": []map[string]interface{}{
				{"ingredient_id": "flour", "quantity": 10},
				{"ingredient_id": "blueberries", "quantity": 10},
				{"ingredient_id": "sugar", "quantity": 10},
			},
		},
	}
}

// Default content for order.json
func DefaultOrders() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"order_id":      "order137",
			"customer_name": "John Doe",
			"items": []map[string]interface{}{
				{"product_id": "espresso", "quantity": 5},
				{"product_id": "shit", "quantity": 5},
			},
			"status":     "closed",
			"created_at": "2024-11-14T17:06:42+05:00",
		},
		{
			"order_id":      "order138",
			"customer_name": "John Doe",
			"items": []map[string]interface{}{
				{"product_id": "espresso", "quantity": 5},
				{"product_id": "shit", "quantity": 5},
			},
			"status":     "closed",
			"created_at": "2024-11-14T17:07:01+05:00",
		},
	}
}

func PrintUsage() {
	fmt.Println("Coffee Shop Management System")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    hot-coffee [--port <N>] [--dir <S>] ")
	fmt.Println("    hot-coffee --help")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help     Show this screen.")
	fmt.Println("  --port N   Port number")
	fmt.Println("  --dir S    Path to the directory")
}
