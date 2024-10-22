package main

import (
	"bufio"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
	"strings"
	"unicode"
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
	case 1:
		switch task {
		case 1:
			task_1_1()
		case 2:
			task_1_2()
		case 3:
			task_1_3()
		case 4:
			task_1_4()
		case 5:
			task_1_5()
		default:
			fmt.Print("Некорректный номер задачи")
		}
	case 2:
		switch task {
		case 1:
			task_2_1()
		case 2:
			task_2_2()
		case 3: task_2_3()
		case 4: task_2_4()
		case 5:
			task_2_5()
		default:
			fmt.Print("Некорректный номер задачи")
		}
	case 3:
		switch task {
		case 1:
			task_3_1()
		case 2:
			task_3_2()
		case 3:
			task_3_3()
		case 4:
			task_3_4()
		case 5:
			task_3_5()
		default:
			fmt.Print("Некорректный номер задачи")
		}
	default:
		fmt.Print("Некорректный номер раздела")
	}
}

func to10(num string, base int) int {
	m := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"A": 10,
		"B": 11,
		"C": 12,
		"D": 13,
		"E": 14,
		"F": 15,
		"G": 16,
		"H": 17,
		"I": 18,
		"J": 19,
		"K": 20,
		"L": 21,
		"M": 22,
		"N": 23,
		"O": 24,
		"P": 25,
		"Q": 26,
		"R": 27,
		"S": 28,
		"T": 29,
		"U": 30,
		"V": 31,
		"W": 32,
		"X": 33,
		"Y": 34,
		"Z": 35,
	}
	var res int
	for index, value := range num {
		place := 1
		for i := 0; i < len(num)-1-index; i++ {
			place *= base
		}
		res += m[string(value)] * place
	}
	return res
}

func from10(num int, base int) string {
	m := map[int]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
		16: "G",
		17: "H",
		18: "I",
		19: "J",
		20: "K",
		21: "L",
		22: "M",
		23: "N",
		24: "O",
		25: "P",
		26: "Q",
		27: "R",
		28: "S",
		29: "T",
		30: "U",
		31: "V",
		32: "W",
		33: "X",
		34: "Y",
		35: "Z",
	}
	res := ""
	for num != 0 {
		res = m[num%base] + res
		num /= base
	}
	return res
}

func task_1_1() {
	fmt.Println("1. Перевод чисел из одной системы счисления в другую")
	fmt.Println("Введите число, исходную систему, конечную систему")
	var (
		num          string
		base1, base2 int
	)
	_, _ = fmt.Scan(&num, &base1, &base2)
	fmt.Println(from10(to10(num, base1), base2))
}

func abc(a, b, c float64) {
	d := b*b - 4*a*c
	sqrt_d := cmplx.Sqrt(complex(d, 0))
	x1 := (-complex(b, 0) + sqrt_d) / (2 * complex(a, 0))
	x2 := (-complex(b, 0) - sqrt_d) / (2 * complex(a, 0))
	if d >= 0 {
		fmt.Println(real(x1), real(x2))
	} else {
		fmt.Println(x1, x2)
	}
}

func task_1_2() {
	fmt.Println("2. Решение квадратного уравнения")
	fmt.Println("коэффициенты a, b, c.")
	var a, b, c float64
	_, _ = fmt.Scan(&a, &b, &c)
	abc(a, b, c)
}

func task_1_3() {
	fmt.Println("3. Сортировка чисел по модулю")
	fmt.Println("Введите количество элементов в массиве")
	var n int
	_, _ = fmt.Scan(&n)
	fmt.Println("Введите массив чисел")
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&a[i])
	}
	sort.Slice(a, func(i, j int) bool {
		return math.Abs(float64(a[i])) < math.Abs(float64(a[j]))
	})
	fmt.Println(a)
}

func mergeSortedArrays(arr1 []int, arr2 []int) []int {
	var merged []int
	i, j := 0, 0

	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			merged = append(merged, arr1[i])
			i++
		} else {
			merged = append(merged, arr2[j])
			j++
		}
	}
	for i < len(arr1) {
		merged = append(merged, arr1[i])
		i++
	}
	for j < len(arr2) {
		merged = append(merged, arr2[j])
		j++
	}
	return merged
}

func task_1_4() {
	var m, n int
	fmt.Println("Введите количество элементов первого массива")
	_, _ = fmt.Scan(&m)
	fmt.Println("Введите первый массив")
	arr1 := make([]int, m)
	for i := 0; i < m; i++ {
		_, _ = fmt.Scan(&arr1[i])
	}
	fmt.Println("Введите количество элементов второго массива")
	_, _ = fmt.Scan(&n)
	fmt.Println("Введите второй массив")
	arr2 := make([]int, n)
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&arr2[i])
	}
	arr := mergeSortedArrays(arr1, arr2)
	fmt.Println(arr)
}

func task_1_5() {
	var s, sub_s string
	fmt.Println("5. Нахождение подстроки в строке без использования встроенных функций ")
	fmt.Println("Введите две строки")
	_, _ = fmt.Scan(&s, &sub_s)
	index := -1
	rune_s := []rune(s)
	rune_sub := []rune(sub_s)
	for i := 0; i < len(rune_s)-len(rune_sub)+1; i++ {
		flag := true
		for j := 0; j < len(rune_sub); j++ {
			if rune_s[j+i] != rune_sub[j] {
				flag = false
			}
		}
		if flag {
			index = i
			break
		}
	}
	fmt.Print(index)
}

func calc(num1, num2 float64, operator string) {
	switch operator {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		switch {
		case num2 == 0:
			fmt.Println("На ноль делить нельзя")
		default:
			fmt.Println(num1 / num2)
		}
	default:
		fmt.Println("Некорректная операция")
	}
}

func task_2_1() {
	var (
		num1, num2 float64
		operator   string
	)
	fmt.Println("1. Калькулятор с расширенными операциями")
	fmt.Println("Введите два числа и оператор")
	fmt.Scan(&num1, &num2, &operator)
	calc(num1, num2, operator)
}

func leapYear(year int) string {
	if year%400 == 0 || year%4 == 0 && year%100 != 0 {
		return "Високосный"
	} else {
		return "Не високосный"
	}
}

func task_2_2() {
	fmt.Println("2. Проверка палиндрома")
	fmt.Println("Введите строку")
	scan := bufio.NewScanner(os.Stdin)
	_ = scan.Scan()
	_ = scan.Scan()
	s := scan.Text()
	letter_s := ""
	for _, l := range s {
		if unicode.IsLetter(l) {
			letter_s += strings.ToLower(string(l))
		}
	}
	r_s := []rune(letter_s)
	flag := true
	for i := 0; i < len(r_s)/2; i++ {
		if r_s[i] != r_s[len(r_s)-i-1] {
			flag = false
		}
	}
	if flag {
		fmt.Println("Палиндром")
	} else {
		fmt.Println("Не палиндром")
	}
}

func intersect(a1, a2, b1, b2 int) bool {
	return a1 < b2 && a1 > b1 || b1 > a1 && b1 < a2
}

func task_2_3() {
	fmt.Println("3. Нахождение пересечения трех отрезков")
	fmt.Println("Введите три пары чисел, задающих отрезки")
	var a1, a2, b1, b2, c1, c2 int
	_, _ = fmt.Scan(&a1, &a2, &b1, &b2, &c1, &c2)
	if intersect(a1, a2, b1, b2) && intersect(a1, a2, c1, c2) && intersect(b1, b2, c1, c2) {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}
	
}

func task_2_4() {
	var max_word, word string
	fmt.Println("4. Выбор самого длинного слова в предложении")
	fmt.Println("Введите строку")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scanner.Scan()
	s := scanner.Text()
	for _, l := range s {
		if unicode.IsLetter(l) {
			word += string(l)
		} else {
			if len(word) > len(max_word) {
				max_word = word
			}
			word = ""
		}
	}
	fmt.Println(max_word)
}

func task_2_5() {
	fmt.Println("5. Проверка высокосного года")
	var year int
	fmt.Println("Введите год")
	_, _ = fmt.Scan(&year)
	fmt.Print(leapYear(year))
}

func fib(n int) {
	n1 := 0
	n2 := 1
	fmt.Print(n1)
	if n > 1 {
		for i := 2; i <= n; i++ {
			n1, n2 = n2, n1+n2
			fmt.Print(", ", n1)
		}
	}
}

func task_3_1() {
	var n int
	fmt.Println("1. Числа Фибоначчи до определенного числа")
	fmt.Println("Введите целое число")
	_, _ = fmt.Scan(&n)
	fib(n)
}

func isPrime(n int) bool {
	flag := true
	for i := 2; i < n; i++ {
		if n%i == 0 {
			flag = false
		}
	}
	return flag
}

func task_3_2() {
	fmt.Println("2. Определение простых чисел в диапазоне")
	fmt.Println("Введите два числа")
	var a, b int
	_, _ = fmt.Scan(&a, &b)
	for i := a; i <= b; i++ {
		if isPrime(i) {
			fmt.Print(i, " ")
		}
	}
}

func sumDigits(n int) int {
	sum := 0
	for n != 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

func countOfDigits(n int) int {
	cnt := 0
	for n != 0 {
		cnt++
		n /= 10
	}
	return cnt
}

func pow(n, p int) int {
	res := 1
	for i := 1; i <= p; i++ {
		res *= n
	}
	return res
}

func arm(n int) bool {
	return n == pow(sumDigits(n), countOfDigits(n))
}

func task_3_3() {
	fmt.Println("3. Числа Армстронга в заданном диапазоне")
	fmt.Println("Введите два числа, задающие диапазон")
	var a, b int
	_, _ = fmt.Scan(&a, &b)
	for i := a; i <= b; i++ {
		if arm(i) {
			fmt.Print(i, " ")
		}
	}
}

func rev(s string) string {
	a := []rune(s)
	l := len(a)
	rev_s := ""
	for i := 0; i < l; i++ {
		rev_s += string(a[l-i-1])
	}
	return rev_s

}

func task_3_4() {
	var s string
	fmt.Println("4. Реверс строки")
	fmt.Println("Введите строку")
	_, _ = fmt.Scan(&s)
	fmt.Print(rev(s))
}

func nod(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func task_3_5() {
	fmt.Println("5. Нахождение наибольшего общего делителя (НОД)")
	fmt.Println("Введите два числа")
	var a, b int
	_, _ = fmt.Scan(&a, &b)
	fmt.Print(nod(a, b))
}
