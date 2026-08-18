# Cross-browser checklist

Manual matrix for PerfCheck public pages and the audit island.

## Browsers

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- iOS Safari
- Android Chrome

## Widths

- 320px
- 768px
- 1024px
- 1440px

## Checks

- [ ] Primary nav works without JavaScript (checkbox toggle)
- [ ] Skip link appears on keyboard focus
- [ ] Audit form: keyboard submit, focus on invalid field path, live region announces results
- [ ] Score gauges expose text labels (not colour alone)
- [ ] Dashboard filter/sort usable by keyboard
- [ ] No horizontal overflow at 320px
- [ ] Focus rings visible on buttons, links and inputs

## Evidence

Run Lighthouse against a local production build (`hugo --minify` + `make api` with `-static public`) and paste scores into the README. Target: 95+ on Performance, Accessibility, Best Practices and SEO for marketing pages.
