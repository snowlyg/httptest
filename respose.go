package httptest

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gavv/httpexpect/v2"
)

type Responses []Response
type Response struct {
	Type  string
	Key   string
	Value interface{}
}

// Keys
func (res Responses) Keys() []string {
	keys := []string{}
	for _, re := range res {
		keys = append(keys, re.Key)
	}
	return keys
}

func IdKeys() Responses {
	return Responses{
		{Key: "id", Value: uint(0)},
	}
}

func (res Responses) Test(object *httpexpect.Object) Responses {
	for _, rs := range res {
		reflectTypeString := reflect.TypeOf(rs.Value).String()
		if rs.Value == nil {
			continue
		}
		switch reflectTypeString {
		case "string":
			if strings.ToLower(rs.Type) == "notempty" {
				object.Value(rs.Key).String().NotEmpty()
			} else {
				object.Value(rs.Key).String().Equal(rs.Value.(string))
			}
		case "float64":
			if strings.ToLower(rs.Type) == "ge" {
				object.Value(rs.Key).Number().Ge(rs.Value.(float64))
			} else {
				object.Value(rs.Key).Number().Equal(rs.Value.(float64))
			}
		case "uint":
			if strings.ToLower(rs.Type) == "ge" {
				object.Value(rs.Key).Number().Ge(rs.Value.(uint))
			} else {
				object.Value(rs.Key).Number().Equal(rs.Value.(uint))
			}
		case "int":
			if strings.ToLower(rs.Type) == "ge" {
				object.Value(rs.Key).Number().Ge(rs.Value.(int))
			} else {
				object.Value(rs.Key).Number().Equal(rs.Value.(int))
			}
		case "[]httptest.Responses":
			object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]Responses)))
			length := int(object.Value(rs.Key).Array().Length().Raw())
			if length > 0 {
				if len(rs.Value.([]Responses)) == length {
					for i := 0; i < length; i++ {
						rs.Value.([]Responses)[i].Test(object.Value(rs.Key).Array().Element(i).Object())
					}
				} else {
					rs.Value.([]Responses)[0].Test(object.Value(rs.Key).Array().First().Object())
				}
			}
		case "map[int][]httptest.Responses":
			values := rs.Value.(map[int][]Responses)
			object.Value(rs.Key).Object().Keys().Length().Equal(len(values))
			for key, v := range values {
				for _, vres := range v {
					vres.Test(object.Value(rs.Key).Object().Value(strconv.FormatInt(int64(key), 10)).Object())
				}
			}
		case "httptest.Responses":
			rs.Value.(Responses).Test(object.Value(rs.Key).Object())
		case "[]uint":
			object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]uint)))
			length := int(object.Value(rs.Key).Array().Length().Raw())
			if length > 0 {
				if len(rs.Value.([]uint)) == length {
					for i := 0; i < length; i++ {
						object.Value(rs.Key).Array().Element(i).Number().Equal(rs.Value.([]uint)[i])
					}
				} else {
					object.Value(rs.Key).Array().First().Number().Equal(rs.Value.([]uint)[0])
				}
			}
		case "[]string":
			if strings.ToLower(rs.Type) == "null" {
				object.Value(rs.Key).Null()
			} else if strings.ToLower(rs.Type) == "notnull" {
				object.Value(rs.Key).NotNull()
			} else {
				object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]string)))
				length := int(object.Value(rs.Key).Array().Length().Raw())
				if length > 0 {
					if len(rs.Value.([]string)) == length {
						for i := 0; i < length; i++ {
							object.Value(rs.Key).Array().Element(i).String().Equal(rs.Value.([]string)[i])
						}
					} else {
						object.Value(rs.Key).Array().First().String().Equal(rs.Value.([]string)[0])
					}
				}
			}
		case "map[int]string":
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
		default:
			continue
		}
	}
	return res.Scan(object)
}

func (res Responses) Scan(object *httpexpect.Object) Responses {
	for k, rk := range res {
		if !Exist(object, rk.Key) {
			continue
		}
		if rk.Value == nil {
			continue
		}
		switch reflect.TypeOf(rk.Value).String() {
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
			object.Value(rk.Key).Array().Length().Equal(len(rk.Value.([]Responses)))
			length := int(object.Value(rk.Key).Array().Length().Raw())
			if length > 0 {
				if len(rk.Value.([]Responses)) == length {
					for i := 0; i < length; i++ {
						rk.Value.([]Responses)[i].Scan(object.Value(rk.Key).Array().Element(i).Object())
					}
				} else {
					rk.Value.([]Responses)[0].Scan(object.Value(rk.Key).Array().First().Object())
				}
			}
		case "httptest.Responses":
			res[k].Value = rk.Value.(Responses).Scan(object.Value(rk.Key).Object())
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
	return res
}

func Exist(object *httpexpect.Object, key string) bool {
	objectKyes := object.Keys().Raw()
	for _, objectKey := range objectKyes {
		if key == objectKey.(string) {
			return true
		}
	}
	return false
}

func (res Responses) GetString(key string) string {
	var keys []string
	if strings.Contains(key, ".") {
		keys = strings.Split(key, ".")
		if len(keys) != 2 {
			return ""
		}
		key = keys[0]
	}
	for _, rk := range res {
		reflectTypeString := reflect.TypeOf(rk.Value).String()
		if key == rk.Key {
			if rk.Value == nil {
				return ""
			}
			switch reflectTypeString {
			case "string":
				return rk.Value.(string)
			case "httptest.Responses":
				return rk.Value.(Responses).GetString(keys[1])
			}
		}
	}
	return ""
}

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

func (rks Responses) GetUint(key string) uint {
	var keys []string
	if strings.Contains(key, ".") {
		keys = strings.Split(key, ".")
		if len(keys) != 2 {
			return 0
		}
		key = keys[0]
	}
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return 0
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "float64":
				return uint(rk.Value.(float64))
			case "int32":
				return uint(rk.Value.(int32))
			case "uint":
				return rk.Value.(uint)
			case "int":
				return uint(rk.Value.(int))
			case "httptest.Responses":
				return rk.Value.(Responses).GetUint(keys[1])
			}
		}
	}
	return 0
}

func (rks Responses) GetInt(key string) int {
	var keys []string
	if strings.Contains(key, ".") {
		keys = strings.Split(key, ".")
		if len(keys) != 2 {
			return 0
		}
		key = keys[0]
	}
	for _, rk := range rks {
		if key == rk.Key {
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
				return rk.Value.(Responses).GetInt(keys[1])
			}
		}
	}
	return 0
}
func (rks Responses) GetInt32(key string) int32 {
	var keys []string
	if strings.Contains(key, ".") {
		keys = strings.Split(key, ".")
		if len(keys) != 2 {
			return 0
		}
		key = keys[0]
	}
	for _, rk := range rks {
		if key == rk.Key {
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
				return rk.Value.(Responses).GetInt32(keys[1])
			}
		}
	}
	return 0
}

func (res Responses) GetId() uint {
	return res.GetUint("data.id")
}
