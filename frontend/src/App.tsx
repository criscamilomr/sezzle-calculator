import { useEffect, type ReactNode } from 'react';
import { useCalculator } from './hooks/useCalculator';
import './App.css';

// Renders expressions with exponents raised as visual superscripts (e.g., "5^2" -> 5²) just like macOS Calculator.
function renderFormattedMath(expr: string): ReactNode[] {
  const parts: ReactNode[] = [];
  const regex = /\^([+-]?\d*\.?\d*)/g;
  let lastIndex = 0;
  let match: RegExpExecArray | null;

  while ((match = regex.exec(expr)) !== null) {
    if (match.index > lastIndex) {
      parts.push(expr.slice(lastIndex, match.index));
    }

    const expValue = match[1];
    if (expValue === '') {
      parts.push(
        <sup
          key={`exp-${match.index}`}
          className="exponent-sup exponent-placeholder"
        >
          □
        </sup>
      );
    } else {
      parts.push(
        <sup key={`exp-${match.index}`} className="exponent-sup">
          {expValue}
        </sup>
      );
    }

    lastIndex = regex.lastIndex;
  }

  if (lastIndex < expr.length) {
    parts.push(expr.slice(lastIndex));
  }

  return parts;
}

export function App() {
  const {
    displayValue,
    activeExpression,
    history,
    isLoading,
    error,
    inputDigit,
    inputDecimal,
    inputParenthesis,
    toggleSign,
    backspace,
    clearAll,
    clearError,
    clearHistory,
    reuseHistoryResult,
    selectOperation,
    executeSqrt,
    evaluate,
  } = useCalculator();

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.metaKey || event.ctrlKey || event.altKey) return;

      const { key } = event;

      if (/^\d$/.test(key)) {
        event.preventDefault();
        inputDigit(key);
      } else if (key === '.' || key === ',') {
        event.preventDefault();
        inputDecimal();
      } else if (key === '(' || key === ')') {
        event.preventDefault();
        inputParenthesis(key);
      } else if (key === '+') {
        event.preventDefault();
        selectOperation('add');
      } else if (key === '-') {
        event.preventDefault();
        selectOperation('subtract');
      } else if (key === '*') {
        event.preventDefault();
        selectOperation('multiply');
      } else if (key === '/') {
        event.preventDefault();
        selectOperation('divide');
      } else if (key === '^') {
        event.preventDefault();
        selectOperation('exponent');
      } else if (key === '%') {
        event.preventDefault();
        selectOperation('percent');
      } else if (key === 'r' || key === 'R') {
        event.preventDefault();
        executeSqrt();
      } else if (key === 'Enter' || key === '=') {
        event.preventDefault();
        void evaluate();
      } else if (key === 'Backspace') {
        event.preventDefault();
        backspace();
      } else if (key === 'Escape' || key === 'c' || key === 'C') {
        event.preventDefault();
        clearAll();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [
    inputDigit,
    inputDecimal,
    inputParenthesis,
    selectOperation,
    executeSqrt,
    evaluate,
    backspace,
    clearAll,
  ]);

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="brand">
          <span className="brand-dot" aria-hidden="true" />
          <h1 className="brand-title">Calculator</h1>
        </div>
        <span className="keyboard-hint">
          Keyboard enabled · Press <kbd className="kbd-key">R</kbd> for √
        </span>
      </header>

      <main className="workspace">
        <section className="calculator-card" aria-label="Calculator">
          <div className="display-container">
            <div className="display-top-row">
              <span
                className="expression-preview"
                data-testid="expression-preview"
              >
                {activeExpression ? (
                  <>{renderFormattedMath(activeExpression)} =</>
                ) : (
                  '\u00A0'
                )}
              </span>
            </div>
            <div
              className="display-value"
              role="status"
              aria-label="Calculator display"
              data-testid="calculator-display"
            >
              {renderFormattedMath(displayValue)}
            </div>
          </div>

          {error && (
            <div className="error-wrapper">
              <div className="error-banner" role="alert">
                <span>{error}</span>
                <button
                  type="button"
                  className="error-dismiss"
                  onClick={clearError}
                  aria-label="Dismiss error"
                >
                  ✕
                </button>
              </div>
            </div>
          )}

          <div className="keypad">
            {/* Row 1: Clear, Backspace, Parentheses */}
            <button
              type="button"
              className="calc-btn btn-fn"
              onClick={clearAll}
              disabled={isLoading}
              aria-label="All Clear"
            >
              AC
            </button>
            <button
              type="button"
              className="calc-btn btn-fn"
              onClick={backspace}
              disabled={isLoading}
              aria-label="Backspace"
            >
              ⌫
            </button>
            <button
              type="button"
              className="calc-btn btn-fn"
              onClick={() => inputParenthesis('(')}
              disabled={isLoading}
              aria-label="Open parenthesis"
            >
              (
            </button>
            <button
              type="button"
              className="calc-btn btn-fn"
              onClick={() => inputParenthesis(')')}
              disabled={isLoading}
              aria-label="Close parenthesis"
            >
              )
            </button>

            {/* Row 2: Advanced Operations + Divide */}
            <button
              type="button"
              className="calc-btn btn-advanced"
              onClick={executeSqrt}
              disabled={isLoading}
              aria-label="Square root"
              title="Square root √() (Key: R)"
            >
              √x
            </button>
            <button
              type="button"
              className="calc-btn btn-advanced"
              onClick={() => selectOperation('exponent')}
              disabled={isLoading}
              aria-label="Power"
              title="Exponentiation (Key: ^)"
            >
              xʸ
            </button>
            <button
              type="button"
              className="calc-btn btn-advanced"
              onClick={() => selectOperation('percent')}
              disabled={isLoading}
              aria-label="Percentage"
              title="Percentage: B% of A (Key: %)"
            >
              %
            </button>
            <button
              type="button"
              className="calc-btn btn-op"
              onClick={() => selectOperation('divide')}
              disabled={isLoading}
              aria-label="Divide"
            >
              ÷
            </button>

            {/* Row 3: 7, 8, 9, Multiply */}
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('7')}
              disabled={isLoading}
            >
              7
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('8')}
              disabled={isLoading}
            >
              8
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('9')}
              disabled={isLoading}
            >
              9
            </button>
            <button
              type="button"
              className="calc-btn btn-op"
              onClick={() => selectOperation('multiply')}
              disabled={isLoading}
              aria-label="Multiply"
            >
              ×
            </button>

            {/* Row 4: 4, 5, 6, Subtract */}
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('4')}
              disabled={isLoading}
            >
              4
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('5')}
              disabled={isLoading}
            >
              5
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('6')}
              disabled={isLoading}
            >
              6
            </button>
            <button
              type="button"
              className="calc-btn btn-op"
              onClick={() => selectOperation('subtract')}
              disabled={isLoading}
              aria-label="Subtract"
            >
              −
            </button>

            {/* Row 5: 1, 2, 3, Add */}
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('1')}
              disabled={isLoading}
            >
              1
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('2')}
              disabled={isLoading}
            >
              2
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('3')}
              disabled={isLoading}
            >
              3
            </button>
            <button
              type="button"
              className="calc-btn btn-op"
              onClick={() => selectOperation('add')}
              disabled={isLoading}
              aria-label="Add"
            >
              +
            </button>

            {/* Row 6: +/-, 0, Decimal, Equals */}
            <button
              type="button"
              className="calc-btn btn-fn"
              onClick={toggleSign}
              disabled={isLoading}
              aria-label="Toggle sign"
            >
              +/−
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={() => inputDigit('0')}
              disabled={isLoading}
            >
              0
            </button>
            <button
              type="button"
              className="calc-btn"
              onClick={inputDecimal}
              disabled={isLoading}
              aria-label="Decimal point"
            >
              .
            </button>
            <button
              type="button"
              className="calc-btn btn-eq"
              onClick={() => void evaluate()}
              disabled={isLoading}
              aria-label="Equals"
            >
              =
            </button>
          </div>
        </section>

        <aside className="history-card" aria-label="Calculation history">
          <div className="history-header">
            <h2 className="history-title">History</h2>
            {history.length > 0 && (
              <button
                type="button"
                className="clear-history-btn"
                onClick={clearHistory}
              >
                Clear
              </button>
            )}
          </div>

          {history.length === 0 ? (
            <p className="history-empty">
              No calculations yet. Results will appear here—click any item to
              reuse it.
            </p>
          ) : (
            <ul className="history-list">
              {history.map((item) => (
                <li key={item.id} className="history-item">
                  <button
                    type="button"
                    className="history-item-btn"
                    onClick={() => reuseHistoryResult(item.result)}
                    title="Click to reuse this result"
                  >
                    <span className="history-expr">
                      {renderFormattedMath(item.expression)} =
                    </span>
                    <span className="history-res">{item.result}</span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </aside>
      </main>
    </div>
  );
}

export default App;
