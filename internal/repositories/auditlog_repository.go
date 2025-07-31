package repositories

import (
	"database/sql"

	"github.com/Severustark/movietracker-backend/internal/models"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
}

type auditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (entity_type, entity_id, action, details, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(query,
		log.EntityType,
		log.EntityID,
		log.Action,
		log.Details,
	)
	return err
}
