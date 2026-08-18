import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { EngineBadge } from './EngineBadge';

describe('EngineBadge', () => {
  it('labels live fetch data', () => {
    render(<EngineBadge engine="fetch" />);
    expect(screen.getByText(/live page fetch/i)).toBeInTheDocument();
  });

  it('falls back to the raw engine name', () => {
    render(<EngineBadge engine="unknown" />);
    expect(screen.getByText('unknown')).toBeInTheDocument();
  });
});
