package cheapjson_test

import (
	"math"
	"testing"

	"github.com/acrazing/cheapjson"
	"github.com/stretchr/testify/assert"
)

func TestNewValue(t *testing.T) {
	object := cheapjson.NewValue()
	object.AsObject(nil)
	assert.Equal(t, true, object.IsObject())
	assert.Equal(t, false, object.IsNull())
	hello := object.AddField("hello")
	hello.AsString("world")
	assert.Equal(t, "world", hello.String())
	assert.Equal(t, "world", object.Get("hello").String())
	integer := object.AddField("integer")
	integer.AsInt(math.MaxInt64)
	assert.Equal(t, int64(math.MaxInt64), integer.Int())
	assert.Equal(t, true, integer.IsInt())
	assert.Equal(t, true, integer.IsNumber())
	double := object.AddField("float")
	double.AsFloat(math.SmallestNonzeroFloat64)
	assert.Equal(t, math.SmallestNonzeroFloat64, double.Float())
	assert.Equal(t, true, double.IsNumber())
	assert.Equal(t, false, double.IsInt())
	assert.Equal(t, false, double.IsString())
	null := object.AddField("null")
	null.AsNull()
	assert.Equal(t, nil, null.Value())
	assert.Equal(t, true, null.IsNull())
	null = object.Get("hello", "world", "empty")
	assert.Nil(t, null)
	null = object.Ensure("hello", "world", "empty")
	assert.NotNil(t, null)
	assert.Nil(t, null.Value())
	array := object.AddField("array")
	array.AsArray(nil)
	array.AddElement().AsInt(1)
	array.AddElement().AsBool(true)
	array.AddElement().AsArray(nil)
	assert.Equal(t, true, array.Array()[1].IsTrue())
	assert.Equal(t, true, array.Array()[2].IsArray())
	assert.Equal(t, true, object.Get("array", "2").IsArray())
	object.Ensure("array", "2", "sub").AsBool(false)
	assert.Equal(t, false, object.Object()["array"].Object()["2"].IsArray())
	assert.Equal(t, true, object.Object()["array"].Object()["2"].Object()["sub"].IsFalse())
}
