# Deep Dive: Search & RAG (Retrieval-Augmented Generation)

## What is RAG?

Large language models are powerful, but they have a critical weakness: they only know what was in their training data. They don't know about your company's documents, your product's features, or your team's conversations.

**Retrieval-Augmented Generation (RAG)** solves this by giving AI access to your data at query time:

```
Without RAG:
  User: "What are our authentication requirements?"
  AI: "Authentication typically involves..." (generic, possibly wrong)

With RAG:
  User: "What are our authentication requirements?"
  System: [searches your docs, finds the auth spec]
  AI: "According to your Feature-Auth-2024 spec, authentication requires
       OAuth2 with PKCE flow, session timeout of 30 minutes, and MFA
       for admin roles." (specific, grounded, accurate)
```

The Search & RAG engine makes this possible by indexing your content, understanding its meaning, and retrieving the most relevant pieces when the AI needs them.

---

## Architecture at a Glance

```
┌──────────────────────────────────────────────────────┐
│                   Search & RAG Engine                 │
│                                                       │
│  ┌────────────────────────────────────────────────┐  │
│  │              Content Sources                    │  │
│  │  Documents | Features | Chats | Tickets | PRs  │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │ (events)                       │
│  ┌───────────────────▼────────────────────────────┐  │
│  │           Indexing Pipeline                      │  │
│  │  Extract → Hash → Embed → Store                 │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌───────────────────▼────────────────────────────┐  │
│  │        PostgreSQL + pgvector                     │  │
│  │  Full-text index | Vector embeddings             │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌───────────────────▼────────────────────────────┐  │
│  │           Hybrid Search                         │  │
│  │  (0.7 × vector score) + (0.3 × text score)     │  │
│  └───────────────────┬────────────────────────────┘  │
│                      │                                │
│  ┌──────────┐  ┌─────▼──────┐  ┌──────────────────┐ │
│  │  Search   │  │    RAG     │  │   AI Agent       │ │
│  │   API     │  │  Retrieve  │  │  Tool Calling    │ │
│  └──────────┘  └────────────┘  └──────────────────┘ │
└──────────────────────────────────────────────────────┘
```

---

## How It Works: Step by Step

### Step 1: Content Extraction

When content is created or modified anywhere in the platform, **extractors** convert it into searchable text. Each entity type has a specialized extractor:

| Entity Type | What Gets Extracted |
|---|---|
| **Documents** | Full document text from the collaborative editor (Yjs binary → plain text) |
| **Features** | Name, status, description, and full document content |
| **Use Cases** | Name, actors, acceptance criteria |
| **Pull Requests** | Title, state, type, branch names, metadata |
| **Assessments** | Name, description, questions, and iteration reviews |
| **Chat Messages** | Thread content (up to 50 messages per thread) |
| **Support Tickets** | Title, status, description, metadata |

This means search covers everything -- not just documents, but conversations, code reviews, tickets, and more.

### Step 2: Change Detection

Before generating a new embedding (which costs money), the system checks if the content actually changed:

1. Compute a **SHA-256 hash** of the extracted text
2. Compare against the stored hash for this entity
3. **If identical**: Skip re-embedding (saves cost)
4. **If different**: Generate a new embedding and update the hash

This ensures you only pay for embedding generation when content genuinely changes.

### Step 3: Embedding Generation

The extracted text is sent to an **embedding model** that converts it into a high-dimensional vector -- a numerical representation of its meaning:

```
"OAuth2 authentication with PKCE flow"
    → [0.023, -0.156, 0.891, 0.045, ..., -0.234]  (1536 dimensions)
```

**Supported embedding models:**

| Model | Dimensions | Best For |
|---|---|---|
| `text-embedding-3-small` | 1536 | Cost-effective general use |
| `text-embedding-3-large` | 3072 | Maximum accuracy |
| `text-embedding-ada-002` | 1536 | Legacy compatibility |

The model is configurable per project, so you can balance cost and quality based on your needs.

### Step 4: Storage

Both the original text and the embedding vector are stored in **PostgreSQL with pgvector**:

- **Full-text index** (`tsvector`): For keyword matching
- **Vector column**: For semantic similarity search
- **Metadata**: Entity type, entity ID, content hash, embedding model, timestamp

Using PostgreSQL for both relational data and vector search means no additional database to manage -- everything lives in one place.

### Step 5: Search

When a query comes in, the system runs **hybrid search** -- combining two different search strategies:

#### Vector Search (Semantic)
- The query is embedded using the same model
- pgvector finds the closest vectors by cosine distance
- This finds content that **means** the same thing, even if different words are used
- Example: searching "login security" matches content about "authentication mechanisms"

#### Full-Text Search (Keyword)
- PostgreSQL's built-in `tsvector` and `plainto_tsquery` find exact and partial word matches
- This catches content that uses the **exact terms** the user searched for
- Example: searching "OAuth2" matches content containing that specific term

#### Combined Scoring

```
final_score = (vector_score × 0.7) + (text_score × 0.3)
```

The 70/30 weighting prioritizes semantic understanding while still rewarding exact keyword matches. These weights are configurable per query.

**Why hybrid?** Pure semantic search misses exact terms (searching "JIRA-1234" should find that exact ticket). Pure keyword search misses meaning (searching "authentication" should find content about "login"). Hybrid gets the best of both.

---

## Event-Driven Indexing

Content is indexed **automatically** the moment it changes. No cron jobs. No manual triggers. No stale results.

The system listens for events from every content module:

| Module | Events |
|---|---|
| **Documents** | Created, Updated, Deleted |
| **Features** | Created, Name/Description/Status Updated, Archived, Deleted |
| **Use Cases** | Created, Updated, Deleted |
| **Pull Requests** | Created, Updated, Deleted |
| **Assessments** | Created, Updated, Deleted |
| **Chat Messages** | Created, Updated, Deleted |
| **Support Tickets** | Created, Updated, Deleted |

When an event fires:
1. The appropriate extractor generates the searchable text
2. The indexing pipeline checks if re-embedding is needed
3. If yes, a new embedding is generated and stored
4. The search index is immediately up-to-date

An **async queue** with configurable concurrency (default: 3) processes indexing jobs without blocking the main application.

---

## RAG Retrieval for AI Agents

The RAG retrieval function is the bridge between search and AI:

```
ragRetrieve(projectId, query, topK = 5)
  → [
      { content: "...", source: "Feature: Auth System", score: 0.92 },
      { content: "...", source: "Use Case: User Login", score: 0.87 },
      { content: "...", source: "Chat: Security Discussion", score: 0.71 },
    ]
```

When an AI agent needs information:
1. It calls the `search_features` tool with a natural language query
2. The tool runs hybrid search against the project's index
3. Top results are returned with content, source attribution, and relevance scores
4. The agent incorporates these results into its response, **citing sources**

This is what prevents hallucination. The agent doesn't guess -- it retrieves actual content from your knowledge base and references it.

---

## Search Modes

The system supports two operating modes:

### Text-Only Mode

- Uses only PostgreSQL full-text search
- No embedding generation (zero AI cost for indexing)
- Keyword matching only
- Good for: Small projects, cost-sensitive environments, or when semantic search isn't needed

### Hybrid Mode

- Combines vector search with full-text search
- Requires an embedding provider (OpenAI)
- Semantic + keyword matching
- Good for: Any project where users need to find content by meaning, not just exact words

Projects can switch between modes at any time. When switching to hybrid mode, a **background job** automatically generates embeddings for all existing text-only records.

---

## Configuration

### Per-Project Settings

Each project can configure:
- **Search mode**: text_only or hybrid
- **Embedding provider**: Which OpenAI provider to use
- **Embedding model**: Which model to generate embeddings with
- **Inheritance**: Use global defaults or override per project

### Global Defaults

Administrators can set system-wide defaults that projects inherit, ensuring consistent behavior unless explicitly overridden.

### Reindexing

Manual reindex can be triggered when:
- Switching embedding models (old embeddings need regeneration)
- Switching from text-only to hybrid mode
- Data integrity concerns

The reindex process:
1. Extracts content from all entities in the project
2. Processes in batches of 10
3. Tracks progress (total entities, indexed entities)
4. Reports errors without stopping the batch
5. Updates configuration status on completion

---

## Search Results

Each search result includes:

| Field | Description |
|---|---|
| **title** | Entity name or derived title |
| **snippet** | Relevant excerpt with context around the query match |
| **score** | Combined relevance score (0 to 1) |
| **vectorScore** | Semantic similarity component |
| **textScore** | Keyword match component |
| **entityType** | What kind of content (feature, use case, chat, etc.) |
| **url** | Direct link to the source entity |
| **metadata** | Entity-specific details (status, owner, etc.) |

Results link directly to the source content, so users can jump from a search result to the actual document, conversation, or ticket.

---

## API Reference

### Search Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| `POST` | `/api/projects/{id}/search` | Execute a hybrid search query |
| `PUT` | `/api/projects/{id}/search` | Trigger a full reindex |
| `GET` | `/api/projects/{id}/search` | Get current index status |

### Configuration Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/projects/{id}/search/config` | Get project search settings |
| `POST` | `/api/projects/{id}/search/config` | Update project search settings |
| `GET` | `/api/admin/search/embedding-config` | Get global defaults |
| `POST` | `/api/admin/search/embedding-config` | Set global defaults |
| `GET` | `/api/admin/search/reindex-status` | Monitor reindex progress |
| `GET` | `/api/admin/search/debug` | Debug search index state |

---

## Client-Side Integration

A React hook provides seamless search integration:

```
useSearch({ projectId, debounceMs: 300, limit: 20 })
```

Features:
- **Debounced queries**: Waits 300ms after the user stops typing before searching
- **Request cancellation**: Previous search is cancelled when a new query arrives
- **Result caching**: 30-second stale time, 60-second garbage collection
- **Prefetch support**: Pre-load results for predictable queries

---

## Technical Decisions

### Why PostgreSQL + pgvector (not a dedicated vector DB)?

- **One database**: No additional infrastructure to manage (no Pinecone, Weaviate, or Qdrant)
- **Transactional consistency**: Embeddings update in the same transaction as content
- **Mature ecosystem**: PostgreSQL's full-text search is battle-tested at any scale
- **Cost**: No additional SaaS subscription for a vector database

### Why Hybrid Search (not pure vector)?

- **Exact matches matter**: Searching for "PROJ-1234" needs keyword matching, not semantic similarity
- **Acronyms and codes**: Domain-specific terms often need exact matching
- **Configurable balance**: The 70/30 weight can be adjusted per query for different use cases

### Why Event-Driven Indexing (not batch)?

- **Real-time freshness**: Content is searchable immediately after creation
- **Cost efficiency**: Only changed content is re-embedded
- **No scheduled jobs**: No cron complexity, no "wait for the next indexing run"

---

## Summary

The Search & RAG engine is not just a search box. It is an **intelligent knowledge layer** that:

1. **Indexes everything** -- documents, conversations, tickets, code reviews, assessments
2. **Understands meaning** through semantic vector embeddings
3. **Stays current** with event-driven, real-time indexing
4. **Grounds AI responses** in actual data through the RAG retrieval API
5. **Balances precision and recall** with configurable hybrid search
6. **Scales simply** using PostgreSQL + pgvector (no additional database infrastructure)
7. **Respects budgets** with content hashing to avoid redundant embedding generation

This is how you make AI reliable in production: by giving it access to your real data, not just its training set.
