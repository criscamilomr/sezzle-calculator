package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/criscamilomr/sezzle-calculator/backend/internal/calculator"
)

const maxBodyBytes = 64 << 10 // 64 KB limit to prevent large payload abuse

// Error codes returned in ErrorResponse for programmatic frontend handling.
const (
	CodeInvalidJSON       = "INVALID_JSON"
	CodeInvalidOperation  = "INVALID_OPERATION"
	CodeInvalidExpression = "INVALID_EXPRESSION"
	CodeMissingNumbers    = "MISSING_NUMBERS"
	CodeDivisionByZero    = "DIVISION_BY_ZERO"
	CodeNegativeSqrt      = "NEGATIVE_SQRT"
	CodeOverflow          = "NUMERIC_OVERFLOW"
	CodeInternalError     = "INTERNAL_ERROR"
)

type CalculateRequest struct {
	Expression string                 `json:"expression,omitempty"`
	Operation  calculator.Operation   `json:"operation,omitempty"`
	Operations []calculator.Operation `json:"operations,omitempty"`
	Numbers    []float64              `json:"numbers,omitempty"`
}

type CalculateResponse struct {
	Result     float64              `json:"result"`
	Operation  calculator.Operation `json:"operation,omitempty"`
	Expression string               `json:"expression,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)
	mux.HandleFunc("GET /api/v1/health", h.Health)
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	defer r.Body.Close()

	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload", CodeInvalidJSON)
		return
	}

	// Ensure the body contains only a single JSON object.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request body must contain a single JSON object", CodeInvalidJSON)
		return
	}

	var (
		result float64
		err    error
		opName calculator.Operation
	)

	switch {
	case req.Expression != "":
		result, err = calculator.EvaluateExpression(req.Expression)
		opName = calculator.Operation("expression")
	case len(req.Operations) > 0:
		result, err = calculator.CalculateMixed(req.Numbers, req.Operations)
		opName = calculator.Operation("mixed")
	default:
		result, err = calculator.Calculate(req.Operation, req.Numbers)
		opName = req.Operation
	}

	if err != nil {
		status, code := mapDomainError(err)
		writeError(w, status, err.Error(), code)
		return
	}

	writeJSON(w, http.StatusOK, CalculateResponse{
		Result:     result,
		Operation:  opName,
		Expression: req.Expression,
	})
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, calculator.ErrInvalidOperation):
		return http.StatusBadRequest, CodeInvalidOperation
	case errors.Is(err, calculator.ErrInvalidExpression):
		return http.StatusBadRequest, CodeInvalidExpression
	case errors.Is(err, calculator.ErrMissingNumbers):
		return http.StatusBadRequest, CodeMissingNumbers
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusUnprocessableEntity, CodeDivisionByZero
	case errors.Is(err, calculator.ErrNegativeSqrt):
		return http.StatusUnprocessableEntity, CodeNegativeSqrt
	case errors.Is(err, calculator.ErrOverflow):
		return http.StatusUnprocessableEntity, CodeOverflow
	default:
		return http.StatusInternalServerError, CodeInternalError
	}
}

func writeError(w http.ResponseWriter, status int, message string, code string) {
	writeJSON(w, status, ErrorResponse{
		Error: message,
		Code:  code,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
