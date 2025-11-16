# Деплой приложения

Этот каталог содержит примеры и конфигурацию для деплоя Go приложения на Render и Heroku.

## Файлы

- `main.go` - пример приложения готового к деплою
- `Dockerfile` - конфигурация Docker для контейнеризации
- `render.yaml` - конфигурация для автоматического деплоя на Render
- `Procfile` - конфигурация для Heroku
- `app.json` - метаданные для Heroku
- `.dockerignore` - файлы, которые игнорируются при сборке Docker образа
- `DEPLOY_GUIDE.md` - подробное руководство по деплою

## Быстрый старт

### Локальный запуск

```bash
# Установи переменную окружения (опционально)
export PORT=8080

# Запусти приложение
go run main.go

# Или скомпилируй и запусти
go build -o server main.go
./server
```

### Деплой на Render

1. Зарегистрируйся на https://render.com
2. Нажми "New +" -> "Web Service"
3. Подключи GitHub репозиторий
4. Настрой:
   - Build Command: `go build -o server ./deploy`
   - Start Command: `./server`
5. Добавь переменные окружения через Dashboard
6. Нажми "Create Web Service"

### Деплой на Heroku

1. Установи Heroku CLI
2. Войди: `heroku login`
3. Создай приложение: `heroku create golang-app`
4. Настрой переменные: `heroku config:set KEY=value`
5. Деплой: `git push heroku main`
6. Открой: `heroku open`

## Проверка работы

После деплоя проверь:

- `https://your-app.onrender.com/health` - должен вернуть "OK"
- `https://your-app.onrender.com/` - главная страница

## Дополнительная информация

Смотри `DEPLOY_GUIDE.md` для подробной информации о деплое.

