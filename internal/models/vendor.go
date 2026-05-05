package models

type VendorApplication struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ShopName  string `json:"shop_name"`
	GSTIN     string `json:"gstin"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}