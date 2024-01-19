package main_test

import (
	"encoding/json"
	"strconv"
	"testing"
)

type DioS struct {
	Dio string `json:"dio"`
}

type Dio struct {
	dio uint16
}

type DacS struct {
	Dac string `json:"dac"`
}

func TestUnmarshall(t *testing.T) {
	s := `{"Dio": "1"}`

	var v DioS
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		t.Error(err)
	}
	var d Dio
	i, err := strconv.Atoi(v.Dio)
	if err != nil {
		t.Error(err)
	}

	if i != 1 {
		t.Errorf("Expected 1, got %d", d.dio)
	}

	d.dio = uint16(i)
	if d.dio != 1 {
		t.Errorf("Expected 1, got %d", d.dio)
	}
}
