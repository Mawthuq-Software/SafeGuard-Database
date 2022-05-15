package db

import "time"

type StandardModel struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Policies struct {
	StandardModel
	Name        string `gorm:"unique"`
	Permissions string
}

type Groups struct {
	StandardModel
	Name string `gorm:"unique"`
}

type GroupPolicies struct {
	StandardModel
	GroupID  int
	Groups   Groups `gorm:"foreignKey:GroupID"`
	PolicyID int
	Policies Policies `gorm:"foreignKey:PolicyID"`
}

type Authentications struct {
	StandardModel
	Username string `gorm:"unique"`
	Password string
	Email    string `gorm:"unique"`
}

type Users struct {
	StandardModel
	GroupID         int
	Groups          Groups `gorm:"foreignKey:GroupID"`
	AuthID          int
	Authentications Authentications `gorm:"foreignKey:AuthID"`
}

type Servers struct {
	StandardModel
	Name       string `gorm:"unique"`
	Region     string
	Country    string
	IPAddress  string
	Online     bool
	LastOnline time.Time
	// GroupID    int
	// Groups     Groups `gorm:"foreignKey:GroupID"`
}
