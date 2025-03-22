package request

// GetRoles ...
type GetRoles struct {
	Search string `query:"search"`
	Length int    `query:"length"`
	Page   int    `query:"page"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
	ID     uint   `query:"id"`
	Source string `query:"source"`
}

type GetMenu struct {
	Search   string `query:"search"`
	Length   int    `query:"length"`
	Page     int    `query:"page"`
	Sort     string `query:"sort"`
	Order    string `query:"order"`
	ParentID uint   `query:"parent_id"`
	Source   string `query:"source"`
}

// GetRoleMenus ...
type GetRoleMenus struct {
	Search string `query:"search"`
	Length int    `query:"length"`
	Page   int    `query:"page"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
	ID     uint   `query:"id"`
	Source string `query:"source"`
}

// PutRoleMenu ...
type PutRoleMenu struct {
	Menu []int `json:"menu" `
}

// GetPersonal ...
type GetPersonal struct {
	Search string `query:"search"`
	Length int    `query:"length"`
	Page   int    `query:"page"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
	ID     uint   `query:"id"`
	Source string `query:"source"`
}

// PutUser ...
type PutUser struct {
	Name             string `json:"name" `
	Email            string `json:"email" `
	EmailVerifiedAt  string `json:"emailVerifiedAt" `
	Password         string `json:"password"`
	Avatar           string `json:"avatar" `
	RememberToken    string `json:"rememberToken" `
	WorkGroupID      uint   `json:"workGroupID" `
	WorkIDNumber     string `json:"workIDNumber" `
	IDNumber         string `json:"idNumber" `
	Address          string `json:"address" `
	Mobile           string `json:"mobile" `
	TerritoryAreaID  uint   `json:"territoryAreaID" `
	TerritoryID      uint   `json:"territoryID" `
	TerritoryIDs     []uint `json:"territoryIDS" `
	Role             []int  `json:"role" `
	ZoneID           uint   `json:"zoneID" `
	SalesChannelID   []uint `json:"salesChannelID"`
	Country          string `json:"country"`
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	SubDistrict      string `json:"subDistrict"`
	PostalCode       string `json:"postalCode"`
	StatusActive     bool   `json:"statusActive"`
	StatusSfa        bool   `json:"statusSfa"`
	StatusAvatar     string `json:"statusAvatar"`
	RegistrationDate string `json:"registrationDate"`
	CodeExternal     string `json:"codeExternal"`
	Source           string `json:"source"`
	WorkPositionID   uint   `json:"workPositionID"`
}

// PostUser ...
type PostUser struct {
	Name             string `json:"name" `
	Email            string `json:"email" `
	EmailVerifiedAt  string `json:"emailVerifiedAt" `
	Password         string `json:"password"`
	Avatar           string `json:"avatar" `
	RememberToken    string `json:"rememberToken" `
	WorkGroupID      uint   `json:"workGroupID" `
	WorkIDNumber     string `json:"workIDNumber" `
	IDNumber         string `json:"idNumber" `
	Address          string `json:"address" `
	Mobile           string `json:"mobile" `
	TerritoryAreaID  uint   `json:"territoryAreaID" `
	TerritoryID      uint   `json:"territoryID" `
	TerritoryIDs     []uint `json:"territoryIDS" `
	Role             []int  `json:"role" `
	ZoneID           uint   `json:"zoneID" `
	SalesChannelID   []uint `json:"salesChannelID"`
	Country          string `json:"country"`
	Province         string `json:"province"`
	City             string `json:"city"`
	District         string `json:"district"`
	SubDistrict      string `json:"subDistrict"`
	PostalCode       string `json:"postalCode"`
	RegistrationDate string `json:"registrationDate"`
	CodeExternal     string `json:"codeExternal"`
	Source           string `json:"source"`
	WorkPositionID   uint   `json:"workPositionID"`
}

// GetWorkPersonal ...
type GetWorkPersonal struct {
	Search string `query:"search"`
	Length int    `query:"length"`
	Page   int    `query:"page"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
}
