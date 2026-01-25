# 📖 Документация КММ Project MKS - Указатель

**Версия:** 1.0 Consolidated + Nginx Reverse Proxy  
**Дата обновления:** 10 января 2026  
**Статус:** ✅ Production Ready с Nginx

---

## 🗺️ Навигация по документации

### 🚀 Начните отсюда

| Документ                                                 | Для кого           | Время чтения |
| -------------------------------------------------------- | ------------------ | ------------ |
| **[README.md](README.md)** ⭐ НОВОЕ                      | Все разработчики   | 15 минут     |
| **[QUICKSTART.md](QUICKSTART.md)** ⭐ НОВОЕ              | Новые разработчики | 10 минут     |
| **[DOCKER_COMPOSE_README.md](DOCKER_COMPOSE_README.md)** | DevOps / Backend   | 10 минут     |

### 📡 API и Интеграция

| Документ                                                       | Содержание                                          |
| -------------------------------------------------------------- | --------------------------------------------------- |
| **[API_DOCUMENTATION.md](API_DOCUMENTATION.md)**               | Proto definitions, примеры вызовов, все RPC методы  |
| **[INTEGRATION_TESTS_README.md](INTEGRATION_TESTS_README.md)** | Запуск интеграционных тестов (88 тестов), seed data |

### 🔍 Observability

| Документ                                                                               | Содержание                                       |
| -------------------------------------------------------------------------------------- | ------------------------------------------------ |
| **[OPENTELEMETRY_TRACING_IMPLEMENTATION.md](OPENTELEMETRY_TRACING_IMPLEMENTATION.md)** | Jaeger tracing, span attributes, troubleshooting |

### 🔐 Безопасность

| Документ                                                                              | Содержание                                                  |
| ------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| **[SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)**                          | JWT, encryption, OWASP compliance, password management      |
| **[JWT_SECRET_ROTATION_IMPLEMENTATION.md](JWT_SECRET_ROTATION_IMPLEMENTATION.md)**    | Автоматическая смена JWT ключей, rotation strategy          |
| **[services/nginx-proxy/README.md](./services/nginx-proxy/README.md)** ⭐ НОВОЕ       | Nginx reverse proxy, rate limiting, SSL/TLS config          |
| **[services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md)** ⭐ НОВОЕ | HTTPS настройка, Let's Encrypt, самоподписанные сертифікаты |

### 📊 Производительность и Тестирование

| Документ                                           | Содержание                                               |
| -------------------------------------------------- | -------------------------------------------------------- |
| **[LOAD_TESTING_GUIDE.md](LOAD_TESTING_GUIDE.md)** | Load testing, stress testing, baseline metrics, сценарии |

### 🚀 Production

| Документ                                                                                           | Содержание                                         |
| -------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| **[PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)**                                           | K8s deployment, health checks, scaling, monitoring |
| **[QUERY_SERVERS_PRIORITIES_COMPLETION_REPORT.md](QUERY_SERVERS_PRIORITIES_COMPLETION_REPORT.md)** | Итоговый отчёт по всем 5 приоритетам               |

---

## 📚 Быстрые ссылки по задачам

### Я хочу...

#### 🏃 **Быстро запустить проект локально**

1. Прочитайте: [QUICKSTART.md](QUICKSTART.md)
2. Выполните: `docker-compose up -d && ./seed-test-data.sh`
3. Проверьте: `http://localhost:16686` (Jaeger)

#### 🔌 **Подключиться к API**

1. Прочитайте: [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
2. Примеры: `grpcurl -plaintext localhost:50051 list`
3. Тесты: [INTEGRATION_TESTS_README.md](INTEGRATION_TESTS_README.md)

#### 📊 **Мониторить систему**

1. Jaeger UI: http://localhost:16686
2. Prometheus: http://localhost:9090
3. Grafana: http://localhost:3000
4. Читайте: [OPENTELEMETRY_TRACING_IMPLEMENTATION.md](OPENTELEMETRY_TRACING_IMPLEMENTATION.md)

#### 🔒 **Обеспечить безопасность**

1. Проверьте: [SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)
2. JWT ротация: [JWT_SECRET_ROTATION_IMPLEMENTATION.md](JWT_SECRET_ROTATION_IMPLEMENTATION.md)
3. Тесты: security/README.md

#### ⚡ **Тестировать производительность**

1. Load tests: `cd load-tests && ./run-load-test.sh`
2. Stress tests: `./stress-test.sh`
3. Читайте: [LOAD_TESTING_GUIDE.md](LOAD_TESTING_GUIDE.md)

#### 🚀 **Развернуть в production**

1. Читайте: [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)
2. **Настройте Nginx**: [services/nginx-proxy/README.md](./services/nginx-proxy/README.md)
3. **Включите HTTPS**: [services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md)
4. K8s manifests: `k8s/` папка (если есть)
5. Jaeger collector: Setup external Jaeger

#### 🔐 **Настроить Nginx Reverse Proxy** ⭐ НОВОЕ

1. **Основное**: [services/nginx-proxy/README.md](./services/nginx-proxy/README.md)
   - Rate limiting, гzip, upstream pooling
   - Health checks, метрики, логирование
2. **SSL/TLS**: [services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md)
   - Let's Encrypt сертификаты
   - Самоподписанные сертификаты для dev
   - Automatic renewal

#### 🧪 **Написать тесты**

1. Integration tests: [INTEGRATION_TESTS_README.md](INTEGRATION_TESTS_README.md)
2. Примеры: `services/*/tests/integration/`
3. Run: `go test -v ./integration/...`

---

## 📂 Структура документации

```
/
├── README.md                                    # 📍 ГЛАВНАЯ документация (обновлена с Nginx)
├── QUICKSTART.md                                # Быстрый старт (обновлен с Nginx)
├── INDEX.md                                     # Этот файл - указатель всей документации
├── INDEX.md                                     # Этот файл (навигация)
│
├── API_DOCUMENTATION.md                         # Proto, методы, примеры
├── INTEGRATION_TESTS_README.md                  # Тесты + seed data
│
├── DOCKER_COMPOSE_README.md                     # Docker Compose setup
├── OPENTELEMETRY_TRACING_IMPLEMENTATION.md      # Jaeger tracing
│
├── SECURITY_BEST_PRACTICES.md                   # Security guidelines
├── JWT_SECRET_ROTATION_IMPLEMENTATION.md        # JWT rotation
│
├── PRODUCTION_DEPLOYMENT.md                     # K8s, production
├── LOAD_TESTING_GUIDE.md                        # Performance testing
│
├── QUERY_SERVERS_PRIORITIES_COMPLETION_REPORT.md # Итоговый отчёт
│
├── security/                                    # Security tools
│   ├── README.md
│   ├── owasp-top10-checklist.md
│   ├── compliance-standards.md
│   └── penetration-testing-guide.md
│
├── monitoring/                                  # Prometheus, Grafana
│   ├── prometheus.yml
│   └── grafana-dashboard-*.json
│
├── load-tests/                                  # Load testing
│   ├── README.md
│   └── *.sh
│
└── proto/                                       # Proto definitions
    └── README.md
```

---

## 🎯 По уровню детализации

### Для менеджеров (5-10 минут)

→ [README.md](README.md) - Обзор, архитектура, статистика

### Для новых разработчиков (30-60 минут)

1. [QUICKSTART.md](QUICKSTART.md)
2. [DOCKER_COMPOSE_README.md](DOCKER_COMPOSE_README.md)
3. [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - первая часть

### Для backend разработчиков (2-3 часа)

1. [README.md](README.md) - полностью
2. [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - полностью
3. [INTEGRATION_TESTS_README.md](INTEGRATION_TESTS_README.md)
4. [SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)

### Для DevOps/SRE (3-4 часа)

1. [DOCKER_COMPOSE_README.md](DOCKER_COMPOSE_README.md)
2. [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)
3. [OPENTELEMETRY_TRACING_IMPLEMENTATION.md](OPENTELEMETRY_TRACING_IMPLEMENTATION.md)
4. [LOAD_TESTING_GUIDE.md](LOAD_TESTING_GUIDE.md)

### Для security специалистов (2-3 часа)

1. [SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)
2. [JWT_SECRET_ROTATION_IMPLEMENTATION.md](JWT_SECRET_ROTATION_IMPLEMENTATION.md)
3. security/README.md и guides

---

## 🔄 Обновление документации

Документация обновляется при:

- Изменении API (обновите API_DOCUMENTATION.md)
- Новых фичах (обновите README.md)
- Изменении инфраструктуры (обновите DOCKER_COMPOSE_README.md)
- Изменении безопасности (обновите SECURITY_BEST_PRACTICES.md)

**Правило:** Каждый PR должен включать обновление соответствующей документации.

---

## 📊 Статистика документации

```
Документ                                    Строк   Размер
────────────────────────────────────────────────────────────
README.md                                   450     28KB
API_DOCUMENTATION.md                        650     38KB
INTEGRATION_TESTS_README.md                 400     25KB
SECURITY_BEST_PRACTICES.md                  550     32KB
DOCKER_COMPOSE_README.md                    300     18KB
JWT_SECRET_ROTATION_IMPLEMENTATION.md       650     38KB
OPENTELEMETRY_TRACING_IMPLEMENTATION.md     550     32KB
PRODUCTION_DEPLOYMENT.md                    550     32KB
LOAD_TESTING_GUIDE.md                       650     38KB
QUICKSTART.md                               400     24KB
QUERY_SERVERS_PRIORITIES_COMPLETION_REPORT  900     52KB
────────────────────────────────────────────────────────────
TOTAL:                                      6100    357KB
```

---

## ✅ Checklist для новых разработчиков

- [ ] Прочитайте [QUICKSTART.md](QUICKSTART.md)
- [ ] Запустите `docker-compose up -d`
- [ ] Выполните `./seed-test-data.sh`
- [ ] Запустите интеграционные тесты
- [ ] Посетите Jaeger UI (http://localhost:16686)
- [ ] Прочитайте [API_DOCUMENTATION.md](API_DOCUMENTATION.md)
- [ ] Изучите [SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)
- [ ] Выполните первый gRPC запрос
- [ ] Запустите load tests

---

## 🆘 Где найти ответ?

| Вопрос                       | Документ                                                                           |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| Как запустить?               | [QUICKSTART.md](QUICKSTART.md)                                                     |
| Какие есть методы API?       | [API_DOCUMENTATION.md](API_DOCUMENTATION.md)                                       |
| Как писать тесты?            | [INTEGRATION_TESTS_README.md](INTEGRATION_TESTS_README.md)                         |
| Как проверить безопасность?  | [SECURITY_BEST_PRACTICES.md](SECURITY_BEST_PRACTICES.md)                           |
| Как тестировать performance? | [LOAD_TESTING_GUIDE.md](LOAD_TESTING_GUIDE.md)                                     |
| Как развернуть в production? | [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)                               |
| Как настроить трассировку?   | [OPENTELEMETRY_TRACING_IMPLEMENTATION.md](OPENTELEMETRY_TRACING_IMPLEMENTATION.md) |
| Как настроить Docker?        | [DOCKER_COMPOSE_README.md](DOCKER_COMPOSE_README.md)                               |
| Как работает JWT?            | [JWT_SECRET_ROTATION_IMPLEMENTATION.md](JWT_SECRET_ROTATION_IMPLEMENTATION.md)     |

---

## 📞 Поддержка

- **Issues:** GitHub Issues (используйте labels: `question`, `bug`, `enhancement`)
- **Docs Issues:** GitHub Issues с label `documentation`
- **Security:** security@example.com
- **General:** support@example.com

---

**Последнее обновление:** 8 января 2026  
**Ответственный:** Backend Team  
**Status:** ✅ Актуально
