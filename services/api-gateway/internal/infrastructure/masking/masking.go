package masking

import (
	"regexp"
	"strings"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

// SensitiveFieldMasker маскирует чувствительные данные в логах
type SensitiveFieldMasker struct {
	patterns map[string]*regexp.Regexp
	metrics  *observability.Metrics
}

// NewSensitiveFieldMasker создает новый masker с дефолтными паттернами
func NewSensitiveFieldMasker(metrics *observability.Metrics) *SensitiveFieldMasker {
	masker := &SensitiveFieldMasker{
		patterns: make(map[string]*regexp.Regexp),
		metrics:  metrics,
	}

	// Регистрируем паттерны для различных типов данных
	masker.RegisterPattern("password", `"password"\s*:\s*"[^"]*"`)
	masker.RegisterPattern("jwt_token", `"token"\s*:\s*"[^"]*"`)
	masker.RegisterPattern("api_key", `"api_key"\s*:\s*"[^"]*"`)
	masker.RegisterPattern("secret", `"secret"\s*:\s*"[^"]*"`)
	masker.RegisterPattern("authorization", `Authorization:\s*Bearer\s+[^\s]+`)
	masker.RegisterPattern("credit_card", `\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`)
	masker.RegisterPattern("ssn", `\b\d{3}-\d{2}-\d{4}\b`)
	masker.RegisterPattern("phone", `\b\d{3}[\s-]?\d{3}[\s-]?\d{4}\b`)
	masker.RegisterPattern("email", `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	masker.RegisterPattern("bearer_token", `Bearer\s+[^\s,]+`)

	return masker
}

// RegisterPattern регистрирует паттерн для маскирования
func (sfm *SensitiveFieldMasker) RegisterPattern(name string, pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	sfm.patterns[name] = regex
	return nil
}

// Mask маскирует все чувствительные данные в строке
func (sfm *SensitiveFieldMasker) Mask(input string) string {
	start := time.Now()
	result := input
	fieldCount := 0

	for fieldType, regex := range sfm.patterns {
		matches := regex.FindAllString(result, -1)
		fieldCount += len(matches)

		result = regex.ReplaceAllStringFunc(result, func(match string) string {
			sfm.metrics.RecordMaskedField(fieldType)
			return maskValue(match)
		})
	}

	if fieldCount > 0 {
		latency := time.Since(start).Seconds()
		for fieldType := range sfm.patterns {
			sfm.metrics.RecordMaskingLatency(fieldType, latency/float64(len(sfm.patterns)))
		}
	}

	return result
}

// MaskJSONField маскирует значение JSON поля по пути
func (sfm *SensitiveFieldMasker) MaskJSONField(jsonStr string, fieldName string) string {
	// Простой поиск паттерна "fieldName": "value"
	pattern := regexp.MustCompile(`"` + regexp.QuoteMeta(fieldName) + `"\s*:\s*"[^"]*"`)
	return pattern.ReplaceAllString(jsonStr, `"`+fieldName+`":"***"`)
}

// maskValue маскирует значение, сохраняя первые и последние символы
func maskValue(value string) string {
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}

	// Находим позицию значения (после : или = или пробела)
	var dataStart, dataEnd int

	// Ищем символ, после которого идет значение
	dataStart = 0
	for i, ch := range value {
		if ch == '"' || ch == '\'' {
			dataStart = i + 1
			break
		}
		if ch == ':' || ch == '=' || ch == ' ' {
			if i+1 < len(value) && (value[i+1] == '"' || value[i+1] == '\'') {
				dataStart = i + 2
			} else {
				dataStart = i + 1
			}
		}
	}

	// Находим конец значения
	dataEnd = len(value)
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == '"' || value[i] == '\'' {
			dataEnd = i
			break
		}
	}

	if dataStart >= dataEnd {
		return strings.Repeat("*", len(value))
	}

	// Извлекаем значение
	actualData := value[dataStart:dataEnd]

	// Если очень короткое значение, маскируем полностью
	if len(actualData) <= 6 {
		return value[:dataStart] + strings.Repeat("*", len(actualData)) + value[dataEnd:]
	}

	// Показываем первые 2 и последние 2 символа
	masked := actualData[:2] + strings.Repeat("*", len(actualData)-4) + actualData[len(actualData)-2:]
	return value[:dataStart] + masked + value[dataEnd:]
}

// MaskFields маскирует список известных чувствительных полей
var sensitiveFields = []string{
	"password", "passwd", "pwd",
	"secret", "api_key", "apikey",
	"token", "access_token", "refresh_token",
	"authorization",
	"credit_card", "card_number",
	"ssn", "social_security",
	"phone", "phone_number",
	"email", "email_address",
	"private_key", "private_key_id",
	"cvv", "cvc", "pin",
}

// MaskSensitiveFields маскирует все известные чувствительные поля в JSON
func (sfm *SensitiveFieldMasker) MaskSensitiveFields(jsonStr string) string {
	start := time.Now()
	result := jsonStr

	for _, field := range sensitiveFields {
		// Маскируем в разных форматах
		result = sfm.MaskJSONField(result, field)
		result = sfm.MaskJSONField(result, strings.ToUpper(field))
		result = sfm.MaskJSONField(result, strings.ToLower(field))
		sfm.metrics.RecordMaskedField(field)
	}

	// Применяем regex паттерны
	result = sfm.Mask(result)

	latency := time.Since(start).Seconds()
	sfm.metrics.RecordMaskingLatency("json_fields", latency)

	return result
}

// MaskInHttpHeader маскирует чувствительные HTTP headers
func (sfm *SensitiveFieldMasker) MaskInHttpHeader(headerName string, headerValue string) string {
	sensitiveHeaders := map[string]bool{
		"Authorization":       true,
		"X-API-Key":           true,
		"X-Auth-Token":        true,
		"Cookie":              true,
		"Set-Cookie":          true,
		"X-CSRF-Token":        true,
		"X-Access-Token":      true,
		"Proxy-Authorization": true,
	}

	if sensitiveHeaders[headerName] {
		return maskValue(headerValue)
	}

	return headerValue
}
