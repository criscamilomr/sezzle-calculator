import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  calculate,
  CalculatorApiError,
  checkHealth,
} from './calculatorApi';

describe('calculatorApi service', () => {
  const fetchMock = vi.fn();

  beforeEach(() => {
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  describe('calculate()', () => {
    it('sends a POST request with expression and returns the parsed result', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          result: 11,
          operation: 'expression',
          expression: '5+2×3',
        }),
      });

      const response = await calculate({ expression: '5+2×3' });

      expect(fetchMock).toHaveBeenCalledTimes(1);
      expect(fetchMock).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/calculate',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ expression: '5+2×3' }),
        }
      );
      expect(response).toEqual({
        result: 11,
        operation: 'expression',
        expression: '5+2×3',
      });
    });

    it('sends a POST request with operation and numbers array', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          result: 24,
          operation: 'add',
        }),
      });

      const response = await calculate({
        operation: 'add',
        numbers: [1, 6, 9, 5, 3],
      });

      expect(response.result).toBe(24);
      expect(response.operation).toBe('add');
    });

    it('throws CalculatorApiError with backend error message and code on 422 response', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: false,
        status: 422,
        json: async () => ({
          error: 'division by zero is not allowed',
          code: 'DIVISION_BY_ZERO',
        }),
      });

      await expect(
        calculate({ expression: '10÷0' })
      ).rejects.toMatchObject({
        name: 'CalculatorApiError',
        message: 'division by zero is not allowed',
        status: 422,
        code: 'DIVISION_BY_ZERO',
      });
    });

    it('throws CalculatorApiError with fallback message if error body is not valid JSON', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => {
          throw new Error('Unexpected token');
        },
      });

      await expect(calculate({ expression: '1+1' })).rejects.toMatchObject({
        name: 'CalculatorApiError',
        message: 'Request failed with status 500',
        status: 500,
        code: 'UNKNOWN_ERROR',
      });
    });

    it('throws CalculatorApiError with NETWORK_ERROR code when fetch rejects', async () => {
      fetchMock.mockRejectedValueOnce(new TypeError('Failed to fetch'));

      try {
        await calculate({ expression: '5+5' });
        expect.fail('Expected calculate to throw CalculatorApiError');
      } catch (err) {
        expect(err).toBeInstanceOf(CalculatorApiError);
        const apiErr = err as CalculatorApiError;
        expect(apiErr.status).toBe(0);
        expect(apiErr.code).toBe('NETWORK_ERROR');
        expect(apiErr.message).toContain(
          'Unable to connect to the calculator service'
        );
      }
    });
  });

  describe('checkHealth()', () => {
    it('returns status ok when backend health check succeeds', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ status: 'ok' }),
      });

      const res = await checkHealth();

      expect(fetchMock).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/health',
        {
          method: 'GET',
          headers: { Accept: 'application/json' },
        }
      );
      expect(res).toEqual({ status: 'ok' });
    });

    it('throws CalculatorApiError when health check returns non-200 status', async () => {
      fetchMock.mockResolvedValueOnce({
        ok: false,
        status: 503,
      });

      await expect(checkHealth()).rejects.toMatchObject({
        name: 'CalculatorApiError',
        message: 'Health check failed',
        status: 503,
        code: 'HEALTH_CHECK_FAILED',
      });
    });
  });
});
