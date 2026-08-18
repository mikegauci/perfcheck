import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ScoreGauge } from './ScoreGauge';

describe('ScoreGauge', () => {
  it('exposes an accessible name with band text', () => {
    render(<ScoreGauge label="Performance" value={92} />);
    expect(
      screen.getByRole('img', { name: /performance: 92 out of 100, good/i }),
    ).toBeInTheDocument();
    expect(screen.getByText('Good')).toBeInTheDocument();
  });
});
