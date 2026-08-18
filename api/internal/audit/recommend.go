package audit

// rule is an open/closed recommendation rule: append to rules to extend behaviour.
type rule struct {
	ID        string
	Category  string
	Severity  string
	Title     string
	Detail    string
	Condition func(Scores) bool
}

var rules = []rule{
	{
		ID:       "img-format",
		Category: "performance",
		Severity: "high",
		Title:    "Serve images in modern formats",
		Detail:   "Converting photographic images to AVIF or WebP would cut transfer size substantially.",
		Condition: func(s Scores) bool {
			return s.Performance < 85
		},
	},
	{
		ID:       "cache-headers",
		Category: "performance",
		Severity: "medium",
		Title:    "Add long-lived cache headers for static assets",
		Detail:   "Fingerprinted CSS and JS can use Cache-Control: public, max-age=31536000, immutable.",
		Condition: func(s Scores) bool {
			return s.Performance < 90
		},
	},
	{
		ID:       "meta-description",
		Category: "seo",
		Severity: "medium",
		Title:    "Ensure every page has a unique meta description",
		Detail:   "Unique descriptions improve snippet quality in search results and social previews.",
		Condition: func(s Scores) bool {
			return s.SEO < 92
		},
	},
	{
		ID:       "heading-order",
		Category: "seo",
		Severity: "low",
		Title:    "Keep a single H1 and avoid skipped heading levels",
		Detail:   "A clear heading outline helps both assistive tech and crawlers understand structure.",
		Condition: func(s Scores) bool {
			return s.SEO < 88
		},
	},
	{
		ID:       "focus-visible",
		Category: "accessibility",
		Severity: "high",
		Title:    "Keep visible focus styles on interactive controls",
		Detail:   "Never remove outline without providing an equivalent :focus-visible treatment.",
		Condition: func(s Scores) bool {
			return s.Accessibility < 90
		},
	},
	{
		ID:       "form-labels",
		Category: "accessibility",
		Severity: "high",
		Title:    "Associate every input with a visible label",
		Detail:   "Use <label for> or wrapping labels so screen readers announce purpose clearly.",
		Condition: func(s Scores) bool {
			return s.Accessibility < 85
		},
	},
	{
		ID:       "overall-pass",
		Category: "performance",
		Severity: "low",
		Title:    "Solid baseline — keep measuring on real traffic",
		Detail:   "Mock scores look healthy. Wire a real lab tool when you need production evidence.",
		Condition: func(s Scores) bool {
			return s.Overall >= 85
		},
	},
}

// Recommend returns recommendations for the given scores.
func Recommend(scores Scores) []Recommendation {
	out := make([]Recommendation, 0, len(rules))
	for _, r := range rules {
		if r.Condition(scores) {
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
