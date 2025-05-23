package convert

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ToString returns the string result converted by src.
// ToString 将任意类型的输入转换为字符串
// 参数:
//	src interface{}: 待转换的输入
// 返回值:
//	string: 转换后的字符串
func ToString(src interface{}) string {
	// 使用类型断言和类型选择来检查src的类型
	switch v := src.(type) {
	// 处理整数类型
	case int, int8, int16, int32, int64:
		// 将整数转换为int64类型，然后使用strconv.FormatInt将其转换为字符串
		return strconv.FormatInt(ToInt64(v), 10)
	// 处理无符号整数类型
	case uint, uint8, uint16, uint32, uint64, uintptr:
		// 将无符号整数转换为uint64类型，然后使用strconv.FormatUint将其转换为字符串
		return strconv.FormatUint(ToUint64(v), 10)
	// 处理浮点数和复数类型
	case float32, float64, complex64, complex128:
		// 将浮点数或复数转换为float64类型，然后使用strconv.FormatFloat将其转换为字符串
		return strconv.FormatFloat(ToFloat64(v), 'f', -1, 64)
	// 处理字符串类型
	case string:
		// 直接返回字符串
		return v
	// 处理字节切片类型
	case []byte:
		// 将字节切片转换为字符串
		return string(v)
	// 处理rune切片类型
	case []rune:
		// 将rune切片转换为字符串
		return string(v)
	// 处理布尔类型
	case bool:
		// 使用strconv.FormatBool将布尔值转换为字符串
		return strconv.FormatBool(v)
	// 处理nil类型
	case nil:
		// 返回空字符串
		return ""
	// 默认情况，处理其他类型
	default:
		// 使用fmt.Sprint将其他类型转换为字符串
		return fmt.Sprint(v)
	}
}

// ToBool returns the bool result converted by src.
// ToBool 函数将任意类型的参数转换为布尔值。
// 参数：
//	src: 接口类型的参数，可以是整数、浮点数、复数、布尔值、字符串、字节切片或rune切片。
// 返回值：
//	如果参数可以转换为布尔值，则返回转换后的布尔值；否则返回false。
func ToBool(src interface{}) bool {
	// 根据src的类型进行类型断言
	switch v := src.(type) {
	// 如果src是整数类型（int, int8, int16, int32, int64）
	case int, int8, int16, int32, int64:
		// 将整数转换为int64类型，判断其是否大于0
		return ToInt64(v) > 0
	// 如果src是无符号整数类型（uint, uint8, uint16, uint32, uint64, uintptr）
	case uint, uint8, uint16, uint32, uint64, uintptr:
		// 将无符号整数转换为uint64类型，判断其是否大于0
		return ToUint64(v) > 0
	// 如果src是浮点数或复数类型（float32, float64, complex64, complex128）
	case float32, float64, complex64, complex128:
		// 将浮点数或复数转换为float64类型，判断其是否大于0
		return ToFloat64(v) > 0
	// 如果src是布尔类型
	case bool:
		// 直接返回src的值
		return v
	// 如果src是字符串、字节切片或rune切片
	case string, []byte, []rune:
		// 将src转换为字符串，然后使用strconv.ParseBool解析布尔值
		result, _ := strconv.ParseBool(ToString(v))
		// 返回解析后的布尔值
		return result
	// 默认情况，返回false
	default:
		return false
	}
}

// ToInt returns the int result converted by src.
func ToInt(src interface{}) int {
	return int(ToInt64(src))
}

// ToInt32 returns the int32 result converted by src.
func ToInt32(src interface{}) int32 {
	return int32(ToInt64(src))
}

// ToInt64 returns the int64 result converted by src.
func ToInt64(src interface{}) int64 {
	switch v := src.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case uintptr:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case complex64:
		return int64(real(v))
	case complex128:
		return int64(real(v))
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		v = strings.TrimSpace(v)
		index := strings.Index(v, ".")
		if index != -1 {
			v = v[:index]
		}
		result, _ := strconv.ParseInt(v, 10, 64)
		return result
	case []byte:
		return BytesToInt64(v)
	default:
		return 0
	}
}

// ToUint returns the uint result converted by src.
func ToUint(src interface{}) uint {
	return uint(ToUint64(src))
}

// ToUint32 returns the uint32 result converted by src.
func ToUint32(src interface{}) uint32 {
	return uint32(ToUint64(src))
}

// ToUint64 returns the uint64 result converted by src.
func ToUint64(src interface{}) uint64 {
	switch v := src.(type) {
	case int:
		return uint64(v)
	case int8:
		return uint64(v)
	case int16:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case uintptr:
		return uint64(v)
	case float32:
		return uint64(v)
	case float64:
		return uint64(v)
	case complex64:
		return uint64(real(v))
	case complex128:
		return uint64(real(v))
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		v = strings.TrimSpace(v)
		index := strings.Index(v, ".")
		if index != -1 {
			v = v[:index]
		}
		result, _ := strconv.ParseUint(v, 10, 64)
		return result
	case []byte:
		return BytesToUint64(v)
	default:
		return 0
	}
}

// ToFloat returns the float64 result converted by src.
func ToFloat(src interface{}) float64 {
	return ToFloat64(src)
}

// ToFloat32 returns the float32 result converted by src.
func ToFloat32(src interface{}) float32 {
	return float32(ToFloat64(src))
}

// ToFloat64 returns the float64 result converted by src.
func ToFloat64(src interface{}) float64 {
	switch v := src.(type) {
	case int, int8, int16, int32, int64:
		return float64(ToInt64(v))
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return float64(ToUint64(v))
	case float32:
		return float64(v)
	case float64:
		return v
	case complex64:
		return float64(real(v))
	case complex128:
		return real(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		v = strings.TrimSpace(v)
		result, _ := strconv.ParseFloat(v, 64)
		return result
	case []byte:
		return BytesToFloat64(v)
	default:
		return 0
	}
}

// BytesToInt64 returns the int64 result converted by byte slice bytes.
func BytesToInt64(bytes []byte) int64 {
	return int64(BytesToUint64(bytes))
}

// Int64ToBytes returns the byte slice result converted by int64 i.
func Int64ToBytes(i int64) []byte {
	return Uint64ToBytes(uint64(i))
}

// BytesToUint64 returns the uint64 result converted by byte slice bytes.
func BytesToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

// Uint64ToBytes returns the byte slice result converted by uint64 i.
func Uint64ToBytes(i uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	return bytes
}

// BytesToFloat64 returns the float64 result converted by byte slice bytes.
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// Float64ToBytes returns the byte slice result converted by float64 f.
func Float64ToBytes(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

// BytesToRunes returns the rune slice result converted by byte slice bytes.
func BytesToRunes(bytes []byte) []rune {
	size, count := utf8.RuneCount(bytes), 0
	runes := make([]rune, size)
	for i := 0; i < size; i++ {
		r, c := utf8.DecodeRune(bytes[count:])
		runes[i], count = r, count+c
	}
	return runes
}

// RunesToBytes returns the byte slice result converted by rune slice runes.
func RunesToBytes(runes []rune) []byte {
	size := 0
	for _, r := range runes {
		size += utf8.RuneLen(r)
	}
	bytes := make([]byte, size)
	count := 0
	for _, r := range runes {
		count += utf8.EncodeRune(bytes[count:], r)
	}
	return bytes
}

// BytesEncodeHex returns the hexadecimal encoding string result of byte slice bytes.
func BytesEncodeHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HexDecodeBytes returns the byte slice result represented by the hexadecimal string h.
func HexDecodeBytes(h string) []byte {
	bytes, _ := hex.DecodeString(h)
	return bytes
}

// BytesEncodeHexs returns the hexadecimal encoding byte slice result of byte slice bytes.
func BytesEncodeHexs(bytes []byte) []byte {
	hexs := make([]byte, hex.EncodedLen(len(bytes)))
	n := hex.Encode(hexs, bytes)
	return hexs[:n]
}

// HexsDecodeBytes returns the byte slice result represented by the hexadecimal byte slice hs.
func HexsDecodeBytes(hs []byte) []byte {
	bytes := make([]byte, hex.DecodedLen(len(hs)))
	n, err := hex.Decode(bytes, hs)
	if err != nil {
		return nil
	}
	return bytes[:n]
}

// ToBase returns the toBase string result converted by fromBase string src.
func ToBase(src string, fromBase, toBase int) string {
	i, err := strconv.ParseInt(src, fromBase, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatInt(i, toBase)
}

// DecToBin returns the binary string result converted by decimal int64 dec.
func DecToBin(dec int64) string {
	return strconv.FormatInt(dec, 2)
}

// BinToDec returns the decimal int64 result converted by binary string bin.
func BinToDec(bin string) int64 {
	dec, _ := strconv.ParseInt(strings.TrimPrefix(bin, "0b"), 2, 64)
	return dec
}

// HexToBin returns the binary string result converted by hexadecimal string hex.
func HexToBin(hex string) string {
	i, _ := strconv.ParseInt(strings.TrimPrefix(hex, "0x"), 16, 64)
	return strconv.FormatInt(i, 2)
}

// BinToHex returns the hexadecimal string result converted by binary string bin.
func BinToHex(bin string) string {
	i, _ := strconv.ParseInt(strings.TrimPrefix(bin, "0b"), 2, 64)
	return strconv.FormatInt(i, 16)
}

// DecToHex returns the hexadecimal string result converted by decimal int64 dec.
func DecToHex(dec int64) string {
	return strconv.FormatInt(dec, 16)
}

// HexToDec returns the decimal int64 result converted by hexadecimal string hex.
func HexToDec(hex string) int64 {
	dec, _ := strconv.ParseInt(strings.TrimPrefix(hex, "0x"), 16, 64)
	return dec
}

// isStruct reports whether i is struct.
func isStruct(i interface{}) (reflect.Value, bool) {
	if i == nil {
		return reflect.Value{}, false
	}

	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		return reflect.Value{}, false
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	return v, true
}

// getStructFieldName returns the struct field name.
func getStructFieldName(sf reflect.StructField) string {
	name := strings.SplitN(sf.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	} else if name == "" {
		return sf.Name
	}
	return name
}

// StructToInterfaceMap returns the map[string]interface{} result converted by struct s.
func StructToInterfaceMap(s interface{}, ignoreZeroValue ...bool) map[string]interface{} {
	m := make(map[string]interface{})
	v, ok := isStruct(s)
	if !ok {
		return m
	}

	ignore := false
	if len(ignoreZeroValue) != 0 {
		ignore = ignoreZeroValue[0]
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		name := getStructFieldName(t.Field(i))
		if name == "" {
			continue
		}
		for f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		if f.IsValid() {
			if ignore && f.IsZero() {
				continue
			}
			m[name] = f.Interface()
		} else {
			if ignore {
				continue
			}
			m[name] = nil
		}
	}

	return m
}

// StructToStringMap returns the map[string]string result converted by struct s.
func StructToStringMap(s interface{}, ignoreZeroValue ...bool) map[string]string {
	m := make(map[string]string)
	v, ok := isStruct(s)
	if !ok {
		return m
	}

	ignore := false
	if len(ignoreZeroValue) != 0 {
		ignore = ignoreZeroValue[0]
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		name := getStructFieldName(t.Field(i))
		if name == "" {
			continue
		}
		for f.Kind() == reflect.Ptr {
			f = f.Elem()
		}
		if f.IsValid() {
			if ignore && f.IsZero() {
				continue
			}
			m[name] = ToString(f.Interface())
		} else {
			if ignore {
				continue
			}
			m[name] = ""
		}
	}

	return m
}
