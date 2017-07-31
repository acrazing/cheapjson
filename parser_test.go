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
	"strconv"
	"io/ioutil"
	"os"
	"flag"
	"github.com/a8m/djson"
	"github.com/acrazing/json-test-suite"
)

func initData() {
	var normalData interface{} = map[string]interface{}{
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
	var normalInput, _ = json.MarshalIndent(normalData, "", "  ")
	var bigData interface{} = map[string]interface{}{}
	var bigInput []byte
	var deepData interface{} = map[string]interface{}{}
	var deepInput []byte
	tempData := bigData
	size := 40
	for i := 0; i < size; i++ {
		if temp1, ok := tempData.(map[string]interface{}); ok {
			key1 := strings.Repeat("中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文", i + 1)
			temp2 := map[string]map[string]interface{}{}
			temp1[key1] = temp2
			for j := 0; j < size; j++ {
				key2 := strconv.Itoa(math.MaxInt64 - j)
				temp3 := map[string]interface{}{}
				temp2[key2] = temp3
				for k := 0; k < size; k++ {
					key3 := strings.Repeat("中文中文\r\n\t\f\b中文中😍", k + 1)
					temp3[key3] = normalData
				}
			}
		}
	}
	tempData = deepData
	for i := 0; i < 1000; i++ {
		if temp, ok := tempData.(map[string]interface{}); ok {
			tempData = map[string]interface{}{}
			temp["中文中文\r\n\t\f\b中文中😍😘😢😓文中文"] = tempData
		}
	}
	if temp, ok := tempData.(map[string]interface{}); ok {
		temp["1"] = normalData
	}
	bigInput, _ = json.MarshalIndent(bigData, "", "  ")
	deepInput, _ = json.MarshalIndent(deepData, "", "  ")
	os.MkdirAll("./data", 0777)
	ioutil.WriteFile("./data/normal.json", normalInput, 0777)
	ioutil.WriteFile("./data/big.json", bigInput, 0777)
	ioutil.WriteFile("./data/deep.json", deepInput, 0777)
}

var normalInput []byte
var bigInput []byte
var deepInput []byte
var profileCount int

func init() {
	flag.IntVar(&profileCount, "run-profile", 0, "specify run profile test")
	flag.Parse()
	if _, err := os.Stat("./data/normal.json"); os.IsNotExist(err) {
		initData()
	}
	normalInput, _ = ioutil.ReadFile("./data/normal.json")
	bigInput, _ = ioutil.ReadFile("./data/big.json")
	deepInput, _ = ioutil.ReadFile("./data/deep.json")
	log.Printf("big input size: %d, normal input size: %d, deep input size: %d", len(bigInput), len(normalInput), len(deepInput))
}

func TestUnmarshal(t *testing.T) {
	value, err := cheapjson.Unmarshal(normalInput)
	assert.Nil(t, err)
	jsonOutput, err := json.MarshalIndent(value.Value(), "", "  ")
	assert.Nil(t, err)
	assert.Equal(t, normalInput, jsonOutput)
	value, err = cheapjson.Unmarshal(bigInput)
	assert.Nil(t, err)
	jsonOutput, err = json.MarshalIndent(value.Value(), "", "  ")
	assert.Nil(t, err)
	assert.Equal(t, bigInput, jsonOutput)
	value, err = cheapjson.Unmarshal(deepInput)
	assert.Nil(t, err)
	jsonOutput, err = json.MarshalIndent(value.Value(), "", "  ")
	assert.Nil(t, err)
	assert.Equal(t, deepInput, jsonOutput)
	value, err = cheapjson.Unmarshal([]byte("\"\\ud83d\\ude02\\ud83d\\ude03\\u4e2d\\u56fd\\u4ebA\""))
	assert.Nil(t, err)
	assert.Equal(t, "😂😃中国人", value.String())
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

func TestProfileUnmarshal(t *testing.T) {
	for i := 0; i < profileCount; i++ {
		value, err := cheapjson.Unmarshal(bigInput)
		assert.NotNil(t, value)
		assert.Nil(t, err)
	}
}

func TestProfileSimpleJson(t *testing.T) {
	for i := 0; i < profileCount; i++ {
		value, err := simplejson.NewJson(bigInput)
		assert.NotNil(t, value)
		assert.Nil(t, err)
	}
}

func TestProfileDjson(t *testing.T) {
	for i := 0; i < profileCount; i++ {
		value, err := djson.Decode(bigInput)
		assert.Nil(t, err)
		assert.NotNil(t, value)
	}
}

func BenchmarkCompareMarshal(b *testing.B) {
	json_test_suite.CompareUnmarshal(map[string]func(data []byte) (interface{}, error){
		"cheapjson": func(data []byte) (interface{}, error) {
			return cheapjson.Unmarshal(data)
		},
		"djson": func(data []byte) (interface{}, error) {
			return djson.Decode(data)
		},
		"go-simplejson": func(data []byte) (interface{}, error) {
			return simplejson.NewJson(data)
		},
	}, "./json-test-suite/correct")
}
