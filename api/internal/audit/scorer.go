package audit

import "context"

// Engine identifiers returned on every audit so the UI can label provenance.
const (
	EngineMock  = "mock"
	EngineFetch = "fetch"
	EnginePSI   = "psi"
)

// Scorer produces scores for a normalised URL.
type Scorer interface {
	Score(ctx context.Context, url string) (Result, error)
}

// Result is the output of a scoring engine.
type Result struct {
	Scores  Scores
	Engine  string
	Signals *Signals
}

// Signals are optional evidence collected by live engines.
// MockScorer leaves this nil.
type Signals struct {
	StatusCode         int    `json:"statusCode,omitempty"`
	TTFBMs             int64  `json:"ttfbMs,omitempty"`
	Bytes              int    `json:"bytes,omitempty"`
	ContentEncoding    string `json:"contentEncoding,omitempty"`
	CacheControl       string `json:"cacheControl,omitempty"`
	HTTPS              bool   `json:"https"`
	Title              string `json:"title,omitempty"`
	TitleLength        int    `json:"titleLength,omitempty"`
	HasTitle           bool   `json:"hasTitle"`
	MetaDescription    string `json:"metaDescription,omitempty"`
	HasMetaDescription bool   `json:"hasMetaDescription"`
	H1Count            int    `json:"h1Count"`
	HasCanonical       bool   `json:"hasCanonical"`
	HasOpenGraph       bool   `json:"hasOpenGraph"`
	HasJSONLD          bool   `json:"hasJsonLd"`
	HasLang            bool   `json:"hasLang"`
	ImagesMissingAlt   int    `json:"imagesMissingAlt"`
	HasViewport        bool   `json:"hasViewport"`
	InputsWithoutLabel int    `json:"inputsWithoutLabel"`
	ScriptCount        int    `json:"scriptCount"`
	StylesheetCount    int    `json:"stylesheetCount"`
	InlineStyleBytes   int    `json:"inlineStyleBytes"`
	PSIStrategy        string `json:"psiStrategy,omitempty"`
}
