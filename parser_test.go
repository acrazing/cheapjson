package cheapjson_test

import (
	"encoding/json"
	"testing"
	"math"
	"github.com/stretchr/testify/assert"
	"."
	"strings"
	"github.com/bitly/go-simplejson"
	"log"
)

var normalBlock = map[string]interface{}{
	"string": "string",
	"true": true,
	"false": false,
	"null": nil,
	"int": int64(math.MaxInt64),
	"-int": -int64(math.MaxInt64),
	"float": math.MaxFloat64,
	"float2": math.SmallestNonzeroFloat64,
	"-float": -math.MaxFloat64,
	"-float2": -math.SmallestNonzeroFloat64,
	"文\r\n\t\f\b中文中😍😘😢😓文":"中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文",
	"文\r\n\t\f\b中中😍😘😢😓文": []interface{}{
		"string",
		true,
		false,
		nil,
		int64(math.MaxInt64),
		-int64(math.MaxInt64),
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		-math.MaxFloat64,
		-math.SmallestNonzeroFloat64,
		"中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文",
		[]interface{}{"😍😘😢😓", "\r\t\f\n\b"},
	},
}

var normalInput, _ = json.MarshalIndent(normalBlock, "", "  ")

var bigData = map[string]map[int]map[string]interface{}{}
var bigInput []byte
var deepData interface{} = map[string]interface{}{}
var deepInput []byte

func init() {
	for i := 0; i < 10; i++ {
		key1 := strings.Repeat("中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文", i)
		bigData[key1] = map[int]map[string]interface{}{}
		for j := 0; j < 30; j++ {
			key2 := math.MaxInt64 - j
			bigData[key1][key2] = map[string]interface{}{}
			for k := 0; k < 100; k++ {
				key3 := strings.Repeat("中文中文\r\n\t\f\b中文中😍", k)
				bigData[key1][key2][key3] = normalBlock
			}
		}
	}
	tempData := deepData
	for i := 0; i < 1000; i++ {
		if temp, ok := tempData.(map[string]interface{}); ok {
			tempData = map[string]interface{}{}
			temp[strings.Repeat("中文中文\r\n\t\f\b中文中😍😘😢😓文中文", i)] = tempData
		}
	}
	bigInput, _ = json.MarshalIndent(bigData, "", "  ")
	deepInput, _ = json.MarshalIndent(deepData, "", "  ")
	log.Printf("big input size: %d, normal input size: %d, deep input size: %d", len(bigInput), len(normalInput), len(deepInput))
}

func TestUnmarshal(t *testing.T) {
	value, err := cheapjson.Unmarshal(normalInput)
	assert.Nil(t, err, "should not throw error")
	assert.Equal(t, normalBlock, value.Value(), "strict same")
	value, err = cheapjson.Unmarshal([]byte("\"\\ud83d\\ude02\\ud83d\\ude03\\u4e2d\\u56fd\\u4ebA\""))
	assert.Nil(t, err)
	assert.Equal(t, "😂😃中国人", value.String())
}

func TestSimpleJson(t *testing.T) {
	value := simplejson.New()
	err := value.UnmarshalJSON(bigInput)
	assert.Nil(t, err)
}

func BenchmarkUnmarshalBigInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, _ := cheapjson.Unmarshal(bigInput)
			assert.NotNil(b, value)
		}
	})
}

func BenchmarkSimpleJsonBigInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value := simplejson.New()
			err := value.UnmarshalJSON(bigInput)
			assert.Nil(b, err)
		}
	})
}

func BenchmarkUnmarshalNormalInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, _ := cheapjson.Unmarshal(normalInput)
			assert.NotNil(b, value)
		}
	})
}

func BenchmarkSimpleJsonNormalInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value := simplejson.New()
			err := value.UnmarshalJSON(normalInput)
			assert.Nil(b, err)
		}
	})
}

func BenchmarkUnmarshalDeepInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, _ := cheapjson.Unmarshal(deepInput)
			assert.NotNil(b, value)
		}
	})
}

func BenchmarkSimpleJsonDeepInput(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value := simplejson.New()
			err := value.UnmarshalJSON(deepInput)
			assert.Nil(b, err)
		}
	})
}

