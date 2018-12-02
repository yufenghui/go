package cmap

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
)

// randElement 会生成并返回一个伪随机元素值。
func randElement() interface{} {
	if i := rand.Int31(); i%3 != 0 {
		return i
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}

// randString 会生成并返回一个伪随机字符串。
func randString() string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rand.Int31())
	return hex.EncodeToString(buf.Bytes())
}
