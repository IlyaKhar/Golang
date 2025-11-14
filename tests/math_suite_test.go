package tests

import (
    "testing"
    "github.com/stretchr/testify/suite"
)
//
// // Простая функция для тестирования (в реальном проекте это была бы функция из другого пакета)
func Sum(a, b int) int {
    return a + b
}

// // ТЕОРИЯ: MathSuite - это структура для группировки тестов
// // - suite.Suite встраивается - получаем доступ к s.Assert(), s.Require() и т.д.
// // - a, b - это наши поля для хранения данных


type MathSuite struct {
    suite.Suite
    a, b int
}

// // ТЕОРИЯ: SetupTest() вызывается ПЕРЕД КАЖДЫМ тестом автоматически
// // - Здесь мы готовим данные для тестов
// // - Гарантирует что каждый тест начинается с "чистого листа"


func (s *MathSuite) SetupTest() {
    s.a = 2  // инициализируем a = 2
    s.b = 3  // инициализируем b = 3
    // Теперь в TestSum() будут доступны s.a = 2 и s.b = 3
}

// // ТЕОРИЯ: TestSum - это тест
// // - Использует данные из SetupTest() (s.a и s.b)
// // - s.Require() - жёсткая проверка (если не пройдёт - тест упадёт сразу)


func (s *MathSuite) TestSum() {
    result := Sum(s.a, s.b)  // Sum(2, 3) = 5
    s.Require().Equal(5, result)  // проверяем что результат = 5
}

// // ТЕОРИЯ: TestMathSuite - точка входа для запуска suite
// // - suite.Run() запускает все тесты из MathSuite
// // - Автоматически вызывает SetupTest() перед каждым тестом

func TestMathSuite(t *testing.T) {
    suite.Run(t, new(MathSuite))
}