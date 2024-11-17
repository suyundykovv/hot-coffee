package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"hot-coffee/config"
	"hot-coffee/logging"
	"hot-coffee/utils"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func isRestrictedDir(dir string) bool {
	for _, restricted := range config.RestrictedDirs {
		if strings.EqualFold(dir, restricted) {
			return true
		}
	}
	return false
}

func isPortInRange(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return portNum >= 1024 && portNum <= 65535
}

func isPortAvailable(port string) bool {
	if !isPortInRange(port) {
		fmt.Println("Port out of range. Valid range is 1024-65535.")
		return false
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func Setup() {
	defer utils.CatchCriticalPoint()

	defaultPort := "6002"
	defaultStorageDir := "./"
	help := flag.Bool("help", false, "Display help information")

	flag.StringVar(&config.Port, "port", defaultPort, "Port to run the server on")
	flag.StringVar(&config.StorageDir, "directory", defaultStorageDir, "Directory for file storage")
	flag.Parse()
	if !isPortAvailable(config.Port) {
		logging.Error("The specified port is already in use", nil, "port", config.Port)
		log.Fatalf("The specified port '%s' is already in use. Please choose a different port.", config.Port)
	}
	logging.Info("Using port for the server", "port", config.Port)
	if *help {
		config.PrintUsage()
		logging.Info("Help command used", "help", true)
		os.Exit(0)
	}

	if isRestrictedDir(config.StorageDir) {
		logging.Error("The specified directory is restricted", fmt.Errorf("Error"), config.StorageDir)
		log.Fatalf("The specified directory '%s' is restricted. Please choose a different name.", config.StorageDir)
	}
	logging.Info("Using port for the server", "port", config.Port)

	if err := os.MkdirAll(config.StorageDir, os.ModePerm); err != nil {
		logging.Error("Failed to create storage directory", nil, config.StorageDir, "error", err)
		log.Fatalf("Failed to create storage directory %s: %v", config.StorageDir, err)
	}
	logging.Info("Using storage directory", "directory", config.StorageDir)

	dataDir := filepath.Join(config.StorageDir, "data")
	config.StorageDir = dataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		logging.Error("Failed to create 'data' directory", nil, dataDir, "error", err)
		log.Fatalf("Failed to create 'data' directory %s: %v", dataDir, err)
	}

	logFile := filepath.Join(dataDir, "app.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			logging.Error("Failed to create app.log file", nil, logFile, "error", err)
			log.Fatalf("Failed to create app.log file: %v", err)
		}
		defer file.Close()
		logging.Info("Created app.log file", "file", logFile)
	}
	config.InventoryFile = filepath.Join(config.StorageDir, "inventory.json")
	createJSONFileIfNotExists(filepath.Join(dataDir, "inventory.json"), config.DefaultInventory())
	config.MenuFile = filepath.Join(config.StorageDir, "menu_item.json")
	createJSONFileIfNotExists(filepath.Join(dataDir, "menu_item.json"), config.DefaultMenuItems())
	config.OrdersFile = filepath.Join(config.StorageDir, "order.json")
	createJSONFileIfNotExists(filepath.Join(dataDir, "order.json"), config.DefaultOrders())
	config.LogFile = filepath.Join(config.StorageDir, "app.log")
	config.AggregationFile = filepath.Join(config.StorageDir, "aggregation.json")

	logging.Info("Using data directory", "data_directory", dataDir)
}

// Helper function to create a JSON file if it doesn't exist
func createJSONFileIfNotExists(filePath string, defaultData interface{}) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create and write default data to the file
		file, err := os.Create(filePath)
		if err != nil {
			logging.Error("Failed to create JSON file", nil, filePath, "error", err)
			log.Fatalf("Failed to create JSON file %s: %v", filePath, err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultData); err != nil {
			logging.Error("Failed to write default data to JSON file", nil, filePath, "error", err)
			log.Fatalf("Failed to write default data to JSON file %s: %v", filePath, err)
		}
		logging.Info("Created JSON file with default data", "file", filePath)
	}
}
