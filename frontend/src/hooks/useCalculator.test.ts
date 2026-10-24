import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import * as api from '../services/calculatorApi';
import { useCalculator } from './useCalculator';

vi.mock('../services/calculatorApi', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../services/calculatorApi')>();
  return {
    ...actual,
    calculate: vi.fn(),
  };
});

describe('useCalculator hook', () => {
  const calculateMock = vi.mocked(api.calculate);

  beforeEach(() => {
    vi.clearAllMocks();
    window.localStorage.clear();
  });

  it('initializes with default state', () => {
    const { result } = renderHook(() => useCalculator());

    expect(result.current.displayValue).toBe('0');
    expect(result.current.activeExpression).toBeNull();
    expect(result.current.history).toEqual([]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  describe('digit & decimal input', () => {
    it('replaces initial zero and concatenates subsequent digits', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.inputDigit('2');
      });

      expect(result.current.displayValue).toBe('52');
    });

    it('ignores non-digit characters in inputDigit', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('a');
        result.current.inputDigit('12');
      });

      expect(result.current.displayValue).toBe('0');
    });

    it('replaces leading zero after an operator with the next digit', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.selectOperation('add');
        result.current.inputDigit('0');
        result.current.inputDigit('8');
      });

      expect(result.current.displayValue).toBe('5+8');
    });

    it('handles decimal point and prevents multiple decimals in the same number token', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDecimal();
        result.current.inputDigit('5');
        result.current.inputDecimal(); // duplicate in first operand ignored
        result.current.selectOperation('add');
        result.current.inputDecimal(); // starts second operand as 0.
        result.current.inputDigit('2');
      });

      expect(result.current.displayValue).toBe('0.5+0.2');
    });
  });

  describe('mixed operations, exponentiation, and percentage', () => {
    it('builds multi-number mixed expressions (5+2×3 and 1+6+9+5+3)', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.selectOperation('add');
        result.current.inputDigit('2');
        result.current.selectOperation('multiply');
        result.current.inputDigit('3');
      });

      expect(result.current.displayValue).toBe('5+2×3');
    });

    it('replaces the last binary operator if another operator is selected consecutively', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.selectOperation('add');
        result.current.selectOperation('multiply');
      });

      expect(result.current.displayValue).toBe('5×');
    });

    it('preserves postfix % when followed by a binary operator (50×10%+20)', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.inputDigit('0');
        result.current.selectOperation('multiply');
        result.current.inputDigit('1');
        result.current.inputDigit('0');
        result.current.selectOperation('percent');
        result.current.selectOperation('add');
        result.current.inputDigit('2');
        result.current.inputDigit('0');
      });

      expect(result.current.displayValue).toBe('50×10%+20');
    });
  });

  describe('square root √() interactive flow', () => {
    it('starts with √() on initial display and inserts digits inside the parentheses', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.executeSqrt();
      });
      expect(result.current.displayValue).toBe('√()');

      act(() => {
        result.current.inputDigit('1');
        result.current.inputDigit('4');
        result.current.inputDigit('4');
      });
      expect(result.current.displayValue).toBe('√(144)');
    });

    it('handles decimals and sign toggle inside √()', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.executeSqrt();
        result.current.inputDecimal();
        result.current.inputDigit('2');
        result.current.inputDecimal(); // ignored duplicate
        result.current.inputDigit('5');
      });
      expect(result.current.displayValue).toBe('√(0.25)');

      act(() => {
        result.current.toggleSign();
      });
      expect(result.current.displayValue).toBe('√(-0.25)');

      act(() => {
        result.current.toggleSign();
      });
      expect(result.current.displayValue).toBe('√(0.25)');
    });

    it('steps outside √(...) when a binary operator or closing parenthesis is pressed', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.executeSqrt();
        result.current.inputDigit('9');
        result.current.selectOperation('add');
        result.current.inputDigit('5');
      });
      expect(result.current.displayValue).toBe('√(9)+5');
    });

    it('allows stepping out of √(...) with ) without duplicating the parenthesis', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.executeSqrt();
        result.current.inputDigit('1');
        result.current.inputDigit('6');
        result.current.inputParenthesis(')');
        result.current.selectOperation('multiply');
        result.current.inputDigit('2');
      });
      expect(result.current.displayValue).toBe('√(16)×2');
    });

    it('removes empty √() on backspace and deletes inner digits when non-empty', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.executeSqrt();
        result.current.inputDigit('9');
      });
      expect(result.current.displayValue).toBe('√(9)');

      act(() => {
        result.current.backspace();
      });
      expect(result.current.displayValue).toBe('√()');

      act(() => {
        result.current.backspace();
      });
      expect(result.current.displayValue).toBe('0');
    });

    it('chains √() with implicit multiplication when pressed after a number and direct √() after an operator', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('2');
        result.current.executeSqrt();
        result.current.inputDigit('0');
        result.current.inputDigit('9');
      });
      expect(result.current.displayValue).toBe('2×√(9)');

      act(() => {
        result.current.selectOperation('add');
        result.current.executeSqrt();
        result.current.selectOperation('subtract'); // negative sign inside empty √()
        result.current.inputDigit('4');
      });
      expect(result.current.displayValue).toBe('2×√(9)+√(-4)');
    });
  });

  describe('parentheses, toggleSign, backspace, and clear', () => {
    it('supports parentheses, backspace on multi-char numbers, and toggling sign on numbers and expressions', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.toggleSign(); // on '0' does nothing
      });
      expect(result.current.displayValue).toBe('0');

      act(() => {
        result.current.inputParenthesis('(');
        result.current.inputDigit('5');
        result.current.inputDigit('8');
        result.current.backspace();
        result.current.selectOperation('add');
        result.current.inputDigit('3');
        result.current.inputParenthesis(')');
        result.current.toggleSign();
      });
      expect(result.current.displayValue).toBe('-((5+3))');
    });

    it('resets everything when clearAll is called', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('9');
        result.current.selectOperation('add');
        result.current.inputDigit('1');
        result.current.clearAll();
      });

      expect(result.current.displayValue).toBe('0');
      expect(result.current.activeExpression).toBeNull();
    });
  });

  describe('evaluate() and history management', () => {
    it('evaluates expression via backend API, updates display and history, and resets on next digit', async () => {
      calculateMock.mockResolvedValueOnce({
        result: 11,
        operation: 'expression',
        expression: '5+2×3',
      });

      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.selectOperation('add');
        result.current.inputDigit('2');
        result.current.selectOperation('multiply');
        result.current.inputDigit('3');
      });

      await act(async () => {
        await result.current.evaluate();
      });

      expect(calculateMock).toHaveBeenCalledWith({ expression: '5+2×3' });
      expect(result.current.displayValue).toBe('11');
      expect(result.current.activeExpression).toBe('5+2×3');
      expect(result.current.history).toHaveLength(1);
      expect(result.current.history[0]).toMatchObject({
        expression: '5+2×3',
        result: 11,
      });

      // Next digit after evaluation starts a fresh calculation
      act(() => {
        result.current.inputDigit('7');
      });
      expect(result.current.displayValue).toBe('7');
      expect(result.current.activeExpression).toBeNull();
    });

    it('reuses a history result and clears history', async () => {
      calculateMock.mockResolvedValueOnce({
        result: 25,
        operation: 'expression',
        expression: '5^2',
      });

      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.selectOperation('exponent');
        result.current.inputDigit('2');
      });

      await act(async () => {
        await result.current.evaluate();
      });

      act(() => {
        result.current.clearAll();
        result.current.reuseHistoryResult(25);
      });
      expect(result.current.displayValue).toBe('25');

      act(() => {
        result.current.clearHistory();
      });
      expect(result.current.history).toHaveLength(0);
    });

    it('captures and clears CalculatorApiError when evaluation fails', async () => {
      calculateMock.mockRejectedValueOnce(
        new api.CalculatorApiError(
          'division by zero is not allowed',
          422,
          'DIVISION_BY_ZERO'
        )
      );

      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('1');
        result.current.inputDigit('0');
        result.current.selectOperation('divide');
        result.current.inputDigit('0');
      });

      await act(async () => {
        await result.current.evaluate();
      });

      expect(result.current.error).toBe('division by zero is not allowed');

      act(() => {
        result.current.clearError();
      });
      expect(result.current.error).toBeNull();
    });
  });
});
