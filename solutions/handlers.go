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
	meta := ui.PageMeta{
		Title:       "AI Platform Solutions",
		Description: "Production-grade, embeddable AI platform components: AI Engine, Multi Modal Chat, Collaborative Editing, and Search & RAG.",
		Path:        "/solutions",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/overview", nil, meta)
}

func AIEngineHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-ai-engine.md")
	meta := ui.PageMeta{
		Title:       "AI Engine — Agent Runtime & Orchestration",
		Description: "A complete AI orchestration layer: multi-provider LLM support, tool calling, MCP protocol, voice I/O, and streaming for your product.",
		Path:        "/solutions/ai-engine",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"Title":        "AI Engine",
		"Subtitle":     "From chat assistants to autonomous agents -- a complete AI orchestration layer",
		"Content":      content,
		"Slug":         "ai-engine",
		"GradientFrom": "from-violet-600",
		"GradientTo":   "to-indigo-700",
		"Icon":         "ai-engine",
	}, meta)
}

func ChatCollaborationHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-chat-collaboration.md")
	meta := ui.PageMeta{
		Title:       "Multi Modal Chat",
		Description: "A unified communication layer where humans and AI agents collaborate naturally with rich media, voice messages, and interactive workflows.",
		Path:        "/solutions/chat-collaboration",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"Title":        "Multi Modal Chat",
		"Subtitle":     "A complete messaging backbone where humans and AI agents communicate as equals",
		"Content":      content,
		"Slug":         "chat-collaboration",
		"GradientFrom": "from-blue-600",
		"GradientTo":   "to-cyan-600",
		"Icon":         "chat",
	}, meta)
}

func CollaborativeEditingHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-collaborative-editing.md")
	meta := ui.PageMeta{
		Title:       "Collaborative Document Editing",
		Description: "Multi-user real-time editing with AI assistance, CRDT sync, live cursors, and version history — embeddable into any product.",
		Path:        "/solutions/collaborative-editing",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"Title":        "Collaborative Document Editing",
		"Subtitle":     "Multi-user real-time editing with AI assistance -- like Google Docs, built for your product",
		"Content":      content,
		"Slug":         "collaborative-editing",
		"GradientFrom": "from-emerald-600",
		"GradientTo":   "to-teal-600",
		"Icon":         "editor",
	}, meta)
}

func SearchRAGHandler(w http.ResponseWriter, r *http.Request) {
	content := readMarkdown("docs/deep-dive-search-rag.md")
	meta := ui.PageMeta{
		Title:       "Search & RAG",
		Description: "Semantic search and intelligent retrieval that makes your entire knowledge base AI-accessible with hybrid vector and keyword search.",
		Path:        "/solutions/search-rag",
		OGType:      "website",
	}
	ui.RenderPageWithMeta(w, r, "solutions/detail", map[string]interface{}{
		"Title":        "Search & RAG",
		"Subtitle":     "Semantic search and intelligent retrieval that makes your entire knowledge base AI-accessible",
		"Content":      content,
		"Slug":         "search-rag",
		"GradientFrom": "from-amber-500",
		"GradientTo":   "to-orange-600",
		"Icon":         "search",
	}, meta)
}
