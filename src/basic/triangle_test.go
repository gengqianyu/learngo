/*
测试
表格驱动测试
代码覆盖率 go test . -coverProfile
性能优化工具
http测试
文档以及示例代码
*/
package main

import "testing"

func TestTriangle(t *testing.T) {
	// 定义一个切片，里面是一个结构体
	tests := []struct{ a, b, c int }{
		{3, 4, 5},
		{5, 12, 13},
		{8, 15, 17},
		{12, 35, 37},
		{30000, 40000, 50000},
	}

	for _, tt := range tests {
		if actual := calcTriangle(tt.a, tt.b); actual != tt.c {
			//                                  得到     预期
			t.Errorf("calcTriganle(%d,%d);got %d;expected %d", tt.a, tt.b, actual, tt.c)
		}
	}
}
