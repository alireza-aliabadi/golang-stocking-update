package models

type StockUpdatePayload struct {
	StockUnit string `json:"stock_unit"`
	NewStock int `json:"new_stock"`
	StoreIDs []string `json:"store_ids"`
}