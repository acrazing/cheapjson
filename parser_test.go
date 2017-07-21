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

type Block struct {
	String    string
	True      bool
	False     bool
	Null      interface{}
	Int       int64
	NegInt    int64
	Float     float64
	NegFloat  float64
	Float2    float64
	NegFloat2 float64
	String2   string
	Array     []interface{}
}

var block = Block{
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
}

var bigData = map[string]map[int]map[string]Block{}
var bigInput []byte

func init() {
	for i := 0; i < 10; i++ {
		key1 := strings.Repeat("中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文", i)
		bigData[key1] = map[int]map[string]Block{}
		for j := 0; j < 10; j++ {
			key2 := math.MaxInt64 - j
			bigData[key1][key2] = map[string]Block{}
			for k := 0; k < 10; k++ {
				key3 := strings.Repeat("中文中文\r\n\t\f\b中文中😍", i)
				bigData[key1][key2][key3] = block
			}
		}
	}
	bigInput, _ = json.MarshalIndent(bigData, "", "  ")
	log.Printf("big input size: %d", len(bigInput))
}

var testBlock = map[string]interface{}{
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

var testInput, _ = json.MarshalIndent(testBlock, "", "  ")

func TestUnmarshal(t *testing.T) {
	value, err := cheapjson.Unmarshal(testInput)
	assert.Nil(t, err, "should not throw error")
	assert.Equal(t, testBlock, value.Value(), "strict same")
	value, err = cheapjson.Unmarshal([]byte("\"\\ud83d\\ude02\\ud83d\\ude03\\u4e2d\\u56fd\\u4ebA\""))
	assert.Nil(t, err)
	assert.Equal(t, "😂😃中国人", value.String())
}

func TestSimpleJson(t *testing.T) {
	value := simplejson.New()
	err := value.UnmarshalJSON(bigInput)
	assert.Nil(t, err)
}

func BenchmarkUnmarshal(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, _ := cheapjson.Unmarshal(bigInput)
			assert.NotNil(b, value)
		}
	})
}

func BenchmarkSimpleJson(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value := simplejson.New()
			err := value.UnmarshalJSON(bigInput)
			assert.Nil(b, err)
		}
	})
}
