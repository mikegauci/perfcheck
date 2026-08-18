export type Scores = {
  overall: number;
  performance: number;
  seo: number;
  accessibility: number;
};

export type Recommendation = {
  id: string;
  category: string;
  severity: string;
  title: string;
  detail: string;
};

export type Audit = {
  id: string;
  url: string;
  createdAt: string;
  scores: Scores;
  recommendations: Recommendation[];
};

export type AuditListResponse = {
  items: Audit[];
};

export type ApiErrorCode =
  | 'invalid_url'
  | 'invalid_json'
  | 'not_found'
  | 'rate_limited'
  | 'internal'
  | 'network'
  | 'timeout'
  | 'unknown';

export type ApiErrorBody = {
  error: {
    code: ApiErrorCode | string;
    message: string;
    field?: string;
  };
};

export class ApiError extends Error {
  readonly code: ApiErrorCode | string;
  readonly field?: string;
  readonly status: number;

  constructor(code: ApiErrorCode | string, message: string, status: number, field?: string) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
    this.field = field;
  }
}
