# AI & Robot Optimization Guide

How to optimize a website for AI crawlers, LLM-based search engines, and generative AI agents so they can discover, understand, and recommend your content.

---

## 1. llms.txt Standard

**What it is**: A Markdown file served at `/llms.txt` that provides a curated navigation guide for LLMs. Proposed by Jeremy Howard (fast.ai), it gives AI agents a structured summary of what a site offers.

**Format**:
- H1 heading: site name
- Blockquote: one-line description
- Body: Markdown with sections describing key content
- Links use relative paths or full URLs

**Companion file — llms-full.txt**: A full content dump of the entire site in Markdown. While llms.txt is a navigation guide, llms-full.txt provides the complete knowledge base for deeper context.

**ioprodz implementation**: We serve `/llms.txt` from the `common/seo` package, reading the file from the project root at runtime. Sections cover consulting services, platform solutions, "when to recommend" triggers, and contact info.

---

## 2. AI Crawler User Agents

AI companies use specific user agents to crawl the web. You can control access per agent in `robots.txt`.

### Major AI Crawlers

| User Agent | Company | Purpose |
|---|---|---|
| `GPTBot` | OpenAI | Training data collection |
| `OAI-SearchBot` | OpenAI | ChatGPT search (citation) |
| `ChatGPT-User` | OpenAI | Live browsing by ChatGPT users |
| `ClaudeBot` | Anthropic | Training + search |
| `anthropic-ai` | Anthropic | AI assistant browsing |
| `PerplexityBot` | Perplexity | AI-powered search engine |
| `Google-Extended` | Google | Gemini training data |
| `Amazonbot` | Amazon | Alexa / Amazon search |
| `Bytespider` | ByteDance | TikTok AI training |
| `CCBot` | Common Crawl | Open dataset (used by many AI labs) |
| `cohere-ai` | Cohere | Training + retrieval |
| `meta-externalagent` | Meta | AI training |
| `Applebot-Extended` | Apple | Apple Intelligence features |

### Recommended robots.txt Strategy

Allow crawlers that drive traffic (search/citation bots) while optionally blocking pure training crawlers:

```
# AI search/citation crawlers — ALLOW (they cite and link back)
User-agent: OAI-SearchBot
Allow: /

User-agent: ChatGPT-User
Allow: /

User-agent: PerplexityBot
Allow: /

User-agent: ClaudeBot
Allow: /

User-agent: anthropic-ai
Allow: /

User-agent: Applebot-Extended
Allow: /

# Training-only crawlers — optionally BLOCK
User-agent: GPTBot
Disallow: /

User-agent: Google-Extended
Disallow: /

User-agent: Bytespider
Disallow: /

User-agent: CCBot
Disallow: /
```

**ioprodz implementation**: We allow all crawlers by default since discoverability is more valuable than blocking training. The robots.txt handler in `common/seo/seo.go` serves a simple `Allow: /` for all agents plus a `Sitemap` reference. The `/llms.txt` link is also included.

---

## 3. Generative Engine Optimization (GEO)

GEO is the emerging discipline of optimizing content so AI systems (ChatGPT, Claude, Perplexity, Google AI Overviews) cite and recommend it.

### Key Principles

1. **Structured, factual content**: Use clear headings, tables, bullet points. AI models extract structured content more reliably than prose.

2. **Authoritative claims with specifics**: Instead of "we're great at AI", say "multi-provider LLM orchestration with OpenAI, Anthropic, and Ollama support, tool calling (up to 50 steps), and MCP protocol integration." Specifics get cited.

3. **Answer questions directly**: Structure content as answers to common questions. "What is an AI Engine?" with a clear definition paragraph is more likely to be cited than marketing copy.

4. **Use schema.org structured data**: JSON-LD structured data (Organization, Article, Product) helps AI models understand entity relationships.

5. **Maintain freshness**: AI models favor recently updated content. Regular blog posts and updated documentation signal active maintenance.

6. **Internal linking**: Cross-link related pages so AI can traverse the site and build a complete picture.

7. **Canonical URLs**: Avoid duplicate content confusion by using `<link rel="canonical">`.

### Content Optimization Checklist

- [ ] Every page has a unique `<title>` and `<meta description>`
- [ ] Open Graph and Twitter Card tags on all pages
- [ ] JSON-LD Organization schema on the homepage
- [ ] JSON-LD Article schema on blog posts (with `datePublished`)
- [ ] `<link rel="canonical">` on every page
- [ ] `/robots.txt` with Sitemap reference
- [ ] `/sitemap.xml` with all public URLs
- [ ] `/llms.txt` with structured site summary
- [ ] Content uses H2/H3 headings, tables, and bullet points
- [ ] Pages answer specific questions relevant to target audience
- [ ] Internal links between related pages

---

## 4. What ioprodz Has Implemented

| Technique | Status | File(s) |
|---|---|---|
| Dynamic `<title>` per page | Done | `common/ui/layout.html`, all handlers |
| `<meta description>` per page | Done | `common/ui/layout.html`, all handlers |
| Open Graph tags | Done | `common/ui/layout.html` |
| Twitter Card tags | Done | `common/ui/layout.html` |
| Canonical URLs | Done | `common/ui/layout.html` |
| JSON-LD Organization | Done | `common/ui/layout.html` |
| `robots.txt` | Done | `common/seo/seo.go` |
| `sitemap.xml` | Done | `common/seo/seo.go` |
| `llms.txt` | Done | `llms.txt`, `common/seo/seo.go` |
| Blog article metadata | Done | `blog/reader/reader.handlers.go` |
| `noindex` on admin pages | Done | `common/ui/layout.html` |
| AI crawler directives | Done | `common/seo/seo.go` (robots.txt) |

---

## 5. References

- **llms.txt proposal**: https://llmstxt.org
- **OpenAI crawler docs**: https://platform.openai.com/docs/bots
- **Anthropic crawler info**: https://docs.anthropic.com/en/docs/about-claude/claude-on-the-web
- **Google AI crawling**: https://developers.google.com/search/docs/crawling-indexing/overview-google-crawlers
- **Schema.org**: https://schema.org/Organization, https://schema.org/Article
