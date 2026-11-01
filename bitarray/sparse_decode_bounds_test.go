package bitarray_test

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/bitarray"
)

func decodeWithoutPanic(t *testing.T, data []byte) {
	t.Helper()
	defer func() {
		if v := recover(); v != nil {
			t.Errorf("Unmarshal panicked: %v", v)
		}
	}()
	if value, err := bitarray.Unmarshal(data); err == nil || value != nil {
		t.Errorf("Unmarshal accepted incomplete or inconsistent sparse data: %#v", data)
	}
}

func TestUnmarshalRejectsTruncatedSparseEncoding(t *testing.T) {
	value := bitarray.NewSparseBitArray()
	for _, n := range []uint64{1, 64, 130} {
		value.SetBit(n)
	}
	encoded, err := bitarray.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	for cut := 0; cut < len(encoded); cut++ {
		t.Run(fmt.Sprint(cut), func(t *testing.T) {
			t.Run("retained capacity", func(t *testing.T) { decodeWithoutPanic(t, encoded[:cut]) })
			t.Run("limited capacity", func(t *testing.T) { decodeWithoutPanic(t, encoded[:cut:cut]) })
		})
	}
	decoded, err := bitarray.Unmarshal(encoded)
	if err != nil || !reflect.DeepEqual(decoded.ToNums(), []uint64{1, 64, 130}) {
		t.Fatalf("complete frame=%v,%v", decoded, err)
	}
}

func TestUnmarshalRejectsInconsistentSparseLengths(t *testing.T) {
	for _, tc := range []struct {
		name            string
		blocks, indices uint64
		body            []byte
	}{
		{"block count exceeds payload", 1000, 0, nil},
		{"index count exceeds payload", 0, 1000, nil},
		{"index count differs from blocks", 1, 0, make([]byte, 8)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := []byte{'S'}
			var count [8]byte
			binary.LittleEndian.PutUint64(count[:], tc.blocks)
			b = append(b, count[:]...)
			b = append(b, tc.body...)
			binary.LittleEndian.PutUint64(count[:], tc.indices)
			b = append(b, count[:]...)
			decodeWithoutPanic(t, b)
		})
	}
}

func TestUnmarshalSparseCompleteControls(t *testing.T) {
	for _, bits := range [][]uint64{nil, {0}, {1, 63, 64, 128, 4097}} {
		value := bitarray.NewSparseBitArray()
		for _, n := range bits {
			value.SetBit(n)
		}
		encoded, err := bitarray.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		got, err := bitarray.Unmarshal(append(encoded, 0, 0))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got.ToNums(), value.ToNums()) {
			t.Fatal("sparse roundtrip changed")
		}
	}
}

func TestDeserializeSparseRejectsEmptyInput(t *testing.T) {
	value := bitarray.NewSparseBitArray().(interface{ Deserialize([]byte) error })
	defer func() {
		if v := recover(); v != nil {
			t.Errorf("Deserialize panicked: %v", v)
		}
	}()
	if err := value.Deserialize(nil); err == nil {
		t.Fatal("Deserialize accepted missing input")
	}
}
