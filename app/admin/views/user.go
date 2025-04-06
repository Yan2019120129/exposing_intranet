package views

import (
	"github.com/gin-gonic/gin"
	models "my-base/app/model"
	"my-base/core/view"
)

func init() {
	ViewFuncList = append(ViewFuncList, UserView)
}

// UserView 用户视图
func UserView(engine *gin.Engine) {
	userView := view.NewView("用户")
	userTable := userView.BodyTable(&models.User{})
	userTable.Field("id", view.Number).SetDesc("用户ID")
	userTable.Field("name", view.Text).SetDesc("用户名")
	userView.Register(engine)
}
