package models

import "time"

const (
	AdjustmentReasonAdjustment = "adjustment"
	AdjustmentReasonCorrection = "correction"
	AdjustmentReasonRestock    = "restock"
	AdjustmentReasonSale       = "sale"
	AdjustmentReasonReturn     = "return"
)

type InventoryAdjustment struct {
	ID                 int64     `gorm:"type:bigint;primary_key;autoIncrement" json:"id"`
	VariantInventoryID int64     `gorm:"type:bigint;not null" json:"variant_inventory_id"`
	AdjustedBy         int64     `gorm:"type:bigint;not null" json:"adjusted_by"`
	PreviousQuantity   int       `gorm:"not null" json:"previous_quantity"`
	NewQuantity        int       `gorm:"not null" json:"new_quantity"`
	AdjustmentAmount   int       `gorm:"not null" json:"adjustment_amount"`
	Reason             string    `gorm:"type:varchar(50);not null" json:"reason"`
	Notes              string    `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
}

func (InventoryAdjustment) TableName() string { return "inventory_adjustments" }
