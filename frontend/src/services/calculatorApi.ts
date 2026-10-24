export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'exponent'
  | 'percent'
  | 'sqrt'
  | 'mixed'
  | 'expression';

export interface CalculateRequest {
  expression?: string;
  operation?: Operation;
  operations?: Operation[];
  numbers?: number[];
}

export interface CalculateResponse {
  result: number;
  operation?: Operation;
  expression?: string;
}

export interface ErrorResponse {
  error: string;
  code: string;
}

export interface HealthResponse {
  status: string;
}

export class CalculatorApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(message: string, status: number, code: string) {
    super(message);
    this.name = 'CalculatorApiError';
    this.status = status;
    this.code = code;
  }
}

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1';

export async function calculate(
  request: CalculateRequest
): Promise<CalculateResponse> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE_URL}/calculate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
  } catch {
    throw new CalculatorApiError(
      'Unable to connect to the calculator service. Please verify the backend is running.',
      0,
      'NETWORK_ERROR'
    );
  }

  if (!response.ok) {
    let errorPayload: ErrorResponse | null = null;
    try {
      errorPayload = (await response.json()) as ErrorResponse;
    } catch {
      // Fallback if response body is not valid JSON
    }

    throw new CalculatorApiError(
      errorPayload?.error ?? `Request failed with status ${response.status}`,
      response.status,
      errorPayload?.code ?? 'UNKNOWN_ERROR'
    );
  }

  return (await response.json()) as CalculateResponse;
}

export async function checkHealth(): Promise<HealthResponse> {
  const response = await fetch(`${API_BASE_URL}/health`, {
    method: 'GET',
    headers: {
      Accept: 'application/json',
    },
  });

  if (!response.ok) {
    throw new CalculatorApiError(
      'Health check failed',
      response.status,
      'HEALTH_CHECK_FAILED'
    );
  }

  return (await response.json()) as HealthResponse;
}
