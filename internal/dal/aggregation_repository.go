// aggregationRepository.go
package dal

import (
	"encoding/json"
	"hot-coffee/config"
	"hot-coffee/logging"
	"hot-coffee/models"
	"os"
)

type AggregationRepository interface {
	SaveAggregationData(aggregationData models.AggregationData) error
}

type AggregationService struct{}

// SaveAggregationData saves the aggregated results (e.g., total sales, popular items, daily item) to a file.
func (a *AggregationService) SaveAggregationData(aggregationData models.AggregationData) error {
	// Marshal the aggregation data into JSON
	data, err := json.MarshalIndent(aggregationData, "", "  ")
	if err != nil {
		logging.Error("Failed to marshal aggregation data", err)
		return err
	}

	// Open the aggregation file (create it if it doesn't exist, overwrite if it does)
	file, err := os.OpenFile(config.AggregationFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		logging.Error("Failed to open aggregation file for writing", err)
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		logging.Error("Failed to write aggregation data to file", err)
		return err
	}

	logging.Info("Successfully saved aggregation data", "file", config.AggregationFile)
	return nil
}
