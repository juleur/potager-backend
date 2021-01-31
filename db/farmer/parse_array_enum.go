package farmer

import (
	"npp_backend/entity/enums"
	"strings"
)

func ParsePGSystemeEchangeEnumArray(enumArr string) []enums.SystemeEchange {
	systemEchangeEnum := []enums.SystemeEchange{}
	if enumArr[0] != '{' && enumArr[len(enumArr)-1] != '}' {
		return systemEchangeEnum
	}
	arr := strings.Split(enumArr[1:len(enumArr)-1], ",")
	for _, v := range arr {
		for _, enum := range enums.AllSystemeEchange {
			if v == enum.String() {
				systemEchangeEnum = append(systemEchangeEnum, enum)
			}
		}
	}
	return systemEchangeEnum
}
