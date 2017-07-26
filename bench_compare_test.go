package cheapjson_test

import (
	"."
	"testing"
	"github.com/acrazing/json-test-suite"
	"github.com/a8m/djson"
	"github.com/bitly/go-simplejson"
)

func TestCompareMarshal(t *testing.T) {
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
