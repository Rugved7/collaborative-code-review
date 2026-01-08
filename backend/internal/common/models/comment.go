package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	ReviewID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"review_id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	FilePath        string     `gorm:"size:500" json:"file_path"`
	LineNumber      int        `gorm:"default:0" json:"line_number"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	ParentCommentID *uuid.UUID `gorm:"type:uuid;index" json:"parent_comment_id,omitempty"`
	IsResolved      bool       `gorm:"default:false" json:"is_resolved"`
	CreatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Review        Review    `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentComment *Comment  `gorm:"foreignKey:ParentCommentID" json:"parent_comment,omitempty"`
	Replies       []Comment `gorm:"foreignKey:ParentCommentID" json:"replies,omitempty"`
}

// BeforeCreate hook to generate UUID
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Comment) TableName() string {
	return "comments"
}
