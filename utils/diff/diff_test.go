package diff

import "testing"

func Test_diff(t *testing.T) {
	type k struct {
		X string
		Y string
	}

	a := k{X: "1", Y: "2"}
	b := k{X: "1", Y: "3"}

	c, err := Diff(a, b)
	if err != nil {
		t.Errorf("diff error %v", err)
		return
	}

	if c == nil {
		t.Errorf("response cant be nil")
		return
	}

	if c["y"] != "3" {
		t.Errorf("Expected 3, got: %v -> %v", c["y"], c)
		return
	}
}

func Test_diff_equal(t *testing.T) {
	type k struct {
		X string
	}

	type l struct {
		Y string
	}

	a := k{X: "1"}
	b := l{Y: "3"}

	_, err := Diff(a, b)
	if err == nil {
		t.Errorf("failed on check type of a and b values")
		return
	}
}

func Test_diff_void(t *testing.T) {
	type k struct {
		X string
		Y string
	}

	a := k{X: "1", Y: "2"}
	b := k{X: "1", Y: ""}

	c, err := Diff(a, b)
	if err != nil {
		t.Errorf("diff error %v", err)
		return
	}

	if c == nil {
		t.Errorf("response cant be nil")
		return
	}

	if c["y"] == "3" {
		t.Errorf("Expected void, got: %v -> %v", c["y"], c)
		return
	}
}

func Test_diff_pointer(t *testing.T) {
	type k struct {
		X string
		Y string
	}

	type l struct {
		X string
		K *k
	}

	a := k{X: "1", Y: "2"}
	b := k{X: "1", Y: "3"}
	z := l{X: "1", K: &a}
	h := l{X: "2", K: &b}

	c, err := Diff(&z, &h)
	if err != nil {
		t.Errorf("diff error %v", err)
		return
	}

	if c == nil {
		t.Errorf("response cant be nil")
		return
	}

	if c["x"] != "2" {
		t.Errorf("Expected 2, got: %v -> %v", c["x"], c)
		return
	}
}
