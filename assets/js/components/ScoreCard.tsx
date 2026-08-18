import type { Scores } from '../api/types';
import { ScoreGauge } from './ScoreGauge';

type ScoreCardProps = {
  scores: Scores;
  url: string;
};

export function ScoreCard({ scores, url }: ScoreCardProps) {
  return (
    <section className="perfcheck-score" aria-labelledby="audit-results-heading">
      <h2 id="audit-results-heading">Results for {url}</h2>
      <div className="perfcheck-score__summary">
        <ScoreGauge label="Overall" value={scores.overall} />
        <div className="perfcheck-score__gauges">
          <ScoreGauge label="Performance" value={scores.performance} />
          <ScoreGauge label="SEO" value={scores.seo} />
          <ScoreGauge label="Accessibility" value={scores.accessibility} />
        </div>
      </div>
    </section>
  );
}
