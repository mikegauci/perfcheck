import { describe, expect, it } from 'vitest';
import { render, screen } from '@testing-library/react';
import { SignalList } from './SignalList';

describe('SignalList', () => {
  it('renders evidence from signals', () => {
    render(
      <SignalList
        signals={{
          https: true,
          hasTitle: true,
          title: 'Home',
          titleLength: 4,
          hasMetaDescription: false,
          h1Count: 1,
          hasCanonical: true,
          hasOpenGraph: false,
          hasJsonLd: false,
          hasLang: true,
          imagesMissingAlt: 2,
          hasViewport: true,
          inputsWithoutLabel: 0,
          scriptCount: 3,
          stylesheetCount: 1,
          inlineStyleBytes: 0,
        }}
      />,
    );
    expect(screen.getByText('Evidence')).toBeInTheDocument();
    expect(screen.getByText('Images missing alt')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
  });
});
