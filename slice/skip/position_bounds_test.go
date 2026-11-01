package skip_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
)

type positionKey int

func (k positionKey) Compare(other common.Comparator) int {
	o := other.(positionKey)
	if k < o {
		return -1
	}
	if k > o {
		return 1
	}
	return 0
}
func positionList() *skip.SkipList {
	list := skip.New(uint64(0))
	list.Insert(positionKey(10), positionKey(20), positionKey(30))
	return list
}
func positionValues(list *skip.SkipList) []common.Comparator {
	values := make([]common.Comparator, 0, list.Len())
	for i := uint64(0); i < list.Len(); i++ {
		values = append(values, list.ByPosition(i))
	}
	return values
}
func TestSplitAtOutOfRangePreservesList(t *testing.T) {
	for _, index := range []uint64{2, 3, 4, ^uint64(0) - 1, ^uint64(0)} {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			list := positionList()
			want := positionValues(list)
			left, right := list.SplitAt(index)
			if left != list || right != nil || !reflect.DeepEqual(positionValues(list), want) {
				t.Fatalf("out-of-range split changed list: left=%v right=%v values=%v", left, right, positionValues(list))
			}
		})
	}
}
func TestReplaceAtPositionOutOfRangeIsNoOp(t *testing.T) {
	for _, index := range []uint64{3, 4, ^uint64(0) - 1, ^uint64(0)} {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			list := positionList()
			want := positionValues(list)
			list.ReplaceAtPosition(index, positionKey(99))
			if !reflect.DeepEqual(positionValues(list), want) {
				t.Fatal("out-of-range replacement changed entries")
			}
			if got := list.ByPosition(^uint64(0)); got != nil {
				t.Fatalf("invalid replacement changed sentinel: %v", got)
			}
			iter := list.IterAtPosition(^uint64(0))
			if iter.Next() || iter.Value() != nil {
				t.Fatal("invalid replacement created out-of-range iterator value")
			}
		})
	}
}
func TestSkipPositionBoundsControls(t *testing.T) {
	t.Run("valid split", func(t *testing.T) {
		list := positionList()
		left, right := list.SplitAt(0)
		if left != list || right == nil || !reflect.DeepEqual(positionValues(left), []common.Comparator{positionKey(10)}) || !reflect.DeepEqual(positionValues(right), []common.Comparator{positionKey(20), positionKey(30)}) {
			t.Fatal("valid split changed")
		}
	})
	t.Run("valid replacement", func(t *testing.T) {
		list := positionList()
		list.ReplaceAtPosition(1, positionKey(21))
		if list.ByPosition(1) != positionKey(21) {
			t.Fatal("valid replacement changed")
		}
		iter := list.IterAtPosition(1)
		if !iter.Next() || iter.Value() != positionKey(21) || !iter.Next() || iter.Value() != positionKey(30) || iter.Next() {
			t.Fatal("valid iteration changed")
		}
	})
	t.Run("empty", func(t *testing.T) {
		list := skip.New(uint64(0))
		list.ReplaceAtPosition(^uint64(0), positionKey(99))
		left, right := list.SplitAt(^uint64(0))
		iter := list.IterAtPosition(^uint64(0))
		if left != list || right != nil || list.Len() != 0 || list.ByPosition(^uint64(0)) != nil || iter.Next() || iter.Value() != nil {
			t.Fatal("empty position operations changed")
		}
	})
}
