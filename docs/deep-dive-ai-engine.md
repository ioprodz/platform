# Deep Dive: AI Engine

## What is an AI Engine?

An AI Engine is the layer between your product and large language models (LLMs). Without it, integrating AI means writing custom code for every provider, managing API keys, handling streaming, building tool orchestration, and tracking costs -- all from scratch.

The AI Engine abstracts all of this into a single, configurable platform. You define **agents** (who the AI is), give them **tools** (what the AI can do), and connect them to **providers** (which model powers them). Everything else -- streaming, retries, cost tracking, multi-tenancy -- is handled for you.

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────┐
│                   Your Application                  │
├─────────────────────────────────────────────────────┤
│                                                     │
│   ┌──────────┐   ┌──────────┐   ┌──────────┐        │
│   │  Agent A │   │  Agent B │   │  Agent C │        │
│   │ "Helper" │   │ "Analyst"│   │ "Coder"  │        │
│   └────┬─────┘   └─────┬────┘   └──────┬───┘        │
│        │               │               │            │
│   ┌────▼───────────────▼───────────────▼────┐       │
│   │          Tool Calling Loop              │       │
│   │  (reason → call tool → observe → repeat)│       │
│   └────┬────────────────────────────────┬───┘       │
│        │                                │           │
│   ┌────▼──────┐                  ┌──────▼─────┐     │
│   │ Built-in  │                  │     MCP    │     │
│   │   Tools   │                  │ (external) │     │
│   └───────────┘                  └────────────┘     │
│                                                     │
├─────────────────────────────────────────────────────┤
│            Provider Abstraction Layer               │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│   │  OpenAI  │  │  Claude  │  │  Ollama  │          │
│   │  (cloud) │  │  (cloud) │  │  (local) │          │
│   └──────────┘  └──────────┘  └──────────┘          │
├─────────────────────────────────────────────────────┤
│   Usage Tracking  │  Cost Management  │  Analytics  │
└─────────────────────────────────────────────────────┘
```

---

## Key Concepts

### 1. Providers

A **provider** is a connection to an LLM service. Each provider has:
- An API endpoint and credentials
- A list of available models
- Pricing rules for cost estimation

**Supported providers:**

| Provider | Type | Models | Tool Calling | Streaming |
|---|---|---|---|---|
| **OpenAI** | Cloud | GPT-4o, GPT-4, GPT-3.5, etc. | Yes | Yes |
| **Anthropic** | Cloud | Claude 3.5 Sonnet, Claude 3, etc. | Yes | Yes |
| **Ollama** | Local | Mistral 7B, Llama, etc. | Limited | Yes |

The provider abstraction means your application code never references a specific model. You configure which provider an agent uses, and the Engine handles the translation between your request and each provider's API format.

**Why this matters:** If a new model launches tomorrow, you switch one configuration value. No code changes. No redeployment.

### 2. Agents

An **agent** is a named AI persona with:
- **Name and avatar** -- how it appears in the UI
- **System prompt** -- instructions that define its behavior and personality
- **Provider assignment** -- which LLM powers it (can override the project default)
- **Tool access** -- which tools it can use
- **MCP server connections** -- external tool servers it can reach

Agents are project-scoped, meaning each project can have its own set of agents with different configurations.

**Example:** A "Product Analyst" agent might use Claude for nuanced reasoning, have access to feature search tools and documentation tools, and be instructed to always cite sources. A "Quick Helper" agent might use a local Ollama model for fast, cost-free responses to simple questions.

### 3. Tool Calling

Tool calling is what makes AI agents useful beyond conversation. Instead of just generating text, an agent can:

1. **Reason** about what information it needs
2. **Call a tool** to retrieve data or perform an action
3. **Observe** the tool's result
4. **Repeat** until it has enough information to respond

This loop runs for up to 50 iterations per request, allowing agents to perform complex multi-step workflows autonomously.

**Built-in tools include:**
- `search_features` -- Search the knowledge base for relevant content
- `send_message` -- Send messages to users or channels
- `label_activity` -- Tag and categorize the agent's own work
- `get_tool_guide` -- Retrieve detailed instructions for complex tool usage

### 4. MCP (Model Context Protocol)

MCP is a standard protocol for connecting AI agents to external tool servers. Think of it as "USB for AI" -- any system that speaks MCP can instantly give an agent new capabilities.

**How it works:**
1. An MCP server exposes a set of tools via JSON-RPC 2.0
2. The agent discovers available tools at connection time
3. During conversation, the agent calls these tools just like built-in ones
4. Results flow back into the agent's reasoning loop

**Supported transports:**
- **stdio** -- Local process communication (for tools running on the same server)
- **HTTP** -- Remote server communication (for tools hosted elsewhere)

**Why this matters:** You can extend agent capabilities without modifying the AI Engine itself. Need an agent to query your CRM? Deploy an MCP server that wraps your CRM API. Need it to check inventory? Deploy another MCP server. The agent automatically discovers and uses the new tools.

---

## Activation Triggers

Agents don't just sit idle. They activate in response to events:

| Trigger | How It Works |
|---|---|
| **Chat Mention** | A user @-mentions an agent in chat. The agent receives the conversation context and responds. |
| **Survey Response** | A user completes a question set. The agent processes the answers and takes action. |
| **Calendar Event** | A scheduled event fires. The agent runs its workflow at the configured time. |
| **Coding Completion** | An external coding agent finishes a task. The AI agent processes the results and reports back. |

This event-driven design means agents integrate naturally into existing workflows. Users don't need to learn a new interface -- they just tag an agent in the conversation they're already having.

---

## Voice Capabilities

The AI Engine includes built-in voice processing:

- **Speech-to-Text**: Voice messages are automatically transcribed in the background. The transcript is stored alongside the audio, making voice content searchable and AI-readable.
- **Text-to-Speech**: When a user sends a voice message, the agent can respond with generated audio. The system prompt automatically adjusts for voice context (concise responses, no markdown formatting, 30-60 second duration).

---

## Usage Tracking & Cost Management

Every AI request is tracked with:

| Metric | Description |
|---|---|
| **Input tokens** | How many tokens were sent to the model |
| **Output tokens** | How many tokens the model generated |
| **Estimated cost** | Calculated from configurable pricing rules per provider/model |
| **Latency** | End-to-end response time |
| **Status** | Success, error, timeout, or rate-limited |
| **Feature type** | Which capability triggered the request (chat, title generation, system prompt, etc.) |
| **User attribution** | Which user initiated the request |

This data enables:
- **Budget controls** per project
- **Provider comparison** (cost vs. quality vs. speed)
- **Usage reports** for billing and capacity planning
- **Anomaly detection** for runaway costs

---

## Activity Tracking

Every agent interaction is logged as an **activity record** containing:
- The full system prompt and conversation context
- All tool calls with inputs, outputs, and iteration numbers
- Duration and completion status
- Associated feature or ticket context

This provides a complete audit trail of what the AI did, why, and how long it took -- essential for debugging, compliance, and trust.

---

## Multi-Tenancy

The AI Engine is designed for multi-tenant products from day one:

- **Project-scoped agents**: Each project has its own agents with independent configurations
- **Provider sharing**: Administrators can create system-level providers and grant access to specific projects
- **Usage isolation**: Token counts and costs are tracked per project
- **Configuration inheritance**: Projects can use global defaults or override them

---

## Integration Points

The AI Engine connects to other platform components:

| Component | Integration |
|---|---|
| **Chat** | Agents receive messages, stream responses, and use interactive elements |
| **Search & RAG** | Agents call `search_features` to ground responses in actual data |
| **Collaborative Editor** | Agents edit documents programmatically via streaming insert and diff application |
| **Calendar** | Scheduled triggers activate agent workflows |
| **File Storage** | Voice recordings and attachments flow through S3 |

---

## Summary

The AI Engine is not a chatbot. It is a complete agent runtime that:

1. **Abstracts providers** so you never lock into one model
2. **Orchestrates tools** so agents can take real actions
3. **Extends via MCP** so capabilities grow without code changes
4. **Tracks everything** so you understand cost, performance, and behavior
5. **Isolates tenants** so each project operates independently

This is the AI infrastructure layer that every modern product needs -- and building it from scratch would take months of dedicated engineering.
