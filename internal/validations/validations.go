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

func CheckPassword(password, passwordrepeat string) bool {
	if password == "" {
		return false
	}
	if len(password) < 4 {
		return false
	}
	if len(password) > 50 {
		return false
	}
	if PasswordRepeat(password, passwordrepeat) {
		return true
	}
	return false
}

func PasswordRepeat(password, passwordrepeat string) bool {
	if passwordrepeat == "" {
		return false
	}
	if password != passwordrepeat {
		return false
	}
	return true
}

func IsUniqueTwoValue(tableName, firstfieldName, secondfieldName string, firstValue, secondValue uint) bool {
	var count int64

	result := initializers.DB.Table(tableName).
		Where(firstfieldName+" = ? AND "+secondfieldName+" = ?", firstValue, secondValue).
		Count(&count)

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
