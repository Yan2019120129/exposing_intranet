package service

type User struct {
}

// Get 获取用户信息
func (u *User) Get() (any, error) {
	return "user", nil
}

// GetPage 获取用户信息列表
func (u *User) GetPage() (any, error) {
	return "user list", nil
}
