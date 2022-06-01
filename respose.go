package httptest

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gavv/httpexpect/v2"
)

type TestFunc interface {
}

type Responses []Response
type Response struct {
	Type   string                // httpest type , if empty use  Equal() function to test
	Key    string                // httptest data's key
	Value  interface{}           // httptest data's value
	Length int                   // httptest data's length,when the data are array or map
	Func   func(obj interface{}) // httpest func, you can add your test logic ,can be empty
}

// Keys return Responses object key array
func (res Responses) Keys() []string {
	keys := []string{}
	for _, re := range res {
		keys = append(keys, re.Key)
	}
	return keys
}

// IdKeys return Responses with id
func IdKeys() Responses {
	return Responses{
		{Key: "id", Value: 0, Type: "ge"},
	}
}

func Test(object *httpexpect.Object, reses ...Responses) {
	if len(reses) == 0 {
		return
	}

	//return once
	if len(reses) == 1 {
		reses[0].Test(object.Value("data").Object())
		return
	}

	array := object.Value("data").Array()
	length := int(array.Length().Raw())
	if length < len(reses) {
		fmt.Println("Return data not equal keys length")
		array.Length().Equal(len(reses))
		return
	}

	// return array
	for m, ks := range reses {
		if ks == nil {
			return
		}
		ks.Test(object.Value("data").Array().Element(m).Object())
	}
}

func Scan(object *httpexpect.Object, reses ...Responses) {
	if len(reses) == 0 {
		return
	}

	//return once
	if len(reses) == 1 {
		reses[0].Scan(object.Value("data").Object())
		return
	}

	array := object.Value("data").Array()
	length := int(array.Length().Raw())
	if length < len(reses) {
		fmt.Println("Return data not equal keys length")
		array.Length().Equal(len(reses))
		return
	}

	// return array
	for m, res := range reses {
		if res == nil {
			return
		}
		res.Scan(object.Value("data").Array().Element(m).Object())
	}
}

// Test Test Responses object
func (res Responses) Test(object *httpexpect.Object) {
	for _, rs := range res {
		if rs.Value == nil {
			continue
		}
		reflectTypeString := reflect.TypeOf(rs.Value).String()
		switch reflectTypeString {
		case "string":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "notempty" {
					object.Value(rs.Key).String().NotEmpty()
				} else {
					object.Value(rs.Key).String().Equal(rs.Value.(string))
				}
			}
		case "float64":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "ge" {
					object.Value(rs.Key).Number().Ge(rs.Value.(float64))
				} else {
					object.Value(rs.Key).Number().Equal(rs.Value.(float64))
				}
			}
		case "uint":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "ge" {
					object.Value(rs.Key).Number().Ge(rs.Value.(uint))
				} else {
					object.Value(rs.Key).Number().Equal(rs.Value.(uint))
				}
			}
		case "int":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "ge" {
					object.Value(rs.Key).Number().Ge(rs.Value.(int))
				} else {
					object.Value(rs.Key).Number().Equal(rs.Value.(int))
				}
			}
		case "[]httptest.Responses":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				valueLen := len(rs.Value.([]Responses))
				length := int(object.Value(rs.Key).Array().Length().Raw())
				if rs.Length == 0 {
					object.Value(rs.Key).Array().Length().Equal(valueLen)
				}
				if length > 0 {
					if valueLen == length {
						max := length
						if rs.Length > 0 {
							max = rs.Length
						}
						for i := 0; i < max; i++ {
							rs.Value.([]Responses)[i].Test(object.Value(rs.Key).Array().Element(i).Object())
						}
					} else {
						rs.Value.([]Responses)[0].Test(object.Value(rs.Key).Array().First().Object())
					}
				}
			}
		case "map[int][]httptest.Responses":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				values := rs.Value.(map[int][]Responses)
				length := len(values)
				if length > 0 {
					object.Value(rs.Key).Object().Keys().Length().Equal(length)
					for key, v := range values {
						for _, vres := range v {
							vres.Test(object.Value(rs.Key).Object().Value(strconv.FormatInt(int64(key), 10)).Object())
						}
					}
				}
			}
		case "httptest.Responses":
			rs.Value.(Responses).Test(object.Value(rs.Key).Object())
		case "[]uint":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				valueLen := len(rs.Value.([]uint))
				if rs.Length == 0 {
					object.Value(rs.Key).Array().Length().Equal(valueLen)
				}
				length := int(object.Value(rs.Key).Array().Length().Raw())
				if length > 0 {
					if valueLen == length {
						max := length
						if rs.Length > 0 {
							max = rs.Length
						}
						for i := 0; i < max; i++ {
							object.Value(rs.Key).Array().Element(i).Number().Equal(rs.Value.([]uint)[i])
						}
					} else {
						object.Value(rs.Key).Array().First().Number().Equal(rs.Value.([]uint)[0])
					}
				}
			}
		case "[]string":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "null" {
					object.Value(rs.Key).Null()
				} else if strings.ToLower(rs.Type) == "notnull" {
					object.Value(rs.Key).NotNull()
				} else {
					valueLen := len(rs.Value.([]string))
					if rs.Length == 0 {
						object.Value(rs.Key).Array().Length().Equal(valueLen)
					}
					length := int(object.Value(rs.Key).Array().Length().Raw())
					if length > 0 {
						if valueLen == length {
							max := length
							if rs.Length > 0 {
								max = rs.Length
							}
							for i := 0; i < max; i++ {
								object.Value(rs.Key).Array().Element(i).String().Equal(rs.Value.([]string)[i])
							}
						} else {
							object.Value(rs.Key).Array().First().String().Equal(rs.Value.([]string)[0])
						}
					}
				}
			}
		case "map[int]string":
			if rs.Func != nil {
				rs.Func(object.Value(rs.Key))
			} else {
				if strings.ToLower(rs.Type) == "null" {
					object.Value(rs.Key).Null()
				} else if strings.ToLower(rs.Type) == "notnull" {
					object.Value(rs.Key).NotNull()
				} else {
					values := rs.Value.(map[int]string)
					object.Value(rs.Key).Object().Keys().Length().Equal(len(values))
					for key, v := range values {
						object.Value(rs.Key).Object().Value(strconv.FormatInt(int64(key), 10)).Equal(v)
					}
				}
			}
		default:
			continue
		}
	}
	res.Scan(object)
}

// Scan Scan response data to Responses object.
func (res Responses) Scan(object *httpexpect.Object) {
	for k, rk := range res {
		if !Exist(object, rk.Key) {
			continue
		}
		if rk.Value == nil {
			continue
		}
		valueTypeName := reflect.TypeOf(rk.Value).String()
		switch valueTypeName {
		case "string":
			res[k].Value = object.Value(rk.Key).String().Raw()
		case "uint":
			res[k].Value = uint(object.Value(rk.Key).Number().Raw())
		case "int":
			res[k].Value = int(object.Value(rk.Key).Number().Raw())
		case "int32":
			res[k].Value = int32(object.Value(rk.Key).Number().Raw())
		case "float64":
			res[k].Value = object.Value(rk.Key).Number().Raw()
		case "[]httptest.Responses":
			valueLen := len(res[k].Value.([]Responses))
			if rk.Length > 0 {
				valueLen = rk.Length
			}
			if rk.Length == 0 {
				object.Value(rk.Key).Array().Length().Equal(valueLen)
			}
			length := int(object.Value(rk.Key).Array().Length().Raw())
			if length > 0 {
				if valueLen == length {
					for i := 0; i < valueLen; i++ {
						res[k].Value.([]Responses)[i].Scan(object.Value(rk.Key).Array().Element(i).Object())
					}
				} else {
					res[k].Value.([]Responses)[0].Scan(object.Value(rk.Key).Array().First().Object())
				}
			}
		case "httptest.Responses":
			rk.Value.(Responses).Scan(object.Value(rk.Key).Object())
		case "[]string":
			if strings.ToLower(rk.Type) == "null" {
				res[k].Value = []string{}
			} else if strings.ToLower(rk.Type) == "notnull" {
				continue
			} else {
				length := int(object.Value(rk.Key).Array().Length().Raw())
				if length == 0 {
					continue
				}
				reskey, ok := res[k].Value.([]string)
				if ok {
					var strings []string
					for i := 0; i < length; i++ {
						strings = append(reskey, object.Value(rk.Key).Array().Element(i).String().Raw())
					}
					res[k].Value = strings
				}
			}
		default:
			continue
		}
	}
}

// Exist Check object keys if the key is in the keys array.
func Exist(object *httpexpect.Object, key string) bool {
	objectKyes := object.Keys().Raw()
	for _, objectKey := range objectKyes {
		if key == objectKey.(string) {
			return true
		}
	}
	return false
}

// GetString return string value.
func (res Responses) GetString(key ...string) string {
	if len(key) == 0 {
		return ""
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return ""
			}
			key = keys
		}
	}

	for i := 0; i < len(key); i++ {
		for m, rk := range res {
			if rk.Value == nil {
				return ""
			}
			reflectTypeString := reflect.TypeOf(rk.Value).String()
			if key[i] == rk.Key {
				switch reflectTypeString {
				case "string":
					return rk.Value.(string)
				case "httptest.Responses":
					return res[m].Value.(Responses).GetString(key[i+1:]...)
				}
			}
		}

	}
	return ""
}

// GetStrArray return string array value.
func (rks Responses) GetStrArray(key string) []string {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return nil
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "[]string":
				return rk.Value.([]string)
			}
		}
	}
	return nil
}

// GetResponses return Resposnes Array value
func (rks Responses) GetResponses(key string) []Responses {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return nil
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "[]httptest.Responses":
				return rk.Value.([]Responses)
			}
		}
	}
	return nil
}

// GetResponsereturn Resposnes value
func (rks Responses) GetResponse(key string) Responses {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return nil
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "httptest.Responses":
				return rk.Value.(Responses)
			}
		}
	}
	return nil
}

// GetUint return uint value
func (rks Responses) GetUint(key ...string) uint {

	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	for i := 0; i < len(key); i++ {
		for m, rk := range rks {
			if key[i] == rk.Key {
				if rk.Value == nil {
					return 0
				}
				valueTypeName := reflect.TypeOf(rk.Value).String()
				switch valueTypeName {
				case "float64":
					return uint(rk.Value.(float64))
				case "int32":
					return uint(rk.Value.(int32))
				case "uint":
					return rk.Value.(uint)
				case "int":
					return uint(rk.Value.(int))
				case "httptest.Responses":
					return rks[m].Value.(Responses).GetUint(key[i:]...)
				}
			}
		}
	}

	return 0
}

// GetInt return int value
func (rks Responses) GetInt(key ...string) int {
	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	for i := 0; i < len(key); i++ {
		for m, rk := range rks {
			if key[i] == rk.Key {
				if rk.Value == nil {
					return 0
				}
				switch reflect.TypeOf(rk.Value).String() {
				case "float64":
					return int(rk.Value.(float64))
				case "int":
					return rk.Value.(int)
				case "int32":
					return int(rk.Value.(int32))
				case "uint":
					return int(rk.Value.(uint))
				case "httptest.Responses":
					return rks[m].Value.(Responses).GetInt(key[i+1:]...)
				}
			}
		}
	}

	return 0
}

// GetInt32 return int32.
func (rks Responses) GetInt32(key ...string) int32 {
	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	for i := 0; i < len(key); i++ {
		for m, rk := range rks {
			if key[i] == rk.Key {
				if rk.Value == nil {
					return 0
				}
				switch reflect.TypeOf(rk.Value).String() {
				case "float64":
					return int32(rk.Value.(float64))
				case "int32":
					return rk.Value.(int32)
				case "int":
					return int32(rk.Value.(int))
				case "uint":
					return int32(rk.Value.(uint))
				case "httptest.Responses":
					return rks[m].Value.(Responses).GetInt32(key[i+1:]...)
				}
			}
		}
	}

	return 0
}

func (rks Responses) GetFloat64(key ...string) float64 {
	if len(key) == 0 {
		return 0
	}

	if len(key) == 1 {
		k := key[0]
		if strings.Contains(k, ".") {
			keys := strings.Split(k, ".")
			if len(keys) == 0 {
				return 0
			}
			key = keys
		}
	}

	for i := 0; i < len(key); i++ {
		for m, rk := range rks {
			if key[i] == rk.Key {
				if rk.Value == nil {
					return 0
				}
				switch reflect.TypeOf(rk.Value).String() {
				case "float64":
					return rk.Value.(float64)
				case "int":
					return float64(rk.Value.(int))
				case "int32":
					return float64(rk.Value.(int32))
				case "uint":
					return float64(rk.Value.(uint))
				case "httptest.Responses":
					return rks[m].Value.(Responses).GetFloat64(key[i+1:]...)
				}
			}
		}
	}

	return 0
}

// GetId return id.
func (res Responses) GetId(key ...string) uint {
	if len(key) == 0 {
		key = append(key, "data", "id")
	}
	return res.GetUint(key...)
}
