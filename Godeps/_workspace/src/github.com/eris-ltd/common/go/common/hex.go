package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

//-------------------------------------------------------
// hex and ints

// keeps N bytes of the conversion
func NumberToBytes(num interface{}, N int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		// TODO: get this guy a return error?
	}
	//fmt.Println("btyes!", buf.Bytes())
	if buf.Len() > N {
		return buf.Bytes()[buf.Len()-N:]
	}
	return buf.Bytes()
}

// s can be string, hex, or int.
// returns properly formatted 32byte hex value
func Coerce2Hex(s string) string {
	//fmt.Println("coercing to hex:", s)
	// is int?
	i, err := strconv.Atoi(s)
	if err == nil {
		return "0x" + hex.EncodeToString(NumberToBytes(int32(i), i/256+1))
	}
	// is already prefixed hex?
	if len(s) > 1 && s[:2] == "0x" {
		if len(s)%2 == 0 {
			return s
		}
		return "0x0" + s[2:]
	}
	// is unprefixed hex?
	if len(s) > 32 {
		return "0x" + s
	}
	pad := strings.Repeat("\x00", (32-len(s))) + s
	ret := "0x" + hex.EncodeToString([]byte(pad))
	//fmt.Println("result:", ret)
	return ret
}

func IsHex(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[:2] == "0x" {
		return true
	}
	return false
}

func AddHex(s string) string {
	if len(s) < 2 {
		return "0x" + s
	}

	if s[:2] != "0x" {
		return "0x" + s
	}

	return s
}

func StripHex(s string) string {
	if len(s) > 1 {
		if s[:2] == "0x" {
			s = s[2:]
			if len(s)%2 != 0 {
				s = "0" + s
			}
			return s
		}
	}
	return s
}

func StripZeros(s string) string {
	i := 0
	for ; i < len(s); i++ {
		if s[i] != '0' {
			break
		}
	}
	return s[i:]
}

// hex and ints
//---------------------------------------------------------------------------
// reflection and json

func WriteJson(config interface{}, config_file string) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config_file, out.Bytes(), 0600)
	return err
}

func ReadJson(config interface{}, config_file string) error {
	b, err := ioutil.ReadFile(config_file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		fmt.Println("error unmarshalling config from file:", err)
		return err
	}
	return nil
}

func NewInvalidKindErr(kind, k reflect.Kind) error {
	return fmt.Errorf("Invalid kind. Expected %s, received %s", kind, k)
}

func FieldFromTag(v reflect.Value, field string) (string, error) {
	iv := v.Interface()
	st := reflect.TypeOf(iv)
	for i := 0; i < v.NumField(); i++ {
		tag := st.Field(i).Tag.Get("json")
		if tag == field {
			return st.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("Invalid field name")
}

// Set a field in a struct value
// Field can be field name or json tag name
// Values can be strings that can be cast to int or bool
//  only handles strings, ints, bool
func SetProperty(cv reflect.Value, field string, value interface{}) error {
	f := cv.FieldByName(field)
	if !f.IsValid() {
		name, err := FieldFromTag(cv, field)
		if err != nil {
			return err
		}
		f = cv.FieldByName(name)
	}
	kind := f.Kind()

	k := reflect.ValueOf(value).Kind()
	if k != kind && k != reflect.String {
		return NewInvalidKindErr(kind, k)
	}

	if kind == reflect.String {
		f.SetString(value.(string))
	} else if kind == reflect.Int {
		if k != kind {
			v, err := strconv.Atoi(value.(string))
			if err != nil {
				return err
			}
			f.SetInt(int64(v))
		} else {
			f.SetInt(int64(value.(int)))
		}
	} else if kind == reflect.Bool {
		if k != kind {
			v, err := strconv.ParseBool(value.(string))
			if err != nil {
				return err
			}
			f.SetBool(v)
		} else {
			f.SetBool(value.(bool))
		}
	}
	return nil
}

// reflection and json
//---------------------------------------------------------------------------
