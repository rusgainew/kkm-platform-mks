package conversion

import (
	"fmt"
	"math"
)

// SafeIntToInt32 безопасно конвертирует int в int32, проверяя переполнение
func SafeIntToInt32(value int) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d exceeds int32 range", value)
	}
	return int32(value), nil
}

// SafeInt64ToInt32 безопасно конвертирует int64 в int32, проверяя переполнение
func SafeInt64ToInt32(value int64) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d exceeds int32 range", value)
	}
	return int32(value), nil
}

// SafeIntToInt32WithDefault безопасно конвертирует int в int32, возвращая значение по умолчанию при переполнении
func SafeIntToInt32WithDefault(value int, defaultValue int32) int32 {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return defaultValue
	}
	return int32(value)
}

// SafeInt64ToInt32WithDefault безопасно конвертирует int64 в int32, возвращая значение по умолчанию при переполнении
func SafeInt64ToInt32WithDefault(value int64, defaultValue int32) int32 {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return defaultValue
	}
	return int32(value)
}
