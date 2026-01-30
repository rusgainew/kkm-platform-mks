// Файл api-gateway/internal/interfaces/http/httputil/conversion.go содержит реализацию пакета httputil.
package httputil

import (
	"math"
	"strconv"
)

// ParseInt32FromString парсит строку в int32 с проверкой границ
func ParseInt32FromString(s string, defaultValue int32) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue, err
	}

	if val > math.MaxInt32 || val < math.MinInt32 {
		return defaultValue, nil
	}

	return int32(val), nil // #nosec G109 - range checked above
}

// ParseInt32FromStringOrDefault парсит строку в int32, возвращая значение по умолчанию при ошибке
func ParseInt32FromStringOrDefault(s string, defaultValue int32) int32 {
	val, err := ParseInt32FromString(s, defaultValue)
	if err != nil {
		return defaultValue
	}
	return val
}

// SafeInt64ToInt32 безопасно конвертирует int64 в int32
func SafeInt64ToInt32(val int64, defaultValue int32) int32 {
	if val > math.MaxInt32 || val < math.MinInt32 {
		return defaultValue
	}
	return int32(val) // #nosec G109 - range checked above
}
