package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("Задачи для практической работы на языке Go")
	fmt.Println("Введите номер раздела")
	fmt.Println("1. Задачи на линейное программирование ")
	fmt.Println("2. Задачи с условным оператором")
	fmt.Println("3. Задачи на циклы")
	var part, task int
	_, _ = fmt.Scan(&part)
	fmt.Println("Введите номер задачи от 1 до 5")
	_, _ = fmt.Scan(&task)
	switch part {
	case 1: switch task {
	case 1: task_1_1()
	case 2: task_1_2()
	case 3: task_1_3()
	case 4: task_1_4()
	case 5: task_1_5()
	default: fmt.Print("Некорректный номер задачи")
	}
	case 2: switch task {
	case 1: task_2_1()
	case 2: task_2_2()
	case 3: task_2_3()
	case 4: task_2_4() //0-12 - ребенок, 13-18 - подросток, 19-65 - взрослый, >65 - пожилой
	case 5: task_2_5()
	default: fmt.Print("Некорректный номер задачи")
	}
	case 3: switch task {
	case 1: task_3_1()
	case 2: task_3_2()
	case 3: task_3_3()
	case 4: task_3_4()
	case 5: task_3_5()
	default: fmt.Print("Некорректный номер задачи")
	}
	default: fmt.Print("Некорректный номер раздела")
	}
}

func summator(n int) int {
	if n == 0 {
		return 0
	}
	return n % 10 + summator(n / 10)
}

func task_1_1() {
	fmt.Println("1. Сумма цифр числа")
	fmt.Println("Введите число")
	var n int
	_, _ = fmt.Scan(&n)
	fmt.Print(summator(n))
}

func toFahrenheit(n float64) float64 {
	return n * 9 / 5 + 32
}

func toCelsius(n float64) float64 {
	return (n - 32) * 5 / 9
}

func task_1_2() {
	fmt.Println("2. Преобразование температуры")
	fmt.Println("Введите температуру и единицы измерения (Celsius)/(Fahrenheit)")
	var (
		n float64
		s string
	)
	fmt.Scan(&n, &s)
	switch s {
	case "(Celsius)": 
		fmt.Print(toFahrenheit(n), " (Fahrenheit)")
	case "(Fahrenheit)":
		fmt.Print(toCelsius(n), " (Celsius)")
	}
}

func doubleElement(slice []int) {
	slice[0] = slice[0] * 2
	if len(slice) > 1 {
		doubleElement(slice[1:])
	}
}

func task_1_3() {
	fmt.Println("3. Удвоение каждого элемента массива")
	fmt.Println("Введите количество чисел в массиве , затем введите сами числа")
	var n int
	_, _ = fmt.Scan(&n)
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&slice[i])		
	}
	doubleElement(slice)
	fmt.Print(slice)
}

func stringConcat(slice []string) string {
	if len(slice) == 1 {
		return slice[0]
	}
	return slice[0] + " " + stringConcat(slice[1:])
}

func task_1_4() {
	fmt.Println("4. Объединение строк")
	fmt.Println("Введите количество строк, затем сами строки")
	var n int
	_, _ = fmt.Scan(&n)
	slice := make([]string, n)
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&slice[i])		
	}
	fmt.Print(stringConcat(slice))

}

func distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)))
}

func task_1_5() {
	fmt.Println("5. Расчет расстояния между двумя точками")
	fmt.Println("Введите координаты")
	var x1, y1, x2, y2 int
	_, _ = fmt.Scan(&x1, &y1, &x2, &y2)
	fmt.Print(distance(x1, y1, x2, y2))
}

func isEven(n int) string {
	if n % 2 == 0 {
		return "Четное"
	} else {
		return "Нечетное"
	}
}

func task_2_1() {
	fmt.Println("1. Проверка на четность/нечетность")
	var n int
	fmt.Println("Введите число")
	_, _ = fmt.Scan(&n)
	fmt.Print(isEven(n))
}

func leapYear(year int) string{
	if year % 400 == 0 || year % 4 == 0 && year % 100 != 0 {
		return "Високосный"
	} else {
		return "Не високосный"
	}
}

func task_2_2() {
	fmt.Println("2. Проверка высокосного года")
	var year int
	fmt.Println("Введите год")
	_, _ = fmt.Scan(&year)
	fmt.Print(leapYear(year))
}

func biggestOfThreeNumbers(a, b, c int) int {
	if a > b && a > c {
		return a
	} else if b > c {
		return b
	} else {
		return c
	}
}

func task_2_3() {
	fmt.Println("3. Определение наибольшего из трех чисел")
	var a, b, c int
	fmt.Println("Введите три числа")
	_, _ = fmt.Scan(&a, &b, &c)
	fmt.Print(biggestOfThreeNumbers(a, b, c))
}

func categoryOfAge(n int) string {
	if n < 12 {
		return "Ребенок"
	} else if n < 19 {
		return "Подросток"
	} else if n < 66 {
		return "Взрослый"
	} else {
		return "Пожилой"
	}
}

func task_2_4() {
	fmt.Print("4. Категория возраста")
	var n int
	fmt.Println("Введите возраст")
	_, _ = fmt.Scan(&n)
	fmt.Print(categoryOfAge(n))
}

func divisibleBy3And5(n int) string {
	if n % 3 == 0 && n % 5 == 0 {
		return "Делится"
	} else {
		return "Не делится"
	}
}

func task_2_5() {
	fmt.Println("5. Проверка делимости на 3 и 5")
	var n int
	fmt.Println("Введите число")
	_, _ = fmt.Scan(&n)
	fmt.Println(divisibleBy3And5(n))
}

func factorial(n int) int {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return result
}

func task_3_1() {
	fmt.Println("1. Факториал числа")
	var n int
	fmt.Println("Введите число")
	_, _ = fmt.Scan(&n)
	fmt.Print(factorial(n))
}

func fib(n int) {
	var fib1, fib2 int = 0, 1
	if n == 1 {
		fmt.Print(fib1)
	} else if n == 2 {
		fmt.Print(fib1, ", ", fib2)
	} else {
		fmt.Print(fib1, ", ", fib2)
		for i := 3; i <= n; i++ {
			fib1, fib2 = fib2, fib1 + fib2
			fmt.Print(", ", fib2)
		}
	}
}

func task_3_2() {
	fmt.Println("2. Числа Фибоначчи")
	var n int
	fmt.Println("Введите число")
	_, _ = fmt.Scan(&n)
	fib(n)
}

func reverseArray(slice []int) {
	l := len(slice)
	for i := 0; i < l / 2; i++ {
		slice[i], slice[l - i - 1] = slice[l - i - 1], slice[i]
	}
}

func task_3_3() {
	fmt.Println("3. Реверс массива")
	var n int
	fmt.Println("Введите количество элементов в массиве")
	_, _ = fmt.Scan(&n)
	slice := make([]int, n)
	fmt.Println("Введите элементы массива")
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&slice[i])
	}
	reverseArray(slice)
	fmt.Print(slice)

}

func isPrime(n int) bool {
	result := true
	for i := 2; i < n; i++ {
		if n % i == 0 {
			result = false
			break
		}
	}
	return result
}

func task_3_4() {
	fmt.Println("4. Поиск простых чисел")
	fmt.Println("Введите число")
	var n int
	_, _ = fmt.Scan(&n)
	for i := 2; i <= n; i++ {
		if isPrime(i) {
			fmt.Print(i, " ")
		}
	}
}

func sumOfArr(slice []int) int {
	sum := 0
	l := len(slice)
	for i := 0; i < l; i++ {
		sum += slice[i]
	}
	return sum
}

func task_3_5() {
	fmt.Println("5. Сумма чисел в массиве")
	var n int
	fmt.Println("Введите количество элементов в массиве")
	_, _ = fmt.Scan(&n)
	slice := make([]int, n)
	fmt.Println("Введите элементы массива")
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&slice[i])
	}
	fmt.Print(sumOfArr(slice))
}
