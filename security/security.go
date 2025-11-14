package security

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
)

type User struct {
	ID    int
	Email string
	Name  string
	Age   int
}

// ============================================================================
// ТЕОРИЯ: Безопасность в веб-приложениях
// ============================================================================
//
// ОСНОВНЫЕ УГРОЗЫ:
// 1. SQL Injection - внедрение SQL кода через входные данные
// 2. XSS (Cross-Site Scripting) - внедрение JavaScript через входные данные
// 3. Невалидированные данные - атаки через некорректные входные данные
// 4. CSRF (Cross-Site Request Forgery) - подделка запросов
//
// ПРИНЦИПЫ БЕЗОПАСНОСТИ:
// 1. Валидируй ВСЕ входные данные (не доверяй клиенту!)
// 2. Используй параметризованные запросы для SQL
// 3. Экранируй HTML/JavaScript в выводе
// 4. Используй HTTPS для передачи данных
// 5. Ограничивай права доступа (принцип наименьших привилегий)
//
// ============================================================================

// ============================================================================
// 1. ВАЛИДАЦИЯ ДАННЫХ
// ============================================================================
//
// ТЕОРИЯ: Валидация - это проверка входных данных на корректность
// - Проверяй данные на сервере (клиент может обойти валидацию!)
// - Валидируй тип, формат, длину, диапазон значений
// - Используй whitelist (разрешённые значения) вместо blacklist
//
// ЧТО ВАЛИДИРОВАТЬ:
// - Email адреса (формат)
// - Пароли (длина, сложность)
// - Числа (диапазон, тип)
// - Строки (длина, формат)
// - URL (формат, протокол)
//

// ValidateEmail проверяет корректность email адреса
func ValidateEmail(email string) error {
	// ТЕОРИЯ: Валидация email
	// - Проверяем что email не пустой
	// - Проверяем формат через регулярное выражение
	// - Проверяем длину (максимум обычно 254 символа)

	// TODO: проверь что email не пустой
	// TODO: проверь длину email (максимум 254 символа)
	// TODO: используй регулярное выражение для проверки формата
	// TODO: пример паттерна: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Проверка на пустоту
	if email == "" {
		return fmt.Errorf("email не может быть пустым")
	}

	// Проверка длины
	if len(email) > 254 {
		return fmt.Errorf("email слишком длинный (максимум 254 символа)")
	}

	// Проверка формата через регулярное выражение
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("неверный формат email")
	}

	return nil
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string) error {
	// ТЕОРИЯ: Валидация пароля
	// - Минимальная длина (обычно 8+ символов)
	// - Наличие заглавных букв
	// - Наличие строчных букв
	// - Наличие цифр
	// - Наличие специальных символов (опционально)

	// TODO: проверь минимальную длину (8 символов)
	// TODO: проверь наличие заглавных букв
	// TODO: проверь наличие строчных букв
	// TODO: проверь наличие цифр
	// TODO: верни понятную ошибку если пароль не соответствует требованиям

	if len(password) < 8 {
		return fmt.Errorf("пароль должен содержать минимум 8 символов")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper {
		return fmt.Errorf("пароль должен содержать заглавные буквы")
	}
	if !hasLower {
		return fmt.Errorf("пароль должен содержать строчные буквы")
	}
	if !hasDigit {
		return fmt.Errorf("пароль должен содержать цифры")
	}

	return nil
}

// ValidateAge проверяет корректность возраста
func ValidateAge(age int) error {
	// ТЕОРИЯ: Валидация числовых значений
	// - Проверяем диапазон (обычно 0-150 для возраста)
	// - Проверяем что значение не отрицательное

	// TODO: проверь что возраст в допустимом диапазоне (0-150)
	// TODO: верни ошибку если возраст некорректен

	if age < 0 {
		return fmt.Errorf("возраст не может быть отрицательным")
	}

	if age > 150 {
		return fmt.Errorf("возраст слишком большой (максимум 150 лет)")
	}

	return nil
}

// SanitizeInput очищает входные данные от опасных символов
func SanitizeInput(input string) string {
	// ТЕОРИЯ: Санитизация (очистка) входных данных
	// - Удаляем или экранируем опасные символы
	// - Обрезаем пробелы в начале и конце
	// - Ограничиваем длину

	// TODO: обрежь пробелы в начале и конце (strings.TrimSpace)
	// TODO: ограничь длину строки (например, максимум 1000 символов)
	// TODO: удали или замени опасные символы (например, нулевые байты)

	// Обрезаем пробелы
	sanitized := strings.TrimSpace(input)

	// Ограничиваем длину
	maxLength := 1000
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength]
	}

	// Удаляем нулевые байты (опасны для некоторых систем)
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")

	return sanitized
}

// ============================================================================
// 2. ЗАЩИТА ОТ SQL INJECTION
// ============================================================================
//
// ТЕОРИЯ: SQL Injection - это атака через внедрение SQL кода
// - Злоумышленник передаёт SQL код через входные данные
// - Если запрос формируется через конкатенацию строк - уязвимость!
//
// ПРИМЕР АТАКИ:
//   username := "admin' OR '1'='1"
//   query := "SELECT * FROM users WHERE username = '" + username + "'"
//   // Получится: SELECT * FROM users WHERE username = 'admin' OR '1'='1'
//   // Это вернёт ВСЕХ пользователей!
//
// ЗАЩИТА:
// 1. Используй параметризованные запросы (prepared statements)
// 2. НИКОГДА не конкатенируй SQL через строки!
// 3. Используй ORM (GORM, SQLX) - они автоматически защищают
// 4. Валидируй входные данные перед использованием
//

// GetUserByEmail НЕБЕЗОПАСНЫЙ способ (уязвим к SQL injection)
func GetUserByEmailUnsafe(db *sql.DB, email string) (*User, error) {
	// ТЕОРИЯ: НЕПРАВИЛЬНЫЙ способ - конкатенация строк
	// - Это уязвимо к SQL injection!
	// - НИКОГДА так не делай!

	// ПЛОХОЙ КОД (НЕ ИСПОЛЬЗУЙ!):
	// query := "SELECT id, email, name FROM users WHERE email = '" + email + "'"
	// rows, err := db.Query(query)

	// TODO: покажи НЕПРАВИЛЬНЫЙ способ (закомментированный)
	// TODO: объясни почему это опасно

	return nil, nil
}

// GetUserByEmail БЕЗОПАСНЫЙ способ (защищён от SQL injection)
func GetUserByEmailSafe(db *sql.DB, email string) (*User, error) {
	// ТЕОРИЯ: ПРАВИЛЬНЫЙ способ - параметризованный запрос
	// - Используем плейсхолдеры (?) вместо конкатенации
	// - БД сама экранирует параметры
	// - Это защищает от SQL injection

	// TODO: используй параметризованный запрос
	// query := "SELECT id, email, name FROM users WHERE email = ?"
	// TODO: передай email как параметр
	// row := db.QueryRow(query, email)
	// TODO: прочитай данные из row

	// ВАЛИДАЦИЯ перед запросом
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	// ПАРАМЕТРИЗОВАННЫЙ ЗАПРОС - защита от SQL injection
	query := "SELECT id, email, name FROM users WHERE email = ?"
	row := db.QueryRow(query, email) // email передаётся как параметр

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser создаёт пользователя с защитой от SQL injection
func CreateUser(ctx context.Context, db *sql.DB, user *User) error {
	// ТЕОРИЯ: Создание записи с параметризованным запросом
	// - Используем ExecContext с параметрами
	// - Все значения передаём как параметры, не через конкатенацию

	// TODO: валидируй данные перед вставкой (ValidateEmail, ValidatePassword)
	// TODO: используй параметризованный запрос для INSERT
	// query := "INSERT INTO users (email, name, age) VALUES (?, ?, ?)"
	// TODO: выполни запрос с параметрами
	// _, err := db.ExecContext(ctx, query, user.Email, user.Name, user.Age)

	return nil
}

// SearchUsers безопасный поиск пользователей
func SearchUsers(db *sql.DB, searchTerm string) ([]User, error) {
	// ТЕОРИЯ: Поиск с LIKE и защитой от SQL injection
	// - Используем параметризованный запрос даже для LIKE
	// - Экранируем специальные символы SQL (%, _)

	// TODO: санитизируй searchTerm (SanitizeInput)
	// TODO: экранируй специальные символы для LIKE (% и _)
	// TODO: используй параметризованный запрос
	// query := "SELECT * FROM users WHERE name LIKE ?"
	// searchPattern := "%" + escapedTerm + "%"
	// rows, err := db.Query(query, searchPattern)

	sanitized := SanitizeInput(searchTerm)

	// Экранирование специальных символов для LIKE
	// % и _ - это специальные символы в SQL LIKE
	escaped := strings.ReplaceAll(sanitized, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")

	// Параметризованный запрос
	query := "SELECT id, email, name FROM users WHERE name LIKE ?"
	pattern := "%" + escaped + "%"
	rows, err := db.Query(query, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// ============================================================================
// 3. ЗАЩИТА ОТ XSS (Cross-Site Scripting)
// ============================================================================
//
// ТЕОРИЯ: XSS - это внедрение JavaScript через входные данные
// - Злоумышленник передаёт JavaScript код в форме
// - Если код выводится без экранирования - он выполнится в браузере
//
// ПРИМЕР АТАКИ:
//   userInput := "<script>alert('XSS')</script>"
//   // Если вывести без экранирования: <div>{{userInput}}</div>
//   // Браузер выполнит JavaScript!
//
// ЗАЩИТА:
// 1. Экранируй HTML при выводе (html.EscapeString)
// 2. Используй Content Security Policy (CSP)
// 3. Валидируй и санитизируй входные данные
// 4. Используй шаблоны с автоматическим экранированием (html/template)
//

// EscapeHTML экранирует HTML символы
func EscapeHTML(input string) string {
	// ТЕОРИЯ: Экранирование HTML
	// - html.EscapeString заменяет опасные символы на HTML entities
	// - < становится &lt;
	// - > становится &gt;
	// - & становится &amp;

	// TODO: используй html.EscapeString для экранирования
	// return html.EscapeString(input)

	return html.EscapeString(input)
}

// RenderUserHTML безопасно выводит данные пользователя в HTML
func RenderUserHTML(user *User) string {
	// ТЕОРИЯ: Безопасный вывод в HTML
	// - Экранируем все данные перед выводом
	// - Это предотвращает XSS атаки

	// TODO: экранируй все поля пользователя
	// escapedName := html.EscapeString(user.Name)
	// escapedEmail := html.EscapeString(user.Email)
	// TODO: собери безопасный HTML
	// return fmt.Sprintf("<div>Name: %s, Email: %s</div>", escapedName, escapedEmail)

	// ЭКРАНИРОВАНИЕ всех данных перед выводом
	escapedName := html.EscapeString(user.Name)
	escapedEmail := html.EscapeString(user.Email)

	// Безопасный HTML
	return fmt.Sprintf(
		"<div class='user'>"+
			"<p>Name: %s</p>"+
			"<p>Email: %s</p>"+
			"</div>",
		escapedName, escapedEmail)
}

// HandleUserInput обрабатывает пользовательский ввод с защитой от XSS
func HandleUserInput(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// ТЕОРИЯ: Обработка пользовательского ввода
	// - Получаем данные из формы
	// - Валидируем данные
	// - Санитизируем данные
	// - Экранируем при выводе

	// TODO: получи данные из формы (r.FormValue)
	// userInput := r.FormValue("comment")
	// TODO: санитизируй данные (SanitizeInput)
	// TODO: валидируй данные (проверь длину, формат)
	// TODO: сохрани в БД (с параметризованным запросом!)
	// TODO: при выводе экранируй HTML (EscapeHTML)

	comment := r.FormValue("comment")

	// Валидация
	if len(comment) == 0 {
		http.Error(w, "комментарий не может быть пустым", http.StatusBadRequest)
		return
	}

	if len(comment) > 1000 {
		http.Error(w, "комментарий слишком длинный", http.StatusBadRequest)
		return
	}

	// Санитизация
	sanitized := SanitizeInput(comment)

	// Сохранение в БД (с параметризованным запросом!)
	query := "INSERT INTO comments (text) VALUES (?)"
	_, err := db.Exec(query, sanitized)
	if err != nil {
		http.Error(w, "ошибка сохранения", http.StatusInternalServerError)
		return
	}

	// Вывод (с экранированием!)
	escaped := html.EscapeString(sanitized)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<p>Комментарий: %s</p>", escaped)
}

// ============================================================================
// 4. ДОПОЛНИТЕЛЬНЫЕ МЕРЫ БЕЗОПАСНОСТИ
// ============================================================================

// ValidateURL проверяет корректность URL
func ValidateURL(url string) error {
	// ТЕОРИЯ: Валидация URL
	// - Проверяем формат URL
	// - Проверяем протокол (только http/https)
	// - Проверяем домен (whitelist разрешённых доменов)

	// TODO: проверь что URL не пустой
	// TODO: используй regexp для проверки формата URL
	// TODO: проверь что протокол http или https
	// TODO: проверь домен (можно использовать whitelist)

	return nil
}

// SanitizeSQLIdentifier очищает идентификатор SQL (имя таблицы, колонки)
func SanitizeSQLIdentifier(identifier string) string {
	// ТЕОРИЯ: Санитизация SQL идентификаторов
	// - Используется когда нужно динамически формировать имена таблиц/колонок
	// - Разрешаем только буквы, цифры и подчёркивания
	// - НИКОГДА не используй пользовательский ввод для имён таблиц/колонок!

	// TODO: проверь что identifier содержит только разрешённые символы
	// TODO: разрешённые символы: буквы (a-z, A-Z), цифры (0-9), подчёркивание (_)
	// TODO: если содержит недопустимые символы - верни ошибку или замени их

	return identifier
}

// RateLimit проверяет лимит запросов (защита от брутфорса)
func RateLimit(ip string, maxRequests int, windowSeconds int) error {
	// ТЕОРИЯ: Rate Limiting (ограничение частоты запросов)
	// - Защита от брутфорс атак
	// - Ограничиваем количество запросов с одного IP
	// - Обычно используется для логина, регистрации, API endpoints

	// TODO: используй кэш (например, Redis) для хранения счётчика запросов
	// TODO: ключ: "ratelimit:" + ip
	// TODO: увеличивай счётчик при каждом запросе
	// TODO: проверяй лимит (если превышен - верни ошибку)
	// TODO: устанавливай TTL на ключ (windowSeconds)

	return nil
}

// КЛЮЧЕВЫЕ ПРАВИЛА:
// 1. ВСЕГДА валидируй входные данные
// 2. ВСЕГДА используй параметризованные SQL запросы
// 3. ВСЕГДА экранируй HTML при выводе
// 4. НИКОГДА не доверяй данным от клиента
// 5. Используй HTTPS для передачи данных
// 6. Ограничивай права доступа (принцип наименьших привилегий)
//
// Полезные команды для проверки:
//   go test ./security/...                    - запустить тесты безопасности
//   go vet ./security/...                    - проверить код на ошибки
//   go run ./security/...                    - запустить примеры
//
// Следующий шаг:
//   1. Реализуй все TODO функции
//   2. Создай тесты для проверки безопасности
//   3. Протестируй защиту от SQL injection и XSS
//   4. Используй эти паттерны в реальном проекте
//
// Profit: защищённое приложение + доверие пользователей

// Подавляем предупреждения о неиспользуемых импортах
// (они используются в TODO функциях и примерах кода)
var (
	_ = fmt.Errorf
	_ = html.EscapeString
	_ = http.StatusOK
	_ = regexp.MustCompile
	_ = strings.TrimSpace
)
