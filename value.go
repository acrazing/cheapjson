package cheapjson

import "strconv"

type Value struct {
	value interface{}
}

type null int

var NULL = null(0)

func NewValue() *Value {
	return &Value{nil}
}

func (v *Value) AsObject(value map[string]*Value) {
	if value == nil {
		v.value = map[string]*Value{}
	} else {
		v.value = value
	}
}

func (v *Value) AsArray(value []*Value) {
	if value == nil {
		v.value = []*Value{}
	} else {
		v.value = value
	}
}

func (v *Value) AsInt(value int64) {
	v.value = value
}

func (v *Value) AsFloat(value float64) {
	v.value = value
}

func (v *Value) AsBool(ok bool) {
	v.value = ok
}

func (v *Value) AsNull() {
	v.value = NULL
}

func (v *Value) AsString(value string) {
	v.value = value
}

func (v *Value) AddField(key string) *Value {
	if values, ok := v.value.(map[string]*Value); ok {
		value := NewValue()
		values[key] = value
		return value
	}
	panic("not a object value")
}

func (v *Value) AddElement() *Value {
	if values, ok := v.value.([]*Value); ok {
		value := NewValue()
		v.value = append(values, value)
		return value
	}
	panic("not a array value")
}

func (v *Value) IsObject() bool {
	switch v.value.(type) {
	case map[string]*Value:
		return true
	default:
		return false
	}
}

func (v *Value) IsArray() bool {
	switch v.value.(type) {
	case []*Value:
		return true
	default:
		return false
	}
}

func (v *Value) IsInt() bool {
	switch v.value.(type) {
	case int64:
		return true
	default:
		return false
	}
}

func (v *Value) IsNumber() bool {
	switch v.value.(type) {
	case int64, float64:
		return true
	default:
		return false
	}
}

func (v *Value) IsBool() bool {
	switch v.value.(type) {
	case bool:
		return true
	default:
		return false
	}
}

func (v *Value) IsTrue() bool {
	switch v.value.(type) {
	case bool:
		return v.value == true
	default:
		return false
	}
}

func (v *Value) IsFalse() bool {
	switch v.value.(type) {
	case bool:
		return v.value == false
	default:
		return false
	}
}

func (v *Value) IsNull() bool {
	switch v.value.(type) {
	case null:
		return true
	default:
		return false
	}
}

func (v *Value) IsString() bool {
	switch v.value.(type) {
	case string:
		return true
	default:
		return false
	}
}

// return the value of the specified path
// if path not exist, will return nil
// if some path is array, will covert the
// path to integer, if covert error, will
// return nil rather than panic.
func (v *Value) Get(path... string) *Value {
	value := v
	index := 0
	var obj map[string]*Value
	var arr []*Value
	var err error
	var ok bool
	for _, key := range path {
		if obj, ok = value.value.(map[string]*Value); ok {
			if value, ok = obj[key]; !ok {
				return nil
			}
		} else if arr, ok = value.value.([]*Value); ok {
			index, err = strconv.Atoi(key)
			if err != nil {
				return nil
			}
			if index < 0 || len(arr) < index + 1 {
				return nil
			}
			value = arr[index]
		} else {
			return nil
		}
	}
	return value
}

// This will force add a path to a value
// requires all the values on the path is an object
// if not, will force covert to an object
func (v *Value) Ensure(path... string) *Value {
	temp := v
	var ok bool
	var obj map[string]*Value
	for _, field := range path {
		if !temp.IsObject() {
			temp.AsObject(nil)
		}
		if obj, ok = temp.value.(map[string]*Value); ok {
			if temp, ok = obj[field]; !ok {
				temp = NewValue()
				obj[field] = temp
			}
		} else {
			panic("any thing do not want")
		}
	}
	return temp
}

func (v *Value) Object() map[string]*Value {
	if value, ok := v.value.(map[string]*Value); ok {
		return value
	}
	panic("not a object value")
}

func (v *Value) Array() []*Value {
	if value, ok := v.value.([]*Value); ok {
		return value
	}
	panic("not a array value")
}

func (v *Value) Int() int64 {
	if value, ok := v.value.(int64); ok {
		return value
	}
	panic("not a int value")
}

func (v *Value) Float() float64 {
	if value, ok := v.value.(int64); ok {
		return float64(value)
	}
	if value, ok := v.value.(float64); ok {
		return float64(value)
	}
	panic("not a number value")
}

func (v *Value) String() string {
	if value, ok := v.value.(string); ok {
		return value
	}
	panic("not a string value")
}

func (v *Value) Value() interface{} {
	if v == nil {
		return nil
	}
	switch v.value.(type) {
	case null, nil:
		return nil
	case string, bool, int64, float64:
		return v.value
	case map[string]*Value:
		if values, ok := v.value.(map[string]*Value); ok {
			out := map[string]interface{}{}
			for key, value := range values {
				out[key] = value.Value()
			}
			return out
		}
		return nil
	case []*Value:
		if values, ok := v.value.([]*Value); ok {
			out := []interface{}{}
			for _, value := range values {
				out = append(out, value.Value())
			}
			return out
		}
		return nil
	default:
		return nil
	}
}
