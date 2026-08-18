type EngineBadgeProps = {
  engine: string;
};

const COPY: Record<string, { label: string; title: string }> = {
  fetch: {
    label: 'Live page fetch',
    title: 'Scores come from one HTTP GET and HTML inspection of the URL.',
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
