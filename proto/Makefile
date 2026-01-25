# Конфигурация сборки Proto-файлов
# Использование: make proto

# Базовые директории
PROTO_DIR := .
ROOT_DIR := $(abspath ..)
PROTO_LIB_DIR := $(ROOT_DIR)/proto-lib
PROTO_MODULE := github.com/rusgainew/kkm-project-mks

# База команд для новых (организованных) proto — генерируем файлы
# с путями, относительными к исходным файлам
PROTOC_BASE := \
	protoc --proto_path=$(PROTO_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

# Для legacy-файлов уважаем go_package пути, чтобы не смешивать пакеты
PROTOC_BASE_LEGACY := \
	protoc --proto_path=$(PROTO_DIR) \
		--go_opt=module=$(PROTO_MODULE) --go-grpc_opt=module=$(PROTO_MODULE)

# Legacy proto-файлы (расположены в каталоге legacy/)
PROTO_LEGACY_AUTH := $(wildcard legacy/auth*.proto)
PROTO_LEGACY_COMPANY := $(wildcard legacy/company*.proto)
PROTO_LEGACY_DOCUMENT := $(wildcard legacy/document*.proto)
PROTO_LEGACY_COMMON := $(wildcard legacy/common.proto)

# Новые организованные proto-файлы
PROTO_DICTIONARIES := $(wildcard dictionaries/*.proto)
PROTO_ENTITIES := $(wildcard entities/*.proto)
PROTO_API := $(wildcard api/*.proto)
PROTO_ERRORS := $(wildcard errors/*.proto)

.PHONY: proto
# Полная сборка: очистка, создание директорий, сборка legacy и новых proto
proto: proto-clean proto-dirs proto-legacy proto-new

.PHONY: proto-legacy
proto-legacy:
	@echo "Сборка legacy proto файлов..."
	# ВАЖНО: для legacy уважаем go_package (в путях есть 'proto-lib/...'), поэтому
	# выводим в корень репозитория, а не в proto-lib, чтобы не получить 'proto-lib/proto-lib'
	$(PROTOC_BASE_LEGACY) --go_out=$(ROOT_DIR) --go-grpc_out=$(ROOT_DIR) legacy/common.proto
	$(PROTOC_BASE_LEGACY) --go_out=$(ROOT_DIR) --go-grpc_out=$(ROOT_DIR) legacy/auth.proto legacy/auth_service.proto
	$(PROTOC_BASE_LEGACY) --go_out=$(ROOT_DIR) --go-grpc_out=$(ROOT_DIR) legacy/company.proto legacy/company_service.proto
	$(PROTOC_BASE_LEGACY) --go_out=$(ROOT_DIR) --go-grpc_out=$(ROOT_DIR) legacy/document.proto legacy/document_service.proto

.PHONY: proto-new
proto-new:
	@echo "Сборка организованных proto файлов..."
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_DICTIONARIES)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_ENTITIES)
	# Включает и устаревший api/requests.proto для обратной совместимости
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_API)
	$(PROTOC_BASE) --go_out=$(PROTO_LIB_DIR) --go-grpc_out=$(PROTO_LIB_DIR) $(PROTO_ERRORS)

.PHONY: proto-dirs
# Создание директорий для сгенерированных файлов
proto-dirs:
	mkdir -p $(PROTO_LIB_DIR)/common $(PROTO_LIB_DIR)/auth $(PROTO_LIB_DIR)/company $(PROTO_LIB_DIR)/document
	mkdir -p $(PROTO_LIB_DIR)/dictionaries $(PROTO_LIB_DIR)/entities $(PROTO_LIB_DIR)/api $(PROTO_LIB_DIR)/errors

.PHONY: proto-clean
# Очистка ранее сгенерированных .pb.go файлов
proto-clean:
	@[ -d $(PROTO_LIB_DIR) ] && find $(PROTO_LIB_DIR) -name "*.pb.go" -delete || true
	find $(PROTO_DIR) -maxdepth 1 -name "*.pb.go" -delete
	# Удаляем ошибочно созданную вложенную структуру 'proto-lib/proto-lib'
	@[ -d $(PROTO_LIB_DIR)/proto-lib ] && rm -rf $(PROTO_LIB_DIR)/proto-lib || true

.PHONY: proto-install-tools
# Установка компилятора protobuf и Go-плагинов
proto-install-tools:
	@echo "Устанавливаю protobuf компилятор и Go-плагины..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Инструменты для proto успешно установлены!"

.PHONY: proto-check
# Проверка версий установленных инструментов
proto-check:
	@echo "Проверяю версии инструментов для proto..."
	protoc --version
	protoc-gen-go --version
	protoc-gen-go-grpc --version

.PHONY: help
help:
	@echo "Управление Proto-файлами"
	@echo "========================"
	@echo ""
	@echo "Доступные цели:"
	@echo "  make proto               - Скомпилировать все .proto файлы"
	@echo "  make proto-clean         - Удалить все сгенерированные .pb.go файлы"
	@echo "  make proto-install-tools - Установить компилятор protobuf и Go-плагины"
	@echo "  make proto-check         - Проверить версии установленных инструментов"
	@echo "  make help                - Показать это сообщение"
	@echo ""
	@echo "Расположение proto-файлов: proto/ (текущий каталог)"
	@echo "Сгенерированный код: в каталоге proto-lib/* по соответствующим подпапкам"
