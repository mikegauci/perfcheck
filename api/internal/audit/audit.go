package audit

import "time"

// Scores holds the four numeric ratings returned by an audit.
type Scores struct {
	Overall       int `json:"overall"`
	Performance   int `json:"performance"`
	SEO           int `json:"seo"`
	Accessibility int `json:"accessibility"`
}

// Recommendation is a single actionable finding.
type Recommendation struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
}

// Audit is a completed mock audit result.
type Audit struct {
	ID              string           `json:"id"`
	URL             string           `json:"url"`
	CreatedAt       time.Time        `json:"createdAt"`
	Scores          Scores           `json:"scores"`
	Recommendations []Recommendation `json:"recommendations"`
}
