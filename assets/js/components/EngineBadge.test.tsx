import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { EngineBadge } from './EngineBadge';

describe('EngineBadge', () => {
  it('labels mock data', () => {
    render(<EngineBadge engine="mock" />);
    expect(screen.getByText(/simulated data/i)).toBeInTheDocument();
  });

  it('labels PageSpeed data', () => {
    render(<EngineBadge engine="psi" />);
    expect(screen.getByText(/live pagespeed data/i)).toBeInTheDocument();
  });
});
