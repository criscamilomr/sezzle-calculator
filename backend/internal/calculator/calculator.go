package calculator

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Operation string

const (
	Add      Operation = "add"
	Subtract Operation = "subtract"
	Multiply Operation = "multiply"
	Divide   Operation = "divide"
	Exponent Operation = "exponent"
	Percent  Operation = "percent"
	Sqrt     Operation = "sqrt"
)

var (
	ErrDivisionByZero    = errors.New("division by zero is not allowed")
	ErrNegativeSqrt      = errors.New("cannot calculate square root of a negative number")
	ErrInvalidOperation  = errors.New("invalid operation")
	ErrMissingNumbers    = errors.New("invalid number of operands for operation")
	ErrInvalidExpression = errors.New("invalid mathematical expression")
	ErrOverflow          = errors.New("result is out of range or undefined")
)

const floatPrecision = 1e10

// Calculate performs a single operation across one (for Sqrt or unary Percent) or two-or-more numbers.
func Calculate(operation Operation, numbers []float64) (float64, error) {
	if operation == Sqrt {
		if len(numbers) != 1 {
			return 0, ErrMissingNumbers
		}
		if numbers[0] < 0 {
			return 0, ErrNegativeSqrt
		}
		res := math.Sqrt(numbers[0])
		if err := validateNumber(res); err != nil {
			return 0, err
		}
		return normalizeFloat(res), nil
	}

	if operation == Percent && len(numbers) == 1 {
		res := numbers[0] / 100
		if err := validateNumber(res); err != nil {
			return 0, err
		}
		return normalizeFloat(res), nil
	}

	if !isBinaryOperation(operation) {
		return 0, ErrInvalidOperation
	}

	if len(numbers) < 2 {
		return 0, ErrMissingNumbers
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		var err error
		result, err = applyBinaryOperation(operation, result, numbers[i])
		if err != nil {
			return 0, err
		}
	}

	return normalizeFloat(result), nil
}

// CalculateMixed evaluates a sequence of numbers and binary operations (e.g., numbers: [5, 2, 3], operations: [Add, Multiply] -> 5 + 2 * 3 = 11)
// respecting standard mathematical operator precedence (PEMDAS/BODMAS).
func CalculateMixed(numbers []float64, operations []Operation) (float64, error) {
	if len(numbers) < 2 || len(operations) != len(numbers)-1 {
		return 0, ErrMissingNumbers
	}

	for _, op := range operations {
		if !isBinaryOperation(op) {
			return 0, ErrInvalidOperation
		}
	}

	nums := make([]float64, len(numbers))
	copy(nums, numbers)
	ops := make([]Operation, len(operations))
	copy(ops, operations)

	// Pass 1: Exponentiation (highest binary precedence)
	var err error
	nums, ops, err = reducePrecedenceLevel(nums, ops, func(op Operation) bool {
		return op == Exponent
	})
	if err != nil {
		return 0, err
	}

	// Pass 2: Multiplication, Division, Percentage
	nums, ops, err = reducePrecedenceLevel(nums, ops, func(op Operation) bool {
		return op == Multiply || op == Divide || op == Percent
	})
	if err != nil {
		return 0, err
	}

	// Pass 3: Addition, Subtraction
	nums, _, err = reducePrecedenceLevel(nums, ops, func(op Operation) bool {
		return op == Add || op == Subtract
	})
	if err != nil {
		return 0, err
	}

	return normalizeFloat(nums[0]), nil
}

// EvaluateExpression parses and evaluates a mathematical expression string (e.g., "50×10%", "200-15%", "200+15%", "5+2×3", "1+6+9+5+3", "√(144)+5²")
// with full operator precedence, parentheses, unary operators, and percentage calculations.
func EvaluateExpression(expr string) (float64, error) {
	normalized := normalizeExpressionSymbols(expr)
	if strings.TrimSpace(normalized) == "" {
		return 0, ErrInvalidExpression
	}

	p := &exprParser{input: normalized}
	result, err := p.parseExpression()
	if err != nil {
		return 0, err
	}

	p.skipWhitespace()
	if p.pos < len(p.input) {
		return 0, ErrInvalidExpression
	}

	if err := validateNumber(result); err != nil {
		return 0, err
	}

	return normalizeFloat(result), nil
}

func isBinaryOperation(op Operation) bool {
	switch op {
	case Add, Subtract, Multiply, Divide, Exponent, Percent:
		return true
	default:
		return false
	}
}

func applyBinaryOperation(op Operation, a, b float64) (float64, error) {
	var res float64
	switch op {
	case Add:
		res = a + b
	case Subtract:
		res = a - b
	case Multiply:
		res = a * b
	case Divide:
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		res = a / b
	case Exponent:
		res = math.Pow(a, b)
	case Percent:
		// Calculates b% of a (e.g., 50 * (10 / 100) = 5)
		res = (a * b) / 100
	default:
		return 0, ErrInvalidOperation
	}

	if err := validateNumber(res); err != nil {
		return 0, err
	}
	return res, nil
}

func reducePrecedenceLevel(
	nums []float64,
	ops []Operation,
	matches func(Operation) bool,
) ([]float64, []Operation, error) {
	nextNums := []float64{nums[0]}
	var nextOps []Operation

	for i, op := range ops {
		right := nums[i+1]
		if matches(op) {
			left := nextNums[len(nextNums)-1]
			val, err := applyBinaryOperation(op, left, right)
			if err != nil {
				return nil, nil, err
			}
			nextNums[len(nextNums)-1] = val
		} else {
			nextNums = append(nextNums, right)
			nextOps = append(nextOps, op)
		}
	}

	return nextNums, nextOps, nil
}

func validateNumber(val float64) error {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		return ErrOverflow
	}
	return nil
}

// normalizeFloat eliminates standard IEEE 754 floating-point noise (e.g., 0.1 + 0.2 = 0.30000000000000004).
func normalizeFloat(val float64) float64 {
	if val == 0 {
		return 0 // Avoid returning -0
	}
	return math.Round(val*floatPrecision) / floatPrecision
}

var superscriptDigits = map[rune]byte{
	'⁰': '0',
	'¹': '1',
	'²': '2',
	'³': '3',
	'⁴': '4',
	'⁵': '5',
	'⁶': '6',
	'⁷': '7',
	'⁸': '8',
	'⁹': '9',
}

// normalizeExpressionSymbols replaces visual math symbols (×, ÷, −, √, and superscripts like 5²) with ASCII equivalents for parsing.
func normalizeExpressionSymbols(expr string) string {
	replacer := strings.NewReplacer(
		"×", "*",
		"÷", "/",
		"−", "-",
		"√", "sqrt",
	)
	replaced := replacer.Replace(expr)

	var b strings.Builder
	b.Grow(len(replaced) + 4)
	inSuperscript := false

	for _, r := range replaced {
		if digit, ok := superscriptDigits[r]; ok {
			if !inSuperscript {
				b.WriteByte('^')
				inSuperscript = true
			}
			b.WriteByte(digit)
		} else {
			inSuperscript = false
			b.WriteRune(r)
		}
	}

	return b.String()
}

// Recursive descent parser for mathematical expressions.
// Grammar:
//   Expression := Term ( ('+' | '-') Term )*
//                 Note: If a right-hand Term is a standalone percentage (e.g. 200 - 15% or 200 + 15%),
//                 it computes a discount/increment relative to the left value: left ± (left * right).
//   Term       := Power ( ('*' | '/') Power )*
//   Power      := Unary ( '^' Unary )*
//   Unary      := ('+' | '-') Unary | 'sqrt' Unary | Postfix
//   Postfix    := Primary ('%')*   -> converts X% into X / 100 (e.g. 10% -> 0.1, so 50 * 10% = 5)
//   Primary    := Number | '(' Expression ')'
type exprParser struct {
	input string
	pos   int
}

func (p *exprParser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *exprParser) parseExpression() (float64, error) {
	left, _, err := p.parseTerm()
	if err != nil {
		return 0, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}
		ch := p.input[p.pos]
		if ch != '+' && ch != '-' {
			break
		}
		p.pos++

		right, rightIsBarePercent, err := p.parseTerm()
		if err != nil {
			return 0, err
		}

		// If the user wrote "A + B%" or "A - B%" (discount / markup), B% is relative to A (i.e. A * (B/100)).
		operand := right
		if rightIsBarePercent {
			operand = left * right
		}

		op := Add
		if ch == '-' {
			op = Subtract
		}
		left, err = applyBinaryOperation(op, left, operand)
		if err != nil {
			return 0, err
		}
	}

	return left, nil
}

func (p *exprParser) parseTerm() (float64, bool, error) {
	left, isBarePercent, err := p.parsePower()
	if err != nil {
		return 0, false, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) {
			break
		}
		ch := p.input[p.pos]
		if ch != '*' && ch != '/' {
			break
		}
		p.pos++

		right, _, err := p.parsePower()
		if err != nil {
			return 0, false, err
		}

		var op Operation
		if ch == '*' {
			op = Multiply
		} else {
			op = Divide
		}

		left, err = applyBinaryOperation(op, left, right)
		if err != nil {
			return 0, false, err
		}
		// Once combined in multiplication/division (e.g., 50 * 10% = 5), the result is already a proportional value, not a bare %
		isBarePercent = false
	}

	return left, isBarePercent, nil
}

func (p *exprParser) parsePower() (float64, bool, error) {
	left, isBarePercent, err := p.parseUnary()
	if err != nil {
		return 0, false, err
	}

	for {
		p.skipWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != '^' {
			break
		}
		p.pos++

		right, _, err := p.parseUnary()
		if err != nil {
			return 0, false, err
		}

		left, err = applyBinaryOperation(Exponent, left, right)
		if err != nil {
			return 0, false, err
		}
		isBarePercent = false
	}

	return left, isBarePercent, nil
}

func (p *exprParser) parseUnary() (float64, bool, error) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return 0, false, ErrInvalidExpression
	}

	if p.input[p.pos] == '+' {
		p.pos++
		val, isPct, err := p.parseUnary()
		return val, isPct, err
	}

	if p.input[p.pos] == '-' {
		p.pos++
		val, isPct, err := p.parseUnary()
		if err != nil {
			return 0, false, err
		}
		return -val, isPct, nil
	}

	if strings.HasPrefix(p.input[p.pos:], "sqrt") {
		p.pos += 4
		val, _, err := p.parseUnary()
		if err != nil {
			return 0, false, err
		}
		if val < 0 {
			return 0, false, ErrNegativeSqrt
		}
		res := math.Sqrt(val)
		if err := validateNumber(res); err != nil {
			return 0, false, err
		}
		return res, false, nil
	}

	return p.parsePostfix()
}

func (p *exprParser) parsePostfix() (float64, bool, error) {
	val, err := p.parsePrimary()
	if err != nil {
		return 0, false, err
	}

	isPercent := false
	for {
		p.skipWhitespace()
		if p.pos < len(p.input) && p.input[p.pos] == '%' {
			p.pos++
			val = val / 100.0
			isPercent = true
			if err := validateNumber(val); err != nil {
				return 0, false, err
			}
		} else {
			break
		}
	}

	return val, isPercent, nil
}

func (p *exprParser) parsePrimary() (float64, error) {
	p.skipWhitespace()
	if p.pos >= len(p.input) {
		return 0, ErrInvalidExpression
	}

	if p.input[p.pos] == '(' {
		p.pos++
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return 0, ErrInvalidExpression
		}
		p.pos++
		return val, nil
	}

	start := p.pos
	hasDot := false
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == '.' {
			if hasDot {
				return 0, ErrInvalidExpression
			}
			hasDot = true
			p.pos++
		} else if ch >= '0' && ch <= '9' {
			p.pos++
		} else {
			break
		}
	}

	if start == p.pos || p.input[start:p.pos] == "." {
		return 0, ErrInvalidExpression
	}

	val, err := strconv.ParseFloat(p.input[start:p.pos], 64)
	if err != nil {
		return 0, ErrInvalidExpression
	}

	return val, nil
}
