import type { Signals } from '../api/types';

type SignalListProps = {
  signals?: Signals;
};

function Row({ term, value }: { term: string; value: string }) {
  return (
    <div className="perfcheck-signals__row">
      <dt>{term}</dt>
      <dd>{value}</dd>
    </div>
  );
}

function yn(v: boolean) {
  return v ? 'Yes' : 'No';
}

export function SignalList({ signals }: SignalListProps) {
  if (!signals) {
    return null;
  }
  return (
    <details className="perfcheck-signals">
      <summary>Evidence</summary>
      <dl>
        {signals.statusCode ? <Row term="HTTP status" value={String(signals.statusCode)} /> : null}
        {signals.ttfbMs != null ? <Row term="TTFB" value={`${signals.ttfbMs} ms`} /> : null}
        {signals.bytes ? <Row term="Bytes transferred" value={String(signals.bytes)} /> : null}
        <Row term="HTTPS" value={yn(signals.https)} />
        {signals.cacheControl ? <Row term="Cache-Control" value={signals.cacheControl} /> : null}
        <Row term="Title" value={signals.hasTitle ? signals.title || 'Present' : 'Missing'} />
        <Row term="Meta description" value={yn(signals.hasMetaDescription)} />
        <Row term="H1 count" value={String(signals.h1Count)} />
        <Row term="Canonical" value={yn(signals.hasCanonical)} />
        <Row term="Open Graph" value={yn(signals.hasOpenGraph)} />
        <Row term="JSON-LD" value={yn(signals.hasJsonLd)} />
        <Row term="html[lang]" value={yn(signals.hasLang)} />
        <Row term="Images missing alt" value={String(signals.imagesMissingAlt)} />
        <Row term="Viewport meta" value={yn(signals.hasViewport)} />
        <Row term="Inputs without label" value={String(signals.inputsWithoutLabel)} />
        <Row term="Scripts" value={String(signals.scriptCount)} />
        <Row term="Stylesheets" value={String(signals.stylesheetCount)} />
        {signals.psiStrategy ? <Row term="PSI strategy" value={signals.psiStrategy} /> : null}
      </dl>
    </details>
  );
}
