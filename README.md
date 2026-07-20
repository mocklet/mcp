# Mocklet MCP Server

Mocklet MCP Server — это сервер на базе [Model Context Protocol (MCP)](https://modelcontextprotocol.io/), который позволяет AI-ассистентам (таким как Claude Desktop, Cursor, Windsurf и др.) нативно взаимодействовать с платформой мокирования Mocklet (Harmockery). 

С его помощью ИИ может "из коробки" загружать HAR-файлы, создавать шаблоны и разворачивать одноразовые (эфемеровые) mock-серверы для автономного прототипирования фронтенда и E2E-тестирования.

## Зачем это нужно?

- **Автономная генерация тестов:** ИИ может самостоятельно поднять мок-сервер из HAR-файла, написать Cypress/Playwright тесты для проверки UI и затем удалить мок.
- **Прототипирование Frontend:** ИИ-агент может за долю секунды создать реалистичный backend из шаблона Mocklet и начать верстать дашборд.
- **Дебаггинг API:** Получение статистики использования мока (hit/miss) и отладка отсутствующих роутов прямо в чате с AI.

---

## 🛠 Установка

### Требования
- Установленный [Go](https://go.dev/) версии 1.21 или новее.

### Сборка из исходников
Склонируйте репозиторий и соберите бинарный файл:

```bash
cd mcp
go mod tidy
go build -o mocklet-mcp .
```

Это создаст исполняемый файл `mocklet-mcp` в текущей директории. Убедитесь, что вы запомнили или скопировали абсолютный путь к этому файлу (например, `/home/user/coding/mcp/mocklet-mcp`), он понадобится для настройки клиентов.

---

## ⚙️ Конфигурация (Переменные окружения)

Для работы сервера требуются следующие переменные окружения:

| Переменная | Описание | Пример |
| --- | --- | --- |
| `MOCKLET_API_URL` | Базовый URL вашего Mocklet API. Если не указан, используется `http://localhost:8080` по умолчанию. | `https://api.mocklet.dev` |
| `MOCKLET_SERVICE_TOKEN` | Сервисный токен (Bearer Token) для аутентификации в API Mocklet. | `mckt_123456789...` |

---

## 🚀 Настройка клиентов

### Claude Desktop
Откройте конфигурационный файл Claude Desktop (обычно находится в `~/Library/Application Support/Claude/claude_desktop_config.json` на macOS или `%APPDATA%\Claude\claude_desktop_config.json` на Windows) и добавьте секцию `mcpServers`:

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/абсолютный/путь/к/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "ваш_сервисный_токен"
      }
    }
  }
}
```
*После изменения файла перезапустите Claude Desktop.*

### Cursor
1. Перейдите в **Settings > Features > MCP**.
2. Нажмите **+ Add New MCP Server**.
3. Тип (Type): `command`
4. Имя (Name): `mocklet`
5. Команда (Command): `MOCKLET_API_URL="http://localhost:8080" MOCKLET_SERVICE_TOKEN="ваш_токен" /абсолютный/путь/к/mocklet-mcp`

### Google Antigravity (AGY)
Для интеграции с Antigravity 2.0 (AGY / IDE), добавьте сервер в конфигурационный файл MCP (`~/.gemini/antigravity-cli/mcp.json`):

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/абсолютный/путь/к/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "ваш_сервисный_токен"
      }
    }
  }
}
```
*После изменения файла, перезапустите AGY.*

### Codex / Cline / VS Code
Если вы используете расширения вроде Cline (ранее Claude Dev) или Codex в VS Code, вам необходимо отредактировать файл настроек MCP (например, `cline_mcp_settings.json` для Cline, обычно расположен в `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json` на macOS):

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/абсолютный/путь/к/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "ваш_сервисный_токен"
      }
    }
  }
}
```

### Zed IDE
Для локального использования в [Zed](https://zed.dev/) откройте файл настроек (`~/.config/zed/settings.json`) и добавьте блок `context_servers`:

```json
{
  "context_servers": {
    "mocklet": {
      "command": "/абсолютный/путь/к/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "ваш_сервисный_токен"
      }
    }
  }
}
```

---

## 🧰 Доступные инструменты (Tools)

После подключения ИИ-ассистент получит доступ к следующим операциям:

- `mocklet_validate_har` — валидация HAR-файла перед деплоем.
- `mocklet_create_mock` — создание одноразового мок-сервера из HAR-файла.
- `mocklet_list_mocks` — получение списка активных моков.
- `mocklet_get_mock_stats` — получение статистики (хиты/промахи) конкретного мока.
- `mocklet_delete_mock` — остановка и удаление мок-сервера.
- `mocklet_create_template` — загрузка HAR-файла для создания переиспользуемого шаблона.
- `mocklet_list_templates` — поиск существующих шаблонов.
- `mocklet_spawn_mock` — быстрый запуск эфемерового мока на базе готового шаблона.
- `mocklet_upload_template_revision` — обновление логики шаблона новым HAR-файлом.

## 💬 Встроенные промпты (Prompts)

Сервер также предоставляет готовые сценарии (prompts), доступные в клиентах с поддержкой этой функции (например, в Claude), для автоматизации популярных задач:
- **Spawn Dependency Mock** (Запуск зависимости)
- **Debug Mock Usage** (Дебаг ошибок 404)
- **Generate Frontend with Mock Data** (Создание UI с тестовыми данными)
- **Integration Testing Setup** (Настройка Cypress/Playwright тестов)
- **HAR Validation & Cleanup** (Очистка и проверка HAR)
