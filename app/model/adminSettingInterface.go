package models

// AdminSettingInviteRegisterReward 管理邀请注册奖励
type AdminSettingInviteRegisterReward struct {
	Register float64 `json:"register"` //	注册奖励
	Invite   float64 `json:"invite"`   //	邀请奖励
}

// AdminSettingDistributionReward 管理分销奖励
type AdminSettingDistributionReward struct {
	Level int8    `json:"level"` //	分销级数
	Type  int8    `json:"type"`  //	账单类型
	Rate  float64 `json:"rate"`  //	收益比例
}

// AdminSettingWalletWithdrawNums 管理钱包提现次数
type AdminSettingWalletWithdrawNums struct {
	Day  int8 `json:"day"`  //	天数
	Nums int8 `json:"nums"` //	次数
}

// AdminSettingStoreOrderTimeout 店铺订单超时设置
type AdminSettingStoreOrderTimeout struct {
	Payment      float64 `json:"payment"`      //	付款超时
	Receipt      float64 `json:"receipt"`      //	收货超时
	Comment      float64 `json:"comment"`      //	评论超时
	Delivery     float64 `json:"delivery"`     //	发货超时
	DeliveryNums float64 `json:"deliveryNums"` //	扣除店铺信用分
}
