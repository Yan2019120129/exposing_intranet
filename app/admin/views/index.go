package views

import "github.com/gin-gonic/gin"

var (
	ViewFuncList []func(e *gin.Engine)
)

func InitView(engine *gin.Engine) {
	for _, v := range ViewFuncList {
		v(engine)
	}
	ViewFuncList = nil
}
