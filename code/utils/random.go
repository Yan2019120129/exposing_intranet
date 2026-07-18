package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

// Random 随机
type Random struct {
	runes []rune
}

// NewRandom 随机数
func NewRandom() *Random {
	return &Random{
		runes: []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
	}
}

// OrderSn 随机生成订单号
func (c *Random) OrderSn() string {
	firstCode := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	nowTime := time.Now()
	return firstCode[nowTime.Year()-2020] +
		nowTime.Format("20060102150405") +
		strconv.Itoa(int(uuid.New().ID()))
}

// SetLetterRunes 设置字母基数
func (c *Random) SetLetterRunes() *Random {
	c.runes = []rune("abcdefghijklmnopqrstuvwxyz")
	return c
}

// SetNumberRunes 设置数字基数
func (c *Random) SetNumberRunes() *Random {
	c.runes = []rune("0123456789")
	return c
}

// String 随机字符串
func (c *Random) String(n int) string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = c.runes[rand.Intn(len(c.runes))]
	}
	return string(b)
}

// Intn 随机数
func (c *Random) Intn(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
