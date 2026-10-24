package calculator

import (
	"errors"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		numbers   []float64
		want      float64
		wantErr   error
	}{
		// --- Addition (2 and multiple numbers) ---
		{
			name:      "add two positive numbers",
			operation: Add,
			numbers:   []float64{10, 5},
			want:      15,
			wantErr:   nil,
		},
		{
			name:      "add multiple numbers (1 + 6 + 9 + 5 + 3 = 24)",
			operation: Add,
			numbers:   []float64{1, 6, 9, 5, 3},
			want:      24,
			wantErr:   nil,
		},
		{
			name:      "add negative and positive numbers",
			operation: Add,
			numbers:   []float64{-12.5, 7.5, 2},
			want:      -3,
			wantErr:   nil,
		},
		{
			name:      "add floating point precision (0.1 + 0.2 = 0.3)",
			operation: Add,
			numbers:   []float64{0.1, 0.2},
			want:      0.3,
			wantErr:   nil,
		},
		{
			name:      "add with missing operand returns ErrMissingNumbers",
			operation: Add,
			numbers:   []float64{10},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Subtraction (2 and multiple numbers) ---
		{
			name:      "subtract two numbers",
			operation: Subtract,
			numbers:   []float64{20, 8},
			want:      12,
			wantErr:   nil,
		},
		{
			name:      "subtract multiple numbers (100 - 20 - 15 - 5 = 60)",
			operation: Subtract,
			numbers:   []float64{100, 20, 15, 5},
			want:      60,
			wantErr:   nil,
		},
		{
			name:      "subtract resulting in negative number",
			operation: Subtract,
			numbers:   []float64{5, 12},
			want:      -7,
			wantErr:   nil,
		},
		{
			name:      "subtract floating point precision (0.3 - 0.1 = 0.2)",
			operation: Subtract,
			numbers:   []float64{0.3, 0.1},
			want:      0.2,
			wantErr:   nil,
		},
		{
			name:      "subtract with missing operand returns ErrMissingNumbers",
			operation: Subtract,
			numbers:   []float64{10},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Multiplication (2 and multiple numbers) ---
		{
			name:      "multiply two positive numbers",
			operation: Multiply,
			numbers:   []float64{6, 7},
			want:      42,
			wantErr:   nil,
		},
		{
			name:      "multiply multiple numbers (2 * 3 * 4 * 5 = 120)",
			operation: Multiply,
			numbers:   []float64{2, 3, 4, 5},
			want:      120,
			wantErr:   nil,
		},
		{
			name:      "multiply negative by zero normalizes -0 to 0",
			operation: Multiply,
			numbers:   []float64{-15, 0},
			want:      0,
			wantErr:   nil,
		},
		{
			name:      "multiply floating point precision (0.1 * 0.2 = 0.02)",
			operation: Multiply,
			numbers:   []float64{0.1, 0.2},
			want:      0.02,
			wantErr:   nil,
		},
		{
			name:      "multiply with missing operand returns ErrMissingNumbers",
			operation: Multiply,
			numbers:   []float64{6},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Division (2 and multiple numbers) ---
		{
			name:      "divide two numbers evenly",
			operation: Divide,
			numbers:   []float64{20, 4},
			want:      5,
			wantErr:   nil,
		},
		{
			name:      "divide multiple numbers (100 / 2 / 5 = 10)",
			operation: Divide,
			numbers:   []float64{100, 2, 5},
			want:      10,
			wantErr:   nil,
		},
		{
			name:      "divide resulting in decimal",
			operation: Divide,
			numbers:   []float64{7, 2},
			want:      3.5,
			wantErr:   nil,
		},
		{
			name:      "divide repeating decimal rounds cleanly to 10 decimal places",
			operation: Divide,
			numbers:   []float64{1, 3},
			want:      0.3333333333,
			wantErr:   nil,
		},
		{
			name:      "divide zero by non-zero number",
			operation: Divide,
			numbers:   []float64{0, 5},
			want:      0,
			wantErr:   nil,
		},
		{
			name:      "divide by zero returns ErrDivisionByZero",
			operation: Divide,
			numbers:   []float64{10, 0},
			want:      0,
			wantErr:   ErrDivisionByZero,
		},
		{
			name:      "divide by zero in third operand returns ErrDivisionByZero",
			operation: Divide,
			numbers:   []float64{100, 5, 0},
			want:      0,
			wantErr:   ErrDivisionByZero,
		},
		{
			name:      "divide with missing operand returns ErrMissingNumbers",
			operation: Divide,
			numbers:   []float64{10},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Exponentiation ---
		{
			name:      "exponent positive base and power",
			operation: Exponent,
			numbers:   []float64{2, 8},
			want:      256,
			wantErr:   nil,
		},
		{
			name:      "exponent multiple numbers ((2 ^ 3) ^ 2 = 64)",
			operation: Exponent,
			numbers:   []float64{2, 3, 2},
			want:      64,
			wantErr:   nil,
		},
		{
			name:      "exponent to the power of zero",
			operation: Exponent,
			numbers:   []float64{99, 0},
			want:      1,
			wantErr:   nil,
		},
		{
			name:      "exponent with negative power",
			operation: Exponent,
			numbers:   []float64{2, -2},
			want:      0.25,
			wantErr:   nil,
		},
		{
			name:      "exponent with missing operand returns ErrMissingNumbers",
			operation: Exponent,
			numbers:   []float64{2},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Percentage ---
		{
			name:      "convert single number into fraction of 100 (50% = 0.5)",
			operation: Percent,
			numbers:   []float64{50},
			want:      0.5,
			wantErr:   nil,
		},
		{
			name:      "calculate proportional percentage of a number (10% of 50 = 5)",
			operation: Percent,
			numbers:   []float64{50, 10},
			want:      5,
			wantErr:   nil,
		},
		{
			name:      "calculate fractional percentage (12.5% of 80 = 10)",
			operation: Percent,
			numbers:   []float64{80, 12.5},
			want:      10,
			wantErr:   nil,
		},
		{
			name:      "percent with zero operands returns ErrMissingNumbers",
			operation: Percent,
			numbers:   []float64{},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Square Root ---
		{
			name:      "sqrt of perfect square",
			operation: Sqrt,
			numbers:   []float64{81},
			want:      9,
			wantErr:   nil,
		},
		{
			name:      "sqrt of zero",
			operation: Sqrt,
			numbers:   []float64{0},
			want:      0,
			wantErr:   nil,
		},
		{
			name:      "sqrt of non-perfect square rounded to precision",
			operation: Sqrt,
			numbers:   []float64{2},
			want:      1.4142135624,
			wantErr:   nil,
		},
		{
			name:      "sqrt of negative number returns ErrNegativeSqrt",
			operation: Sqrt,
			numbers:   []float64{-9},
			want:      0,
			wantErr:   ErrNegativeSqrt,
		},
		{
			name:      "sqrt with zero operands returns ErrMissingNumbers",
			operation: Sqrt,
			numbers:   []float64{},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},
		{
			name:      "sqrt with two operands returns ErrMissingNumbers",
			operation: Sqrt,
			numbers:   []float64{9, 2},
			want:      0,
			wantErr:   ErrMissingNumbers,
		},

		// --- Invalid Operation & Overflow / NaN Edge Cases ---
		{
			name:      "unknown operation returns ErrInvalidOperation",
			operation: Operation("modulo"),
			numbers:   []float64{10, 3},
			want:      0,
			wantErr:   ErrInvalidOperation,
		},
		{
			name:      "empty operation returns ErrInvalidOperation",
			operation: Operation(""),
			numbers:   []float64{10, 3},
			want:      0,
			wantErr:   ErrInvalidOperation,
		},
		{
			name:      "multiplication overflow (+Inf) returns ErrOverflow",
			operation: Multiply,
			numbers:   []float64{1e308, 10},
			want:      0,
			wantErr:   ErrOverflow,
		},
		{
			name:      "exponent overflow (+Inf) returns ErrOverflow",
			operation: Exponent,
			numbers:   []float64{10, 1000},
			want:      0,
			wantErr:   ErrOverflow,
		},
		{
			name:      "zero to negative power (+Inf) returns ErrOverflow",
			operation: Exponent,
			numbers:   []float64{0, -1},
			want:      0,
			wantErr:   ErrOverflow,
		},
		{
			name:      "negative base to fractional power (NaN) returns ErrOverflow",
			operation: Exponent,
			numbers:   []float64{-4, 0.5},
			want:      0,
			wantErr:   ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.operation, tt.numbers)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Calculate(%q, %v) error = %v, wantErr %v", tt.operation, tt.numbers, err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("Calculate(%q, %v) = %v, want %v", tt.operation, tt.numbers, got, tt.want)
			}
		})
	}
}

func TestCalculateMixed(t *testing.T) {
	tests := []struct {
		name       string
		numbers    []float64
		operations []Operation
		want       float64
		wantErr    error
	}{
		{
			name:       "mixed addition and multiplication respects precedence (5 + 2 * 3 = 11)",
			numbers:    []float64{5, 2, 3},
			operations: []Operation{Add, Multiply},
			want:       11,
			wantErr:    nil,
		},
		{
			name:       "chained additions (1 + 6 + 9 + 5 + 3 = 24)",
			numbers:    []float64{1, 6, 9, 5, 3},
			operations: []Operation{Add, Add, Add, Add},
			want:       24,
			wantErr:    nil,
		},
		{
			name:       "all precedence levels combined (20 - 12 / 3 + 2 ^ 3 = 24)",
			numbers:    []float64{20, 12, 3, 2, 3},
			operations: []Operation{Subtract, Divide, Add, Exponent},
			want:       24,
			wantErr:    nil,
		},
		{
			name:       "mixed with percentage (200 % 15 + 10 = 40)",
			numbers:    []float64{200, 15, 10},
			operations: []Operation{Percent, Add},
			want:       40,
			wantErr:    nil,
		},
		{
			name:       "floating-point precision in mixed operations (0.1 + 0.2 * 2 = 0.5)",
			numbers:    []float64{0.1, 0.2, 2},
			operations: []Operation{Add, Multiply},
			want:       0.5,
			wantErr:    nil,
		},
		{
			name:       "fewer than 2 numbers returns ErrMissingNumbers",
			numbers:    []float64{5},
			operations: []Operation{},
			want:       0,
			wantErr:    ErrMissingNumbers,
		},
		{
			name:       "mismatched operations count returns ErrMissingNumbers",
			numbers:    []float64{5, 2, 3},
			operations: []Operation{Add},
			want:       0,
			wantErr:    ErrMissingNumbers,
		},
		{
			name:       "unary sqrt in binary operations list returns ErrInvalidOperation",
			numbers:    []float64{5, 9},
			operations: []Operation{Sqrt},
			want:       0,
			wantErr:    ErrInvalidOperation,
		},
		{
			name:       "division by zero in mixed chain returns ErrDivisionByZero",
			numbers:    []float64{5, 10, 0},
			operations: []Operation{Add, Divide},
			want:       0,
			wantErr:    ErrDivisionByZero,
		},
		{
			name:       "exponent overflow in mixed chain returns ErrOverflow",
			numbers:    []float64{5, 10, 1000},
			operations: []Operation{Add, Exponent},
			want:       0,
			wantErr:    ErrOverflow,
		},
		{
			name:       "addition overflow in mixed chain returns ErrOverflow",
			numbers:    []float64{1e308, 1e308},
			operations: []Operation{Add},
			want:       0,
			wantErr:    ErrOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateMixed(tt.numbers, tt.operations)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CalculateMixed(%v, %v) error = %v, wantErr %v", tt.numbers, tt.operations, err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("CalculateMixed(%v, %v) = %v, want %v", tt.numbers, tt.operations, got, tt.want)
			}
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		want    float64
		wantErr error
	}{
		// --- User Screenshot Scenarios & Multi-Number Mixed Operations ---
		{
			name:    "mixed precedence with unicode multiplication (5+2×3 = 11)",
			expr:    "5+2×3",
			want:    11,
			wantErr: nil,
		},
		{
			name:    "mixed precedence with ASCII multiplication (5 + 2 * 3 = 11)",
			expr:    "5 + 2 * 3",
			want:    11,
			wantErr: nil,
		},
		{
			name:    "multi-number addition chain (1+6+9+5+3 = 24)",
			expr:    "1+6+9+5+3",
			want:    24,
			wantErr: nil,
		},
		{
			name:    "parentheses override precedence ((5+2)*3 = 21)",
			expr:    "(5+2)*3",
			want:    21,
			wantErr: nil,
		},
		{
			name:    "complex mixed expression with unicode symbols (20 − 12 ÷ 3 + 2^3 = 24)",
			expr:    "20 − 12 ÷ 3 + 2^3",
			want:    24,
			wantErr: nil,
		},
		{
			name:    "square root symbol and function in expression (√144 + sqrt(81) * 2 = 30)",
			expr:    "√144 + sqrt(81) * 2",
			want:    30,
			wantErr: nil,
		},
		{
			name:    "superscript exponent notation (5² = 25)",
			expr:    "5²",
			want:    25,
			wantErr: nil,
		},
		{
			name:    "multi-digit superscript exponent and sqrt with parentheses (√(144) + 2¹⁰ = 1036)",
			expr:    "√(144) + 2¹⁰",
			want:    1036,
			wantErr: nil,
		},
		{
			name:    "convert number to fraction of 100 (10% = 0.1)",
			expr:    "10%",
			want:    0.1,
			wantErr: nil,
		},
		{
			name:    "proportional percentage multiplication as in screenshot (50×10% = 5)",
			expr:    "50×10%",
			want:    5,
			wantErr: nil,
		},
		{
			name:    "proportional percentage division (50 ÷ 10% = 500)",
			expr:    "50 ÷ 10%",
			want:    500,
			wantErr: nil,
		},
		{
			name:    "percentage discount calculation (200 − 15% = 170)",
			expr:    "200 − 15%",
			want:    170,
			wantErr: nil,
		},
		{
			name:    "percentage increment/markup calculation (200 + 15% = 230)",
			expr:    "200 + 15%",
			want:    230,
			wantErr: nil,
		},
		{
			name:    "mixed expression combining proportional percentage and addition (50×10% + 20 = 25)",
			expr:    "50×10% + 20",
			want:    25,
			wantErr: nil,
		},
		{
			name:    "unary minus and plus (-5 + +12 - -3 = 10)",
			expr:    "-5 + +12 - -3",
			want:    10,
			wantErr: nil,
		},
		{
			name:    "floating-point precision in expression (0.1 + 0.2 = 0.3)",
			expr:    "0.1 + 0.2",
			want:    0.3,
			wantErr: nil,
		},

		// --- Edge Cases & Error Handling ---
		{
			name:    "division by zero in expression returns ErrDivisionByZero",
			expr:    "5 + 10 / (3 - 3)",
			want:    0,
			wantErr: ErrDivisionByZero,
		},
		{
			name:    "negative square root in expression returns ErrNegativeSqrt",
			expr:    "10 + √(-9)",
			want:    0,
			wantErr: ErrNegativeSqrt,
		},
		{
			name:    "exponent overflow in expression returns ErrOverflow",
			expr:    "10 ^ 1000",
			want:    0,
			wantErr: ErrOverflow,
		},
		{
			name:    "addition overflow in expression returns ErrOverflow",
			expr:    "1e308 + 1e308",
			want:    0,
			wantErr: ErrInvalidExpression, // scientific notation not used in keypad
		},
		{
			name:    "empty expression returns ErrInvalidExpression",
			expr:    "   ",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "trailing operator returns ErrInvalidExpression",
			expr:    "5 + 2 *",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "unclosed parenthesis returns ErrInvalidExpression",
			expr:    "(5 + 2 * 3",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "unmatched closing parenthesis returns ErrInvalidExpression",
			expr:    "5 + 2) * 3",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "multiple decimal points in number returns ErrInvalidExpression",
			expr:    "5.2.3 + 1",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "lone decimal point returns ErrInvalidExpression",
			expr:    ". + 1",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "invalid operand inside sqrt returns ErrInvalidExpression",
			expr:    "sqrt(*)",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "invalid operand after unary minus returns ErrInvalidExpression",
			expr:    "-*",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
		{
			name:    "invalid operand after power operator returns ErrInvalidExpression",
			expr:    "2 ^ *",
			want:    0,
			wantErr: ErrInvalidExpression,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvaluateExpression(tt.expr)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("EvaluateExpression(%q) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("EvaluateExpression(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
