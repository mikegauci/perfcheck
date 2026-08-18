export function ReportSkeleton() {
  return (
    <div aria-busy="true">
      <div className="perfcheck-island-skel" aria-hidden="true">
        <div className="perfcheck-island-skel__badge" />
        <div className="perfcheck-island-skel__title" />
        <div className="perfcheck-island-skel__gauge perfcheck-island-skel__gauge--lg" />
        <div className="perfcheck-island-skel__gauge" />
        <div className="perfcheck-island-skel__gauge" />
        <div className="perfcheck-island-skel__gauge" />
      </div>
      <p className="perfcheck-sr-only">Loading result…</p>
    </div>
  );
}
