import { scoreBand, scoreBandLabel } from '../lib/scoreBand';

type ScoreGaugeProps = {
  label: string;
  value: number;
  size?: 'md' | 'lg';
};

export function ScoreGauge({ label, value, size = 'md' }: ScoreGaugeProps) {
  const band = scoreBand(value);
  const clamped = Math.min(100, Math.max(0, value));
  const bandText = scoreBandLabel(band);
  const name = `${label}: ${value} out of 100, ${bandText}`;

  return (
    <div
      className={`perfcheck-meter perfcheck-meter--${size} perfcheck-meter--${band}`}
      role="meter"
      aria-label={name}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={clamped}
    >
      <div className="perfcheck-meter__head">
        <div className="perfcheck-meter__label">{label}</div>
        <div className="perfcheck-meter__value">{value}</div>
      </div>
      <div className="perfcheck-meter__track">
        <div className="perfcheck-meter__fill" style={{ width: `${clamped}%` }} />
      </div>
      <div className={`perfcheck-score__band perfcheck-score__band--${band}`}>{bandText}</div>
    </div>
  );
}
