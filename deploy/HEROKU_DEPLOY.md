# Инструкция по деплою на Heroku

## Шаг 1: Установка Heroku CLI

Heroku CLI уже установлен! ✅

Проверь установку:
```bash
heroku --version
```

## Шаг 2: Вход в Heroku

```bash
heroku login
```

Это откроет браузер для входа. Войди в свой аккаунт Heroku (или зарегистрируйся на https://heroku.com).

## Шаг 3: Создание приложения на Heroku

```bash
cd /Users/maloy/Desktop/Изучение/Go/day2
heroku create golang-day2-app
```

**ТЕОРИЯ:** 
- `heroku create` создаёт новое приложение на Heroku
- Имя должно быть уникальным (может быть занято, попробуй другое)
- Heroku создаст два remote: `heroku` и `origin`
- Ты получишь URL типа: https://golang-day2-app.herokuapp.com

## Шаг 4: Настройка переменных окружения (опционально)

```bash
heroku config:set ENVIRONMENT=production
```

**ТЕОРИЯ:**
- Переменные окружения настраиваются через `heroku config:set`
- `PORT` устанавливается автоматически Heroku (не нужно настраивать)
- Можно добавить другие переменные: `DATABASE_URL`, `JWT_SECRET` и т.д.

Посмотреть все переменные:
```bash
heroku config
```

## Шаг 5: Убедись что Procfile в корне проекта

Heroku ищет `Procfile` в корне репозитория. Если его нет - скопируй:

```bash
cp deploy/Procfile Procfile
```

Или создай симлинк:
```bash
ln -s deploy/Procfile Procfile
```

## Шаг 6: Деплой на Heroku

```bash
git add .
git commit -m "feat: подготовка к деплою на Heroku"
git push heroku main
```

**ТЕОРИЯ:**
- `git push heroku main` отправляет код на Heroku
- Heroku автоматически:
  1. Определяет что это Go приложение
  2. Устанавливает Go buildpack
  3. Компилирует приложение
  4. Запускает согласно Procfile
- Процесс может занять 2-5 минут

## Шаг 7: Проверка работы

```bash
# Открыть приложение в браузере
heroku open

# Или проверь логи
heroku logs --tail

# Проверь статус
heroku ps
```

## Шаг 8: Полезные команды

```bash
# Логи в реальном времени
heroku logs --tail

# Показать все переменные окружения
heroku config

# Добавить переменную
heroku config:set KEY=value

# Удалить переменную
heroku config:unset KEY

# Перезапустить приложение
heroku restart

# Открыть приложение
heroku open

# Показать информацию о приложении
heroku info

# Показать запущенные процессы
heroku ps

# Выполнить команду в контейнере
heroku run bash
```

## Решение проблем

### Ошибка: "No Procfile found"
**Решение:** Убедись что `Procfile` находится в корне проекта (не в подпапке)

### Ошибка: "App name already taken"
**Решение:** Используй другое имя:
```bash
heroku create другое-имя-приложения
```

### Ошибка: "Build failed"
**Решение:** 
1. Проверь логи: `heroku logs --tail`
2. Убедись что `go.mod` существует
3. Проверь что все зависимости указаны в `go.mod`

### Приложение не запускается
**Решение:**
1. Проверь логи: `heroku logs --tail`
2. Убедись что порт читается из `PORT`: `os.Getenv("PORT")`
3. Проверь что Procfile правильный

## Что дальше?

После успешного деплоя:
1. ✅ Приложение доступно по URL: https://твоё-приложение.herokuapp.com
2. ✅ Health check: https://твоё-приложение.herokuapp.com/health
3. ✅ Главная страница: https://твоё-приложение.herokuapp.com/

**Profit: твоё приложение работает в интернете! 🚀**

