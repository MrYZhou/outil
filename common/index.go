package common

import (
	"encoding/hex"
	"math/rand"
	"time"
)

/*
获取随机字符

num 字符长度
*/
func Random(num int) string {
	rand.Seed(time.Now().UnixNano())
	uLen := 10
	b := make([]byte, uLen)
	rand.Read(b)

	rand_str := hex.EncodeToString(b)[0:uLen]
	return rand_str
}
