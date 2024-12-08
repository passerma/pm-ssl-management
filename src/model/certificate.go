package model

import (
	"time"

	"gorm.io/gorm"
)

type Certificate struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CertificateId *int64         `json:"certificateId"`
	Domain        string         `json:"domain"`
	StartTime     *time.Time     `json:"startTime"`
	EndTime       *time.Time     `json:"endTime"`
	KeyPath       string         `json:"keyPath"`
	CertPath      string         `json:"certPath"`
	KeyContent    string         `json:"keyContent"`
	CertContent   string         `json:"certContent"`
	Command       string         `json:"command"`
	AutoRenew     bool           `json:"autoRenew"`
	RenewTime     int64          `json:"renewTime"`
	OrderId       *int64         `json:"orderId"`
	State         string         `json:"state"`
	Validate      bool           `json:"validate"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
