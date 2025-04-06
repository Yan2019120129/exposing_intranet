package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my-base/app/admin/service"
	"my-base/app/common/api"
)

type User struct {
	api.Api
	user service.User
}

func (u *User) GetPage(c *gin.Context) {
	path := c.FullPath()
	fmt.Println("---------------------", path)
	u.MakeContext(c).
		SetService(u.user.GetPage())
}

func (u *User) Get(c *gin.Context) {
	u.MakeContext(c).
		SetService(u.user.Get())
}
