// handler/report_handler.go
package handler

import (
	"encoding/json"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/service"
	"hot-coffee/logging"
	"hot-coffee/utils"
	"net/http"
)

var reportService service.ReportService

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetReports(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Invalid HTTP method")
	}
}

func handleGetReports(w http.ResponseWriter, r *http.Request) {
	defer utils.CatchCriticalPoint()

	// Log the GET request
	logging.Info("Handling GET request for reports", "url", r.URL.Path)
	orderRepo := &dal.OrderService{}
	menuitemRepo := &dal.MenuItemService{}
	aggRepo := &dal.AggregationService{}
	reportService = service.NewReportService(menuitemRepo, aggRepo, orderRepo)

	// Check the URL for specific report
	switch r.URL.Path {
	case "/reports/total-sales":
		handleTotalSales(w)
	case "reports/popular-items":
		handlePopularItems(w)
	case "/reports/daily-item":
		handleDailyItem(w)
	default:
		writeJSONError(w, http.StatusNotFound, "Report not found")
	}
}

func handleTotalSales(w http.ResponseWriter) {
	defer utils.CatchCriticalPoint()

	totalSales, err := reportService.TotalSalesAmount()
	if err != nil {
		logging.Error("Failed to fetch total sales", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch total sales")
		return
	}

	// Return total sales as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{"total_sales": totalSales})
}

func handlePopularItems(w http.ResponseWriter) {
	defer utils.CatchCriticalPoint()

	popularItems, err := reportService.GetMostPopularItem()
	if err != nil {
		logging.Error("Failed to fetch popular items", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch popular items")
		return
	}

	// Return popular items as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(popularItems)
}

func handleDailyItem(w http.ResponseWriter) {
	defer utils.CatchCriticalPoint()

	dailyItem, err := reportService.GetDailyItem()
	if err != nil {
		logging.Error("Failed to fetch daily item", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch daily item")
		return
	}

	// Return the daily item as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dailyItem)
}
