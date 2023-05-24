package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var IsDigits = regexp.MustCompile(`^[0-9]+$`).MatchString
var IsLetters = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
var userVars = make(map[string]int)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.Trim(scanner.Text(), " ")
		switch {
		case strings.Contains(input, "/exit"):
			fmt.Println("Bye!")
			return
		case strings.Contains(input, "/help"):
			fmt.Println("This program performs the addition and subtraction of numbers, use + or -")
		case strings.Trim(input, " ") == "":
			continue
		case strings.HasPrefix(input, "/"):
			fmt.Println("Unknown Command")
		case validateExpression(input):
			fmt.Println("Invalid Expression")
		case IsLetters(input):
			if userVars[input] == 0 {
				fmt.Println("Unknown variable")
			} else {
				fmt.Println(userVars[input])
			}
		case strings.Contains(input, "="):
			x := checkVariables(input)
			if x != "" {
				fmt.Println(x)
			}
		default:
			fmt.Println(calculator(input))
		}
	}
}

func validateExpression(input string) bool {
	letrs := strings.Split(input, " ")
	open, close := 0, 0
	for _, c := range input {
		if c == '(' {
			open++
		} else if c == ')' {
			close++
		}
	}
	if open != close {
		return true
	} else if strings.HasSuffix(input, "+") || strings.HasSuffix(input, "-") {
		return true
	}
	for _, char := range letrs {
		if len(char) > 1 {
			if strings.ContainsAny(char, "*/") {
				return true
			}
		}
	}
	return false
}

func checkVariables(input string) string {
	if strings.Count(input, "=") > 1 {
		return "Invalid identifier"
	}
	varSlice := strings.Split(input, "=")
	varName := strings.Trim(varSlice[0], " ")
	varVal := strings.Trim(varSlice[len(varSlice)-1], " ")
	varValNum, err := strconv.Atoi(varVal)
	if _, ok := userVars[varVal]; !ok && IsLetters(varVal) {
		return "Unknown variable"
	} else if IsLetters(varName) && err == nil {
		userVars[varName] = varValNum
	} else if IsLetters(varName) && (IsLetters(varVal) || IsDigits(varVal)) {
		userVars[varName] = userVars[varVal]
	} else {
		return "Invalid identifier"
	}
	return ""
}

func calculator(formula string) int {
	formula = strings.ReplaceAll(formula, " ", "")
	formula = strings.ReplaceAll(formula, "(-", "(0-")
	if formula[0] == '-' {
		formula = "0" + formula
	}
	for k, v := range userVars {
		formula = strings.ReplaceAll(formula, k, strconv.Itoa(v))
	}
	return evaluate(formula)
}

func evaluate(formula string) int {
	if !strings.ContainsAny(formula, "+-*/") {
		num, _ := strconv.Atoi(formula)
		return num
	}
	if strings.Contains(formula, "(") {
		start := strings.LastIndex(formula, "(")
		end := strings.Index(formula[start:], ")") + start
		return evaluate(formula[:start] + strconv.Itoa(evaluate(formula[start+1:end])) + formula[end+1:])
	}
	for _, op := range []string{"+-", "*/"} {
		for i := len(formula) - 1; i >= 0; i-- {
			if strings.ContainsRune(op, rune(formula[i])) {
				left := evaluate(formula[:i])
				right := evaluate(formula[i+1:])
				switch formula[i] {
				case '+':
					return left + right
				case '-':
					return left - right
				case '*':
					return left * right
				case '/':
					return left / right
				}
			}
		}
	}
	return 0
}
