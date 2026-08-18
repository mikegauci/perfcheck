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

// Audit is a completed audit result.
type Audit struct {
	ID              string           `json:"id"`
	URL             string           `json:"url"`
	CreatedAt       time.Time        `json:"createdAt"`
	Engine          string           `json:"engine"`
	Scores          Scores           `json:"scores"`
	Signals         *Signals         `json:"signals,omitempty"`
	Recommendations []Recommendation `json:"recommendations"`
}
