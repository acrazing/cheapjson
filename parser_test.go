package cheapjson_test

import (
	"encoding/json"
	"testing"
	"math"
	"github.com/stretchr/testify/assert"
	"."
	"strings"
	"github.com/bitly/go-simplejson"
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

var block = &Block{
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

var bigData = map[string]map[int]map[string]*Block{}

func init() {
	for i := 0; i < 30; i++ {
		temp := map[int]map[string]*Block{}
		bigData[strings.Repeat("中文中文中文\r\n\t\f\b中文中😍😘😢😓文中文", i)] = temp
		for j := 0; j < 30; j++ {
			temp2 := map[string]*Block{}
			temp[math.MaxInt64 - j] = temp2
			for k := 0; k < 500; k++ {
				temp2[strings.Repeat("中文中文\r\n\t\f\b中文中😍", i)] = block
			}
		}
	}
}

var bigInput, _ = json.MarshalIndent(bigData, "", "  ")

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
}

func TestSimpleJson(t *testing.T) {
	value := simplejson.New()
	err := value.UnmarshalJSON(bigInput)
	assert.Nil(t, err)
}

func TestUnmarshal2(t *testing.T) {
	value, err := cheapjson.Unmarshal(bigInput)
	assert.Nil(t, err)
	output, err := json.MarshalIndent(value.Value(), "", "  ")
	assert.Nil(t, err)
	assert.Equal(t, bigInput, output)
}

func BenchmarkUnmarshal(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			value, err := cheapjson.Unmarshal(bigInput)
			assert.Nil(b, err)
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
