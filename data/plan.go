package data

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Plan struct {
	gorm.Model
	SiteId string `json:"siteId"`
	StageId string `json:"stageId"`
	OperId string `json:"operId"`
	ResourceId string `json:"resourceId"`
	ProductId string `json:"productId"`
	PlanQty float32 `json:"planQty"`
	StartTime time.Time `json:"startTime"`
	EndTime time.Time `json:"endTime"`
}
