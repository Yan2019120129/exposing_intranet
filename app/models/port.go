package models

func init() {
	ModelManage.SetModel("port", &Port{}, "映射端口")
}

// Port 映射的端口
type Port struct {
	BaseModel
	ClientId int    `json:"clientId" gorm:"type:int unsigned not null;comment:客户端Id"` // 客户端id
	Server   string `json:"server" gorm:"type:varchar(128) not null; comment:服务端端口"`  // 服务端端口
	Local    string `json:"local" gorm:"type:varchar(128) not null; comment:客户端端口"`   // 本地端口
	Comment  string `json:"comment" gorm:"type:varchar(1024); comment:端口注释"`          // 端口注释
}

func (Port) TableName() string { return "port" }
