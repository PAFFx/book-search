package models

type Book struct {
	ID          uint   `gorm:"primaryKey"`      // Primary key
	HTMLHash    string `gorm:"unique;not null"` // Unique hash of the HTML content
	URL         string `gorm:"not null"`        // URL of the product page may not be unique
	ImageURL    string // URL of the product image
	Title       string // Title of the product
	ISBN        string // ISBN of the product
	Description string // Description of the product
}

type Author struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type AuthorBook struct {
	AuthorID uint   `gorm:"primaryKey"`
	BookID   uint   `gorm:"primaryKey"`
	Author   Author `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Book     Book   `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
