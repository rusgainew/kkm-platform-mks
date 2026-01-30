// Файл api-gateway/internal/domain/models/models.go содержит реализацию пакета models.
package models

// Company модель компании (Organization)
type Company struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	MemberCount int    `json:"member_count"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	Status      string `json:"status"`
}

// Invoice модель счета-фактуры
type Invoice struct {
	ID            string  `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	DeliveryDate  string  `json:"delivery_date"`
	TotalAmount   float64 `json:"total_amount"`
	IsResident    bool    `json:"is_resident"`
	Note          string  `json:"note"`
	Status        string  `json:"status"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

// Document модель документа
type Document struct {
	ID              string          `json:"id"`
	OrganizationID  string          `json:"organization_id"`
	Title           string          `json:"title"`
	Content         string          `json:"content"`
	Status          string          `json:"status"`
	CreatedBy       string          `json:"created_by"`
	AssignedTo      string          `json:"assigned_to,omitempty"`
	CreatedAt       int64           `json:"created_at"`
	UpdatedAt       int64           `json:"updated_at"`
	StatusChangedAt int64           `json:"status_changed_at"`
	Version         int32           `json:"version"`
	Entries         []DocumentEntry `json:"entries,omitempty"`
}

// DocumentEntry модель записи документа
type DocumentEntry struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// Catalog модель элемента каталога
type Catalog struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Number      string  `json:"number"`      // Номер по каталогу
	Description string  `json:"description"` // Описание (опционально)
	TnvedCode   string  `json:"tnved_code"`  // Код ТНВЭД
	Category    string  `json:"category"`    // Категория (deprecated, использовать TnvedCode)
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Unit        string  `json:"unit"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}

// BankAccount модель банковского счета
type BankAccount struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
	OwnerID       string `json:"owner_id"`
	IsActive      bool   `json:"is_active"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

// User модель пользователя
type User struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Status    string `json:"status,omitempty"`
}

// AuthResponse ответ аутентификации
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshExp   int64  `json:"refresh_expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	User         User   `json:"user"`
	Timestamp    int64  `json:"timestamp"`
}

// TokenResponse ответ только с токенами (refresh flow)
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in,omitempty"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	Timestamp        int64  `json:"timestamp,omitempty"`
}

// PageInfo описывает пагинацию
type PageInfo struct {
	Page       int32 `json:"page"`
	Size       int32 `json:"size"`
	TotalCount int32 `json:"total_count"`
}

// ListUsersResponse HTTP обертка для списка пользователей
type ListUsersResponse struct {
	Users []User   `json:"users"`
	Page  PageInfo `json:"page_info"`
}

// DocumentListResponse HTTP обертка для списка документов
type DocumentListResponse struct {
	Documents []Document `json:"documents"`
	Page      PageInfo   `json:"page_info"`
}

// ForeignCompany модель иностранной компании
type ForeignCompany struct {
	ID          int64  `json:"id"`
	PIN         string `json:"pin"`
	FullName    string `json:"full_name"`
	CountryCode string `json:"country_code"`
	Address     string `json:"address,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// APIResponse общая структура ответа API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

// APIError структура для ошибок API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// MetaData метаданные для пагинации
type MetaData struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// PaginationRequest запрос с пагинацией
type PaginationRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// Employee модель сотрудника организации
type Employee struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	Position       string `json:"position,omitempty"`
	Department     string `json:"department,omitempty"`
	JoinedAt       int64  `json:"joined_at"`
	LastActiveAt   int64  `json:"last_active_at,omitempty"`
}

// AddMemberRequest запрос на добавление члена в организацию
type AddMemberRequest struct {
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id" binding:"required"`
	Role           string `json:"role" binding:"required"`
}

// InvoiceDetail модель детализации счета
type InvoiceDetail struct {
	ID           string  `json:"id"`
	InvoiceUUID  string  `json:"invoice_uuid"`
	CatalogCode  string  `json:"catalog_code"`
	CatalogName  string  `json:"catalog_name"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
	VATRate      float64 `json:"vat_rate"`
	VATAmount    float64 `json:"vat_amount"`
	ExciseRate   float64 `json:"excise_rate,omitempty"`
	ExciseAmount float64 `json:"excise_amount,omitempty"`
	TurnoverSize float64 `json:"turnover_size,omitempty"`
	CreatedAt    int64   `json:"created_at"`
}
