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

type UserGroups struct {
	StandardModel
	UserID  int
	Users   Users `gorm:"foreignKey:UserID"`
	GroupID int
	Groups  Groups `gorm:"foreignKey:GroupID"`
}

type Users struct {
	StandardModel
	AuthID          int
	Authentications Authentications `gorm:"foreignKey:AuthID"`
}

type UserKeys struct {
	StandardModel
	UserID int
	Users  Users `gorm:"foreignKey:UserID"`
	KeyID  int
	Keys   Keys `gorm:"foreignKey:KeyID"`
}

type Keys struct {
	StandardModel
	ServerID       int
	Servers        Servers `gorm:"foreignKey:ServerID"`
	ServerKeyID    int     `gorm:"unqiue"`
	UsedBandwidth  int
	TotalBandwidth int
	PublicKey      string
	PresharedKey   string
}

type Servers struct {
	StandardModel
	Name       string `gorm:"unique"`
	Region     string
	Country    string
	IPAddress  string
	Online     bool
	LastOnline time.Time
}

type KeyIPv4 struct {
	StandardModel
	KeyID  int
	Keys   Keys `gorm:"foreignKey:KeyID"`
	IPv4ID int
	IPv4   IPv4 `gorm:"foreignKey:IPv4ID"`
}

type IPv4 struct {
	StandardModel
	Address string
}

type IPv4Interfaces struct {
	StandardModel
	InterfaceID         int
	WireguardInterfaces WireguardInterfaces `gorm:"foreignKey:InterfaceID"`
	IPv4ID              int
	IPv4                IPv4 `gorm:"foreignKey:IPv4ID"`
}

type KeyIPv6 struct {
	StandardModel
	KeyID  int
	Keys   Keys `gorm:"foreignKey:KeyID"`
	IPv6ID int
	IPv6   IPv6 `gorm:"foreignKey:IPv6ID"`
}

type IPv6 struct {
	StandardModel
	Address string
}

type IPv6Interfaces struct {
	StandardModel
	WireguardInterfaces WireguardInterfaces `gorm:"foreignKey:InterfaceID"`
	IPv6ID              int
	IPv6                IPv6 `gorm:"foreignKey:IPv6ID"`
}

type WireguardInterfaces struct {
	StandardModel
	ListenPort int
	PublicKey  string
}

type ServerInterfaces struct {
	StandardModel
	ServerID            int
	Servers             Servers `gorm:"foreignKey:ServerID"`
	InterfaceID         int
	WireguardInterfaces WireguardInterfaces `gorm:"foreignKey:InterfaceID"`
}

type ServerTokens struct {
	StandardModel
	ServerID int
	Servers  Servers `gorm:"foreignKey:ServerID"`
	TokenID  int
	Tokens   Tokens `gorm:"foreignKey:TokenID"`
}

type Tokens struct {
	StandardModel
	AccessToken string `gorm:"unique"`
}
