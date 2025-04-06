package view

import "github.com/gin-gonic/gin"

type Operate struct {
	Url     string  `json:"url"`
	AskType AskType `json:"askType"`
	fun     func(r *gin.Context)
}

func NewOperate(u string, askType AskType, fun func(r *gin.Context)) *Operate {
	return &Operate{Url: u, AskType: askType, fun: fun}
}
