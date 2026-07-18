package models

const (
	StatusOn            = -1   // 断开
	StatusActive        = iota // 活跃
	StatusDisable              // 禁用
	LegacyStatusDisable = 3    // 兼容历史/文档中的禁用值
)

func IsDisabledStatus(status int) bool {
	return status == StatusDisable || status == LegacyStatusDisable
}

func init() {
	ModelManage.SetModel("client", &Client{}, "客户端")
}

// Client 客户端模型
type Client struct {
	BaseModel
	Name   string `json:"name" gorm:"type:varchar(128) not null; comment:客户端名称"`                // 客户端名
	Symbol string `json:"symbol" gorm:"type:varchar(256) not null; unique; comment:客户端标识"`      // 客户端标识
	Status int    `json:"status" gorm:"type:tinyint not null; comment:状态 -1断开 1活跃 2禁用，3为历史兼容值"` // 状态 -1断开 1活跃 2禁用，3为历史兼容值
}

func (Client) TableName() string { return "client" }

// ClientAndPort  关联模型
type ClientAndPort struct {
	Client   `json:"client"`
	PortList []*Port `json:"portList" gorm:"foreignKey:ClientId"`
}

// TableName 实现对应表名
func (cp *ClientAndPort) TableName() string {
	return "client"
}
