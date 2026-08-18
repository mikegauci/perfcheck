type EngineBadgeProps = {
  engine: string;
};

const COPY: Record<string, { label: string; title: string }> = {
  psi: {
    label: 'Live PageSpeed data',
    title: 'Scores come from the Google PageSpeed Insights API (Lighthouse, mobile).',
  },
  fetch: {
    label: 'Live page fetch',
    title: 'Scores come from one HTTP GET and HTML inspection of the URL.',
  },
  mock: {
    label: 'Simulated data',
    title: 'Scores are deterministic mocks derived from the URL. No live fetch ran.',
  },
};

export function EngineBadge({ engine }: EngineBadgeProps) {
  const meta = COPY[engine] ?? {
    label: engine,
    title: 'Scoring engine used for this result.',
  };
  return (
    <span className={`perfcheck-badge perfcheck-badge--${engine}`} title={meta.title}>
      {meta.label}
    </span>
  );
}
