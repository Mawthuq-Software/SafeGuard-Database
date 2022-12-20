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
	UserID  int
	Users   Users `gorm:"foreignKey:UserID"`
	KeyID   int
	VPNKeys VPNKeys `gorm:"foreignKey:KeyID"`
}

type VPNKeys struct {
	StandardModel
	ServerID       int
	Servers        Servers `gorm:"foreignKey:ServerID"`
	UsedBandwidth  int
	TotalBandwidth int
	PublicKey      string `gorm:"unique"`
	PresharedKey   string `gorm:"unique"`
	Enabled        bool
}

type Servers struct {
	StandardModel
	Name       string `gorm:"unique"`
	Region     string
	Country    string
	IPAddress  string `gorm:"unique"`
	Online     bool
	LastOnline time.Time
}

type KeyIPv4 struct {
	StandardModel
	KeyID       int
	VPNKeys     VPNKeys `gorm:"foreignKey:KeyID"`
	IPv4ID      int
	PrivateIPv4 PrivateIPv4 `gorm:"foreignKey:IPv4ID"`
}

type PrivateIPv4 struct {
	StandardModel
	Address string
}

type KeyIPv6 struct {
	StandardModel
	KeyID       int
	VPNKeys     VPNKeys `gorm:"foreignKey:KeyID"`
	IPv6ID      int
	PrivateIPv6 PrivateIPv6 `gorm:"foreignKey:IPv6ID"`
}

type PrivateIPv6 struct {
	StandardModel
	Address string
}

type WireguardInterfaces struct {
	StandardModel
	ListenPort  int
	PublicKey   string `gorm:"unique"`
	IPv4Address string `gorm:"unique"`
	IPv6Address string `gorm:"unique"`
}

type ServerInterfaces struct {
	StandardModel
	ServerID            int     `gorm:"unique"`
	Servers             Servers `gorm:"foreignKey:ServerID"`
	InterfaceID         int
	WireguardInterfaces WireguardInterfaces `gorm:"foreignKey:InterfaceID"`
}

type ServerTokens struct {
	StandardModel
	ServerID int     `gorm:"unique"`
	Servers  Servers `gorm:"foreignKey:ServerID"`
	TokenID  int     `gorm:"unique"`
	Tokens   Tokens  `gorm:"foreignKey:TokenID"`
}

type Tokens struct {
	StandardModel
	Name        string
	AccessToken string `gorm:"unique"`
}

type ServerConfigurations struct {
	StandardModel
	ServerID       int
	Servers        Servers `gorm:"foreignKey:ServerID"`
	ConfigID       int
	Configurations Configurations `gorm:"foreignKey:ConfigID"`
}

type Configurations struct {
	StandardModel
	Name string `gorm:"unique"`
	DNS  string
	Mask int
}

type Subscriptions struct {
	StandardModel
	Name           string
	NumberOfKeys   int
	TotalBandwidth int
}

type UserSubscriptions struct {
	StandardModel
	UserID         int
	Users          Users `gorm:"foreignKey:UserID"`
	SubscriptionID int
	Subscriptions  Subscriptions `gorm:"foreignKey:SubscriptionID"`
	UsedBandwidth  int
	Expiry         time.Time
}
