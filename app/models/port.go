package models

import "gorm.io/gorm"

func init() {
	ModelManage.SetModel("port", &Port{}, "映射端口")
}

// Port 映射的端口
type Port struct {
	gorm.Model
	ClientId uint   `json:"clientId" gorm:"type:int unsigned not null;uniqueIndex:uidx_port_client_local,priority:1;comment:客户端Id"` // 客户端id
	Server   string `json:"server" gorm:"type:varchar(128) not null;uniqueIndex:uidx_port_server;comment:服务端端口"`                    // 服务端端口
	Local    string `json:"local" gorm:"type:varchar(128) not null;uniqueIndex:uidx_port_client_local,priority:2;comment:客户端端口"`    // 本地端口
	Comment  string `json:"comment" gorm:"type:varchar(1024); comment:端口注释"`                                                        // 端口注释
}

func (Port) TableName() string { return "port" }
