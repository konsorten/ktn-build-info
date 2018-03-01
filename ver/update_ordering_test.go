package ver

import (
	"strconv"
	"testing"
)

func TestOrderUpdates(t *testing.T) {
	actions := make(UpdateActions)

	actions["5"] = "a"
	actions["2"] = "$create$"
	actions["1"] = "$delete$"
	actions["6"] = "b"
	actions["0"] = "$delete$"
	actions["4"] = "$___"
	actions["3"] = "$create$"

	ordered := orderUpdates(actions)

	if len(actions) != len(ordered) {
		t.Fatal("Element count mismatch")
	}

	for i, u := range ordered {
		j, _ := strconv.Atoi(u.Path)

		if i != j {
			t.Fatalf("Action in wrong position: %v: %v", i, u)
		}
	}
}
