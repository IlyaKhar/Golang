package tests

// ТЕОРИЯ ПО TESTIFY И GOMEGA
//
// ПОЧЕМУ НУЖНЫ ФРЕЙМВОРКИ:
// 1. Стандартный testing.T подходит для базовых проверок, но не даёт выразительных матчеров.
// 2. Testify и Gomega закрывают разные боли:
//    - Testify = набор assert/require + test suites (структурированный подход).
//    - Gomega = BDD-style матчеры + fluent-синтаксис (часто используется вместе с Ginkgo).
//
// РАЗБИВКА ПО ФРЕЙМВОРКАМ
//
// --- Testify ---------------------------------------------------------------
// Основные компоненты:
//   • assert.*  – мягкие проверки (тест продолжится после ошибки).
//   • require.* – жёсткие проверки (тест сразу упадёт).
//   • suite.Suite – структура для группировки связанных тестов + хуки Setup/TearDown.
//
// Интеграция:
//   go get github.com/stretchr/testify
//   import (
//       "testing"
//       "github.com/stretchr/testify/assert"
//       "github.com/stretchr/testify/require"
//       "github.com/stretchr/testify/suite"
//   )
//
// TODO-БЛОКИ ДЛЯ ТЕБЯ (Testify):
// 1) // TODO testify/suite: создай math_suite_test.go и определи структуру:
//
// ТЕОРИЯ: Test Suite (тестовый набор) - это способ группировать связанные тесты
// - suite.Suite встраивается в структуру
// - Можно добавить свои поля для хранения данных между тестами
// - SetupTest() вызывается ПЕРЕД КАЖДЫМ тестом - там готовим данные
// - TearDownTest() вызывается ПОСЛЕ КАЖДОГО теста - там очищаем данные
//
// ПРАВИЛЬНАЯ СТРУКТУРА:
//      type MathSuite struct {
//          suite.Suite  // встраиваем suite.Suite - получаем доступ к assert/require
//          a, b int     // свои поля для хранения данных
//      }
//
// ТЕОРИЯ: SetupTest() - это хук (хук = функция которая вызывается автоматически)
// - Вызывается ПЕРЕД каждым тестом (TestSum, TestMultiply и т.д.)
// - Используется для подготовки данных: инициализация переменных, создание моков, подключение к БД
// - Гарантирует что каждый тест начинается с "чистого листа"
//
//      func (s *MathSuite) SetupTest() {
//          // ТЕОРИЯ: s - это указатель на MathSuite
//          // - s.a и s.b - это поля структуры MathSuite
//          // - s.Require() - это метод из suite.Suite для проверок
//          s.a = 2  // инициализируем поле a
//          s.b = 3  // инициализируем поле b
//          // Теперь в каждом тесте s.a = 2 и s.b = 3
//      }
//
// ТЕОРИЯ: TestSum() - это тест
// - Название должно начинаться с Test
// - Использует данные из SetupTest() (s.a и s.b)
// - s.Require() - жёсткая проверка (если не пройдёт - тест сразу упадёт)
//
//      func (s *MathSuite) TestSum() {
//          result := Sum(s.a, s.b)  // используем данные из SetupTest()
//          s.Require().Equal(5, result)  // проверяем результат
//      }
//
// ТЕОРИЯ: TestMathSuite - это точка входа для запуска suite
// - suite.Run() запускает все тесты из MathSuite
// - Автоматически вызывает SetupTest() перед каждым тестом
// - Автоматически вызывает TearDownTest() после каждого теста
//
//      func TestMathSuite(t *testing.T) {
//          suite.Run(t, new(MathSuite))  // запускаем все тесты из MathSuite
//      }
//
// 2) // TODO testify/assert: напиши обычный тест без suite:
//      func TestDivide(t *testing.T) {
//          t.Parallel()
//          res, err := Divide(10, 2)
//          require.NoError(t, err)
//          assert.Equal(t, 5, res)
//      }
//
// 3) // TODO testify/mock: создай интерфейс Storage и мок через mock.Mock:
//      type Storage interface { Save(context.Context, User) error }
//      type storageMock struct { mock.Mock }
//      func (m *storageMock) Save(ctx context.Context, u User) error {
//          args := m.Called(ctx, u)
//          return args.Error(0)
//      }
//      // В тесте:
//      st := new(storageMock)
//      st.On("Save", mock.Anything, user).Return(nil)
//      err := service.Create(ctx, user, st)
//      st.AssertExpectations(t)
//
// 4) // TODO testify/require: проверь ошибку:
//      require.ErrorIs(t, err, ErrInvalidInput)
//
// --- Gomega ---------------------------------------------------------------
// Особенности:
//   • Поставляется с набором матчеров (Equal, HaveOccurred, ContainSubstring, etc.).
//   • Поддерживает асинхронные ожидания (Eventually/Consistently).
//   • Можно использовать поверх стандартного testing.T (через gomega.NewWithT) или вместе с Ginkgo.
//
// Интеграция:
//   go get github.com/onsi/gomega
//   import (
//       "testing"
//       . "github.com/onsi/gomega"
//   )
//
// TODO-БЛОКИ ДЛЯ ТЕБЯ (Gomega):
// 1) // TODO gomega/basic: создай math_gomega_test.go:
//      func TestMultiply(t *testing.T) {
//          g := NewWithT(t)
//          g.Expect(Multiply(2, 4)).To(Equal(8))
//          g.Expect(Multiply(2, 4)).NotTo(Equal(7))
//      }
//
// 2) // TODO gomega/matchers: покажи более «человеческие» сообщения:
//      g.Expect(len(users)).To(BeNumerically("==", 3))
//      g.Expect(user.Email).To(ContainSubstring("@"))
//
// 3) // TODO gomega/eventually: протестируй асинхронную операцию:
//      g.Eventually(func() []Order {
//          return repo.PendingOrders()
//      }, 2*time.Second, 100*time.Millisecond).Should(HaveLen(0))
//
// 4) // TODO gomega/combined: комбинируй матчеры:
//      g.Expect(response.Code).To(And(
//          BeNumerically(">=", 200),
//          BeNumerically("<", 300),
//      ))
//
//
// ============================================================================
// ПОЛНЫЙ РАБОЧИЙ ПРИМЕР ДЛЯ ПРОВЕРКИ
// ============================================================================
//
// Скопируй этот код в файл tests/math_suite_test.go и запусти: go test ./tests
//
// ПРИМЕР 1: Testify Suite с SetupTest()
// ```go
// package tests
//
// import (
//     "testing"
//     "github.com/stretchr/testify/suite"
// )
//
// // Простая функция для тестирования (в реальном проекте это была бы функция из другого пакета)
// func Sum(a, b int) int {
//     return a + b
// }
//
// // ТЕОРИЯ: MathSuite - это структура для группировки тестов
// // - suite.Suite встраивается - получаем доступ к s.Assert(), s.Require() и т.д.
// // - a, b - это наши поля для хранения данных
// type MathSuite struct {
//     suite.Suite
//     a, b int
// }
//
// // ТЕОРИЯ: SetupTest() вызывается ПЕРЕД КАЖДЫМ тестом автоматически
// // - Здесь мы готовим данные для тестов
// // - Гарантирует что каждый тест начинается с "чистого листа"
// func (s *MathSuite) SetupTest() {
//     s.a = 2  // инициализируем a = 2
//     s.b = 3  // инициализируем b = 3
//     // Теперь в TestSum() будут доступны s.a = 2 и s.b = 3
// }
//
// // ТЕОРИЯ: TestSum - это тест
// // - Использует данные из SetupTest() (s.a и s.b)
// // - s.Require() - жёсткая проверка (если не пройдёт - тест упадёт сразу)
// func (s *MathSuite) TestSum() {
//     result := Sum(s.a, s.b)  // Sum(2, 3) = 5
//     s.Require().Equal(5, result)  // проверяем что результат = 5
// }
//
// // ТЕОРИЯ: TestMathSuite - точка входа для запуска suite
// // - suite.Run() запускает все тесты из MathSuite
// // - Автоматически вызывает SetupTest() перед каждым тестом
// func TestMathSuite(t *testing.T) {
//     suite.Run(t, new(MathSuite))
// }
// ```
//
// ПРИМЕР 2: Обычный тест с Testify (без suite)
// ```go
// package tests
//
// import (
//     "errors"
//     "testing"
//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/require"
// )
//
// func Divide(a, b float64) (float64, error) {
//     if b == 0 {
//         return 0, errors.New("деление на ноль")
//     }
//     return a / b, nil
// }
//
// func TestDivide(t *testing.T) {
//     t.Parallel()  // тест может выполняться параллельно с другими
//
//     // ТЕОРИЯ: require.NoError - жёсткая проверка
//     // - Если ошибка есть - тест упадёт сразу
//     res, err := Divide(10, 2)
//     require.NoError(t, err)  // проверяем что ошибки нет
//
//     // ТЕОРИЯ: assert.Equal - мягкая проверка
//     // - Если не пройдёт - тест отметится как проваленный, но продолжит выполнение
//     assert.Equal(t, 5.0, res)  // проверяем результат
// }
// ```
//
// ПРИМЕР 3: Gomega
// ```go
// package tests
//
// import (
//     "testing"
//     . "github.com/onsi/gomega"
// )
//
// func Multiply(a, b int) int {
//     return a * b
// }
//
// func TestMultiply_WithGomega(t *testing.T) {
//     // ТЕОРИЯ: NewWithT(t) создаёт Gomega для работы с testing.T
//     // - g.Expect() - это начало цепочки проверок
//     g := NewWithT(t)
//
//     // ТЕОРИЯ: Fluent-синтаксис Gomega
//     // - g.Expect(значение).To(матчер) - читается как "ожидаю что значение соответствует матчеру"
//     g.Expect(Multiply(2, 4)).To(Equal(8))  // ожидаю что 2*4 = 8
//     g.Expect(Multiply(2, 4)).NotTo(Equal(7))  // ожидаю что 2*4 НЕ равно 7
// }
// ```
//
// Полезные команды:
//   go test ./...                               — прогнать все тесты.
//   go test ./... -run TestName                 — запустить конкретный тест.
//   go test ./tests -run Test.*Gomega           — запустить группу Gomega-тестов.
//
// Следующий шаг:
//   1. Установи зависимости: go get github.com/stretchr/testify github.com/onsi/gomega
//   2. Создай реальные файлы *_test.go с указанными примерами.
//   3. Расширь примеры под свои кейсы (HTTP, репозитории, domain сервисы).
//
// Profit: выразительные тесты + лучшая диагностика падений.
