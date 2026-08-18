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
      <ul className="perfcheck-recommendations">
        {items.map((item) => (
          <li key={item.id} className="perfcheck-recommendations__item">
            <div className="perfcheck-recommendations__meta">
              <span>{item.category}</span>
              <span>·</span>
              <span>{item.severity}</span>
            </div>
            <h3 className="perfcheck-recommendations__title">{item.title}</h3>
            <p className="perfcheck-recommendations__detail">{item.detail}</p>
          </li>
        ))}
      </ul>
    </section>
  );
}
