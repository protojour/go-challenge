package main

import (
	"strconv"
	"testing"

	"example.com/m/api"
)

func TestHashWorker(t *testing.T) {
	chnl := make(chan [2]string)
	expectedIndex1 := 0
	expectedIndex2 := 1
	expectedIndex3 := 2
	expectedHash1 := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	expectedHash2 := "cb8379ac2098aa165029e3938a51da0bcecfc008fd6795f401178647f96c5b34"
	expectedHash3 := "3608bca1e44ea6c4d268eb6db02260269892c0b42b86bbf1e77a6fa16c3c9282"

	// test 1
	go api.HashWorker("abc", chnl, expectedIndex1)
	result := <-chnl
	ind, err := strconv.Atoi(result[0]) // converting index back to integer
	if err != nil {
		t.Errorf("strnconv failed in test 1 with error: (%s)", err)
	}
	if ind != expectedIndex1 {
		t.Errorf("wrong index: expected (%d)"+" got (%d)", expectedIndex1, ind)
	}
	if result[1] != expectedHash1 {
		t.Errorf("wrong index: expected (%s)"+"got (%s)", expectedHash1, result[1])
	}

	// test 2
	go api.HashWorker("def", chnl, expectedIndex2)
	result = <-chnl
	ind, err = strconv.Atoi(result[0]) // converting index back to integer
	if err != nil {
		t.Errorf("strnconv failed in test 2 with error: (%s)", err)
	}
	if ind != expectedIndex2 {
		t.Errorf("wrong index: expected (%d)"+" got (%d)", expectedIndex2, ind)
	}
	if result[1] != expectedHash2 {
		t.Errorf("wrong hash: expected (%s)"+"got (%s)", expectedHash2, result[1])
	}

	// test 3
	go api.HashWorker("xyz", chnl, expectedIndex3)
	result = <-chnl
	ind, err = strconv.Atoi(result[0]) // converting index back to integer
	if err != nil {
		t.Errorf("strnconv failed in test 3 with error: (%s)", err)
	}
	if ind != expectedIndex3 {
		t.Errorf("wrong index: expected (%d)"+" got (%d)", expectedIndex3, ind)
	}
	if result[1] != expectedHash3 {
		t.Errorf("wrong hash: expected (%s)"+"got (%s)", expectedHash3, result[1])
	}
}
