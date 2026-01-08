package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusInReview ReviewStatus = "in_review"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type Review struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string       `gorm:"not null;size:200" json:"title"`
	Description   string       `gorm:"type:text" json:"description"`
	RepositoryURL string       `gorm:"size:500" json:"repository_url"`
	BranchName    string       `gorm:"size:100" json:"branch_name"`
	Status        ReviewStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	S3DiffKey     string       `gorm:"size:500" json:"s3_diff_key"`
	CreatedBy     uuid.UUID    `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt     time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Creator  User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Comments []Comment `gorm:"foreignKey:ReviewID" json:"comments,omitempty"`
}

// BeforeCreate hook to generate UUID
func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.Status == "" {
		r.Status = ReviewStatusPending
	}
	return nil
}

// TableName specifies the table name
func (Review) TableName() string {
	return "reviews"
}
