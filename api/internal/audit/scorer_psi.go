package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	psiEndpoint = "https://www.googleapis.com/pagespeedonline/v5/runPagespeed"
	psiTimeout  = 60 * time.Second
)

// PSIScorer calls the PageSpeed Insights API.
type PSIScorer struct {
	Client *http.Client
	APIKey string
	Base   string // override for tests
}

func (s PSIScorer) client() *http.Client {
	if s.Client != nil {
		return s.Client
	}
	return &http.Client{Timeout: psiTimeout}
}

func (s PSIScorer) Score(ctx context.Context, pageURL string) (Result, error) {
	endpoint := s.Base
	if endpoint == "" {
		endpoint = psiEndpoint
	}
	q := url.Values{}
	q.Set("url", pageURL)
	q.Set("strategy", "mobile")
	q.Add("category", "PERFORMANCE")
	q.Add("category", "SEO")
	q.Add("category", "ACCESSIBILITY")
	if s.APIKey != "" {
		q.Set("key", s.APIKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("User-Agent", fetchUserAgent)

	res, err := s.client().Do(req)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("pagespeed status %d", res.StatusCode)
	}

	var payload psiResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return Result{}, fmt.Errorf("decode pagespeed: %w", err)
	}

	cats := payload.LighthouseResult.Categories
	perf := scalePSI(cats.Performance.Score)
	seo := scalePSI(cats.SEO.Score)
	a11y := scalePSI(cats.Accessibility.Score)
	if cats.Performance.Score == nil && cats.SEO.Score == nil && cats.Accessibility.Score == nil {
		return Result{}, fmt.Errorf("pagespeed response missing categories")
	}

	audits := payload.LighthouseResult.Audits
	sig := Signals{
		HTTPS:              true,
		PSIStrategy:        "mobile",
		HasTitle:           auditPassed(audits, "document-title"),
		HasMetaDescription: auditPassed(audits, "meta-description"),
		HasLang:            auditPassed(audits, "html-has-lang"),
		ImagesMissingAlt:   missingIfFailed(audits, "image-alt"),
		Bytes:              int(auditNumeric(audits, "total-byte-weight")),
	}

	return Result{
		Engine: EnginePSI,
		Scores: Scores{
			Performance:   perf,
			SEO:           seo,
			Accessibility: a11y,
			Overall:       (perf*40 + seo*30 + a11y*30) / 100,
		},
		Signals: &sig,
	}, nil
}

type psiResponse struct {
	LighthouseResult struct {
		Categories struct {
			Performance   psiCategory `json:"performance"`
			SEO           psiCategory `json:"seo"`
			Accessibility psiCategory `json:"accessibility"`
		} `json:"categories"`
		Audits map[string]psiAudit `json:"audits"`
	} `json:"lighthouseResult"`
}

type psiCategory struct {
	Score *float64 `json:"score"`
}

type psiAudit struct {
	Score         *float64 `json:"score"`
	NumericValue  float64  `json:"numericValue"`
}

func scalePSI(score *float64) int {
	if score == nil {
		return 0
	}
	return clampScore(int(*score*100 + 0.5))
}

func auditPassed(audits map[string]psiAudit, id string) bool {
	a, ok := audits[id]
	if !ok || a.Score == nil {
		return false
	}
	return *a.Score >= 1
}

func missingIfFailed(audits map[string]psiAudit, id string) int {
	if auditPassed(audits, id) {
		return 0
	}
	if _, ok := audits[id]; ok {
		return 1
	}
	return 0
}

func auditNumeric(audits map[string]psiAudit, id string) float64 {
	if a, ok := audits[id]; ok {
		return a.NumericValue
	}
	return 0
}
