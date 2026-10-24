import { useState, useCallback } from 'react';
import {
  calculate,
  CalculatorApiError,
  type Operation,
} from '../services/calculatorApi';

export type BinaryOperation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'exponent'
  | 'percent';

export interface HistoryItem {
  id: string;
  expression: string;
  result: number;
  operation: Operation;
}

export const OPERATION_SYMBOLS: Record<BinaryOperation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  exponent: '^',
  percent: '%',
};

const OPERATOR_CHARS = new Set(['+', '−', '-', '×', '*', '÷', '/', '^']);

function isOperatorChar(ch: string): boolean {
  return OPERATOR_CHARS.has(ch);
}

// Extracts the trailing number token at the end of the given string (e.g., "5+2×3.5" -> "3.5").
function getCurrentNumberSegment(expr: string): string {
  let i = expr.length - 1;
  while (i >= 0 && /[\d.]/.test(expr[i])) {
    i--;
  }
  return expr.slice(i + 1);
}

interface ExpressionState {
  displayValue: string;
  insideSqrtParens: boolean;
  justEvaluated: boolean;
  lastSubmittedExpression: string | null;
}

const INITIAL_EXPRESSION_STATE: ExpressionState = {
  displayValue: '0',
  insideSqrtParens: false,
  justEvaluated: false,
  lastSubmittedExpression: null,
};

const HISTORY_STORAGE_KEY = 'sezzle_calculator_history';

function loadStoredHistory(): HistoryItem[] {
  try {
    const raw = window.localStorage.getItem(HISTORY_STORAGE_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as unknown;
    return Array.isArray(parsed) ? (parsed as HistoryItem[]) : [];
  } catch {
    return [];
  }
}

function saveStoredHistory(items: HistoryItem[]): void {
  try {
    window.localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(items));
  } catch {
    // Ignore storage quota or privacy-mode errors
  }
}

export function useCalculator() {
  const [exprState, setExprState] = useState<ExpressionState>(
    INITIAL_EXPRESSION_STATE
  );
  const [history, setHistory] = useState<HistoryItem[]>(loadStoredHistory);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const appendHistory = useCallback(
    (expression: string, result: number, op: Operation = 'expression') => {
      const newItem: HistoryItem = {
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
        expression,
        result,
        operation: op,
      };
      setHistory((prev) => {
        const next = [newItem, ...prev.slice(0, 19)];
        saveStoredHistory(next);
        return next;
      });
    },
    []
  );

  const handleApiError = useCallback((err: unknown) => {
    if (err instanceof CalculatorApiError) {
      setError(err.message);
    } else if (err instanceof Error) {
      setError(err.message);
    } else {
      setError('An unexpected error occurred.');
    }
  }, []);

  const inputDigit = useCallback((digit: string) => {
    if (!/^\d$/.test(digit)) return;
    setError(null);

    setExprState((prev) => {
      if (prev.justEvaluated) {
        return {
          displayValue: digit,
          insideSqrtParens: false,
          justEvaluated: false,
          lastSubmittedExpression: null,
        };
      }

      const { displayValue, insideSqrtParens } = prev;

      // When typing inside √(), insert the digit right before the closing ')'
      if (insideSqrtParens && displayValue.endsWith(')')) {
        const inner = displayValue.slice(0, -1);
        const currentSeg = getCurrentNumberSegment(inner);
        if (currentSeg === '0') {
          return {
            ...prev,
            displayValue: `${inner.slice(0, -1)}${digit})`,
          };
        }
        return {
          ...prev,
          displayValue: `${inner}${digit})`,
        };
      }

      if (displayValue === '0') {
        return {
          ...prev,
          displayValue: digit,
        };
      }

      const currentSeg = getCurrentNumberSegment(displayValue);
      if (currentSeg === '0') {
        return {
          ...prev,
          displayValue: displayValue.slice(0, -1) + digit,
        };
      }

      return {
        ...prev,
        displayValue: displayValue + digit,
      };
    });
  }, []);

  const inputDecimal = useCallback(() => {
    setError(null);

    setExprState((prev) => {
      if (prev.justEvaluated) {
        return {
          displayValue: '0.',
          insideSqrtParens: false,
          justEvaluated: false,
          lastSubmittedExpression: null,
        };
      }

      const { displayValue, insideSqrtParens } = prev;

      if (insideSqrtParens && displayValue.endsWith(')')) {
        const inner = displayValue.slice(0, -1);
        const currentSeg = getCurrentNumberSegment(inner);
        if (currentSeg.includes('.')) {
          return prev;
        }
        if (currentSeg === '') {
          return {
            ...prev,
            displayValue: `${inner}0.)`,
          };
        }
        return {
          ...prev,
          displayValue: `${inner}.)`,
        };
      }

      const currentSeg = getCurrentNumberSegment(displayValue);
      if (currentSeg.includes('.')) {
        return prev;
      }
      if (currentSeg === '') {
        return {
          ...prev,
          displayValue: `${displayValue}0.`,
        };
      }
      return {
        ...prev,
        displayValue: `${displayValue}.`,
      };
    });
  }, []);

  const inputParenthesis = useCallback((paren: '(' | ')') => {
    setError(null);

    setExprState((prev) => {
      if (paren === ')' && prev.insideSqrtParens) {
        // Step out of the current √(...) parentheses without duplicating ')'
        return {
          ...prev,
          insideSqrtParens: false,
        };
      }

      if (prev.justEvaluated && paren === '(') {
        return {
          displayValue: '(',
          insideSqrtParens: false,
          justEvaluated: false,
          lastSubmittedExpression: null,
        };
      }

      if (prev.displayValue === '0' && paren === '(') {
        return {
          ...prev,
          displayValue: '(',
          justEvaluated: false,
        };
      }

      return {
        ...prev,
        displayValue: prev.displayValue + paren,
        justEvaluated: false,
      };
    });
  }, []);

  const toggleSign = useCallback(() => {
    setError(null);

    setExprState((prev) => {
      const { displayValue, insideSqrtParens } = prev;

      if (insideSqrtParens && displayValue.endsWith(')')) {
        const openIdx = displayValue.lastIndexOf('√(');
        if (openIdx !== -1) {
          const prefix = displayValue.slice(0, openIdx + 2);
          const content = displayValue.slice(openIdx + 2, -1);
          const toggled = content.startsWith('-')
            ? content.slice(1)
            : `-${content}`;
          return {
            ...prev,
            displayValue: `${prefix}${toggled})`,
            justEvaluated: false,
          };
        }
      }

      if (displayValue === '0' || displayValue === '0.') {
        return prev;
      }

      if (/^-?\d+(\.\d+)?$/.test(displayValue)) {
        return {
          ...prev,
          displayValue: displayValue.startsWith('-')
            ? displayValue.slice(1)
            : `-${displayValue}`,
          justEvaluated: false,
        };
      }

      return {
        ...prev,
        displayValue: `-(${displayValue})`,
        justEvaluated: false,
      };
    });
  }, []);

  const backspace = useCallback(() => {
    setError(null);

    setExprState((prev) => {
      const { displayValue, insideSqrtParens } = prev;

      if (insideSqrtParens && displayValue.endsWith(')')) {
        if (displayValue.endsWith('√()')) {
          const remaining = displayValue.slice(0, -3);
          return {
            ...prev,
            displayValue: remaining === '' ? '0' : remaining,
            insideSqrtParens: false,
            justEvaluated: false,
          };
        }
        return {
          ...prev,
          displayValue: `${displayValue.slice(0, -2)})`,
          justEvaluated: false,
        };
      }

      if (
        displayValue.length <= 1 ||
        (displayValue.length === 2 && displayValue.startsWith('-'))
      ) {
        return {
          ...prev,
          displayValue: '0',
          justEvaluated: false,
        };
      }

      return {
        ...prev,
        displayValue: displayValue.slice(0, -1),
        justEvaluated: false,
      };
    });
  }, []);

  const clearAll = useCallback(() => {
    setExprState(INITIAL_EXPRESSION_STATE);
    setError(null);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const clearHistory = useCallback(() => {
    setHistory([]);
    saveStoredHistory([]);
  }, []);

  const reuseHistoryResult = useCallback((value: number) => {
    setError(null);
    setExprState({
      displayValue: String(value),
      insideSqrtParens: false,
      justEvaluated: false,
      lastSubmittedExpression: null,
    });
  }, []);

  const selectOperation = useCallback(
    (op: BinaryOperation) => {
      if (isLoading) return;
      setError(null);

      const symbol = OPERATION_SYMBOLS[op];

      setExprState((prev) => {
        const { displayValue, insideSqrtParens } = prev;

        if (insideSqrtParens && displayValue.endsWith(')')) {
          if (displayValue.endsWith('√()')) {
            if (op === 'subtract') {
              return {
                ...prev,
                displayValue: `${displayValue.slice(0, -1)}-)`,
                justEvaluated: false,
              };
            }
            return prev;
          }
          return {
            ...prev,
            displayValue: displayValue + symbol,
            insideSqrtParens: false,
            justEvaluated: false,
          };
        }

        const lastChar = displayValue[displayValue.length - 1] ?? '';
        if (isOperatorChar(lastChar)) {
          return {
            ...prev,
            displayValue: displayValue.slice(0, -1) + symbol,
            justEvaluated: false,
          };
        }

        return {
          ...prev,
          displayValue: displayValue + symbol,
          justEvaluated: false,
        };
      });
    },
    [isLoading]
  );

  const executeSqrt = useCallback(() => {
    if (isLoading) return;
    setError(null);

    setExprState((prev) => {
      if (prev.justEvaluated || prev.displayValue === '0') {
        return {
          displayValue: '√()',
          insideSqrtParens: true,
          justEvaluated: false,
          lastSubmittedExpression: null,
        };
      }

      const { displayValue, insideSqrtParens } = prev;

      if (insideSqrtParens && displayValue.endsWith(')')) {
        if (displayValue.endsWith('√()')) {
          return prev;
        }
        return {
          ...prev,
          displayValue: `${displayValue}×√()`,
          insideSqrtParens: true,
          justEvaluated: false,
        };
      }

      const lastChar = displayValue[displayValue.length - 1] ?? '';
      if (isOperatorChar(lastChar) || lastChar === '(') {
        return {
          ...prev,
          displayValue: `${displayValue}√()`,
          insideSqrtParens: true,
          justEvaluated: false,
        };
      }

      return {
        ...prev,
        displayValue: `${displayValue}×√()`,
        insideSqrtParens: true,
        justEvaluated: false,
      };
    });
  }, [isLoading]);

  const evaluate = useCallback(async () => {
    if (isLoading) return;
    setError(null);

    const expr = exprState.displayValue.trim();
    if (!expr) return;

    setIsLoading(true);
    try {
      const response = await calculate({
        expression: expr,
      });
      appendHistory(expr, response.result, 'expression');
      setExprState({
        displayValue: String(response.result),
        insideSqrtParens: false,
        justEvaluated: true,
        lastSubmittedExpression: expr,
      });
    } catch (err) {
      handleApiError(err);
    } finally {
      setIsLoading(false);
    }
  }, [isLoading, exprState.displayValue, appendHistory, handleApiError]);

  return {
    displayValue: exprState.displayValue,
    activeExpression: exprState.lastSubmittedExpression,
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
  };
}
