# 📚 Документация проекта KKM Project MKS

Полная коллекция технической документации микросервисной платформы.

## 🗂️ Основные документы

### 🏗️ Архитектура

| Документ                                                    | Описание                                        | Формат             |
| ----------------------------------------------------------- | ----------------------------------------------- | ------------------ |
| **[SERVICES_ARCHITECTURE.md](../SERVICES_ARCHITECTURE.md)** | Полная архитектурная документация с диаграммами | Markdown + Mermaid |
| **[diagrams/](diagrams/)**                                  | PNG версии всех архитектурных схем              | PNG Images         |
| **[README.md](../README.md)**                               | Основная документация проекта                   | Markdown           |
| **[QUICKSTART.md](../QUICKSTART.md)**                       | Production deployment guide                     | Markdown           |

### 📊 Визуальные схемы

#### 1. Общая архитектура системы (3984x1606px)

![Overall Architecture](diagrams/overall-architecture.png)

[Открыть в полном размере](diagrams/overall-architecture.png) | [Исходник Mermaid](diagrams/overall-architecture.mmd)

---

#### 2. CQRS Pattern - Sequence диаграмма (1539x1117px)

![CQRS Pattern](diagrams/cqrs-pattern.png)

[Открыть в полном размере](diagrams/cqrs-pattern.png) | [Исходник Mermaid](diagrams/cqrs-pattern.mmd)

---

#### 3. RabbitMQ Event Flow (1457x452px)

![RabbitMQ Events](diagrams/rabbitmq-events.png)

[Открыть в полном размере](diagrams/rabbitmq-events.png) | [Исходник Mermaid](diagrams/rabbitmq-events.mmd)

---

#### 4. Аутентификация и авторизация (1321x1402px)

![Auth Flow](diagrams/auth-flow.png)

[Открыть в полном размере](diagrams/auth-flow.png) | [Исходник Mermaid](diagrams/auth-flow.mmd)

---

## 📖 Детальная документация

### Сервисы

| Сервис                    | Документация                                                                                | Статус |
| ------------------------- | ------------------------------------------------------------------------------------------- | ------ |
| API Gateway               | [services/api-gateway/README.md](../services/api-gateway/README.md)                         | ✅     |
| Analytics Server (Python) | [services/analytics-server-python/README.md](../services/analytics-server-python/README.md) | ✅     |
| User Server               | [services/user-server/](../services/user-server/)                                           | ✅     |
| Company Server            | [services/company-server/README.md](../services/company-server/README.md)                   | ✅     |
| Catalog Server            | [services/catalog-server/README.md](../services/catalog-server/README.md)                   | ✅     |
| Invoice Server            | [services/invoice-server/README.md](../services/invoice-server/README.md)                   | ✅     |
| Bank Account Server       | [services/bank-account-server/README.md](../services/bank-account-server/README.md)         | ✅     |
| Foreign Company Server    | [services/foreign-company-server/](../services/foreign-company-server/)                     | ✅     |
| Document Server           | [services/document-server/README.md](../services/document-server/README.md)                 | ✅     |

### gRPC Контракты

| Документ                                                                    | Описание                          |
| --------------------------------------------------------------------------- | --------------------------------- |
| [proto/README.md](../proto/README.md)                                       | gRPC Proto файлы и генерация кода |
| [proto/SERVICES_AND_DEPENDENCIES.md](../proto/SERVICES_AND_DEPENDENCIES.md) | Граф зависимостей сервисов        |
| [proto/PROTO_STRUCTURE.md](../proto/PROTO_STRUCTURE.md)                     | Структура proto файлов            |

### Мониторинг и наблюдаемость

| Документ                                                                     | Описание                                   |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [infrastructure/PROMETHEUS_METRICS.md](infrastructure/PROMETHEUS_METRICS.md) | Метрики и их использование                 |
| [infrastructure/HEALTH_CHECKS.md](infrastructure/HEALTH_CHECKS.md)           | Health check endpoints всех сервисов       |
| [monitoring/](../monitoring/)                                                | Grafana дашборды и Prometheus конфигурация |

### Docker и развертывание

| Документ                                              | Описание                          |
| ----------------------------------------------------- | --------------------------------- |
| [docker-compose.yml](../docker-compose.yml)           | Локальная разработка (22 сервиса) |
| [docker-compose.prod.yml](../docker-compose.prod.yml) | Production конфигурация           |
| [DOCKER_COMPOSE_GUIDE.md](DOCKER_COMPOSE_GUIDE.md)    | Гайд по Docker Compose            |
| [DOCKER_CHEATSHEET.md](DOCKER_CHEATSHEET.md)          | Шпаргалка Docker команд           |

### Безопасность и инфраструктура

| Документ                                                                                                   | Описание                        |
| ---------------------------------------------------------------------------------------------------------- | ------------------------------- |
| [infrastructure/SECURITY_CONFIGURATION.md](infrastructure/SECURITY_CONFIGURATION.md)                       | Конфигурация безопасности       |
| [infrastructure/NGINX_PROXY_CONFIGURATION.md](infrastructure/NGINX_PROXY_CONFIGURATION.md)                 | Nginx настройки и rate limiting |
| [infrastructure/DATABASE_PER_SERVICE_VERIFICATION.md](infrastructure/DATABASE_PER_SERVICE_VERIFICATION.md) | Проверка изоляции БД            |
| [LOGIN_INSTRUCTIONS.md](LOGIN_INSTRUCTIONS.md)                                                             | Инструкции по входу в систему   |

### Frontend

| Документ                                                                  | Описание                    |
| ------------------------------------------------------------------------- | --------------------------- |
| [kkm-platform/README.md](../kkm-platform/README.md)                       | Next.js фронтенд приложение |
| [kkm-platform/PROJECT_STRUCTURE.md](../kkm-platform/PROJECT_STRUCTURE.md) | Структура frontend проекта  |

## 🔍 Быстрая навигация по темам

### Для разработчиков

- **Начало работы**: [README.md](../README.md) → [QUICKSTART.md](../QUICKSTART.md)
- **Архитектура**: [SERVICES_ARCHITECTURE.md](../SERVICES_ARCHITECTURE.md)
- **gRPC разработка**: [proto/README.md](../proto/README.md)
- **API Gateway**: [services/api-gateway/README.md](../services/api-gateway/README.md)

### Для DevOps

- **Docker**: [docker-compose.yml](../docker-compose.yml) → [DOCKER_COMPOSE_GUIDE.md](DOCKER_COMPOSE_GUIDE.md)
- **Мониторинг**: [PROMETHEUS_METRICS.md](infrastructure/PROMETHEUS_METRICS.md)
- **Health Checks**: [HEALTH_CHECKS.md](infrastructure/HEALTH_CHECKS.md)
- **Безопасность**: [SECURITY_CONFIGURATION.md](infrastructure/SECURITY_CONFIGURATION.md)

### Для аналитиков/архитекторов

- **Визуальные схемы**: [diagrams/](diagrams/)
- **Архитектурные решения**: [SERVICES_ARCHITECTURE.md](../SERVICES_ARCHITECTURE.md)
- **База данных**: [DATABASE_PER_SERVICE_VERIFICATION.md](infrastructure/DATABASE_PER_SERVICE_VERIFICATION.md)

## 📊 Статистика проекта

- **Микросервисов**: 14 (7 Command + 6 Query + 1 Analytics)
- **Баз данных**: 7 PostgreSQL instances (DB per service)
- **Языки**: Go, Python, TypeScript/Next.js
- **Протоколы**: gRPC, HTTP/REST, AMQP (RabbitMQ)
- **Паттерны**: CQRS, Event-Driven, API Gateway, Database per Service
- **Observability**: Prometheus, Grafana, Jaeger

## 🚀 Ключевые технологии

- **Backend**: Go 1.23+, Python 3.12
- **Frontend**: Next.js 15, React 19, TypeScript
- **Базы данных**: PostgreSQL 15, Redis 7
- **Message Broker**: RabbitMQ 3
- **API**: gRPC, REST (HTTP/JSON)
- **Мониторинг**: Prometheus, Grafana, Jaeger
- **Контейнеризация**: Docker, Docker Compose
- **Reverse Proxy**: Nginx

## 📝 Обновление документации

### Добавление новых диаграмм

1. Создайте файл `.mmd` в [diagrams/](diagrams/)
2. Используйте [Mermaid Live Editor](https://mermaid.live/) для предпросмотра
3. Генерируйте PNG:
   ```bash
   cd docs/diagrams
   mmdc -i your-diagram.mmd -o your-diagram.png -b transparent -w 2000
   ```
4. Обновите [diagrams/README.md](diagrams/README.md)

### Стандарты документации

- Все Markdown файлы должны использовать UTF-8 encoding
- Используйте относительные пути для ссылок
- Добавляйте эмодзи для визуальной навигации
- Включайте примеры кода где возможно
- Обновляйте дату в конце документа

---

**Создано:** 30 января 2026  
**Поддерживается:** KKM Project Team
