package convert

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	// 测试用例集合
	cases := []struct {
		src    interface{}
		expect string
	}{
		// 整数测试用例
		{src: 123456, expect: "123456"},
		{src: int64(13579), expect: "13579"},
		{src: -123456, expect: "-123456"},
		{src: int64(-13579), expect: "-13579"},
		{src: uint(123456), expect: "123456"},
		{src: uint64(13579), expect: "13579"},

		// 浮点数测试用例
		{src: 1234.5678, expect: "1234.5678"},
		{src: float32(1234.5), expect: "1234.5"},
		{src: 12345.12345 + 6666i, expect: "12345.12345"},
		{src: complex64(1234.5 + 6666i), expect: "1234.5"},
		{src: -1234.5678, expect: "-1234.5678"},
		{src: float32(-1234.5), expect: "-1234.5"},
		{src: -12345.12345 + 6666i, expect: "-12345.12345"},
		{src: complex64(-1234.5 + 6666i), expect: "-1234.5"},

		// 布尔值测试用例
		{src: true, expect: "true"},
		{src: false, expect: "false"},

		// 字符串测试用例
		{src: "to string", expect: "to string"},
		{src: []byte("to string"), expect: "to string"},
		{src: []rune("to string"), expect: "to string"},

		// 字符串切片测试用例
		{src: []string{"a", "b", "c", "d"}, expect: "[a b c d]"},

		// 空值测试用例
		{src: nil, expect: ""},
	}

	// 遍历测试用例
	for _, c := range cases {
		// 调用 ToString 函数
		get := ToString(c.src)
		// 断言预期值和实际值相等
		assert.Equal(t, c.expect, get)
	}
}

func TestToBool(t *testing.T) {
	// 测试用例列表
	cases := []struct {
		src    interface{}
		expect bool
	}{
		{src: 123456, expect: true},
		{src: int64(13579), expect: true},
		{src: -123456, expect: false},
		{src: int64(-13579), expect: false},
		{src: uint(123456), expect: true},
		{src: uint64(13579), expect: true},
		{src: 1234.5678, expect: true},
		{src: float32(1234.5), expect: true},
		{src: 12345.12345 + 6666i, expect: true},
		{src: complex64(1234.5 + 6666i), expect: true},
		{src: -1234.5678, expect: false},
		{src: float32(-1234.5), expect: false},
		{src: -12345.12345 + 6666i, expect: false},
		{src: complex64(-1234.5 + 6666i), expect: false},
		{src: true, expect: true},
		{src: false, expect: false},
		{src: "TOO BOOL", expect: false},
		{src: "TRUE", expect: true},
		{src: "false", expect: false},
		{src: []byte("true"), expect: true},
		{src: []byte("to string"), expect: false},
		{src: []rune("to string"), expect: false},
		{src: []string{"a", "b", "c", "d"}, expect: false},
		{src: nil, expect: false},
	}

	// 遍历测试用例
	for _, c := range cases {
		// 调用 ToBool 函数获取结果
		get := ToBool(c.src)
		// 使用 assert.Equal 验证结果是否符合预期
		assert.Equal(t, c.expect, get)
	}
}

func TestToInt(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect int
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: uint64(1<<64 - 1), expect: -1}, // overflow
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: -1234.5678, expect: -1234},
		{src: float32(-1234.5), expect: -1234},
		{src: -12345.12345 + 6666i, expect: -12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Int64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToInt(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToInt32(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect int32
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: uint64(1<<64 - 1), expect: -1}, // overflow
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: -1234.5678, expect: -1234},
		{src: float32(-1234.5), expect: -1234},
		{src: -12345.12345 + 6666i, expect: -12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Int64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToInt32(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToInt64(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect int64
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: uint64(1<<64 - 1), expect: -1}, // overflow
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: -1234.5678, expect: -1234},
		{src: float32(-1234.5), expect: -1234},
		{src: -12345.12345 + 6666i, expect: -12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Int64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToInt64(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToUint(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect uint
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Uint64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToUint(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToUint32(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect uint32
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: uint32(1<<32-1) - 123456 + 1},      // reverse
		{src: int64(-13579), expect: uint32(1<<32-1) - 13579 + 1}, // reverse
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: -1234.5678, expect: uint32(1<<32-1) - 1234 + 1},                 // reverse
		{src: float32(-1234.5), expect: uint32(1<<32-1) - 1234 + 1},           // reverse
		{src: -12345.12345 + 6666i, expect: uint32(1<<32-1) - 12345 + 1},      // reverse
		{src: complex64(-1234.5 + 6666i), expect: uint32(1<<32-1) - 1234 + 1}, // reverse
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Uint64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToUint32(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToUint64(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect uint64
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: uint64(1<<64-1) - 123456 + 1},      // reverse
		{src: int64(-13579), expect: uint64(1<<64-1) - 13579 + 1}, // reverse
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: uint64(1<<64 - 1), expect: uint64(1<<64 - 1)},
		{src: 1234.5678, expect: 1234},
		{src: float32(1234.5), expect: 1234},
		{src: 12345.12345 + 6666i, expect: 12345},
		{src: complex64(1234.5 + 6666i), expect: 1234},
		{src: -1234.5678, expect: uint64(1<<64-1) - 1234 + 1},                 // reverse
		{src: float32(-1234.5), expect: uint64(1<<64-1) - 1234 + 1},           // reverse
		{src: -12345.12345 + 6666i, expect: uint64(1<<64-1) - 12345 + 1},      // reverse
		{src: complex64(-1234.5 + 6666i), expect: uint64(1<<64-1) - 1234 + 1}, // reverse
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234},
		{src: Uint64ToBytes(12345678), expect: 12345678},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToUint64(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToFloat(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect float64
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: 1234.5678, expect: 1234.5678},
		{src: float32(1234.5), expect: 1234.5},
		{src: 12345.12345 + 6666i, expect: 12345.12345},
		{src: complex64(1234.5 + 6666i), expect: 1234.5},
		{src: -1234.5678, expect: -1234.5678},
		{src: float32(-1234.5), expect: -1234.5},
		{src: -12345.12345 + 6666i, expect: -12345.12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234.5},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234.5678},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234.5678},
		{src: Float64ToBytes(123.456), expect: 123.456},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToFloat(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToFloat32(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect float32
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: 1234.5678, expect: 1234.5678},
		{src: float32(1234.5), expect: 1234.5},
		{src: 12345.12345 + 6666i, expect: 12345.12345},
		{src: complex64(1234.5 + 6666i), expect: 1234.5},
		{src: -1234.5678, expect: -1234.5678},
		{src: float32(-1234.5), expect: -1234.5},
		{src: -12345.12345 + 6666i, expect: -12345.12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234.5},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234.5678},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234.5678},
		{src: Float64ToBytes(123.456), expect: 123.456},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToFloat32(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestToFloat64(t *testing.T) {
	cases := []struct {
		src    interface{}
		expect float64
	}{
		{src: 123456, expect: 123456},
		{src: int64(13579), expect: 13579},
		{src: -123456, expect: -123456},
		{src: int64(-13579), expect: -13579},
		{src: uint(123456), expect: 123456},
		{src: uint64(13579), expect: 13579},
		{src: 1234.5678, expect: 1234.5678},
		{src: float32(1234.5), expect: 1234.5},
		{src: 12345.12345 + 6666i, expect: 12345.12345},
		{src: complex64(1234.5 + 6666i), expect: 1234.5},
		{src: -1234.5678, expect: -1234.5678},
		{src: float32(-1234.5), expect: -1234.5},
		{src: -12345.12345 + 6666i, expect: -12345.12345},
		{src: complex64(-1234.5 + 6666i), expect: -1234.5},
		{src: true, expect: 1},
		{src: false, expect: 0},
		{src: "to string", expect: 0},
		{src: "1234", expect: 1234},
		{src: "1234.5678", expect: 1234.5678},
		{src: "  123456 ", expect: 123456},
		{src: "  1234.5678 ", expect: 1234.5678},
		{src: Float64ToBytes(123.456), expect: 123.456},
		{src: []rune("to string"), expect: 0},
		{src: []string{"a", "b", "c", "d"}, expect: 0},
		{src: nil, expect: 0},
	}

	for _, c := range cases {
		get := ToFloat64(c.src)
		assert.Equal(t, c.expect, get)
	}
}

func TestInt64BytesConversion(t *testing.T) {
	cases := []struct {
		src int64
	}{
		{src: 1<<63 - 1},
		{src: -1 << 63},
		{src: 1<<31 - 1},
		{src: -1 << 31},
		{src: 0},
		{src: 1},
	}

	for _, c := range cases {
		get := Int64ToBytes(c.src)
		assert.Equal(t, c.src, BytesToInt64(get))
		t.Log(get)
	}
}

func TestUint64BytesConversion(t *testing.T) {
	cases := []struct {
		src uint64
	}{
		{src: 1<<64 - 1},
		{src: 1<<32 - 1},
		{src: 1<<16 - 1},
		{src: 1<<8 - 1},
		{src: 0},
		{src: 1},
	}

	for _, c := range cases {
		get := Uint64ToBytes(c.src)
		assert.Equal(t, c.src, BytesToUint64(get))
		t.Log(get)
	}
}

func TestFloat64BytesConversion(t *testing.T) {
	cases := []struct {
		src float64
	}{
		{src: math.MaxFloat32},
		{src: math.SmallestNonzeroFloat32},
		{src: math.MaxFloat64},
		{src: math.SmallestNonzeroFloat64},
		{src: 123.123454678090000},
		{src: 0},
		{src: 1},
	}

	for _, c := range cases {
		get := Float64ToBytes(c.src)
		assert.Equal(t, c.src, BytesToFloat64(get))
		t.Log(get)
	}
}

func TestRunesBytesConversion(t *testing.T) {
	cases := []struct {
		src string
	}{
		{src: "Hello，世界。"},
		{src: "ABCDEFGH"},
		{src: "0123456789"},
		{src: "abcdefgh"},
		{src: "0123456789abcdefghABCDEFGH"},
		{src: ""},
	}

	for _, c := range cases {
		b1 := []byte(c.src)
		r1 := []rune(c.src)
		b2 := RunesToBytes(r1)
		r2 := BytesToRunes(b1)
		assert.Equal(t, b1, b2)
		assert.Equal(t, r1, r2)
		assert.Equal(t, c.src, string(r1))
		assert.Equal(t, c.src, string(r2))
		t.Log(b2, r2)
	}
}

func TestHexBytesConversion(t *testing.T) {
	cases := []struct {
		src string
	}{
		{src: "Hello，世界。"},
		{src: "ABCDEFGH"},
		{src: "0123456789"},
		{src: "abcdefgh"},
		{src: "0123456789abcdefghABCDEFGH"},
		{src: ""},
	}

	for _, c := range cases {
		b1 := []byte(c.src)
		hex := BytesEncodeHex(b1)
		b2 := HexDecodeBytes(hex)
		assert.Equal(t, b1, b2)
		t.Log(hex)
	}
}

func TestHexsBytesConversion(t *testing.T) {
	cases := []struct {
		src string
	}{
		{src: "Hello，世界。"},
		{src: "ABCDEFGH"},
		{src: "0123456789"},
		{src: "abcdefgh"},
		{src: "0123456789abcdefghABCDEFGH"},
		{src: ""},
	}

	for _, c := range cases {
		b1 := []byte(c.src)
		hexs := BytesEncodeHexs(b1)
		b2 := HexsDecodeBytes(hexs)
		assert.Equal(t, b1, b2)
		t.Log(hexs)
	}
}

func TestBaseConversion(t *testing.T) {
	cases := []struct {
		src int64
	}{
		{src: 1<<63 - 1},
		{src: -1 << 63},
		{src: 1<<31 - 1},
		{src: -1 << 31},
		{src: 0},
		{src: 1},
	}

	for _, c := range cases {
		// decimal conversion
		d2b := DecToBin(c.src)
		d2h := DecToHex(c.src)
		// binary conversion
		b2d := BinToDec(d2b)
		b2h := BinToHex(d2b)
		assert.Equal(t, c.src, b2d)
		assert.Equal(t, d2h, b2h)
		// hexadecimal conversion
		h2d := HexToDec(d2h)
		h2b := HexToBin(d2h)
		assert.Equal(t, c.src, h2d)
		assert.Equal(t, d2b, h2b)
		// base conversion
		srcString := strconv.FormatInt(c.src, 10)
		f10t2 := ToBase(srcString, 10, 2)
		f2t16 := ToBase(f10t2, 2, 16)
		f16t10 := ToBase(f2t16, 16, 10)
		assert.Equal(t, f16t10, srcString)
	}
}

func TestStructToMap(t *testing.T) {
	type Student struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
		Age    *int   `json:"age"`
		Room   string `json:"-"`
		Class  string
		School string
	}

	age := 18
	s := Student{
		Name:   "Alice",
		Number: 12345678910,
		Room:   "Room One",
		Age:    &age,
		Class:  "Class One",
	}

	assert.Equal(t, map[string]interface{}{
		"name":   "Alice",
		"number": 12345678910,
		"age":    18,
		"Class":  "Class One",
		"School": "",
	}, StructToInterfaceMap(s))

	assert.Equal(t, map[string]interface{}{
		"name":   "Alice",
		"number": 12345678910,
		"age":    18,
		"Class":  "Class One",
	}, StructToInterfaceMap(s, true))

	assert.Equal(t, map[string]string{
		"name":   "Alice",
		"number": "12345678910",
		"age":    "18",
		"Class":  "Class One",
		"School": "",
	}, StructToStringMap(s))

	assert.Equal(t, map[string]string{
		"name":   "Alice",
		"number": "12345678910",
		"age":    "18",
		"Class":  "Class One",
	}, StructToStringMap(s, true))

	assert.Empty(t, StructToInterfaceMap("err type"))
	assert.Empty(t, StructToStringMap("err type"))
	assert.Empty(t, StructToInterfaceMap(nil))
	assert.Empty(t, StructToStringMap(nil))
}
