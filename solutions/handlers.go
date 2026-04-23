package solutions

import (
	"ioprodz/common/i18n"
	"ioprodz/common/ui"
	"net/http"
	"os"
)

func readMarkdown(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// readMarkdownLocalized tries docs/deep-dive-<slug>.<lang>.md first, then falls
// back to the English docs/deep-dive-<slug>.md. Returns empty string if neither
// exists.
func readMarkdownLocalized(slug string, lang i18n.Lang) string {
	if lang != i18n.DefaultLang {
		if data, err := os.ReadFile("docs/deep-dive-" + slug + "." + string(lang) + ".md"); err == nil {
			return string(data)
		}
	}
	return readMarkdown("docs/deep-dive-" + slug + ".md")
}

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	meta := ui.PageMeta{
		TitleKey: "solutions:meta.overview.title",
		DescKey:  "solutions:meta.overview.description",
		Path:     "/solutions",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/overview", nil, meta)
}

func AIEngineHandler(w http.ResponseWriter, r *http.Request) {
	lang := ui.LangFrom(r)
	content := readMarkdownLocalized("ai-engine", lang)
	meta := ui.PageMeta{
		TitleKey: "solutions:meta.aiEngine.title",
		DescKey:  "solutions:meta.aiEngine.description",
		Path:     "/solutions/ai-engine",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"TitleKey":     "solutions:aiEngine.heroTitle",
		"SubtitleKey":  "solutions:aiEngine.heroSubtitle",
		"Content":      content,
		"Slug":         "ai-engine",
		"GradientFrom": "from-violet-600",
		"GradientTo":   "to-indigo-700",
		"Icon":         "ai-engine",
	}, meta)
}

func ChatCollaborationHandler(w http.ResponseWriter, r *http.Request) {
	lang := ui.LangFrom(r)
	content := readMarkdownLocalized("chat-collaboration", lang)
	meta := ui.PageMeta{
		TitleKey: "solutions:meta.chat.title",
		DescKey:  "solutions:meta.chat.description",
		Path:     "/solutions/chat-collaboration",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"TitleKey":     "solutions:chat.heroTitle",
		"SubtitleKey":  "solutions:chat.heroSubtitle",
		"Content":      content,
		"Slug":         "chat-collaboration",
		"GradientFrom": "from-blue-600",
		"GradientTo":   "to-cyan-600",
		"Icon":         "chat",
	}, meta)
}

func CollaborativeEditingHandler(w http.ResponseWriter, r *http.Request) {
	lang := ui.LangFrom(r)
	content := readMarkdownLocalized("collaborative-editing", lang)
	meta := ui.PageMeta{
		TitleKey: "solutions:meta.collabEditing.title",
		DescKey:  "solutions:meta.collabEditing.description",
		Path:     "/solutions/collaborative-editing",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"TitleKey":     "solutions:collabEditing.heroTitle",
		"SubtitleKey":  "solutions:collabEditing.heroSubtitle",
		"Content":      content,
		"Slug":         "collaborative-editing",
		"GradientFrom": "from-emerald-600",
		"GradientTo":   "to-teal-600",
		"Icon":         "editor",
	}, meta)
}

func SearchRAGHandler(w http.ResponseWriter, r *http.Request) {
	lang := ui.LangFrom(r)
	content := readMarkdownLocalized("search-rag", lang)
	meta := ui.PageMeta{
		TitleKey: "solutions:meta.searchRag.title",
		DescKey:  "solutions:meta.searchRag.description",
		Path:     "/solutions/search-rag",
		OGType:   "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"TitleKey":     "solutions:searchRag.heroTitle",
		"SubtitleKey":  "solutions:searchRag.heroSubtitle",
		"Content":      content,
		"Slug":         "search-rag",
		"GradientFrom": "from-amber-500",
		"GradientTo":   "to-orange-600",
		"Icon":         "search",
	}, meta)
}
