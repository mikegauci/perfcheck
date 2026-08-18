import type { Recommendation } from '../api/types';

type RecommendationListProps = {
  items: Recommendation[];
};

const SEVERITY_TONES = new Set(['low', 'medium', 'high', 'critical']);

function severityTone(severity: string): string | undefined {
  const key = severity.trim().toLowerCase();
  return SEVERITY_TONES.has(key) ? key : undefined;
}

export function RecommendationList({ items }: RecommendationListProps) {
  if (items.length === 0) {
    return (
      <p role="status">No recommendations — this audit looks healthy.</p>
    );
  }

  return (
    <section aria-labelledby="recs-heading">
      <h2 id="recs-heading">Recommendations</h2>
      <ol className="perfcheck-recommendations">
        {items.map((item, index) => {
          const tone = severityTone(item.severity);
          const severityClass = [
            'perfcheck-recommendations__severity',
            tone ? `perfcheck-recommendations__severity--${tone}` : '',
          ]
            .filter(Boolean)
            .join(' ');

          return (
            <li key={item.id} className="perfcheck-recommendations__item">
              <span className="perfcheck-recommendations__index" aria-hidden="true">
                {String(index + 1).padStart(2, '0')}
              </span>
              <div>
                <div className="perfcheck-recommendations__meta">
                  <span>{item.category}</span>
                  <span>·</span>
                  <span className={severityClass}>{item.severity}</span>
                </div>
                <h3 className="perfcheck-recommendations__title">{item.title}</h3>
                <p className="perfcheck-recommendations__detail">{item.detail}</p>
              </div>
            </li>
          );
        })}
      </ol>
    </section>
  );
}
