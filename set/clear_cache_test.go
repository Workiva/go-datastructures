package set_test

import (
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/set"
)

func TestClearInvalidatesFlattenedMembers(t *testing.T) {
	for _, tc := range []struct {
		name  string
		items []interface{}
		cache bool
	}{
		{name: "cached_single", items: []interface{}{"old"}, cache: true},
		{name: "cached_multiple", items: []interface{}{"old", "other"}, cache: true},
		{name: "uncached_members", items: []interface{}{"old", "other"}},
		{name: "cached_empty", cache: true},
		{name: "uncached_empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			members := set.New(tc.items...)
			defer members.Dispose()
			var previous, saved []interface{}
			if tc.cache {
				previous = members.Flatten()
				saved = append(previous[:0:0], previous...)
			}
			members.Clear()
			for i := 0; i < 2; i++ {
				if got := members.Flatten(); len(got) != 0 {
					t.Errorf("Flatten after Clear returned removed members: %v", got)
				}
				if members.Len() != 0 {
					t.Error("Clear left a nonzero length")
				}
				members.Clear()
			}
			for _, item := range tc.items {
				if members.Exists(item) {
					t.Errorf("Clear retained member %v", item)
				}
			}
			if !reflect.DeepEqual(previous, saved) {
				t.Errorf("Clear changed a previously returned slice: %v, want %v", previous, saved)
			}
			members.Add("new")
			if got := members.Flatten(); !reflect.DeepEqual(got, []interface{}{"new"}) {
				t.Errorf("Flatten after reuse = %v, want [new]", got)
			}
		})
	}
}
