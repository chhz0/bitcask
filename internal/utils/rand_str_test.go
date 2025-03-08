package utils

import "testing"

func Test_Rand_Str(t *testing.T) {
	for i := 0; i < 100000; i++ {
		s := RandomString(5, 15) // 长度在5-15之间随机
		t.Log(s)
		if len(s) < 5 || len(s) > 15 {
			t.Fatalf("Invalid length %d", len(s))
		}
	}
}
