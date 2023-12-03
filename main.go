package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Calculator struct {
	romanNumerals  map[string]int
	allowedNumbers map[int]bool
	allowRoman     bool
}

func NewCalculator(allowRoman bool) *Calculator {
	return &Calculator{
		romanNumerals:  map[string]int{"I": 1, "II": 2, "III": 3, "IV": 4, "V": 5, "VI": 6, "VII": 7, "VIII": 8, "IX": 9, "X": 10},
		allowedNumbers: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true},
		allowRoman:     allowRoman,
	}
}

func (calc *Calculator) isRoman(s string) bool {
	for _, char := range s {
		if _, ok := calc.romanNumerals[string(char)]; !ok {
			return false
		}
	}
	return true
}

func (calc *Calculator) toArabic(s string) (int, error) {
	if arabic, ok := calc.romanNumerals[s]; ok {
		return arabic, nil
	}
	return 0, fmt.Errorf("неверное римское число")
}

func (calc *Calculator) toRoman(num int) string {
	if num <= 0 {
		panic("Результат работы с римскими числами не может быть меньше единицы")
	}

	romanNumerals := []struct {
		Value  int
		Symbol string
	}{
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	result := ""
	for _, numeral := range romanNumerals {
		for num >= numeral.Value {
			result += numeral.Symbol
			num -= numeral.Value
		}
	}

	return result
}

func (calc *Calculator) calculate(a, b int, operator string) (int, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("неверная операция")
	}
}

func (calc *Calculator) runCalculator() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите выражение (например, III + V или 3 + 5): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	re := regexp.MustCompile(`^(\d+|[IVX]+)\s*([-+*/])\s*(\d+|[IVX]+)$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		fmt.Println("Неверный формат ввода")
		return
	}

	a, operator, b := matches[1], matches[2], matches[3]

	isARoman := calc.isRoman(a)
	isBRoman := calc.isRoman(b)

	if (isARoman && isBRoman) && !calc.allowRoman {
		fmt.Println("Калькулятор работает только с арабскими числами")
		return
	}

	var operandA, operandB int
	var err error

	if isARoman {
		operandA, err = calc.toArabic(a)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		operandA, err = strconv.Atoi(a)
		if err != nil {
			fmt.Println("Ошибка ввода числа a")
			return
		}
	}

	if isBRoman {
		operandB, err = calc.toArabic(b)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		operandB, err = strconv.Atoi(b)
		if err != nil {
			fmt.Println("Ошибка ввода числа b")
			return
		}
	}

	if (isARoman && !isBRoman) || (!isARoman && isBRoman) {
		fmt.Println("Калькулятор работает только с арабскими или только с римскими числами")
		return
	}

	if _, ok := calc.allowedNumbers[operandA]; !ok {
		fmt.Println("Число 'a' должно быть от 1 до 10 включительно")
		return
	}

	if _, ok := calc.allowedNumbers[operandB]; !ok {
		fmt.Println("Число 'b' должно быть от 1 до 10 включительно")
		return
	}

	result, err := calc.calculate(operandA, operandB, operator)
	if err != nil {
		fmt.Println(err)
		return
	}

	if isARoman || isBRoman {
		fmt.Println("Результат:", calc.toRoman(result))
	} else {
		fmt.Println("Результат:", result)
	}
}

func main() {
	// Если вы хотите работать только с арабскими цифрами, передайте false
	// Если вы хотите работать только с римскими цифрами, передайте true
	calculator := NewCalculator(true)
	calculator.runCalculator()
}
