import type { Audit } from '../api/types';
import { EngineBadge } from './EngineBadge';
import { RecommendationList } from './RecommendationList';
import { ScoreCard } from './ScoreCard';
import { SignalList } from './SignalList';

type AuditReportProps = {
  audit: Audit;
};

export function AuditReport({ audit }: AuditReportProps) {
  return (
    <div>
      <p>
        <EngineBadge engine={audit.engine} />
      </p>
      <ScoreCard scores={audit.scores} url={audit.url} />
      <SignalList signals={audit.signals} />
      <RecommendationList items={audit.recommendations} />
    </div>
  );
}
