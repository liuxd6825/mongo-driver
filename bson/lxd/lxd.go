package lxd

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// GetDbFieldName 驼峰式写法转为下划线写法
func GetFieldName(name string) string {
	if name == "id" {
		return "_id"
	}
	buffer := newBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// Case2camel
//
// @Description:  下划线写法转为驼峰式，并首字母小写  id -> id, _id -> id, _key -> key , key_name -> keyName
// @param fieldName
// @return string
func GetPropertyName(key string) string {
	if key == "" {
		return ""
	}
	if key == "_id" {
		key = "id"
	} else {
		key = camelString(key)
	}
	if strings.HasPrefix(key, "_") {
		key = key[1:]
	}

	return strings.ToLower(key[:1]) + key[1:]
}

// CamelString 蛇形转驼峰
// @Description:
// @param s 要转换的字符串
// @return string
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 内嵌bytes.Buffer，支持连写
type buffer struct {
	*bytes.Buffer
}

func newBuffer() *buffer {
	return &buffer{Buffer: new(bytes.Buffer)}
}

func (b *buffer) Append(i interface{}) *buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *buffer) append(s string) *buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}
