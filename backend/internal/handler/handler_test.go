package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/criscamilomr/sezzle-calculator/backend/internal/calculator"
)

func TestCalculateHandler(t *testing.T) {
	h := NewHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name          string
		payload       string
		wantStatus    int
		wantResult    float64
		wantOperation calculator.Operation
		wantErrorCode string
	}{
		// --- Happy Paths: Single Operation (2 and Multiple Numbers) ---
		{
			name:          "valid addition of two numbers returns 200 OK",
			payload:       `{"operation": "add", "numbers": [10, 5]}`,
			wantStatus:    http.StatusOK,
			wantResult:    15,
			wantOperation: calculator.Add,
		},
		{
			name:          "valid addition of multiple numbers (1+6+9+5+3) returns 24",
			payload:       `{"operation": "add", "numbers": [1, 6, 9, 5, 3]}`,
			wantStatus:    http.StatusOK,
			wantResult:    24,
			wantOperation: calculator.Add,
		},
		{
			name:          "addition with floating-point precision returns 0.3",
			payload:       `{"operation": "add", "numbers": [0.1, 0.2]}`,
			wantStatus:    http.StatusOK,
			wantResult:    0.3,
			wantOperation: calculator.Add,
		},
		{
			name:          "valid subtraction returns 200 OK",
			payload:       `{"operation": "subtract", "numbers": [20, 8.5]}`,
			wantStatus:    http.StatusOK,
			wantResult:    11.5,
			wantOperation: calculator.Subtract,
		},
		{
			name:          "valid multiplication returns 200 OK",
			payload:       `{"operation": "multiply", "numbers": [6, 7]}`,
			wantStatus:    http.StatusOK,
			wantResult:    42,
			wantOperation: calculator.Multiply,
		},
		{
			name:          "valid division returns 200 OK",
			payload:       `{"operation": "divide", "numbers": [25, 4]}`,
			wantStatus:    http.StatusOK,
			wantResult:    6.25,
			wantOperation: calculator.Divide,
		},
		{
			name:          "valid exponentiation returns 200 OK",
			payload:       `{"operation": "exponent", "numbers": [2, 10]}`,
			wantStatus:    http.StatusOK,
			wantResult:    1024,
			wantOperation: calculator.Exponent,
		},
		{
			name:          "valid percentage returns 200 OK",
			payload:       `{"operation": "percent", "numbers": [200, 15]}`,
			wantStatus:    http.StatusOK,
			wantResult:    30,
			wantOperation: calculator.Percent,
		},
		{
			name:          "valid square root returns 200 OK",
			payload:       `{"operation": "sqrt", "numbers": [144]}`,
			wantStatus:    http.StatusOK,
			wantResult:    12,
			wantOperation: calculator.Sqrt,
		},

		// --- Happy Paths: Mixed Operations Array (PEMDAS Precedence) ---
		{
			name:          "mixed operations array (5 + 2 * 3) returns 11",
			payload:       `{"numbers": [5, 2, 3], "operations": ["add", "multiply"]}`,
			wantStatus:    http.StatusOK,
			wantResult:    11,
			wantOperation: calculator.Operation("mixed"),
		},
		{
			name:          "mixed operations array (20 - 12 / 3 + 2 ^ 3) returns 24",
			payload:       `{"numbers": [20, 12, 3, 2, 3], "operations": ["subtract", "divide", "add", "exponent"]}`,
			wantStatus:    http.StatusOK,
			wantResult:    24,
			wantOperation: calculator.Operation("mixed"),
		},

		// --- Happy Paths: Full Mathematical Expression Strings ---
		{
			name:          "expression string 5+2×3 returns 11",
			payload:       `{"expression": "5+2×3"}`,
			wantStatus:    http.StatusOK,
			wantResult:    11,
			wantOperation: calculator.Operation("expression"),
		},
		{
			name:          "expression string 1+6+9+5+3 returns 24",
			payload:       `{"expression": "1+6+9+5+3"}`,
			wantStatus:    http.StatusOK,
			wantResult:    24,
			wantOperation: calculator.Operation("expression"),
		},
		{
			name:          "expression string with parentheses and sqrt returns 30",
			payload:       `{"expression": "(5+5)*2 + √100"}`,
			wantStatus:    http.StatusOK,
			wantResult:    30,
			wantOperation: calculator.Operation("expression"),
		},
		{
			name:          "expression string proportional percentage 50×10% returns 5",
			payload:       `{"expression": "50×10%"}`,
			wantStatus:    http.StatusOK,
			wantResult:    5,
			wantOperation: calculator.Operation("expression"),
		},
		{
			name:          "expression string percentage discount 200-15% returns 170",
			payload:       `{"expression": "200-15%"}`,
			wantStatus:    http.StatusOK,
			wantResult:    170,
			wantOperation: calculator.Operation("expression"),
		},

		// --- JSON & Payload Validation Edge Cases (400 Bad Request) ---
		{
			name:          "empty request body returns 400 and INVALID_JSON",
			payload:       ``,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},
		{
			name:          "malformed JSON syntax returns 400 and INVALID_JSON",
			payload:       `{"operation": "add", "numbers": [10, 5}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},
		{
			name:          "wrong data type in numbers array returns 400 and INVALID_JSON",
			payload:       `{"operation": "add", "numbers": ["ten", 5]}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},
		{
			name:          "unknown field in JSON payload returns 400 and INVALID_JSON",
			payload:       `{"operation": "add", "numbers": [1, 2], "unexpected": true}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},
		{
			name:          "multiple JSON objects in body returns 400 and INVALID_JSON",
			payload:       `{"operation": "add", "numbers": [1, 2]}{"extra": "object"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},
		{
			name:          "payload exceeding 64KB limit returns 400 and INVALID_JSON",
			payload:       `{"operation": "add", "numbers": [1, 2], "pad": "` + strings.Repeat("a", 70000) + `"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidJSON,
		},

		// --- Client / Input Validation Errors (400 Bad Request) ---
		{
			name:          "unsupported operation returns 400 and INVALID_OPERATION",
			payload:       `{"operation": "modulo", "numbers": [10, 3]}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidOperation,
		},
		{
			name:          "invalid syntax in expression string returns 400 and INVALID_EXPRESSION",
			payload:       `{"expression": "5 + * 3"}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeInvalidExpression,
		},
		{
			name:          "missing operand for binary operation returns 400 and MISSING_NUMBERS",
			payload:       `{"operation": "add", "numbers": [10]}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeMissingNumbers,
		},
		{
			name:          "mismatched operations array length returns 400 and MISSING_NUMBERS",
			payload:       `{"numbers": [5, 2, 3], "operations": ["add"]}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeMissingNumbers,
		},
		{
			name:          "too many operands for unary operation returns 400 and MISSING_NUMBERS",
			payload:       `{"operation": "sqrt", "numbers": [9, 2]}`,
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: CodeMissingNumbers,
		},

		// --- Domain / Mathematical Errors (422 Unprocessable Entity) ---
		{
			name:          "division by zero returns 422 and DIVISION_BY_ZERO code",
			payload:       `{"operation": "divide", "numbers": [10, 0]}`,
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrorCode: CodeDivisionByZero,
		},
		{
			name:          "division by zero inside expression returns 422 and DIVISION_BY_ZERO code",
			payload:       `{"expression": "5 + 10 / 0"}`,
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrorCode: CodeDivisionByZero,
		},
		{
			name:          "negative square root returns 422 and NEGATIVE_SQRT code",
			payload:       `{"operation": "sqrt", "numbers": [-25]}`,
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrorCode: CodeNegativeSqrt,
		},
		{
			name:          "numeric overflow (+Inf) returns 422 and NUMERIC_OVERFLOW code",
			payload:       `{"operation": "exponent", "numbers": [10, 1000]}`,
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrorCode: CodeOverflow,
		},
		{
			name:          "indeterminate result (NaN) returns 422 and NUMERIC_OVERFLOW code",
			payload:       `{"operation": "exponent", "numbers": [-4, 0.5]}`,
			wantStatus:    http.StatusUnprocessableEntity,
			wantErrorCode: CodeOverflow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}

			contentType := rec.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				t.Errorf("Content-Type = %q, want application/json", contentType)
			}

			if tt.wantErrorCode != "" {
				var errResp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
					t.Fatalf("failed to decode ErrorResponse: %v", err)
				}
				if errResp.Code != tt.wantErrorCode {
					t.Errorf("error code = %q, want %q", errResp.Code, tt.wantErrorCode)
				}
				if errResp.Error == "" {
					t.Error("expected non-empty error message")
				}
				return
			}

			var resp CalculateResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode CalculateResponse: %v", err)
			}
			if resp.Result != tt.wantResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
			}
			if resp.Operation != tt.wantOperation {
				t.Errorf("operation = %q, want %q", resp.Operation, tt.wantOperation)
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	h := NewHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode HealthResponse: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("status = %q, want %q", resp.Status, "ok")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := NewHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestMapDomainErrorDefaultFallback(t *testing.T) {
	status, code := mapDomainError(errors.New("unexpected internal failure"))
	if status != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", status, http.StatusInternalServerError)
	}
	if code != CodeInternalError {
		t.Errorf("code = %q, want %q", code, CodeInternalError)
	}
}
