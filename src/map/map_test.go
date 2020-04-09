package main

import (
	"testing"
)

// 通过性测试 方法名以Test开头
// 命令行 go test .
// 代码覆盖率测试
// 命令行 go test -coverprofile=c.out
// 再利用工具命令查看 go tool cover -html=c.out
func TestMap(t *testing.T) {
	tests := []struct {
		s   string
		len int
		str string
	}{
		//edge case
		{"abc", 3, "abc"},
		{"abcfabcbacbfdcada", 5, "acbfd"},
		// chinese support
		{"耿洋洋", 2, "耿洋"},
	}

	for _, tt := range tests {
		if l, ms := lengthOfNonRepeatingSubStr(tt.s); l != tt.len || ms != tt.str {
			t.Errorf("lengthOfNonRepeatingSubStr(%s),got maxlen：%d, maxstr:%s;expected maxlen：%d, maxstr:%s", tt.s, l, ms, tt.len, tt.str)
		}
	}
}

//性能测试 方法以Benchmark开头
//命令行 go test -bench .
// 分析性能
//go test-bench . -cpuprofile=cup.out
//go tool pprof cpu.out
func BenchmarkMap(b *testing.B) {
	s := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
	for i := 0; i < 10; i++ {
		s += s
	}
	// len(s)算出来的是字符串s的字节数
	b.Logf("len(s)=%d", len(s))
	len := 8
	str := "会花飞灰化肥挥发"
	// 重置时间，数据准备不算
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l, ms := lengthOfNonRepeatingSubStr(s)
		if l != len || ms != str {
			b.Errorf("lengthOfNonRepeatingSubStr(%s),got maxlen：%d, maxstr:%s;expected maxlen：%d, maxstr:%s", s, l, ms, len, str)
		}
	}

}
