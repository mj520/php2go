package php2go

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
)

type Protocol struct {
	Format []string
}

// Pack 编码
func (p *Protocol) Pack(args ...int64) (ret []byte) {
	la := len(args)
	ls := len(p.Format)
	if ls > 0 && la > 0 && ls == la {
		for i := 0; i < la; i++ {
			if p.Format[i] == "N8" {
				ret = append(ret, IntToBytes8(args[i])...)
			} else if p.Format[i] == "N4" {
				ret = append(ret, IntToBytes4(args[i])...)
			} else if p.Format[i] == "N2" {
				ret = append(ret, IntToBytes2(args[i])...)
			}
		}
	}
	return ret
}

// UnPack 解码
func (p *Protocol) UnPack(data []byte) []int64 {
	la := len(p.Format)
	ret := make([]int64, la)
	if la > 0 {
		for i := 0; i < la; i++ {
			if p.Format[i] == "N8" {
				ret[i] = Bytes8ToInt64(data[0:8])
				data = data[8:]
			} else if p.Format[i] == "N4" {
				ret[i] = Bytes4ToInt64(data[0:4])
				data = data[4:]
			} else if p.Format[i] == "N2" {
				ret[i] = Bytes2ToInt64(data[0:2])
				data = data[2:]
			}
		}
	}
	return ret
}

// Pack16 转成16进制编码字符串
func (p *Protocol) Pack16(args ...int64) (hString string) {
	//变成 byte 码
	hByte := p.Pack(args...)
	//转成 16进制字符串
	hString = p.DecToHexString(hByte)
	return hString
}

// UnPack16 解码16进制字符串
func (p *Protocol) UnPack16(hString string) (unIntList []int64) {
	HSByte := p.HexStringToByte(hString)
	unIntList = p.UnPack(HSByte)
	return unIntList
}

// DecToHexString 10进制转16进制字符串
func (p *Protocol) DecToHexString(decString []byte) (responseStr string) {
	for _, v := range decString {
		hexString := DecHex(int64(v))
		if len(hexString) < 2 {
			hexString = "0" + hexString
		}
		responseStr += hexString
	}
	return strings.ToLower(responseStr)
}

// HexStringToByte 16进制字符串转字节类型
func (p *Protocol) HexStringToByte(hexString string) (responseByte []byte) {
	ls := len(hexString) / 2
	for i := 0; i < ls; i++ {
		hex := hexString[0:2]
		hexString = hexString[2:]
		responseByte = append(responseByte, IntToBytes1(HexDec(hex))...)
	}
	return
}

// IntToBytes8 int64 转 byte8
func IntToBytes8(n int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(n))
	return buf
}

// IntToBytes4 int64 转 byte4
func IntToBytes4(n int64) []byte {
	nb := intToBytes(n, 4)
	return nb
}

// IntToBytes2 int64 转 byte2
func IntToBytes2(n int64) []byte {
	nb := intToBytes(n, 2)
	return nb
}

// IntToBytes1 int64 转 byte1
func IntToBytes1(n int64) []byte {
	nb := intToBytes(n, 1)
	return nb
}

//int64 转 byteN
func intToBytes(n int64, k int) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, n)
	hByte := bytesBuffer.Bytes()
	//c++ 高低位转换
	x := len(hByte)
	nb := make([]byte, k)
	for i := 0; i < k; i++ {
		nb[i] = hByte[x-i-1]
	}
	return nb
}

// Bytes2ToInt64 byte2 转 int64
func Bytes2ToInt64(b []byte) int64 {
	nb := []byte{0, 0, b[1], b[0]}
	bytesBuffer := bytes.NewBuffer(nb)
	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int64(x)
}

// Bytes4ToInt64 byte4 转 int64
func Bytes4ToInt64(b []byte) int64 {
	nb := []byte{b[3], b[2], b[1], b[0]}
	bytesBuffer := bytes.NewBuffer(nb)
	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int64(x)
}

// Bytes8ToInt64 byte8 转 int64
func Bytes8ToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// DecHex dechex()
func DecHex(number int64) string {
	return Dechex(number)
}

// HexDec hexdec()
func HexDec(str string) int64 {
	i, _ := Hexdec(str)
	return i
}

var HexChar = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//HexEncode 进制数转换   n 表示进制， 16 or 36 or 62
func HexEncode(num, hex int64) string {
	code := ""
	for num != 0 {
		residue := num % hex
		code = string(HexChar[residue]) + code
		num = num / hex
	}
	return code
}

//HexDecode 进制数还原   n 表示进制， 16 or 36 or 62
func HexDecode(code string, hex int64) int64 {
	v := 0.0
	length := len(code)
	for i := 0; i < length; i++ {
		s := string(code[i])
		index := strings.Index(HexChar, s)
		v += float64(index) * math.Pow(float64(hex), float64(length-1-i)) // 倒序
	}
	return int64(v)
}
