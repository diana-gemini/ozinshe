package validations

import (
	"fmt"
	"ozinshe/db/initializers"
)

func IsUniqueValue(tableName, fieldName, value string) bool {
	var count int64

	result := initializers.DB.Table(tableName).Where(fieldName+" = ?", value).Count(&count)

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	}

	return count > 0
}

func IsExistValue(tableName, fieldName string, value interface{}) bool {
	var count int64

	result := initializers.DB.Table(tableName).Where(fieldName+" = ?", value).Count(&count)

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	}

	return count > 0
}
