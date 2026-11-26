# Mall Management System (Shopping Center)

Полноценная система управления торговым центром, состоящая из REST API на Go и нативного десктопного клиента на C# WPF. Проект построен с использованием чистой архитектуры (Clean Architecture), паттернов Repository и MVVM.

## Стек технологий

**Backend:**
- Язык: Go (Golang)
- База данных: MySQL (драйвер `go-sql-driver/mysql`)
- Авторизация: JWT (JSON Web Tokens) + Refresh Tokens
- Архитектура: Clean Architecture (Ports & Adapters), Repository, Unit of Work (транзакции)
- Тестирование: Unit-тесты (с моками) + Integration-тесты (с реальной БД)

**Frontend (Client):**
- Платформа: .NET (C# / WPF)
- Паттерн: MVVM (Model-View-ViewModel)
- DI-контейнер: Microsoft.Extensions.DependencyInjection
- Стилизация: Кастомные Control Templates (современный flat-дизайн)

## Основные возможности

- **Ролевая модель (RBAC):** 5 ролей (Admin, Manager, Accountant, Technician, Support). У каждой роли свой набор прав на просмотр и редактирование сущностей.
- **Управление:** Помещения (Premises), Арендаторы (Tenants), Договоры (Lease Contracts), Платежи (Payments), Заявки на обслуживание (Service Requests), Пользователи (Users).
- **Транзакционность:** При создании договора аренды помещение автоматически переводится в статус "Арендовано" в рамках одной транзакции БД.
- **Безопасность:** При регистрации пользователю жестко выдается роль `support` (нельзя зарегистрироваться админом). Пароли хешируются через `bcrypt`.
- **Автоинициализация:** При первом запуске бэкенд сам создает базу данных, таблицы и дефолтного администратора (данные берутся из `.env`).
- **Универсальный UI:** Формы создания/редактирования генерируются динамически на основе схем.

## Скриншоты

<p align="center">
  <img width="350" alt="Окно авторизации" src="https://github.com/user-attachments/assets/c15bf973-2ae4-4801-8576-b813e1292626" />
</p>

<p align="center">
  <img width="800" alt="Главное окно - Помещения" src="https://github.com/user-attachments/assets/19e5a2b1-0c41-4931-b2e3-e13e20590784" />
</p>

<p align="center">
  <img width="800" alt="Валидация при добавлении арендатора" src="https://github.com/user-attachments/assets/e9970a01-fedd-4a29-a55c-62ca4ac940db" />
</p>

<p align="center">
  <img width="800" alt="Подтверждение удаления пользователя" src="https://github.com/user-attachments/assets/4963d2fc-08b1-43db-83f2-27aa8a668448" />
</p>

## Структура проекта

```text
mall-api (Backend)
 ├── cmd/api            # Точка входа (main.go)
 ├── internal
 │   ├── adapter        # Реализации (MySQL, JWT)
 │   ├── db             # Инициализация БД (schema.sql)
 │   ├── domain         # Бизнес-сущности и валидация
 │   ├── handler        # HTTP-контроллеры и middleware
 │   ├── port           # Интерфейсы (контракты)
 │   └── service        # Бизнес-логика
 └── .env               # Конфигурация

MallClient (Frontend)
 ├── Converters         # XAML конвертеры (Роли, Статусы)
 ├── Models             # DTO и модели данных
 ├── Services           # API-клиенты и логика
 ├── ViewModels         # MVVM логика
 ├── Views              # XAML окна
 └── appsettings.json   # URL бэкенда
```

## Установка и запуск

### 1. Подготовка БД
У вас должен быть установлен и запущен MySQL сервер. Создавать базу вручную не нужно — бэкенд сделает это сам при старте.

### 2. Настройка переменных окружения (Backend)
В корне проекта `mall-api` создайте файл `.env`:

```env
DB_USER=root
DB_PASS=root
DB_HOST=127.0.0.1:3306
DB_NAME=shopping_center
JWT_SECRET=ваш_супер_секретный_ключ_минимум_32_символа
PORT=8080

# Данные для автоматического создания админа при первом старте
ADMIN_EMAIL=admin@mall-system.local
ADMIN_PASSWORD=SuperSecretPassword_2026!
```

### 3. Запуск бэкенда
```bash
cd mall-api
go mod tidy
go run ./cmd/main.go
```
*После первого запуска в консоли появятся логи создания БД и админа.*

### 4. Запуск клиента
1. Откройте решение `MallClient.sln` в Visual Studio.
2. Убедитесь, что в `appsettings.json` указан правильный URL бэкенда:

```json
{
  "ApiSettings": {
    "BaseUrl": "http://localhost:8080"
  }
}
```

3. Запустите проект (F5).
4. Войдите под созданным админом (email и пароль из `.env`).

## Матрица прав (RBAC)

| Сущность | Admin | Manager | Accountant | Technician | Support |
| :--- | :---: | :---: | :---: | :---: | :---: |
| Помещения (Premises) | CRUD | R | R | - | - |
| Арендаторы (Tenants) | CRUD | CRU | - | - | - |
| Договоры (Contracts) | CRUD | CRU | R | - | - |
| Платежи (Payments) | CRUD | - | CRU | - | - |
| Заявки (Requests) | CRUD | - | - | RU | CR |
| Пользователи (Users) | CRUD | - | - | - | - |

*(C - Create, R - Read, U - Update, D - Delete)*

## Основные эндпоинты API

Все защищенные эндпоинты требуют заголовок `Authorization: Bearer <token>`.

- `POST /api/auth/login` — Авторизация (возвращает Access + Refresh токены)
- `POST /api/auth/register` — Регистрация (создает пользователя с ролью support)
- `GET /api/premises` — Список помещений
- `POST /api/contracts` — Создание договора (с транзакцией)
- `GET /api/users` — Список пользователей (только Admin)

## Тестирование

В проекте настроены Unit-тесты (с использованием `testify/mock`).
Для запуска интеграционных тестов (требуется БД `shopping_center_test`):

```bash
cd mall-api
go test -tags=integration ./...
```

---
*Разработано в рамках курсового проекта.*
