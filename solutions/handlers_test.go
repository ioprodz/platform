package solutions

import (
	"os"
	"testing"
)

// A /fr/ or /ar/ page whose body falls back to English is a near-duplicate
// that Google crawls but doesn't index.
func TestEverySolutionHasItsDeepDiveInEachLanguage(t *testing.T) {
	for _, slug := range []string{"ai-engine", "chat-collaboration", "collaborative-editing", "search-rag"} {
		for _, lang := range []string{"fr", "ar"} {
			path := "../docs/deep-dive-" + slug + "." + lang + ".md"
			if info, err := os.Stat(path); err != nil || info.Size() < 1000 {
				t.Errorf("missing or empty %s", path)
			}
		}
	}
}
