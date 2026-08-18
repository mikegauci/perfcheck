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
  engine: 'mock' | 'fetch' | 'psi' | string;
  scores: Scores;
  signals?: Signals;
  recommendations: Recommendation[];
};

export type Signals = {
  statusCode?: number;
  ttfbMs?: number;
  bytes?: number;
  contentEncoding?: string;
  cacheControl?: string;
  https: boolean;
  title?: string;
  titleLength?: number;
  hasTitle: boolean;
  metaDescription?: string;
  hasMetaDescription: boolean;
  h1Count: number;
  hasCanonical: boolean;
  hasOpenGraph: boolean;
  hasJsonLd: boolean;
  hasLang: boolean;
  imagesMissingAlt: number;
  hasViewport: boolean;
  inputsWithoutLabel: number;
  scriptCount: number;
  stylesheetCount: number;
  inlineStyleBytes: number;
  psiStrategy?: string;
};

export type AuditListResponse = {
  items: Audit[];
};

export type ApiErrorCode =
  | 'invalid_url'
  | 'invalid_json'
  | 'not_found'
  | 'unauthorized'
  | 'internal'
  | 'network'
  | 'timeout'
  | 'aborted'
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
