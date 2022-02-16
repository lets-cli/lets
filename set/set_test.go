package set

import (
	"reflect"
	"sort"
	"testing"
)

func TestStringSet(t *testing.T) {
	t.Run("add string to set", func(t *testing.T) {
		set := NewStringSet()

		set.Add("a")
		set.Add("b")
		set.Add("a")
		set.Add("c")

		values := set.ToList()
		sort.Strings(values)
		if !reflect.DeepEqual(values, []string{"a", "b", "c"}) {
			t.Errorf("set must contain only unique elements, got: %s", values)
		}
	})
	t.Run("add many strings at once to set", func(t *testing.T) {
		set := NewStringSet()

		set.AddMany([]string{"a", "b", "c"})
		set.Add("c")

		values := set.ToList()
		sort.Strings(values)
		if !reflect.DeepEqual(values, []string{"a", "b", "c"}) {
			t.Errorf("set must contain only unique elements, got: %s", values)
		}
	})

	t.Run("remove string from set", func(t *testing.T) {
		set := NewStringSet()

		set.Add("a")
		set.Add("b")
		set.Add("c")
		set.Remove("c")

		values := set.ToList()
		sort.Strings(values)
		if !reflect.DeepEqual(values, []string{"a", "b"}) {
			t.Errorf("set contains element which must be deleted, got: %s", values)
		}
	})

	t.Run("remove string from set", func(t *testing.T) {
		set := NewStringSet()

		set.Add("a")
		set.Add("b")
		set.Add("c")

		if !set.Contains("c") {
			t.Errorf("set must contain element which was added, got: %s", set.ToList())
		}
	})
}
