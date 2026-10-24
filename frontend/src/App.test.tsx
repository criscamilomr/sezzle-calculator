import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { App } from './App';
import * as api from './services/calculatorApi';

vi.mock('./services/calculatorApi', async (importOriginal) => {
  const actual = await importOriginal<typeof import('./services/calculatorApi')>();
  return {
    ...actual,
    calculate: vi.fn(),
  };
});

describe('App Calculator UI & Integration', () => {
  const calculateMock = vi.mocked(api.calculate);

  beforeEach(() => {
    vi.clearAllMocks();
    window.localStorage.clear();
  });

  it('renders initial calculator interface and empty history state', () => {
    render(<App />);

    expect(
      screen.getByRole('heading', { name: /Calculator/i })
    ).toBeInTheDocument();
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');
    expect(screen.getByText(/No calculations yet/i)).toBeInTheDocument();
  });

  it('allows entering mixed operations via keypad buttons and evaluates result', async () => {
    calculateMock.mockResolvedValueOnce({
      result: 11,
      operation: 'expression',
      expression: '5+2×3',
    });

    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Open parenthesis' }));
    fireEvent.click(screen.getByRole('button', { name: '7' }));
    fireEvent.click(screen.getByRole('button', { name: '8' }));
    fireEvent.click(screen.getByRole('button', { name: '9' }));
    fireEvent.click(screen.getByRole('button', { name: 'Subtract' }));
    fireEvent.click(screen.getByRole('button', { name: '4' }));
    fireEvent.click(screen.getByRole('button', { name: '6' }));
    fireEvent.click(screen.getByRole('button', { name: 'Decimal point' }));
    fireEvent.click(screen.getByRole('button', { name: '3' }));
    fireEvent.click(screen.getByRole('button', { name: 'Close parenthesis' }));
    fireEvent.click(screen.getByRole('button', { name: 'Toggle sign' }));
    fireEvent.click(screen.getByRole('button', { name: 'Backspace' }));
    fireEvent.click(screen.getByRole('button', { name: 'All Clear' }));

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: 'Add' }));
    fireEvent.click(screen.getByRole('button', { name: '2' }));
    fireEvent.click(screen.getByRole('button', { name: 'Multiply' }));
    fireEvent.click(screen.getByRole('button', { name: '3' }));

    expect(screen.getByTestId('calculator-display')).toHaveTextContent('5+2×3');

    fireEvent.click(screen.getByRole('button', { name: 'Equals' }));

    await waitFor(() => {
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('11');
    });

    expect(screen.getByTestId('expression-preview')).toHaveTextContent(
      '5+2×3 ='
    );
    expect(calculateMock).toHaveBeenCalledWith({ expression: '5+2×3' });
  });

  it('renders square root starting with √() and inserts digits inside parentheses', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'Square root' }));
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('√()');

    fireEvent.click(screen.getByRole('button', { name: '1' }));
    fireEvent.click(screen.getByRole('button', { name: '4' }));
    fireEvent.click(screen.getByRole('button', { name: '4' }));

    expect(screen.getByTestId('calculator-display')).toHaveTextContent(
      '√(144)'
    );
  });

  it('renders exponentiation with visual superscript element (5²)', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: 'Power' }));

    const display = screen.getByTestId('calculator-display');
    expect(display.querySelector('sup.exponent-placeholder')).toHaveTextContent(
      '□'
    );

    fireEvent.click(screen.getByRole('button', { name: '2' }));

    const sup = display.querySelector('sup.exponent-sup');
    expect(sup).not.toBeNull();
    expect(sup).toHaveTextContent('2');
    expect(display).toHaveTextContent('52');
  });

  it('evaluates proportional percentage (50×10% = 5) and allows reusing from history', async () => {
    calculateMock.mockResolvedValueOnce({
      result: 5,
      operation: 'expression',
      expression: '50×10%',
    });

    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: '5' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: 'Multiply' }));
    fireEvent.click(screen.getByRole('button', { name: '1' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: 'Percentage' }));

    expect(screen.getByTestId('calculator-display')).toHaveTextContent(
      '50×10%'
    );

    fireEvent.click(screen.getByRole('button', { name: 'Equals' }));

    await waitFor(() => {
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('5');
    });

    // Clear display and click history item to reuse result
    fireEvent.click(screen.getByRole('button', { name: 'All Clear' }));
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');

    const historyItemBtn = screen.getByTitle('Click to reuse this result');
    fireEvent.click(historyItemBtn);
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('5');

    // Clear history
    fireEvent.click(screen.getByRole('button', { name: 'Clear' }));
    expect(screen.getByText(/No calculations yet/i)).toBeInTheDocument();
  });

  it('displays error banner on backend error and allows dismissing it', async () => {
    calculateMock.mockRejectedValueOnce(
      new api.CalculatorApiError(
        'division by zero is not allowed',
        422,
        'DIVISION_BY_ZERO'
      )
    );

    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: '8' }));
    fireEvent.click(screen.getByRole('button', { name: 'Divide' }));
    fireEvent.click(screen.getByRole('button', { name: '0' }));
    fireEvent.click(screen.getByRole('button', { name: 'Equals' }));

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('division by zero is not allowed');

    fireEvent.click(screen.getByRole('button', { name: 'Dismiss error' }));
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('supports physical keyboard navigation for all operations', async () => {
    calculateMock.mockResolvedValueOnce({
      result: 24,
      operation: 'expression',
      expression: '1+6+9+5+3',
    });

    render(<App />);

    fireEvent.keyDown(window, { key: '1' });
    fireEvent.keyDown(window, { key: '+' });
    fireEvent.keyDown(window, { key: '6' });
    fireEvent.keyDown(window, { key: '+' });
    fireEvent.keyDown(window, { key: '9' });
    fireEvent.keyDown(window, { key: '+' });
    fireEvent.keyDown(window, { key: '5' });
    fireEvent.keyDown(window, { key: '+' });
    fireEvent.keyDown(window, { key: '3' });

    expect(screen.getByTestId('calculator-display')).toHaveTextContent(
      '1+6+9+5+3'
    );

    fireEvent.keyDown(window, { key: 'Enter' });

    await waitFor(() => {
      expect(screen.getByTestId('calculator-display')).toHaveTextContent('24');
    });

    // Test Escape, R (sqrt), ^, *, -, /, ., (), Backspace
    fireEvent.keyDown(window, { key: 'Escape' });
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');

    fireEvent.keyDown(window, { key: 'r' });
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('√()');

    fireEvent.keyDown(window, { key: 'Backspace' });
    expect(screen.getByTestId('calculator-display')).toHaveTextContent('0');

    fireEvent.keyDown(window, { key: '(' });
    fireEvent.keyDown(window, { key: '4' });
    fireEvent.keyDown(window, { key: '.' });
    fireEvent.keyDown(window, { key: '5' });
    fireEvent.keyDown(window, { key: ')' });
    fireEvent.keyDown(window, { key: '*' });
    fireEvent.keyDown(window, { key: '2' });
    fireEvent.keyDown(window, { key: '-' });
    fireEvent.keyDown(window, { key: '1' });
    fireEvent.keyDown(window, { key: '/' });
    fireEvent.keyDown(window, { key: '2' });
    fireEvent.keyDown(window, { key: '^' });
    fireEvent.keyDown(window, { key: '3' });
    fireEvent.keyDown(window, { key: '%' });

    expect(screen.getByTestId('calculator-display')).toHaveTextContent(
      '(4.5)×2−1÷23%'
    );
  });
});
