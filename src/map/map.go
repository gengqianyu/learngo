package main

import "fmt"

func main() {
	a, b := lengthOfNonRepeatingSubStr("abcfabcbacbfdcada")
	fmt.Println("a:", a, "b:", b)
}

// 求一个字符串最长不重复字符的子串长度
func lengthOfNonRepeatingSubStr(s string) (int, string) {
	substr := ""
	// 此map容器，存放所有遍历过的字符，以及字符最后出现的位置。
	lastCurred := make(map[rune]int)
	// 当前找到的不含有重复字符子串的开始（不含有重复字符子串的开始索引）
	startIndex := 0
	//记录最大子串长度
	maxLength := 0

	// 以字符串"abcfabcdbacbcfada"为例
	for index, char := range []rune(s) {
		// 去map中获取当前字符，最后出现的索引位置
		lastIndex, ok := lastCurred[char]
		// 如果,当前遍历字符存在于map容器中，说明出现了重复
		// 并且当前遍历字符最后一次出现的索引位置lastIndex,比之前记录的不含有重复字符子串的开始索引startIndex靠后,
		// 则更新子串开始索引位置，以便计算（包含不重复字符）新子串的长度
		//fmt.Println("index:", index)
		//fmt.Println("startIndex:", startIndex)
		if ok && (lastIndex >= startIndex) {
			startIndex = lastIndex + 1
		}
		//fmt.Println("startIndex:", startIndex)
		//fmt.Println("lastIndex:", lastIndex, "ok:", ok)
		// 把当前子串长度算出来，和之前计算的子串长度进行对比
		// 如果当前子串的长度大于之前记录的最大子串长度，则将最大子串长度修改为当前长度
		if subLength := (index + 1) - startIndex; subLength > maxLength {
			maxLength = subLength
			substr = string([]rune(s)[startIndex : index+1])
		}
		// 修改字符最后一次出现的索引位置
		lastCurred[char] = index
		//fmt.Println("lastCurred:", lastCurred)
		//fmt.Println("------------------------------")
	}
	return maxLength, substr
}
