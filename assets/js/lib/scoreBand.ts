export type ScoreBand = 'good' | 'needs-work' | 'poor';

export function scoreBand(value: number): ScoreBand {
  if (value >= 90) return 'good';
  if (value >= 70) return 'needs-work';
  return 'poor';
}

export function scoreBandLabel(band: ScoreBand): string {
  switch (band) {
    case 'good':
      return 'Good';
    case 'needs-work':
      return 'Needs work';
    case 'poor':
      return 'Poor';
  }
}
