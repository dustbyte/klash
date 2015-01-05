package klash_test

import (
	"reflect"
	"testing"

	"github.com/mota/klash"
)

func GetParameter(params interface{}, fieldName string) (error, *klash.Parameter) {
	tsValue := reflect.ValueOf(params).Elem()
	fieldType, _ := tsValue.Type().FieldByName(fieldName)
	field := tsValue.FieldByName(fieldName)

	parameter := klash.NewParameter(fieldName, field)
	err := parameter.DiscoverProperties(fieldType.Tag)
	return err, parameter
}

func TestHelpTag(t *testing.T) {
	type TestStruct struct {
		Version bool `klash-help:"Print version and exit"`
	}
	ts := TestStruct{true}
	err, parameter := GetParameter(&ts, "Version")

	if err != nil {
		t.Error(err)
		return
	}

	if parameter.Help != "Print version and exit" {
		t.Error("Help message is supposed to be filled")
	}
}
