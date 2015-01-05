package klash_test

import (
	"testing"
	"time"

	"github.com/mota/klash"
)

// For convertibility purpose

type Date struct {
	time.Time
}

func (d *Date) FromString(stringval string) error {
	val, err := time.Parse("2006-01-02", stringval)

	if err != nil {
		return err
	}

	d.Time = val
	return nil
}

// Begining of tests

func TestPositional(t *testing.T) {
	var parameters struct{}
	osargs := []string{"dummy", "test"}
	args, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if args[0] != "dummy" && args[1] != "test" {
		t.Error("Arguments are not properly returned")
	}
}

func TestSimple(t *testing.T) {
	type TestStruct struct {
		Dummy string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Dummy != "derp" {
		t.Error("A field is not set")
	}
}

func TestStop(t *testing.T) {
	type TestStruct struct {
		Name     string
		Nickname string
	}
	parameters := TestStruct{}
	osargs := []string{"--name", "Jack", "play", "--nickname", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Nickname != "" {
		t.Error("Parsing didn't stop")
	}
}

func TestNoStop(t *testing.T) {
	type TestStruct struct {
		Name     string
		Nickname string
	}
	parameters := TestStruct{}
	osargs := []string{"--name", "Jack", "play", "--nickname", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, false)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Nickname == "" {
		t.Error("Parsing has stopped")
	}
}

func TestNoValue(t *testing.T) {
	type TestStruct struct {
		Dummy string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err == nil {
		t.Error("Should produce an error")
	}
}

func TestEqual(t *testing.T) {
	type TestStruct struct {
		Dummy string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy=derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Dummy != "derp" {
		t.Error("A field is not set")
	}
}

func TestEqualNoValue(t *testing.T) {
	type TestStruct struct {
		Dummy string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy="}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err == nil {
		t.Error("Should produce an error")
	}
}

func TestUnknown(t *testing.T) {
	type TestStruct struct {
		Dummy string
	}
	parameters := TestStruct{}
	osargs := []string{"--derpy", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err == nil {
		t.Error("Should produce an error")
	}
}

func TestMultiWord(t *testing.T) {
	type TestStruct struct {
		DummyArg string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy-arg", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.DummyArg != "derp" {
		t.Error("A field is not set")
	}
}

func TestMultiWord2(t *testing.T) {
	type TestStruct struct {
		DummyArG string
	}
	parameters := TestStruct{}
	osargs := []string{"--dummy-arg", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.DummyArG != "derp" {
		t.Error("A field is not set")
	}
}

func TestBool(t *testing.T) {
	type TestStruct struct {
		Version bool
	}
	parameters := TestStruct{}
	osargs := []string{"--version"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Version != true {
		t.Error("A field is not set")
	}
}

func TestMultipleBool(t *testing.T) {
	type TestStruct struct {
		V bool
		D bool
	}
	parameters := TestStruct{}
	osargs := []string{"-v", "-d"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.V != true || parameters.D != true {
		t.Error("A field is not set")
	}
}

func TestInt(t *testing.T) {
	type TestStruct struct {
		Temperature int
		Age         uint
	}
	parameters := TestStruct{}
	osargs := []string{"--temperature", "-10", "--age", "27"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Temperature != -10 || parameters.Age != 27 {
		t.Error("A field is not set")
	}
}

func TestFloat(t *testing.T) {
	type TestStruct struct {
		Longitude float32
		Lattitude float64
	}
	parameters := TestStruct{}
	osargs := []string{"--longitude", "-72.7", "--lattitude", "27.4"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Longitude != -72.7 || parameters.Lattitude != 27.4 {
		t.Error("A field is not set")
	}
}

func TestSlice(t *testing.T) {
	type TestStruct struct {
		Grades []string
	}
	parameters := TestStruct{}

	grades := []string{"A-", "C+", "F", "B-"}
	osargs := make([]string, 0, len(grades)*2)
	for _, grade := range grades {
		osargs = append(osargs, "--grades", grade)
	}

	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Grades == nil {
		t.Error("Slice parameter is nil")
		return
	}

	for idx := 0; idx < len(grades); idx++ {
		if parameters.Grades[idx] != grades[idx] {
			t.Error("Slice order mismatch")
			return
		}
	}
}

func TestUnknownType(t *testing.T) {
	type TestStruct struct {
		Duration time.Time
	}
	parameters := TestStruct{}
	osargs := []string{"--duration", "derp"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err == nil {
		t.Error("Should produce an error")
	}
}

func TestConvertible(t *testing.T) {
	type TestStruct struct {
		Start Date
	}

	parameters := TestStruct{}
	osargs := []string{"--start", "2014-07-31"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err != nil {
		t.Error(err)
		return
	}

	if parameters.Start.Year() != 2014 ||
		parameters.Start.Month() != 7 ||
		parameters.Start.Day() != 31 {
		t.Error("A field is not set")
	}
}

func TestConvertibleError(t *testing.T) {
	type TestStruct struct {
		Start Date
	}

	parameters := TestStruct{}
	osargs := []string{"--start", "31/07/2014"}
	_, err := klash.ParseArguments("test", osargs, &parameters, true)

	if err == nil {
		t.Error("Should produce an error")
	}
}
