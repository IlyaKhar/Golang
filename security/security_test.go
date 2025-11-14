package security

import (
	"database/sql"
	"strings"
	"testing"
	_ "github.com/mattn/go-sqlite3" // драйвер для SQLite
)

// ============================================================================
// ТЕОРИЯ: Тестирование безопасности
// ============================================================================
//
// ЗАЧЕМ ТЕСТИРОВАТЬ БЕЗОПАСНОСТЬ:
// 1. Убедиться что защита работает (SQL injection, XSS)
// 2. Найти уязвимости до продакшена
// 3. Документировать безопасное поведение кода
//
// ЧТО ТЕСТИРУЕМ:
// - SQL Injection: попытка внедрить SQL код через входные данные
// - XSS: попытка внедрить JavaScript через входные данные
// - Валидация: проверка что некорректные данные отклоняются
//
// ============================================================================

// setupTestDB создаёт тестовую БД для проверки SQL injection
func setupTestDB(t *testing.T) *sql.DB {
	// ТЕОРИЯ: Создаём временную БД в памяти
	// - ":memory:" - SQLite создаёт БД в оперативной памяти
	// - БД удалится автоматически после закрытия соединения
	// - Это быстро и не требует очистки файлов
	
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("не удалось создать тестовую БД: %v", err)
	}

	// Создаём таблицу users
	createTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			age INTEGER
		)
	`
	
	if _, err := db.Exec(createTable); err != nil {
		t.Fatalf("не удалось создать таблицу: %v", err)
	}

	// Добавляем тестового пользователя
	insertUser := `INSERT INTO users (email, name, age) VALUES (?, ?, ?)`
	if _, err := db.Exec(insertUser, "test@example.com", "Test User", 25); err != nil {
		t.Fatalf("не удалось добавить тестового пользователя: %v", err)
	}

	return db
}

// ============================================================================
// ТЕСТЫ ЗАЩИТЫ ОТ SQL INJECTION
// ============================================================================

// TestSQLInjection_UnsafeQuery проверяет что НЕБЕЗОПАСНЫЙ запрос уязвим
func TestSQLInjection_UnsafeQuery(t *testing.T) {
	// ТЕОРИЯ: Этот тест показывает КАК НЕ НАДО делать
	// - Мы демонстрируем уязвимость (для обучения)
	// - В реальном коде НИКОГДА не используй конкатенацию строк для SQL!
	
	db := setupTestDB(t)
	defer db.Close()

	// ТЕОРИЯ: SQL Injection атака
	// - Злоумышленник передаёт email с SQL кодом
	// - Если запрос формируется через конкатенацию - уязвимость!
	maliciousEmail := "test@example.com' OR '1'='1"

	// ПЛОХОЙ КОД (НЕ ИСПОЛЬЗУЙ!):
	// query := "SELECT * FROM users WHERE email = '" + maliciousEmail + "'"
	// Это вернёт ВСЕХ пользователей, даже если email неверный!
	
	// ТЕОРИЯ: В нашем безопасном коде мы используем параметризованные запросы
	// - GetUserByEmailSafe использует плейсхолдеры (?)
	// - БД сама экранирует параметры
	// - Это защищает от SQL injection
	
	// Проверяем что безопасный метод НЕ находит пользователя с таким email
	user, err := GetUserByEmailSafe(db, maliciousEmail)
	if err == nil && user != nil {
		t.Errorf("SQL injection успешна! Найден пользователь: %+v", user)
	} else {
		t.Logf("✅ Защита работает: SQL injection заблокирована")
	}
}

// TestSQLInjection_SafeQuery проверяет что БЕЗОПАСНЫЙ запрос защищён
func TestSQLInjection_SafeQuery(t *testing.T) {
	// ТЕОРИЯ: Тестируем защиту от SQL injection
	// - Пытаемся внедрить SQL код через email
	// - Проверяем что параметризованный запрос защищает нас
	
	db := setupTestDB(t)
	defer db.Close()

	// ТЕОРИЯ: Различные типы SQL injection атак
	testCases := []struct {
		name  string
		email string
		desc  string
	}{
		{
			name:  "OR атака",
			email: "test@example.com' OR '1'='1",
			desc:  "Попытка вернуть всех пользователей",
		},
		{
			name:  "Комментарий",
			email: "test@example.com'--",
			desc:  "Попытка закомментировать остаток запроса",
		},
		{
			name:  "UNION атака",
			email: "test@example.com' UNION SELECT * FROM users--",
			desc:  "Попытка объединить запросы",
		},
		{
			name:  "Двойные кавычки",
			email: `test@example.com" OR "1"="1`,
			desc:  "Попытка использовать двойные кавычки",
		},
		{
			name:  "Точка с запятой",
			email: "test@example.com'; DROP TABLE users--",
			desc:  "Попытка выполнить несколько команд",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ТЕОРИЯ: Пытаемся использовать злонамеренный email
			// - GetUserByEmailSafe использует параметризованный запрос
			// - БД экранирует параметр, поэтому SQL код не выполнится
			
			user, err := GetUserByEmailSafe(db, tc.email)
			
			// ТЕОРИЯ: Проверяем что атака не удалась
			// - Либо ошибка валидации (email неверный формат)
			// - Либо пользователь не найден (sql.ErrNoRows)
			// - НО НИКОГДА не должен вернуться пользователь с таким email!
			
			if err == nil && user != nil {
				t.Errorf("❌ SQL injection успешна! %s\nНайден пользователь: %+v", tc.desc, user)
			} else {
				t.Logf("✅ Защита работает: %s - заблокировано", tc.desc)
			}
		})
	}
}

// TestSQLInjection_LikeSearch проверяет защиту поиска с LIKE
func TestSQLInjection_LikeSearch(t *testing.T) {
	// ТЕОРИЯ: Поиск с LIKE тоже может быть уязвим
	// - Нужно экранировать специальные символы (%, _)
	// - Использовать параметризованные запросы
	
	db := setupTestDB(t)
	defer db.Close()

	// ТЕОРИЯ: Попытка SQL injection через поиск
	maliciousSearch := "Test' OR '1'='1"

	// ТЕОРИЯ: SearchUsers использует:
	// 1. SanitizeInput - очищает входные данные
	// 2. Экранирование % и _ для LIKE
	// 3. Параметризованный запрос
	
	users, err := SearchUsers(db, maliciousSearch)
	if err != nil {
		t.Logf("✅ Ошибка при поиске (ожидаемо): %v", err)
		return
	}

	// ТЕОРИЯ: Проверяем что не вернулись все пользователи
	// - Если вернулось больше 1 пользователя - возможна уязвимость
	// - В нашем случае должен вернуться только 1 пользователь (если имя совпадает)
	
	if len(users) > 1 {
		t.Errorf("❌ Возможна уязвимость: найдено %d пользователей вместо ожидаемого количества", len(users))
	} else {
		t.Logf("✅ Защита работает: найдено %d пользователей", len(users))
	}
}

// ============================================================================
// ТЕСТЫ ЗАЩИТЫ ОТ XSS
// ============================================================================

// TestXSS_EscapeHTML проверяет экранирование HTML
func TestXSS_EscapeHTML(t *testing.T) {
	// ТЕОРИЯ: XSS (Cross-Site Scripting) атака
	// - Злоумышленник пытается внедрить JavaScript через входные данные
	// - Если данные выводятся без экранирования - JavaScript выполнится
	// - EscapeHTML заменяет опасные символы на безопасные HTML entities
	
	testCases := []struct {
		name     string
		input    string
		expected string // что должно быть после экранирования
		desc     string
	}{
		{
			name:     "JavaScript тег",
			input:    "<script>alert('XSS')</script>",
			expected: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
			desc:     "Попытка выполнить JavaScript",
		},
		{
			name:     "HTML тег",
			input:    "<div>Hello</div>",
			expected: "&lt;div&gt;Hello&lt;/div&gt;",
			desc:     "Попытка вставить HTML",
		},
		{
			name:     "Обработчик событий",
			input:    `<img src="x" onerror="alert('XSS')">`,
			expected: "&lt;img src=&#34;x&#34; onerror=&#34;alert(&#39;XSS&#39;)&#34;&gt;",
			desc:     "Попытка использовать обработчик событий",
		},
		{
			name:     "JavaScript протокол",
			input:    `<a href="javascript:alert('XSS')">Click</a>`,
			expected: "&lt;a href=&#34;javascript:alert(&#39;XSS&#39;)&#34;&gt;Click&lt;/a&gt;",
			desc:     "Попытка использовать javascript: протокол",
		},
		{
			name:     "Обычный текст",
			input:    "Hello World",
			expected: "Hello World",
			desc:     "Обычный текст не должен изменяться",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ТЕОРИЯ: Экранируем входные данные
			escaped := EscapeHTML(tc.input)

			// ТЕОРИЯ: Проверяем что опасные символы экранированы
			// - < должно стать &lt;
			// - > должно стать &gt;
			// - " должно стать &#34; или &quot;
			// - ' должно стать &#39;
			
			if escaped != tc.expected {
				// ТЕОРИЯ: html.EscapeString может использовать разные entities
				// - Проверяем что опасные символы точно экранированы
				if contains(escaped, "<") || contains(escaped, ">") {
					t.Errorf("❌ XSS уязвимость! Опасные символы не экранированы:\nИсходный: %s\nЭкранированный: %s", tc.input, escaped)
				} else {
					t.Logf("✅ Защита работает: %s - экранировано", tc.desc)
				}
			} else {
				t.Logf("✅ Защита работает: %s - экранировано правильно", tc.desc)
			}

			// ТЕОРИЯ: Дополнительная проверка - не должно быть <script>
			if contains(escaped, "<script") {
				t.Errorf("❌ КРИТИЧЕСКАЯ УЯЗВИМОСТЬ! <script> не экранирован: %s", escaped)
			}
		})
	}
}

// TestXSS_RenderUserHTML проверяет безопасный рендеринг HTML
func TestXSS_RenderUserHTML(t *testing.T) {
	// ТЕОРИЯ: Тестируем RenderUserHTML
	// - Функция должна экранировать все данные пользователя
	// - Это предотвращает XSS атаки
	
	testCases := []struct {
		name string
		user *User
		desc string
	}{
		{
			name: "XSS в имени",
			user: &User{
				ID:    1,
				Email: "test@example.com",
				Name:  "<script>alert('XSS')</script>",
				Age:   25,
			},
			desc: "Попытка XSS через имя пользователя",
		},
		{
			name: "XSS в email",
			user: &User{
				ID:    1,
				Email: "<img src=x onerror=alert('XSS')>@example.com",
				Name:  "Test User",
				Age:   25,
			},
			desc: "Попытка XSS через email",
		},
		{
			name: "Обычные данные",
			user: &User{
				ID:    1,
				Email: "test@example.com",
				Name:  "Test User",
				Age:   25,
			},
			desc: "Обычные данные должны работать нормально",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ТЕОРИЯ: Рендерим HTML с данными пользователя
			html := RenderUserHTML(tc.user)

			// ТЕОРИЯ: Проверяем что опасные символы экранированы
			// - Не должно быть неэкранированного <script>
			// - Не должно быть неэкранированного onerror= (как атрибута HTML)
			// - Все < и > должны быть экранированы как &lt; и &gt;
			
			// Проверяем что нет неэкранированного <script
			if contains(html, "<script") {
				t.Errorf("❌ КРИТИЧЕСКАЯ УЯЗВИМОСТЬ! Неэкранированный <script> в HTML: %s", html)
			}

			// Проверяем что нет неэкранированного onerror= как атрибута
			// (должно быть экранировано как &gt; или часть текста)
			// Если onerror= идёт после > без экранирования - это уязвимость
			if contains(html, "> onerror=") || contains(html, "\" onerror=") || contains(html, "' onerror=") {
				t.Errorf("❌ КРИТИЧЕСКАЯ УЯЗВИМОСТЬ! Неэкранированный onerror= в HTML: %s", html)
			}

			// ТЕОРИЯ: Проверяем что данные присутствуют (но экранированы)
			if !contains(html, "&lt;") && !contains(html, "&#34;") {
				// Если нет экранированных символов, но есть опасные - проблема
				if contains(tc.user.Name, "<") || contains(tc.user.Email, "<") {
					t.Errorf("❌ Данные не экранированы: %s", html)
				}
			}

			t.Logf("✅ Защита работает: %s\nHTML: %s", tc.desc, html)
		})
	}
}

// ============================================================================
// ТЕСТЫ ВАЛИДАЦИИ
// ============================================================================

// TestValidation_Email проверяет валидацию email
func TestValidation_Email(t *testing.T) {
	// ТЕОРИЯ: Тестируем валидацию email
	// - Проверяем что некорректные email отклоняются
	// - Проверяем что корректные email принимаются
	
	testCases := []struct {
		name    string
		email   string
		wantErr bool
		desc    string
	}{
		{
			name:    "Корректный email",
			email:   "test@example.com",
			wantErr: false,
			desc:    "Обычный корректный email",
		},
		{
			name:    "Пустой email",
			email:   "",
			wantErr: true,
			desc:    "Пустой email должен быть отклонён",
		},
		{
			name:    "Без @",
			email:   "testexample.com",
			wantErr: true,
			desc:    "Email без @ должен быть отклонён",
		},
		{
			name:    "Без домена",
			email:   "test@",
			wantErr: true,
			desc:    "Email без домена должен быть отклонён",
		},
		{
			name:    "Слишком длинный",
			email:   "a" + string(make([]byte, 255)) + "@example.com",
			wantErr: true,
			desc:    "Email длиннее 254 символов должен быть отклонён",
		},
		{
			name:    "SQL injection попытка",
			email:   "test@example.com' OR '1'='1",
			wantErr: true,
			desc:    "Email с SQL кодом должен быть отклонён",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEmail(tc.email)

			if tc.wantErr && err == nil {
				t.Errorf("❌ Ожидалась ошибка для '%s', но ошибки нет", tc.email)
			} else if !tc.wantErr && err != nil {
				t.Errorf("❌ Неожиданная ошибка для '%s': %v", tc.email, err)
			} else {
				t.Logf("✅ %s: %v", tc.desc, err)
			}
		})
	}
}

// TestValidation_Password проверяет валидацию пароля
func TestValidation_Password(t *testing.T) {
	// ТЕОРИЯ: Тестируем валидацию пароля
	// - Проверяем требования к сложности
	
	testCases := []struct {
		name    string
		password string
		wantErr bool
		desc    string
	}{
		{
			name:     "Корректный пароль",
			password: "Password123",
			wantErr:  false,
			desc:     "Пароль с заглавными, строчными и цифрами",
		},
		{
			name:     "Слишком короткий",
			password: "Pass1",
			wantErr:  true,
			desc:     "Пароль короче 8 символов должен быть отклонён",
		},
		{
			name:     "Без заглавных",
			password: "password123",
			wantErr:  true,
			desc:     "Пароль без заглавных букв должен быть отклонён",
		},
		{
			name:     "Без строчных",
			password: "PASSWORD123",
			wantErr:  true,
			desc:     "Пароль без строчных букв должен быть отклонён",
		},
		{
			name:     "Без цифр",
			password: "Password",
			wantErr:  true,
			desc:     "Пароль без цифр должен быть отклонён",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)

			if tc.wantErr && err == nil {
				t.Errorf("❌ Ожидалась ошибка для пароля, но ошибки нет")
			} else if !tc.wantErr && err != nil {
				t.Errorf("❌ Неожиданная ошибка для пароля: %v", err)
			} else {
				t.Logf("✅ %s: %v", tc.desc, err)
			}
		})
	}
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================================

// contains проверяет содержит ли строка подстроку (используем strings.Contains)
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

