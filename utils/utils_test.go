package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
)

var logFn = log.Panic

//hash test
func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct{ Test string }{Test: "test"}
	t.Run("Hash is always same", func(t *testing.T) {
		x := Hash(s)
		t.Log(x)
		if x != hash {
			t.Errorf("Expected %s, got %s", hash, x)
		}
	})
	t.Run("Hash is hex encoded", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	//Output : e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}

func TestTobytes(t *testing.T) {
	s := "test"
	b := ToBytes(s)
	t.Log(b)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("Tobytes should retun slice of bytes got %s", k)
	}
}

func ExampleToBytes() {
	s := "test"
	b := ToBytes(s)

	fmt.Println(b)
	//Output : [7 12 0 4 116 101 115 116]
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}

	tests := []test{
		{input: "0:0:1", sep: ":", index: 2, output: "1"},
		{input: "0:1:0", sep: ":", index: 1, output: "1"},
		{input: "1:0:0", sep: ":", index: 0, output: "1"},
		{input: "1:1:1", sep: ":", index: 4, output: ""},
		{input: "1:0:0", sep: "/", index: 0, output: "1:0:0"},
	}

	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s and got %s", tc.output, got)
		}
	}
}

func ExampleSplitter() {
	s := "127.0.0.1:4000"
	port := Splitter(s, ":", 1)
	fmt.Println(port)
	//Output : 4000
}

func TestFromBytes(t *testing.T) {
	type testStruct struct {
		Test string
	}
	var restored testStruct
	ts := testStruct{"test"}
	b := ToBytes(ts)
	FromBytes(&restored, b)
	//Deep Eqaul = 깊은 복사 확인?
	if !reflect.DeepEqual(ts, restored) {
		t.Error("FromBytes() should restore struct.")
	}
}

func TestToJson(t *testing.T) {
	type testStruct struct {
		Test string
	}
	s := testStruct{Test: "test"}
	b := ToJSON(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("Expected %v and got %v", reflect.Slice, k)
	}
	var restored testStruct
	json.Unmarshal(b, &restored)
	if !reflect.DeepEqual(restored, s) {
		t.Error("toJson should encode JSON correctly.")
	}
}
