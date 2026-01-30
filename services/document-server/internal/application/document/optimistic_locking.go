// Файл document-server/internal/application/document/optimistic_locking.go содержит реализацию пакета document.
package document

import (
	"context"
	"errors"

	"github.com/rusgainew/kkm-project-mks/document-server/internal/domain"
	"go.uber.org/zap"
)

// UpdateDocumentWithVersion обновляет документ проверяя версию (optimistic locking)
// Возвращает ErrVersionConflict если версия документа в БД отличается от переданной
func (s *Service) UpdateDocumentWithVersion(ctx context.Context, documentID, title, content string, expectedVersion int) error {
	if documentID == "" {
		return domain.ErrInvalidInput
	}

	// Получаем текущую версию документа
	currentDoc, err := s.repo.GetWithVersion(ctx, documentID)
	if err != nil {
		s.logger.Error("Failed to get document for version check",
			zap.Error(err),
			zap.String("document_id", documentID),
		)
		return err
	}

	if currentDoc == nil {
		return domain.ErrDocumentNotFound
	}

	// Проверяем версию
	actualVersion, ok := convertVersion(currentDoc["version"])
	if !ok {
		s.logger.Error("Invalid version type in document",
			zap.String("document_id", documentID),
		)
		return domain.ErrInvalidInput
	}

	if actualVersion != expectedVersion {
		s.logger.Warn("Version conflict detected",
			zap.String("document_id", documentID),
			zap.Int("expected_version", expectedVersion),
			zap.Int("actual_version", actualVersion),
		)
		return domain.ErrVersionConflict
	}

	// Выполняем обновление с проверкой версии
	// Репозиторий должен автоматически увеличить версию на 1
	if err := s.repo.UpdateWithVersion(ctx, documentID, title, content, expectedVersion); err != nil {
		if errors.Is(err, domain.ErrVersionConflict) {
			s.logger.Warn("Update failed due to concurrent modification",
				zap.String("document_id", documentID),
			)
			return domain.ErrConcurrentModification
		}
		s.logger.Error("Failed to update document with version check",
			zap.Error(err),
			zap.String("document_id", documentID),
		)
		return err
	}

	s.logger.Info("Document updated with version check",
		zap.String("document_id", documentID),
		zap.Int("version", expectedVersion+1),
	)

	return nil
}

// convertVersion нормализует версию в int, поддерживая типы int32/int64 из БД
func convertVersion(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetDocumentWithVersion получает документ с информацией о версии
func (s *Service) GetDocumentWithVersion(ctx context.Context, documentID string) (map[string]interface{}, error) {
	if documentID == "" {
		return nil, domain.ErrInvalidInput
	}

	doc, err := s.repo.GetWithVersion(ctx, documentID)
	if err != nil {
		s.logger.Error("Failed to get document with version",
			zap.Error(err),
			zap.String("document_id", documentID),
		)
		return nil, err
	}

	if doc == nil {
		return nil, domain.ErrDocumentNotFound
	}

	return doc, nil
}
