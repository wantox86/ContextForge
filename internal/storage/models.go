// Package storage provides SQLite database access via GORM.
package storage

import "time"

// Template represents a collection of rules for a specific tech stack.
type Template struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"not null"`
	Description string
	Stack       []string  `gorm:"serializer:json"`
	Rules       []Rule    `gorm:"foreignKey:TemplateID"`
	IsBuiltIn   bool      `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rule represents a single coding instruction within a template.
type Rule struct {
	ID         uint   `gorm:"primarykey"`
	TemplateID uint
	Category   string // coding_style | testing | security | naming | architecture
	Title      string `gorm:"not null"`
	Content    string `gorm:"not null"`
	Priority   int
	Enabled    bool `gorm:"default:true"`
}

// Project represents a scanned developer project.
type Project struct {
	ID               uint      `gorm:"primarykey"`
	Name             string
	Path             string    `gorm:"unique"`
	DetectedStack    []string  `gorm:"serializer:json"`
	AppliedTemplates []uint    `gorm:"serializer:json"`
	CreatedAt        time.Time
}
