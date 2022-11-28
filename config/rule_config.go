package config

type MisuRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100"`
	V0    string `gorm:"size:255"` // Subject
	V1    string `gorm:"size:255"` // Domen
	V2    string `gorm:"size:255"` // Object
	V3    string `gorm:"size:255"` // Action
	V4    string `gorm:"size:255"`
	V5    string `gorm:"size:255"`
	V6    string `gorm:"size:255"`
	V7    string `gorm:"size:255"`
}
