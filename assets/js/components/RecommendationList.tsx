import type { Recommendation } from '../api/types';

type RecommendationListProps = {
  items: Recommendation[];
};

export function RecommendationList({ items }: RecommendationListProps) {
  if (items.length === 0) {
    return (
      <p role="status">No recommendations — this mock audit looks healthy.</p>
    );
  }

  return (
    <section aria-labelledby="recs-heading">
      <h2 id="recs-heading">Recommendations</h2>
      <ol className="perfcheck-recommendations">
        {items.map((item, index) => (
          <li key={item.id} className="perfcheck-recommendations__item">
            <span className="perfcheck-recommendations__index" aria-hidden="true">
              {String(index + 1).padStart(2, '0')}
            </span>
            <div>
              <div className="perfcheck-recommendations__meta">
                <span>{item.category}</span>
                <span>·</span>
                <span>{item.severity}</span>
              </div>
              <h3 className="perfcheck-recommendations__title">{item.title}</h3>
              <p className="perfcheck-recommendations__detail">{item.detail}</p>
            </div>
          </li>
        ))}
      </ol>
    </section>
  );
}
