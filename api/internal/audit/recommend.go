package audit

// rule is an open/closed recommendation rule: append to rules to extend behaviour.
type rule struct {
	ID        string
	Category  string
	Severity  string
	Title     string
	Detail    string
	Condition func(Result) bool
}

var rules = []rule{
	{
		ID:       "img-format",
		Category: "performance",
		Severity: "high",
		Title:    "Serve images in modern formats",
		Detail:   "Converting photographic images to AVIF or WebP would cut transfer size substantially.",
		Condition: func(r Result) bool {
			if r.Signals != nil && r.Signals.Bytes > 500_000 {
				return true
			}
			return r.Scores.Performance < 85
		},
	},
	{
		ID:       "cache-headers",
		Category: "performance",
		Severity: "medium",
		Title:    "Add long-lived cache headers for static assets",
		Detail:   "Fingerprinted CSS and JS can use Cache-Control: public, max-age=31536000, immutable.",
		Condition: func(r Result) bool {
			if r.Signals != nil && r.Signals.CacheControl == "" {
				return r.Scores.Performance < 95
			}
			return r.Scores.Performance < 90
		},
	},
	{
		ID:       "meta-description",
		Category: "seo",
		Severity: "medium",
		Title:    "Ensure every page has a unique meta description",
		Detail:   "Unique descriptions improve snippet quality in search results and social previews.",
		Condition: func(r Result) bool {
			if r.Signals != nil && !r.Signals.HasMetaDescription {
				return true
			}
			return r.Scores.SEO < 92
		},
	},
	{
		ID:       "heading-order",
		Category: "seo",
		Severity: "low",
		Title:    "Keep a single H1 and avoid skipped heading levels",
		Detail:   "A clear heading outline helps both assistive tech and crawlers understand structure.",
		Condition: func(r Result) bool {
			if r.Signals != nil && r.Signals.H1Count != 1 {
				return true
			}
			return r.Scores.SEO < 88
		},
	},
	{
		ID:       "focus-visible",
		Category: "accessibility",
		Severity: "high",
		Title:    "Keep visible focus styles on interactive controls",
		Detail:   "Never remove outline without providing an equivalent :focus-visible treatment.",
		Condition: func(r Result) bool {
			return r.Scores.Accessibility < 90
		},
	},
	{
		ID:       "form-labels",
		Category: "accessibility",
		Severity: "high",
		Title:    "Associate every input with a visible label",
		Detail:   "Use <label for> or wrapping labels so screen readers announce purpose clearly.",
		Condition: func(r Result) bool {
			if r.Signals != nil && r.Signals.InputsWithoutLabel > 0 {
				return true
			}
			return r.Scores.Accessibility < 85
		},
	},
	{
		ID:       "missing-alt",
		Category: "accessibility",
		Severity: "high",
		Title:    "Give every image a meaningful alt attribute",
		Detail:   "Decorative images can use alt=\"\"; informative images need a short description.",
		Condition: func(r Result) bool {
			return r.Signals != nil && r.Signals.ImagesMissingAlt > 0
		},
	},
	{
		ID:       "missing-lang",
		Category: "accessibility",
		Severity: "medium",
		Title:    "Set the document language on the html element",
		Detail:   "A lang attribute lets screen readers choose the correct pronunciation.",
		Condition: func(r Result) bool {
			return r.Signals != nil && !r.Signals.HasLang
		},
	},
	{
		ID:       "overall-pass",
		Category: "performance",
		Severity: "low",
		Title:    "Solid baseline — keep measuring on real traffic",
		Detail:   "Scores look healthy. Re-run with a lab tool on production traffic for evidence.",
		Condition: func(r Result) bool {
			return r.Scores.Overall >= 85
		},
	},
}

// Recommend returns recommendations for a scoring result.
func Recommend(res Result) []Recommendation {
	out := make([]Recommendation, 0, len(rules))
	for _, r := range rules {
		if r.Condition(res) {
			out = append(out, Recommendation{
				ID:       r.ID,
				Category: r.Category,
				Severity: r.Severity,
				Title:    r.Title,
				Detail:   r.Detail,
			})
		}
	}
	return out
}
