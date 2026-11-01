package skip_test

import (
	"fmt"
	"github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
	"testing"
)

type lookupKey int

func (k lookupKey) Compare(other common.Comparator) int {
	o := other.(lookupKey)
	if k < o {
		return -1
	}
	if k > o {
		return 1
	}
	return 0
}

func TestGetWithPositionRequiresExactKey(t *testing.T) {
	list := skip.New(uint64(0))
	list.Insert(lookupKey(10), lookupKey(20), lookupKey(30))
	for _, query := range []lookupKey{0, 10, 15, 20, 25, 30, 35} {
		t.Run(fmt.Sprint(query), func(t *testing.T) {
			got, pos := list.GetWithPosition(query)
			if query == 10 || query == 20 || query == 30 {
				if got != query || pos != uint64(query/10-1) {
					t.Fatalf("lookup=%v,%d want%v,%d", got, pos, query, query/10-1)
				}
			} else if got != nil || pos != 0 {
				t.Fatalf("missing key returned successor%v at%d", got, pos)
			}
		})
	}
}
func TestGetWithPositionMissingAfterDeletion(t *testing.T) {
	list := skip.New(uint64(0))
	list.Insert(lookupKey(10), lookupKey(20), lookupKey(30))
	list.Delete(lookupKey(20))
	if got, pos := list.GetWithPosition(lookupKey(20)); got != nil || pos != 0 {
		t.Fatalf("removed key returned%v,%d", got, pos)
	}
}
func TestGetWithPositionEmptyControl(t *testing.T) {
	if got, pos := skip.New(uint64(0)).GetWithPosition(lookupKey(10)); got != nil || pos != 0 {
		t.Fatalf("empty lookup=%v,%d", got, pos)
	}
}
