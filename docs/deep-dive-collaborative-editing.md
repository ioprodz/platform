# Deep Dive: Collaborative Document Editing

## What is a Collaborative Editor?

A collaborative editor lets multiple people -- and AI agents -- edit the same document at the same time, with changes appearing instantly for everyone. No "save and refresh." No "your version vs. my version." No merge conflicts.

This is the same technology behind Google Docs, Notion, and Figma. Polysee provides it as an embeddable component that any product can use.

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────┐
│                  Collaborative Editor                │
│                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │  User A   │  │  User B   │  │ AI Agent  │         │
│  │ (browser) │  │ (browser) │  │ (server)  │         │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘       │
│        │               │               │             │
│  ┌─────▼───────────────▼───────────────▼─────┐      │
│  │        Hocuspocus Sync Server              │      │
│  │   (WebSocket, conflict resolution)         │      │
│  └─────────────────────┬─────────────────────┘      │
│                        │                             │
│  ┌─────────────────────▼─────────────────────┐      │
│  │         Yjs (CRDT Document Model)          │      │
│  │   Automatic conflict-free merging          │      │
│  └─────────────────────┬─────────────────────┘      │
│                        │                             │
│  ┌─────────────────────▼─────────────────────┐      │
│  │           PostgreSQL Storage               │      │
│  │   Binary document state + snapshots        │      │
│  └───────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────┘
```

---

## How Real-Time Sync Works

### The Problem

When two people edit the same paragraph simultaneously, whose edit "wins"? Traditional systems use locking ("someone else is editing this section") or last-write-wins (which silently discards changes). Both are bad.

### The Solution: CRDTs

Polysee uses **Yjs**, a CRDT (Conflict-free Replicated Data Type) library. Here's how it works in plain terms:

1. Each user's editor maintains a local copy of the document
2. Every keystroke generates an **operation** (insert character at position X, delete character at position Y)
3. Operations are sent to the **Hocuspocus sync server** via WebSocket
4. The server broadcasts operations to all other connected editors
5. Yjs's CRDT algorithm **guarantees** that all editors converge to the same document state, regardless of the order operations arrive

**No conflicts. No data loss. No locking.** Two users can type in the same word at the same time, and the result is always consistent.

### Connection Management

- **WebSocket connection** to the Hocuspocus server (port 3001)
- **Automatic reconnection** on disconnect
- **Status tracking**: connecting, connected, disconnected
- **Timeout handling** with configurable defaults (5-10 seconds)

---

## What You Can Write

The editor is built on **Tiptap** (a modern wrapper around ProseMirror) and supports rich content:

| Content Type | Description |
|---|---|
| **Rich text** | Bold, italic, underline, strikethrough |
| **Headings** | H1, H2, H3 hierarchy |
| **Lists** | Bullet lists and numbered/ordered lists |
| **Code blocks** | Syntax highlighting for 50+ programming languages, with language selector |
| **Inline code** | Monospace formatting with custom highlighting |
| **Block quotes** | Indented citation blocks |
| **Tables** | Full table support with resizable columns |
| **Images** | Embedded images via upload, drag-and-drop, or clipboard paste |
| **File attachments** | PDF, Word, Excel, Markdown files with metadata display |
| **Horizontal rules** | Visual section separators |
| **Links** | Clickable URLs with custom styling |

### Slash Commands

Type `/` anywhere in the document to open a command menu:

| Command | Action |
|---|---|
| `/image` | Upload and embed an image |
| `/file` | Attach a document (PDF, Word, Excel, Markdown) |
| `/h1`, `/h2`, `/h3` | Insert heading |
| `/bullet` | Start a bullet list |
| `/numbered` | Start a numbered list |
| `/code` | Insert a code block |
| `/quote` | Insert a block quote |
| `/bold`, `/italic` | Apply text formatting |

Commands filter in real-time as you type, so `/co` shows "Code block" immediately.

### Smart Paste

- **Paste markdown** and it converts to rich text automatically
- **Paste images** from clipboard -- they upload and embed inline
- **Drop files** onto the editor -- images embed, documents attach

---

## Presence & Awareness

When multiple people are in the same document, everyone can see who's there and what they're doing:

### Active Users

A presence indicator in the top corner shows:
- Avatars of all connected users
- Count of active editors
- Deduplication for users with multiple browser tabs

### Live Cursors

Each user's cursor appears in the document in real-time:
- **Color-coded** so you can tell users apart
- **Labeled** with the user's name
- **AI cursors** have a distinct pulsing animation to differentiate bot edits from human edits

This awareness prevents the "are you still editing?" problem and makes collaboration feel natural.

---

## Version History

Every significant change to a document is preserved as a **snapshot**:

### Automatic Snapshots

The system creates snapshots automatically based on:
- **Content hashing** (SHA-256): Only truly different content triggers a new snapshot
- **Word count tracking**: Each snapshot records the document length
- **Creator attribution**: Who made the change (user or AI agent, with name and avatar)

### History Panel

A visual timeline lets users:
1. **Browse** all past versions with timestamps and author info
2. **Preview** any snapshot without affecting the current document
3. **Restore** any version with one click

Restoring creates a new snapshot at the restoration point, so you never lose the current state -- you can always "undo the undo."

### Snapshot Storage

Each snapshot stores:
- Full binary document state (Yjs format)
- Content hash for deduplication
- Word count
- Creator info (user or agent, with avatar)
- Timestamp

---

## AI-Assisted Editing

This is where the collaborative editor becomes truly powerful. AI agents can edit documents programmatically, with the same real-time sync that human users enjoy.

### Streaming Insert

AI can write content into the document in real-time:

```
AI Agent                        Document
   │                               │
   │── start ─────────────────────>│  (cursor appears)
   │── position(paragraph 3) ─────>│  (cursor moves)
   │── chunk("The system") ───────>│  (text appears)
   │── chunk(" provides") ────────>│  (text continues)
   │── chunk(" real-time") ───────>│  (text continues)
   │── complete ──────────────────>│  (cursor disappears)
   │                               │
```

Users watch the AI write in real-time, just like watching a colleague type. They can continue editing other parts of the document simultaneously.

### Diff Application

AI can propose structured changes to existing content:

| Diff Format | Description |
|---|---|
| **Markdown Replace** | Replace the entire document content with new markdown |
| **Unified Diff** | Insert, delete, or replace specific sections |
| **JSON Patch** | Structured patches for precise modifications |

Each diff application records:
- Who applied it (user or agent)
- Source type (ai-bot, user, system)
- Description of what changed
- Timestamp

### Bot Client

For server-side AI editing, a headless bot client connects to the document:
1. Generates a JWT token for authentication
2. Opens a Yjs connection to the Hocuspocus server
3. Applies edits programmatically
4. Appears as a named agent in the presence list

This means AI agents can edit documents even when no human has the document open -- useful for background tasks like documentation generation or content migration.

---

## File & Image Handling

### Image Upload

Three ways to add images:
1. **Slash command** (`/image`) -- opens file picker
2. **Drag and drop** -- drop image files onto the editor
3. **Clipboard paste** -- paste screenshots directly

Images auto-resize to a maximum width of 800px and use the same secure 3-step S3 upload flow as chat attachments.

### File Attachments

Support for non-image files:
- PDF, Word (.docx), Excel (.xlsx), Markdown (.md)
- Size limit: 10MB per file
- Files display as cards with filename, type, and size
- Click to download

### Drop Zone

When dragging files over the editor, a visual overlay appears confirming the drop target. This prevents accidental drops and provides clear feedback.

---

## Export & Conversion

Documents can be extracted in multiple formats:

| Format | Use Case |
|---|---|
| **Plain Text** | Search indexing, simple display |
| **Markdown** | Developer documentation, static site generation |
| **HTML** | Web rendering, email content |
| **JSON** | Programmatic access, ProseMirror schema |

Bidirectional markdown conversion means you can:
- Import existing markdown content into the editor
- Export collaborative documents as markdown for version control
- Round-trip content between the editor and other systems

---

## Access Control

Document access is controlled at multiple levels:

| Level | Control |
|---|---|
| **Authentication** | Session-based, Personal Access Tokens, API keys, JWT (for bots) |
| **Project scope** | Documents are associated with projects; project membership grants access |
| **Role-based** | Owners and maintainers get read/write; other roles get read-only |
| **IDOR protection** | Feature ownership is validated against the project before granting access |

---

## Error Handling & Reliability

### Editor Error Boundary

A React Error Boundary wraps the editor component. If the editor crashes (malformed content, extension error, etc.), a graceful fallback UI appears instead of a blank page. Users can reload without losing data (the server has the latest state).

### Connection Recovery

- Automatic reconnection on WebSocket disconnect
- Document state syncs on reconnect (Yjs handles the merge)
- Cursor and awareness data cleans up on disconnect (no ghost cursors)

### Data Safety

- All edits are persisted to PostgreSQL in binary Yjs format
- Snapshots provide point-in-time recovery
- Content hashing prevents duplicate snapshots from inflating storage

---

## REST API

The editor is fully accessible via API for programmatic workflows:

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/api/features/{id}/documentation` | Fetch current content as markdown |
| `PATCH` | `/api/features/{id}/documentation` | Replace content with new markdown |
| `POST` | `/api/features/{id}/documentation/apply-diff` | Apply a structured diff |
| `POST` | `/api/features/{id}/documentation/stream-insert` | Start a streaming AI insert (SSE) |
| `POST` | `/api/features/{id}/documentation/snapshot` | Create a manual snapshot |
| `GET` | `/api/features/{id}/documentation/history` | List all snapshots |
| `GET` | `/api/features/{id}/documentation/history/{snapshotId}` | Get snapshot details with preview |
| `POST` | `/api/features/{id}/documentation/restore` | Restore to a previous snapshot |

---

## Summary

The Collaborative Editor is not just a text area. It is a **real-time document platform** that:

1. **Enables simultaneous editing** with automatic conflict resolution via CRDTs
2. **Shows who's where** with live cursors and presence indicators
3. **Preserves history** with automatic snapshots and one-click restore
4. **Lets AI participate** as a real-time collaborator -- writing, editing, and proposing changes
5. **Supports rich content** -- text, code, tables, images, and file attachments
6. **Exports everywhere** -- markdown, HTML, JSON, plain text
7. **Stays reliable** with error boundaries, connection recovery, and persistent storage

This is the document infrastructure that products like Notion, Confluence, and Google Docs are built on -- available as an embeddable component for any application.
