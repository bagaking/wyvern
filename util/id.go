package util

import (
	"math/big"
	"time"
)

var (
	sequence  uint8  = 0
	lastTime  uint64 = 0
	partialID uint8  = 0
)

func init() {
	partialID = uint8(time.Now().UnixNano() % 1e8)
}

// GenID 生成一个 64 位整形唯一 id
// 39 位时间戳, 8 位分区 id, 8 位序列号, 8 位保留
func GenID(remain uint8) (uint64, error) {
	// 39 位时间戳
	timestamp := uint64(time.Now().UnixNano() / 1e6)
	if timestamp == lastTime {
		sequence++
	} else {
		sequence = 0
	}
	lastTime = timestamp

	return (timestamp << 24) | (uint64(partialID) << 16) | (uint64(sequence) << 8) | uint64(remain), nil
}

// base58EncodeBigInt encodes a big.Int to a base58 string
func base58EncodeBigInt(i *big.Int) string {
	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	var result []byte
	base := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}
	for i.Cmp(zero) != 0 {
		i.DivMod(i, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}
	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

// Base58EncodeUInt64 encodes a int64 to a base58 string
func Base58EncodeUInt64(i uint64) string {
	// 使用 base58EncodeBigInt
	return base58EncodeBigInt(new(big.Int).SetUint64(i))
}

// Base58Encode encodes a byte slice to a base58 string
func Base58Encode(b []byte) string {
	// 使用 base58EncodeBigInt
	return base58EncodeBigInt(new(big.Int).SetBytes(b))
}
