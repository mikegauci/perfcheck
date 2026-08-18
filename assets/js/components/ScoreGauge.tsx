import { scoreBand, scoreBandLabel } from '../lib/scoreBand';

type ScoreGaugeProps = {
  label: string;
  value: number;
};

const RADIUS = 42;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

export function ScoreGauge({ label, value }: ScoreGaugeProps) {
  const band = scoreBand(value);
  const offset = CIRCUMFERENCE - (Math.min(100, Math.max(0, value)) / 100) * CIRCUMFERENCE;
  const bandText = scoreBandLabel(band);

  return (
    <div className="perfcheck-score__gauge">
      <svg
        width="112"
        height="112"
        viewBox="0 0 112 112"
        role="img"
        aria-label={`${label}: ${value} out of 100, ${bandText}`}
      >
        <circle
          cx="56"
          cy="56"
          r={RADIUS}
          fill="none"
          stroke="#e2ddd4"
          strokeWidth="10"
        />
        <circle
          cx="56"
          cy="56"
          r={RADIUS}
          fill="none"
          stroke={band === 'good' ? '#067647' : band === 'needs-work' ? '#b54708' : '#b42318'}
          strokeWidth="10"
          strokeLinecap="round"
          strokeDasharray={CIRCUMFERENCE}
          strokeDashoffset={offset}
          transform="rotate(-90 56 56)"
        />
        <text
          x="56"
          y="60"
          textAnchor="middle"
          fontSize="22"
          fontWeight="700"
          fill="#1a1a1a"
        >
          {value}
        </text>
      </svg>
      <div className="perfcheck-score__label">{label}</div>
      <div className={`perfcheck-score__band perfcheck-score__band--${band}`}>{bandText}</div>
    </div>
  );
}
