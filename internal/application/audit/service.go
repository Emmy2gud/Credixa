package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `gorm:"index;not null" json:"action"`
	Entity    string    `json:"entity"`
	EntityID  string    `json:"entity_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

type AuditService interface {
	Log(ctx context.Context, userID uint, action, entity, entityID, ip, userAgent string, meta map[string]interface{}) error
	GetUserLogs(ctx context.Context, userID uint, limit int) ([]AuditLog, error)
}

type auditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) AuditService {
	return &auditService{db: db}
}

func (s *auditService) Log(ctx context.Context, userID uint, action, entity, entityID, ip, userAgent string, meta map[string]interface{}) error {
	metaJSON := ""
	if meta != nil {
		bytes, err := json.Marshal(meta)
		if err == nil {
			metaJSON = string(bytes)
		}
	}

	log := AuditLog{
		UserID:    userID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		IPAddress: ip,
		UserAgent: userAgent,
		Metadata:  metaJSON,
		CreatedAt: time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(&log).Error; err != nil {
		return fmt.Errorf("failed to record audit log: %w", err)
	}
	return nil
}

func (s *auditService) GetUserLogs(ctx context.Context, userID uint, limit int) ([]AuditLog, error) {
	if limit <= 0 {
		limit = 20
	}
	var logs []AuditLog
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
