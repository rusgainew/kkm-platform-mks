// Файл company-server/internal/application/company/service_tracing_test.go содержит реализацию пакета company.
package company

import (
	"context"
	"testing"

	"github.com/rusgainew/kkm-project-mks/company-server/internal/domain"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// MockOrganizationRepository - мок репозиторий для организаций
type MockOrganizationRepository struct {
	organizations map[string]*domain.Organization
}

func NewMockOrganizationRepository() *MockOrganizationRepository {
	return &MockOrganizationRepository{
		organizations: make(map[string]*domain.Organization),
	}
}

func (m *MockOrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	m.organizations[org.ID] = org
	return nil
}

func (m *MockOrganizationRepository) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	if org, exists := m.organizations[id]; exists {
		return org, nil
	}
	return nil, domain.ErrOrganizationNotFound
}

func (m *MockOrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	m.organizations[org.ID] = org
	return nil
}

func (m *MockOrganizationRepository) Delete(ctx context.Context, id string) error {
	delete(m.organizations, id)
	return nil
}

func (m *MockOrganizationRepository) List(ctx context.Context, page, perPage int32, ownerID string) ([]*domain.Organization, int32, error) {
	var orgs []*domain.Organization
	for _, org := range m.organizations {
		if ownerID == "" || org.OwnerID == ownerID {
			orgs = append(orgs, org)
		}
	}
	return orgs, int32(len(orgs)), nil
}

func (m *MockOrganizationRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	for _, org := range m.organizations {
		if org.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// MockEmployeeRepository - мок репозиторий для сотрудников
type MockEmployeeRepository struct {
	employees map[string]*domain.Employee
}

func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[string]*domain.Employee),
	}
}

func (m *MockEmployeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	m.employees[emp.ID] = emp
	return nil
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	if emp, exists := m.employees[id]; exists {
		return emp, nil
	}
	return nil, domain.ErrEmployeeNotFound
}

func (m *MockEmployeeRepository) GetByUserAndOrganization(ctx context.Context, userID, organizationID string) (*domain.Employee, error) {
	for _, emp := range m.employees {
		if emp.UserID == userID && emp.OrganizationID == organizationID {
			return emp, nil
		}
	}
	return nil, domain.ErrEmployeeNotFound
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id string) error {
	delete(m.employees, id)
	return nil
}

func (m *MockEmployeeRepository) ListByOrganization(ctx context.Context, organizationID string, page, perPage int32) ([]*domain.Employee, int32, error) {
	var emps []*domain.Employee
	for _, emp := range m.employees {
		if emp.OrganizationID == organizationID {
			emps = append(emps, emp)
		}
	}
	return emps, int32(len(emps)), nil
}

func (m *MockEmployeeRepository) ExistsInOrganization(ctx context.Context, userID, organizationID string) (bool, error) {
	for _, emp := range m.employees {
		if emp.UserID == userID && emp.OrganizationID == organizationID {
			return true, nil
		}
	}
	return false, nil
}

// MockEventPublisher - мок издатель для событий
type MockEventPublisher struct {
	events []string
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		events: make([]string, 0),
	}
}

func (m *MockEventPublisher) Publish(ctx context.Context, eventType string, eventData []byte) error {
	m.events = append(m.events, eventType)
	return nil
}

func (m *MockEventPublisher) Close() error {
	return nil
}

// TestCreateOrganizationWithTracing тестирует создание организации с трассировкой
func TestCreateOrganizationWithTracing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockOrgRepo := NewMockOrganizationRepository()
	mockEmpRepo := NewMockEmployeeRepository()
	mockPublisher := NewMockEventPublisher()

	service := NewService(mockOrgRepo, mockEmpRepo, mockPublisher, logger)

	ctx := context.Background()
	name := "Test Company"
	description := "A test company"
	ownerID := "123e4567-e89b-12d3-a456-426614174000"

	// Выполнение
	org, err := service.CreateOrganization(ctx, name, description, ownerID)

	// Проверка
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if org == nil {
		t.Fatal("Expected organization, got nil")
	}

	if org.Name != name {
		t.Errorf("Expected name %q, got %q", name, org.Name)
	}

	if org.OwnerID != ownerID {
		t.Errorf("Expected owner_id %q, got %q", ownerID, org.OwnerID)
	}

	if len(mockPublisher.events) == 0 {
		t.Fatal("Expected event to be published")
	}
}

// TestGetOrganizationWithTracing тестирует получение организации с трассировкой
func TestGetOrganizationWithTracing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockOrgRepo := NewMockOrganizationRepository()
	mockEmpRepo := NewMockEmployeeRepository()
	mockPublisher := NewMockEventPublisher()

	service := NewService(mockOrgRepo, mockEmpRepo, mockPublisher, logger)

	ctx := context.Background()

	// Создание организации
	org := &domain.Organization{
		ID:      "org1",
		Name:    "Test Company",
		OwnerID: "owner1",
	}
	mockOrgRepo.Create(ctx, org)

	// Выполнение
	retrieved, err := service.GetOrganization(ctx, org.ID)

	// Проверка
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected organization, got nil")
	}

	if retrieved.ID != org.ID {
		t.Errorf("Expected ID %q, got %q", org.ID, retrieved.ID)
	}
}

// TestListOrganizationsWithTracing тестирует получение списка с трассировкой
func TestListOrganizationsWithTracing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockOrgRepo := NewMockOrganizationRepository()
	mockEmpRepo := NewMockEmployeeRepository()
	mockPublisher := NewMockEventPublisher()

	service := NewService(mockOrgRepo, mockEmpRepo, mockPublisher, logger)

	ctx := context.Background()

	// Создание организаций
	org1 := &domain.Organization{
		ID:      "org1",
		Name:    "Company 1",
		OwnerID: "owner1",
	}
	org2 := &domain.Organization{
		ID:      "org2",
		Name:    "Company 2",
		OwnerID: "owner2",
	}

	mockOrgRepo.Create(ctx, org1)
	mockOrgRepo.Create(ctx, org2)

	// Выполнение
	orgs, total, err := service.ListOrganizations(ctx, 1, 10, "")

	// Проверка
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}

	if len(orgs) != 2 {
		t.Errorf("Expected 2 organizations, got %d", len(orgs))
	}
}

// TestAddMemberWithTracing тестирует добавление участника с трассировкой
func TestAddMemberWithTracing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockOrgRepo := NewMockOrganizationRepository()
	mockEmpRepo := NewMockEmployeeRepository()
	mockPublisher := NewMockEventPublisher()

	service := NewService(mockOrgRepo, mockEmpRepo, mockPublisher, logger)

	ctx := context.Background()

	// Создание организации
	org := &domain.Organization{
		ID:      "org1",
		Name:    "Test Company",
		OwnerID: "owner1",
	}
	mockOrgRepo.Create(ctx, org)

	// Выполнение
	emp, err := service.AddMember(ctx, "org1", "user1", "admin")

	// Проверка
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if emp == nil {
		t.Fatal("Expected employee, got nil")
	}

	if emp.UserID != "user1" {
		t.Errorf("Expected user_id user1, got %s", emp.UserID)
	}

	if emp.Role != "admin" {
		t.Errorf("Expected role admin, got %s", emp.Role)
	}
}

// TestTracerInitialization проверяет инициализацию tracer
func TestTracerInitialization(t *testing.T) {
	// Проверка что tracer можно получить
	tracer := otel.Tracer("test-service")

	if tracer == nil {
		t.Fatal("Expected tracer to be initialized")
	}

	// Проверка создания span
	ctx := context.Background()
	_, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	if span == nil {
		t.Fatal("Expected span to be created")
	}
}

// TestServiceWithMocks проверяет интеграцию сервиса с мок репозиториями
func TestServiceWithMocks(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockOrgRepo := NewMockOrganizationRepository()
	mockEmpRepo := NewMockEmployeeRepository()
	mockPublisher := NewMockEventPublisher()

	service := NewService(mockOrgRepo, mockEmpRepo, mockPublisher, logger)

	if service == nil {
		t.Fatal("Expected service to be initialized")
	}

	// Проверка что методы доступны
	if service.orgRepo != mockOrgRepo {
		t.Error("Service should use provided organization repository")
	}

	if service.empRepo != mockEmpRepo {
		t.Error("Service should use provided employee repository")
	}

	if service.eventPublisher != mockPublisher {
		t.Error("Service should use provided event publisher")
	}
}
