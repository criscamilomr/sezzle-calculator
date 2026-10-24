## Unit Tests & Coverage Report

Pre-generated test coverage reports are included directly in the repository for immediate inspection:

- **Backend Coverage Artifacts**:
  - Interactive HTML Report: [`backend/coverage/coverage.html`](./backend/coverage/coverage.html)
  - Function Breakdown Text Report: [`backend/coverage/coverage.txt`](./backend/coverage/coverage.txt)
  - Raw Go Coverprofile: [`backend/coverage/coverage.out`](./backend/coverage/coverage.out)
- **Frontend Coverage Artifacts**:
  - Interactive HTML Report: [`frontend/coverage/index.html`](./frontend/coverage/index.html)
  - Summary Text Report: [`frontend/coverage/coverage.txt`](./frontend/coverage/coverage.txt)

### Backend Coverage Summary (`96.9%` Total Statements)

| Package / File | Function | Statement Coverage |
| :--- | :--- | :---: |
| `internal/calculator/calculator.go` | `Calculate` | `92.0%` |
| `internal/calculator/calculator.go` | `CalculateMixed` | `100.0%` |
| `internal/calculator/calculator.go` | `EvaluateExpression` | `92.3%` |
| `internal/calculator/calculator.go` | `isBinaryOperation` | `100.0%` |
| `internal/calculator/calculator.go` | `applyBinaryOperation` | `92.9%` |
| `internal/calculator/calculator.go` | `reducePrecedenceLevel` | `100.0%` |
| `internal/calculator/calculator.go` | `validateNumber` | `100.0%` |
| `internal/calculator/calculator.go` | `normalizeFloat` | `100.0%` |
| `internal/calculator/calculator.go` | `normalizeExpressionSymbols` | `100.0%` |
| `internal/calculator/calculator.go` | `skipWhitespace` | `100.0%` |
| `internal/calculator/calculator.go` | `parseExpression` | `95.8%` |
| `internal/calculator/calculator.go` | `parseTerm` | `100.0%` |
| `internal/calculator/calculator.go` | `parsePower` | `100.0%` |
| `internal/calculator/calculator.go` | `parseUnary` | `96.0%` |
| `internal/calculator/calculator.go` | `parsePostfix` | `92.9%` |
| `internal/calculator/calculator.go` | `parsePrimary` | `93.5%` |
| `internal/handler/handler.go` | `NewHandler` | `100.0%` |
| `internal/handler/handler.go` | `RegisterRoutes` | `100.0%` |
| `internal/handler/handler.go` | `Health` | `100.0%` |
| `internal/handler/handler.go` | `Calculate` | `100.0%` |
| `internal/handler/handler.go` | `mapDomainError` | `100.0%` |
| `internal/handler/handler.go` | `writeError` | `100.0%` |
| `internal/handler/handler.go` | `writeJSON` | `100.0%` |
| **Total (Backend)** | **All Packages (`internal/...`)** | **`96.9%`** |

### Frontend Coverage Summary (`33/33` Tests Passing — `94.33%` Lines, `100%` Functions)

| File | % Statements | % Branch | % Functions | % Lines |
| :--- | :---: | :---: | :---: | :---: |
| `src/App.tsx` | `98.80%` | `91.07%` | `100.00%` | `100.00%` |
| `src/hooks/useCalculator.ts` | `88.88%` | `75.86%` | `100.00%` | `90.90%` |
| `src/services/calculatorApi.ts` | `100.00%` | `100.00%` | `100.00%` | `100.00%` |
| **Total (Frontend)** | **`92.67%`** | **`81.86%`** | **`100.00%`** | **`94.33%`** |