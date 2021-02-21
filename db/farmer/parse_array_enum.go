package farmer

import (
	"fmt"
	"strings"
)

func ParsePGSystemeEchangeEnumArray(enumArr string) []int {
	fmt.Println("ParsePGSystemeEchangeEnumArray")
	systemEchangeEnum := []int{}
	if enumArr[0] != '{' && enumArr[len(enumArr)-1] != '}' {
		return systemEchangeEnum
	}
	arr := strings.Split(enumArr[1:len(enumArr)-1], ",")
	print(arr)
	return systemEchangeEnum
}
