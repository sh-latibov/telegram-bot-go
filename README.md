# 🤖 Telegram Bot - Go

Мощный Telegram бот на Go с поддержкой погоды и управлением базой данных PostgreSQL.

## 📋 Содержание

- [Описание](#описание)
- [Особенности](#особенности)
- [Требования](#требования)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [Структура проекта](#структура-проекта)
- [API и клиенты](#api-и-клиенты)
- [Использование](#использование)
- [Миграции БД](#миграции-бд)
- [Технологический стек](#технологический-стек)
- [Разработка](#разработка)

## 📝 Описание

Telegram бот, разработанный на Go, предоставляет функциональность получения информации о погоде через OpenWeather API и сохраняет историю запросов пользователей в PostgreSQL.

## ✨ Особенности

- ✅ Интеграция с Telegram Bot API
- ✅ Получение данных о погоде через OpenWeather
- ✅ Сохранение данных пользователей в PostgreSQL
- ✅ Система команд бота
- ✅ Поддержка переменных окружения (.env)
- ✅ Миграции базы данных
- ✅ Логирование операций

## 📦 Требования

- **Go:** версия 1.25.4 или выше
- **PostgreSQL:** 12 или выше
- **Docker** (опционально, для контейнеризации)

## 🚀 Установка

### 1. Клонирование репозитория

```bash
git clone https://github.com/sh-latibov/telegram-bot-go.git
cd telegram-bot-go
```

### 2. Установка зависимостей

```bash
go mod download
go mod tidy
```

### 3. Подготовка базы данных

```bash
# Создайте БД PostgreSQL
createdb telegram_bot_db

# Запустите миграции (если есть инструмент миграции)
# migrate -path ./migrations -database "postgres://user:password@localhost:5432/telegram_bot_db" up
```

## ⚙️ Конфигурация

### Переменные окружения

Создайте файл `.env` в корне проекта:

```env
# Telegram Bot
BOT_TOKEN=your_telegram_bot_token_here

# PostgreSQL
POSTGRES_USER=your_database_user
POSTGRES_PASSWORD=your_database_password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=telegram_bot_db

# OpenWeather API
WEATHER_KEY=your_openweather_api_key_here
```

### Получение необходимых ключей

1. **BOT_TOKEN:** Создайте бот через [@BotFather](https://t.me/botfather) в Telegram
2. **WEATHER_KEY:** Получите API ключ на [OpenWeather](https://openweathermap.org/api)
3. **PostgreSQL:** Установите и настройте локально или используйте облачный сервис

## 📁 Структура проекта

```
telegram-bot-go/
├── main.go                          # Входная точка приложения
├── go.mod                           # Модули Go
├── go.sum                           # Хеши зависимостей
├── .env                             # Переменные окружения (не коммитить!)
├── .gitignore                       # Исключения Git
│
├── clients/                         # Клиенты для внешних сервисов
│   └── openweather/
│       ├── openweather.go          # Основная логика клиента
│       └── models.go               # Модели данных OpenWeather
│
├── handler/                         # Обработчики событий бота
│   └── handler.go                  # Основная логика обработчика
│
├── models/                          # Модели данных приложения
│   └── models.go                   # Определение структур
│
├── repo/                            # Слой доступа к данным (Repository Pattern)
│   └── repo.go                     # SQL запросы и работа с БД
│
└── migrations/                      # Миграции базы данных
    ├── 20260114113629_init.sql          # Начальная схема
    ├── 20260115113936_add_city_column.sql    # Добавление колонки города
    └── 20260120070012_remove_column_user.sql # Удаление колонки пользователя
```

## 🔌 API и клиенты

### OpenWeather Client

Клиент для получения информации о погоде:

```go
owClient := openweather.New(weatherKey)
weather, err := owClient.GetWeather(city)
```

### Repository Pattern

Работа с базой данных через Repository:

```go
userRepo := repo.New(conn)
user, err := userRepo.GetUser(userID)
```

## 💬 Использование

### Запуск бота

```bash
go run main.go
```

### Команды бота

Бот поддерживает различные команды (определены в `handler/handler.go`):

- `/start` - Начать использование бота
- `/weather <город>` - Получить информацию о погоде
- `/help` - Справка по командам

## 🗄️ Миграции БД

Миграции применяются в следующем порядке:

| Файл | Описание | Дата |
|------|---------|------|
| `20260114113629_init.sql` | Создание основных таблиц | 14.01.2026 |
| `20260115113936_add_city_column.sql` | Добавление колонки `city` | 15.01.2026 |
| `20260120070012_remove_column_user.sql` | Удаление устаревшей колонки | 20.01.2026 |

## 🛠️ Технологический стек

| Компонент | Версия | Назначение |
|-----------|--------|-----------|
| **Go** | 1.25.4 | Язык программирования |
| **pgx** | v5.8.0 | PostgreSQL драйвер |
| **Telegram Bot API** | v5.5.1 | Интеграция с Telegram |
| **godotenv** | v1.5.1 | Загрузка .env файлов |
| **OpenWeather** | - | Данные о погоде |

## 👨‍💻 Разработка

### Локальная разработка

```bash
# Установка зависимостей
go mod download

# Форматирование кода
go fmt ./...

# Проверка кода
go vet ./...

# Запуск тестов (если существуют)
go test ./...
```

### Структура кода

- **Clean Architecture:** Разделение на слои (clients, handler, models, repo)
- **Repository Pattern:** Абстракция доступа к данным
- **Error Handling:** Полное обработка ошибок с логированием

### Коммиты

При коммитах используйте следующий формат:

```
<тип>: <описание>

<подробное описание (опционально)>

<ссылка на issue (опционально)>
```

**Типы коммитов:**
- `feat:` - новая функция
- `fix:` - исправление ошибки
- `docs:` - обновление документации
- `style:` - форматирование кода
- `refactor:` - рефакторинг
- `test:` - добавление тестов
- `chore:` - обслуживание проекта

**Пример:**
```
feat: add weather forecast command

Implement 5-day forecast functionality using OpenWeather API.
Display forecast in user-friendly format.

Closes #42
```

## 📋 Лицензия

Проект использует лицензию, указанную в LICENSE файле.

## 👤 Автор

**sh-latibov** - [GitHub](https://github.com/sh-latibov)

## 📞 Контакты

- Telegram: [@sh_latibov](https://t.me/sh_latibov)
- GitHub Issues: [telegram-bot-go/issues](https://github.com/sh-latibov/telegram-bot-go/issues)

---

**Версия:** 1.0.0  
**Последнее обновление:** 20 января 2026 г.  
**Статус:** 🟢 Активная разработка

