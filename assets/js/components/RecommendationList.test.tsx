import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { RecommendationList } from './RecommendationList';
import type { Recommendation } from '../api/types';

const rec = (severity: string, id = severity): Recommendation => ({
  id,
  category: 'seo',
  severity,
  title: `${severity} title`,
  detail: `${severity} detail`,
});

describe('RecommendationList', () => {
  it('announces an empty healthy audit', () => {
    render(<RecommendationList items={[]} />);
    expect(screen.getByRole('status')).toHaveTextContent(/looks healthy/i);
  });

  it.each(['low', 'medium', 'high', 'critical'] as const)(
    'applies a %s severity modifier',
    (severity) => {
      const { container } = render(<RecommendationList items={[rec(severity)]} />);
      expect(container.querySelector(`.perfcheck-recommendations__severity--${severity}`)).not.toBeNull();
      expect(screen.getByText(severity)).toBeInTheDocument();
    },
  );

  it('skips a modifier for unknown severity', () => {
    const { container } = render(<RecommendationList items={[rec('urgent')]} />);
    expect(container.querySelector('[class*="perfcheck-recommendations__severity--"]')).toBeNull();
    expect(screen.getByText('urgent')).toBeInTheDocument();
  });
});
