package parser

import "testing"

func TestMap(t *testing.T) {
	mp := threadSafeMap[int, int]{internal: make(map[int]int)}

	mp.Add(1, 2)
	x, ok := mp.Load(1)
	if !ok || x != 2 {
		t.Fail()
	}

	x, ok = mp.LoadAndDelete(1)
	if !ok || x != 2 {
		t.Fail()
	}

	x, ok = mp.Load(1)
	if ok || x == 2 {
		t.Fail()
	}
}
