# Deep Dive: Real-Time Chat & Collaboration

## What is an AI-Native Chat System?

Traditional chat systems move text between humans. An AI-native chat system treats **AI agents as first-class participants** -- they read messages, understand context, stream responses, call tools, and present interactive elements, all within the same conversation thread that humans use.

This is not a separate "AI chat" bolted onto your product. It is a unified communication layer where humans and AI collaborate naturally.

---

## Architecture at a Glance

```
┌──────────────────────────────────────────────────────┐
│                     Chat System                       │
│                                                       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────┐ │
│  │   Humans    │  │  AI Agents  │  │ CLI Agents   │ │
│  │  (users)    │  │ (assistants)│  │  (devices)   │ │
│  └──────┬──────┘  └──────┬──────┘  └──────┬───────┘ │
│         │                │                 │          │
│  ┌──────▼─────────────────▼─────────────────▼──────┐ │
│  │              Message Stream                      │ │
│  │  text | voice | files | interactive elements     │ │
│  └──────────────────┬──────────────────────────────┘ │
│                     │                                 │
│  ┌──────────────────▼──────────────────────────────┐ │
│  │           Real-Time Delivery (SSE)               │ │
│  │   streaming responses | presence | notifications │ │
│  └──────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

---

## Conversation Types

The Chat system supports four conversation models, covering every communication pattern:

| Type | Description | Example |
|---|---|---|
| **Private (User)** | One-to-one messaging between two users | Direct message to a colleague |
| **Private (Agent)** | One-to-one conversation with an AI agent | Asking the "Product Analyst" agent a question |
| **Group** | Multi-participant conversations | A project team channel with both users and agents |
| **Contextual** | Chat embedded within a specific feature or page | Discussion thread on a feature specification |

Each conversation type shares the same underlying infrastructure -- messages, attachments, real-time delivery, and AI participation all work identically regardless of the conversation type.

---

## Message Types & Content

### Text Messages

Standard text messages with:
- Full markdown rendering (headings, lists, code blocks, links, bold, italic)
- @-mentions with clickable user/agent references
- Copy-to-clipboard for any message

### Voice Messages

Complete voice messaging workflow:

```
Record → Preview → Send → Auto-Transcribe
```

1. **Recording**: Browser-based audio capture with real-time volume visualization and timer
2. **Preview**: Play back before sending, with option to discard and re-record
3. **Delivery**: Audio uploaded to cloud storage, message delivered with playback controls
4. **Transcription**: Background speech-to-text converts audio to searchable text automatically

Voice messages are particularly valuable for:
- Mobile-first users who need hands-free communication
- Field workers reporting from job sites
- Accessibility -- users who find typing difficult
- Capturing nuance that text loses

### File Attachments

Drag-and-drop or click-to-upload with:
- **Image support**: JPEG, PNG, GIF, WebP with inline rendering and hover overlay
- **Document support**: PDF, TXT, Markdown with file cards showing name, type, and size
- **Multiple files** per message
- **Preview strip** before sending -- review and remove files before they go out
- **Download button** on received files

Upload uses a secure 3-step flow:
1. Request a presigned upload URL from the server
2. Upload directly to cloud storage (S3)
3. Confirm the upload and attach metadata to the message

This design means files never pass through the application server, reducing load and latency.

### System Messages

Automated notifications for conversation events:
- Group creation ("Alice created this group")
- Participant changes ("Bob was added to the conversation")
- Custom system events from integrations

---

## Interactive Elements

This is where the Chat system goes beyond traditional messaging. Messages can contain **actionable components** that users interact with directly:

| Element Type | Description | Example |
|---|---|---|
| **Diff Proposal** | AI suggests a code or document change. User can accept or reject inline. | Agent proposes a fix: "Change line 42 from X to Y" -- user clicks Accept |
| **Streaming Insert** | AI streams content into a document in real-time, with completion/rejection controls | Agent writes a feature spec section while the user watches |
| **Form** | Structured data collection within the chat flow | Agent asks for configuration parameters via a form |
| **Confirmation** | Yes/no decision point with context | Agent asks "Should I deploy this to staging?" |
| **Question Set** | Multi-question survey with checkbox and radio options | Onboarding questionnaire with AI-driven follow-up |

Each interactive element tracks its own state (`pending` → `accepted`/`rejected`/`completed`/`failed`), and state changes trigger events that other systems can react to.

**Why this matters:** Users don't need to leave the conversation to take action. The AI proposes, the user decides, and the system executes -- all in one place.

---

## Mentions & Tagging

The @-mention system connects humans and AI:

1. **Type `@`** in the input to open the suggestion menu
2. **Filter** by typing -- suggestions narrow in real time
3. **Select** a user, AI agent, or CLI device
4. **The mention is embedded** in the message as a structured tag

When a message contains an AI agent mention, the system automatically:
- Routes the message to the agent's processing pipeline
- Includes recent conversation history as context
- Starts streaming the agent's response into the chat

### Sticky Agents

For ongoing conversations with a specific agent, the "sticky agent" feature automatically includes the agent in every message without requiring manual @-mentions. This is stored per-conversation and persists across sessions.

---

## Real-Time Delivery

Messages are delivered in real-time using **Server-Sent Events (SSE)**:

```
Client                          Server
  │                                │
  │── GET /api/chat/events ───────>│
  │                                │
  │<── event: message ────────────│  (new message)
  │<── event: message ────────────│  (streaming chunk)
  │<── event: heartbeat ──────────│  (keep-alive)
  │<── event: message ────────────│  (another message)
  │                                │
  │── (connection timeout) ───────>│
  │── GET /api/chat/events ───────>│  (auto-reconnect)
  │                                │
```

**Intelligent polling adapts to context:**
- **During AI streaming**: 1-second polling for smooth token delivery
- **Normal operation**: 3-second polling for efficiency
- **Mobile devices**: Shorter max connection time (2 minutes) to conserve battery
- **Connection timeout**: 10-minute max with automatic reconnection

**Change detection** is optimized:
- Content updates detected for streaming responses
- Context changes detected for metadata updates
- New message detection for standard delivery
- Only changed data is transmitted

---

## Streaming AI Responses

When an AI agent responds, the answer streams into the chat token-by-token:

1. A placeholder message appears immediately with the agent's avatar
2. Text flows in progressively as the model generates it
3. Tool calls are shown in real-time (so users see what the agent is doing)
4. The message finalizes when generation completes

A **stream buffer** manages this flow:
- Chunks are collected and flushed every 200ms
- Atomic database writes prevent race conditions
- Sequential write guarantees maintain message coherence

This creates a natural, responsive experience -- users see the AI "thinking" rather than waiting for a complete response.

---

## Notifications & Unread Tracking

### Unread Counts

Every conversation tracks a **last read timestamp** per user. Unread counts are computed by comparing this timestamp against new messages, excluding system messages.

### Sound Notifications

Configurable sound alerts for new messages with user preference management.

### Real-Time Notification Events

A dedicated SSE channel delivers notification events (new message previews, sender info) so the UI can update badges and indicators without polling conversation lists.

---

## Group Chat Management

Full lifecycle management for multi-participant conversations:

| Operation | Description |
|---|---|
| **Create** | Select participants (users + agents), name the group, auto-add creator |
| **Invite** | Add new participants with system message notification |
| **Leave** | Remove yourself with confirmation dialog |
| **Edit** | Update group name and participant list |

Groups support the same features as private chats: file attachments, voice messages, AI agents, interactive elements, and mentions.

---

## Pagination & History

Chat history uses **cursor-based pagination** for efficient loading:

- **Initial load**: 10 most recent messages
- **Load more**: 15 messages per page, triggered by scrolling near the top
- **Bidirectional**: Load older or newer messages relative to a cursor
- **Scroll management**: Auto-scroll to bottom on new messages, "New messages" badge when scrolled up

---

## Usage Analytics

The Chat system provides built-in analytics:

| Metric | Description |
|---|---|
| **Total messages** | Overall message volume |
| **Unique conversations** | Number of active threads |
| **Active users** | Distinct human participants |
| **Messages by member** | Per-user message counts |
| **Messages by agent** | Per-agent message counts |
| **Daily breakdown** | Members vs. agents activity over time |
| **Per-project breakdown** | Message volume by project |

These metrics enable:
- Understanding team communication patterns
- Measuring AI agent adoption and effectiveness
- Identifying high-activity areas that may need more resources
- Reporting for management and stakeholders

---

## Security & Multi-Tenancy

| Feature | Description |
|---|---|
| **Project isolation** | Messages and conversations are scoped to projects |
| **Authentication** | Session-based (NextAuth), Personal Access Tokens, API keys |
| **Authorization** | Write operations require authentication; read access is configurable |
| **Audit trail** | Every message has sender attribution, timestamps, and context metadata |

---

## Message Context & Extensibility

Every message carries a **context object** that can include:

- Tagged user references
- Streaming state indicators
- CLI agent metadata (device name, model, response time)
- Reply-to references (threading foundation)
- Interactive element arrays
- Attachment metadata
- Voice metadata (MIME type, duration, transcript status)
- Timezone information
- Search indexing state

This extensible context model means new features can be added to messages without schema migrations -- context fields are stored as JSON and interpreted by the UI layer.

---

## Summary

The Chat system is not just messaging. It is a **collaboration runtime** that:

1. **Unifies humans and AI** in the same conversation threads
2. **Supports rich media** -- text, voice, files, and interactive components
3. **Delivers in real-time** with adaptive streaming and notifications
4. **Enables action** through interactive elements that let users decide and execute without leaving the chat
5. **Tracks everything** with analytics, audit trails, and unread management
6. **Scales across projects** with multi-tenant isolation and configurable access

This is the communication layer that modern AI-augmented products need -- where every conversation is an opportunity for intelligent automation.
