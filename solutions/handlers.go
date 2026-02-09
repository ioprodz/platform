package solutions

import (
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

func OverviewHandler(w http.ResponseWriter, r *http.Request) {
	ui.RenderPage(w, r, "solutions/overview", nil)
}

func AIEngineHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-ai-engine.md")
	ui.RenderPage(w, r, "solutions/detail", map[string]interface{}{
		"Title":       "AI Engine",
		"Subtitle":    "From chat assistants to autonomous agents -- a complete AI orchestration layer",
		"Content":     content,
		"Slug":        "ai-engine",
		"GradientFrom": "from-violet-600",
		"GradientTo":   "to-indigo-700",
		"Icon":        "ai-engine",
	})
}

func ChatCollaborationHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-chat-collaboration.md")
	ui.RenderPage(w, r, "solutions/detail", map[string]interface{}{
		"Title":       "Multi Modal Chat",
		"Subtitle":    "A complete messaging backbone where humans and AI agents communicate as equals",
		"Content":     content,
		"Slug":        "chat-collaboration",
		"GradientFrom": "from-blue-600",
		"GradientTo":   "to-cyan-600",
		"Icon":        "chat",
	})
}

func CollaborativeEditingHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-collaborative-editing.md")
	ui.RenderPage(w, r, "solutions/detail", map[string]interface{}{
		"Title":       "Collaborative Document Editing",
		"Subtitle":    "Multi-user real-time editing with AI assistance -- like Google Docs, built for your product",
		"Content":     content,
		"Slug":        "collaborative-editing",
		"GradientFrom": "from-emerald-600",
		"GradientTo":   "to-teal-600",
		"Icon":        "editor",
	})
}

func SearchRAGHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-search-rag.md")
	ui.RenderPage(w, r, "solutions/detail", map[string]interface{}{
		"Title":       "Search & RAG",
		"Subtitle":    "Semantic search and intelligent retrieval that makes your entire knowledge base AI-accessible",
		"Content":     content,
		"Slug":        "search-rag",
		"GradientFrom": "from-amber-500",
		"GradientTo":   "to-orange-600",
		"Icon":        "search",
	})
}
