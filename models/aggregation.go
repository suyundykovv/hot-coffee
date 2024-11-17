package models

type AggregationData struct {
	TotalSales   float64  `json:"total_sales"`
	PopularItems MenuItem `json:"popular_items"`
	DailyItem    MenuItem `json:"daily_item"`
}
