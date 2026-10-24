# Sezzle Full-Stack Calculator

A full-stack calculator web application built with a **React + TypeScript** frontend and a **Go (Golang)** REST API microservice. Designed with clean architecture, strict input validation, comprehensive unit test suites, and full-stack Docker support.

---

## Table of Contents

1. [Features](#features)
2. [Architecture & Project Structure](#architecture--project-structure)
3. [Getting Started (Docker & Local Setup)](#getting-started-docker--local-setup)
4. [REST API Reference & Examples](#rest-api-reference--examples)
5. [Unit Tests & Coverage Report](#unit-tests--coverage-report)
6. [Design Decisions & Assumptions](#design-decisions--assumptions)
7. [AI Prompts Used](#ai-prompts-used)

---

## Features

### Mathematical Capabilities
- **Basic Operations**: Addition (`+`), Subtraction (`−`), Multiplication (`×`), Division (`÷`).
- **Advanced Operations (All Optional Requirements Included)**:
  - **Exponentiation (`xʸ` / `^`)**: Rendered visually with raised superscripts (e.g., `5²`, `2³`) on the display and history.
  - **Square Root (`√x`)**: Interactive `√()` parentheses entry mode allowing numbers or sub-expressions inside the radical (e.g., `√(144)`, `√(9+16)`).
  - **Percentage (`%`)**: Supports direct fraction conversion (`10% = 0.1`), proportional parts (`50 × 10% = 5`), discounts (`200 − 15% = 170`), and markups (`200 + 15% = 230`).
- **Multi-Operand & Mixed Expressions**: Evaluates chained and mixed operations of any length respecting standard **PEMDAS/BODMAS** operator precedence and parentheses `( )` (e.g., `5 + 2 × 3 = 11`, `(5 + 2) × 3 = 21`).
- **IEEE 754 Precision Normalization**: Eliminates floating-point representation noise (e.g., `0.1 + 0.2 = 0.3`) and guards against `NaN` / `±Inf` numeric overflow.

### User Interface & Experience
- **Full Keyboard Support**: Use numbers `0–9`, `.`, `+`, `-`, `*`, `/`, `^`, `%`, `(`, `)`, `R` for `√()`, `Enter` / `=` to evaluate, `Backspace`, and `Esc` / `C` to clear.
- **Persistent Calculation History**: Automatically saves the last 20 calculations in `localStorage` across page reloads, featuring smooth CSS entrance animations, a custom vertical scrollbar, and one-click result reuse.
- **Responsive Layout**: Adapts cleanly from a two-column desktop workspace to a single-column mobile layout.

---

## Architecture & Project Structure

```text
sezzle-calculator/
├── backend/                          # Go REST API Microservice (Zero external dependencies)
│   ├── cmd/api/main.go               # HTTP server entrypoint, CORS middleware, graceful shutdown
│   ├── internal/
│   │   ├── calculator/               # Pure domain layer: PEMDAS parser & math operations
│   │   │   ├── calculator.go
│   │   │   └── calculator_test.go    # Table-driven unit tests (96.3% coverage)
│   │   └── handler/                  # HTTP transport layer: JSON validation & error mapping
│   │       ├── handler.go
│   │       └── handler_test.go       # Endpoint integration tests (100.0% coverage)
│   ├── coverage/                     # Pre-generated Go coverage reports (HTML, TXT, OUT)
│   ├── api.http                      # Ready-to-run requests for VS Code REST Client
│   ├── Dockerfile
│   └── go.mod
├── frontend/                         # React 19 + TypeScript + Vite SPA
│   ├── src/
│   │   ├── services/                 # Typed HTTP client & API contract interfaces
│   │   │   ├── calculatorApi.ts
│   │   │   └── calculatorApi.test.ts # Service unit tests (100% coverage)
│   │   ├── hooks/                    # Calculator state machine & localStorage persistence
│   │   │   ├── useCalculator.ts
│   │   │   └── useCalculator.test.ts # Custom hook unit tests (90.9% line coverage)
│   │   ├── App.tsx                   # Presentational UI, superscript formatter, keyboard handler
│   │   └── App.test.tsx              # UI integration tests (100% line coverage)
│   ├── coverage/                     # Pre-generated Vitest V8 coverage reports (HTML & TXT)
│   ├── nginx.conf                    # Nginx SPA routing & /api reverse proxy for Docker
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml                # Full-stack container orchestration
├── TestsReport.md                    # Pre-generated unit test & coverage report summary
├── PROMPTS.md                        # Log of AI prompts used during development
└── README.md
```

---

## Getting Started (Docker & Local Setup)

### Option 1: Docker Compose (Recommended)

The fastest way to run both the Go backend and React/Nginx frontend together:

```bash
docker compose up --build
```

- **Frontend UI**: [http://localhost:5173](http://localhost:5173)
- **Backend API**: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)

To stop the containers:

```bash
docker compose down
```

### Option 2: Run Locally (Without Docker)

- **Backend** (requires Go `1.23+`):
  ```bash
  cd backend && go run ./cmd/api
  ```
- **Frontend** (requires Node.js `20+`):
  ```bash
  cd frontend && npm install && npm run dev
  ```

---

## REST API Reference & Examples

### 1. Health Check
- **Endpoint**: `GET /api/v1/health`
- **Response (`200 OK`)**:
  ```json
  {
    "status": "ok"
  }
  ```

### 2. Calculate Endpoint
- **Endpoint**: `POST /api/v1/calculate`
- **Headers**: `Content-Type: application/json`

The endpoint supports **three flexible payload modes**:

#### A. Expression Mode (Used by the React UI)
Evaluates complete mathematical strings with PEMDAS precedence, parentheses, `√()`, superscripts/exponents, and percentages:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "5 + 2 × 3"}'
```
**Response (`200 OK`)**:
```json
{
  "result": 11,
  "operation": "expression",
  "expression": "5 + 2 × 3"
}
```

**Percentage Examples (Fraction, Proportional, Discount, Markup)**:
```bash
# Proportional part: 10% of 50 -> 5
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "50×10%"}'

# Discount: 200 minus 15% -> 170
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "200-15%"}'

# Square Root & Exponent: √(144) + 5^2 -> 37
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "√(144)+5^2"}'
```

#### B. Single Operation Mode (2 or More Operands, or Unary `sqrt` / `percent`)
Supported operations: `"add"`, `"subtract"`, `"multiply"`, `"divide"`, `"exponent"`, `"percent"`, `"sqrt"`.

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "numbers": [1, 6, 9, 5, 3]}'
```
**Response (`200 OK`)**:
```json
{
  "result": 24,
  "operation": "add"
}
```

#### C. Structured Mixed Operations Mode
Evaluates an array of numbers and binary operations respecting PEMDAS precedence:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operations": ["add", "multiply"], "numbers": [5, 2, 3]}'
```
**Response (`200 OK`)**:
```json
{
  "result": 11,
  "operation": "mixed"
}
```

#### D. Error Responses (`400 Bad Request` & `422 Unprocessable Entity`)

- **Division by Zero (`422 Unprocessable Entity`)**:
  ```bash
  curl -X POST http://localhost:8080/api/v1/calculate \
    -H "Content-Type: application/json" \
    -d '{"expression": "10 / 0"}'
  ```
  ```json
  {
    "error": "division by zero is not allowed",
    "code": "DIVISION_BY_ZERO"
  }
  ```

- **Negative Square Root (`422 Unprocessable Entity`)**:
  ```json
  {
    "error": "cannot calculate square root of a negative number",
    "code": "NEGATIVE_SQRT"
  }
  ```

- **Invalid Expression / Syntax (`400 Bad Request`)**:
  ```json
  {
    "error": "invalid mathematical expression",
    "code": "INVALID_EXPRESSION"
  }
  ```

---

## Unit Tests & Coverage Report

- **Backend**: `96.9%` total statement coverage (`96.3%` in `internal/calculator`, `100.0%` in `internal/handler`).
- **Frontend**: `33/33` tests passing — `94.33%` line coverage, `100.00%` function coverage.

See [`TestsReport.md`](./TestsReport.md) for the complete function-by-function coverage tables and direct links to the pre-generated HTML/text reports in [`backend/coverage/`](./backend/coverage/) and [`frontend/coverage/`](./frontend/coverage/).

---

## Design Decisions & Assumptions

1. **Clean Layered Architecture**:
   - **Backend**: Domain math logic in `internal/calculator` has zero knowledge of HTTP or JSON. `internal/handler` translates HTTP payloads into domain calls and maps domain errors (`ErrDivisionByZero`, `ErrNegativeSqrt`, `ErrOverflow`, etc.) into structured HTTP responses.
   - **Frontend**: Strictly separates the HTTP service layer (`calculatorApi.ts`), state management hook (`useCalculator.ts`), and declarative view (`App.tsx`).
2. **Recursive Descent Parser for Mixed Expressions**:
   - Rather than restricting users to two numbers at a time, the backend includes a custom recursive descent parser (`EvaluateExpression`) supporting full operator precedence (`()`, `√`, `^`, `%`, `×`/`÷`, `+`/`−`). Visual symbols (`×`, `÷`, `−`, `√`, and Unicode superscripts `⁰–⁹`) are normalized server-side before parsing.
3. **Intuitive Percentage Semantics**:
   - Real-world calculators treat `%` contextually. The parser models `%` as a postfix operator (`X% -> X / 100`) so standalone conversion (`10% = 0.1`) and proportional multiplication (`50 × 10% = 5`) work naturally, while additive/subtractive contexts (`A + B%` or `A - B%`) compute markups and discounts relative to `A` (`200 - 15% = 170`).
4. **Semantic HTTP Error Codes & Security Hardening**:
   - `400 Bad Request` is returned for malformed JSON, unknown fields (`decoder.DisallowUnknownFields()`), trailing JSON objects, invalid syntax, or missing operands.
   - `422 Unprocessable Entity` is returned when the request syntax is valid, but the mathematical operation violates domain rules (division by zero, square root of a negative number, or `Inf`/`NaN` overflow).
   - Request bodies are capped at `64 KB` via `http.MaxBytesReader` to prevent payload abuse.
5. **Atomic Frontend State Updates**:
   - Display state and cursor context (`insideSqrtParens`, `justEvaluated`, `displayValue`) are managed in a single atomic state object inside `useCalculator` to guarantee deterministic updates even when multiple actions are triggered rapidly.

---

## AI Prompts Used

In accordance with the assessment instructions, the prompts used during the development of both the backend and frontend are documented in [`PROMPTS.md`](./PROMPTS.md).
