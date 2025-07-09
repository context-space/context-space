# Context Space

<div align="center">

<!-- ![Context Space Banner](resources/banner.png) -->

[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![Docker](https://img.shields.io/badge/docker-supported-blue.svg)](https://docker.com)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![API Docs](https://img.shields.io/badge/API-documented-green.svg)](http://api.context.space/v1/docs)
[![Contributors](https://img.shields.io/badge/contributors-welcome-orange.svg)]()

**Context Space aims to be a comprehensive context engineering infrastructure for AI agents.** While the industry focuses on prompt engineering, we believe the next frontier lies in Context Engineering that provides AI systems with the right information, at the right time, in the right format.

**Context Space is built to be the backbone of context engineering for AI agents.**
While the industry focuses on prompt engineering, we believe the next frontier lies in Context Engineering that provides AI systems with the right information, at the right time, in the right format.
Join our [Discord](https://discord.gg/Q74Ta5Xv ) and help us build! Have questions? [Open an issue](https://github.com/context-space/context-space/issues) or join our discussions.

**🔗 Integrations** • **🔐 Credential Security** • **🚀 Production Ready** • **🤖 Context Engineering**

[Quick Start](#-quick-start) • [Live Demo](#-live-demo) • [API Documentation](http://api.context.space/v1/docs)

</div>

---

## What is Context Engineering?

**Context engineering** is the systematic design and management of all information surrounding an AI model during inference. While [prompt engineering](https://blog.langchain.com/context-engineering-for-agents/) optimizes *what* you say to the model, context engineering governs *what the model knows* when it generates a response.

**Key Components:**
- **System Instructions** - Rules and examples that guide model behavior
- **Dynamic Memory** - Conversation history and persistent knowledge
- **Retrieved Information** - Real-time data from documents, APIs, and databases  
- **Available Tools** - Functions the model can use (search, send_email, etc.)
- **User State** - Preferences, context, and session information

---

## Why Context Engineering Matters

### The Evolution from Prompt to Context

| **Prompt Engineering** | **Context Engineering** |
|----------------------|-------------------------|
| Single-turn instructions | Multi-turn, stateful workflows |
| Optimizes immediate output | Manages entire AI ecosystem |
| Limited by prompt information | Integrates external knowledge |
| Demo-level solutions | Production-grade systems |

### Critical Problems Context Engineering Solves

- **Hallucinations** → Reduced via grounding in real, external data
- **Statelessness** → Replaced by memory buffers and user modeling
- **Stale Knowledge** → Solved via retrieval pipelines and dynamic updates
- **Weak Personalization** → Addressed by user state tracking

As [Andrej Karpathy recently highlighted](https://x.com/karpathy/status/1937902205765607626), context engineering is becoming the foundation for reliable AI systems that can operate in complex, real-world environments.

---

## Our Story: Building on MCP's Vision for Context Engineering

### Why We Love MCP

When [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) appeared in late 2024, we were *geniuinely excited*. The vision was spot-on: a standardized way for AI tools to securely access external data and services. MCP represented a breakthrough in thinking about how AI agents should interact with the world.

**What MCP got right:**
- **Clear protocol design** for tool-AI communication
- **Security-first approach** to external integrations  
- **Architecture** that can grow with the ecosystem
- **Open standard** that brings the community together

### Our Developer Journey with MCP

As we dove deep into building production systems with MCP, we encountered some implementation challenges that many early adopters face:

- **Complex deployment setup** for enterprise environments
- **OAuth flow gaps** in production-grade security requirements
- **Stateless nature** that needed persistent context management
- **Integration overhead** for each service connection

These weren't flaws in MCP's design - they were opportunities to build the missing infrastructure layer that MCP made possible.

### Our Solution: Context Space

We realized MCP had given us the perfect foundation to build something bigger. So we spent months building **Context Space** - production-ready infrastructure that extends MCP's vision toward full context engineering.

**Today**: A Prouction-ready context provider integration layer with secure credential handling and persistent memory.  
**Tomorrow**: Full context engineering infrastructure built on MCP principles.

### Why We Open-Sourced

When Karpathy recently highlighted "Context Engineering" as the next frontier, we knew MCP was the perfect starting point for this evolution. We open-sourced Context Space because:

1. **MCP and integrations are a great starting point** for context engineering
2. **We're still building toward mature context engineering infrastructure**
3. **The community can help us build something amazing together**

We see ourselves as MCP evangelists building the infrastructure layer that makes MCP's vision accessible to everyone.

---

## 🎯 What We Have Today vs. Where We're Going

### Current Reality: Production-Ready Integration Infrastructure

**We've solved the MCP pain points:**
- **14+ Service Integrations**: GitHub, Slack, Airtable, HubSpot, Notion, Figma, Spotify, Stripe, and more.
- **Secure OAuth Made Simple**: No more editing MCP config files manually. Built-in OAuth flows handle it properly.
- **Enterprise-Grade Security**: HashiCorp Vault integration, automatic token rotation, encrypted credential storage.
- **Production Infrastructure**: Docker, Kubernetes, PostgreSQL, Redis, comprehensive monitoring.
- **RESTful API**: Clean HTTP endpoints that actually work reliably.

**How Context Space enhances MCP for production use:**

| **MCP Foundation** | **Context Space Enhancement** |
|-------------------|------------------------------|
| Clean protocol design | ✅ + Enterprise deployment infrastructure |
| Security-aware architecture | ✅ + OAuth UI flows & Vault integration |
| Extensible tool interface | ✅ + Persistent credential management |
| Open ecosystem vision | ✅ + 14+ ready-to-use integrations |

### 🔍 The Vision: Full Context Engineering Infrastructure

We believe MCP provides the *perfect foundation* for context engineering infrastructure. Here's how we're extending MCP's vision:

**Roadmap Timeline:**

| Phase | Timeline | Key Features | MCP Integration |
|-------|----------|--------------|------------------|
| **Phase 1** | Next 6 months | Native MCP Support, Context Memory, Smart Aggregation | Full MCP protocol compatibility |
| **Phase 2** | 6-12 months | Semantic Retrieval, Context Optimization, Real-time Updates | Enhanced MCP tool capabilities |
| **Phase 3** | 12+ months | Context Synthesis, Predictive Loading, AI Context Reasoning | Advanced MCP ecosystem features |

---

## Supported Services & Context Sources

### Production-Ready Integrations

| Service | Category | Auth | Context Capabilities | Status |
|---------|----------|------|---------------------|--------|
| **GitHub** | Development | OAuth | Code repos, issues, PRs, commit history | ✅ Ready |
| **Slack** | Communication | OAuth | Team conversations, channels, workflows | ✅ Ready |
| **Airtable** | Data Management | OAuth | Structured business data, CRM records | ✅ Ready |
| **HubSpot** | CRM | OAuth | Customer data, sales pipeline, interactions | ✅ Ready |
| **Notion** | Knowledge | OAuth | Documentation, project plans, wikis | ✅ Ready |
| **Spotify** | Personal | OAuth | Music preferences, listening patterns | ✅ Ready |
| **Stripe** | Financial | API Key | Payment data, customer behavior | ✅ Ready |
| **More...** | Various | Various | 5+ additional integrations | ✅ Ready |

**✅ 14+ integrations ready to use • More being added weekly**

**[View All Integrations →](https://context.space/integrations)**

---

## 📖 API Documentation

### Quick API Examples

#### 🔐 Authentication
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     https://api.context.space/v1/users/me
```

#### 🔗 Create OAuth Authorization URL
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     -X POST \
     https://api.context.space/v1/credentials/auth/oauth/github/auth-url
```

#### ⚡ Execute Operations
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     -X POST \
     https://api.context.space/v1/invocations/github/list_repositories
```

**Complete API Documentation**: http://api.context.space/v1/docs

---

## Contributing

**We need your help to build the future of context engineering!**

**Join our [Discord](https://discord.gg/Q74Ta5Xv) to connect with other contributors and get help with your PRs!**

[![Contributors](https://contrib.rocks/image?repo=context-space/context-space)](https://github.com/context-space/context-space/graphs/contributors)

### Quick Contributing Guide

1. **Sign the [CLA](CLA.md)**: Comment "I have read the CLA Document and I hereby sign the CLA" on your first PR
2. **Fork & Branch**: `git checkout -b feat/amazing-feature`
3. **Follow Standards**: Use `make lint` and include tests
4. **Submit PR**: With clear description

**Full Contributing Guide**: [CONTRIBUTING.md](CONTRIBUTING.md)

### Good First Issues

| Type | Difficulty | Examples |
|------|------------|----------|
| **Bug Fixes** | Easy | Fix API response formatting |
| **Documentation** | Easy | Improve API examples |
| **New Integrations** | Medium | Add Discord/Twitter support |
| **Context Features** | Hard | Implement semantic search |

**[See All Issues →](https://github.com/context-space/context-space/issues)**

---

## License

### Current License: AGPL v3 → Apache 2.0 Transition

**Why this approach?**
- **Now**: AGPL v3 protects during our startup phase
- **Future**: Apache 2.0 transition (as community grows) for maximum adoption
- **CLA**: Contributors sign our CLA enabling this transition

| Stakeholder | Today | Tomorrow |
|-------------|-------|----------|
| **👥 Users** | Free production access | Broader ecosystem compatibility |
| **👨‍💻 Contributors** | Protected from exploitation | Maximum community reach |

---

## Community & Support

### Join Our Growing Community

[![Twitter](https://img.shields.io/twitter/follow/contextspace?style=social)](https://twitter.com/contextspace)
[![Discord](https://img.shields.io/discord/1234567890?logo=discord&logoColor=white&label=Discord&color=5865F2)](https://discord.gg/Q74Ta5Xv)

### Resources

- **[API Documentation](https://api.context.space/v1/docs)** - Complete API reference
- **[Discord Community](https://discord.gg/Q74Ta5Xv)** - Real-time chat and collaboration
- **[GitHub Discussions](https://github.com/context-space/context-space/discussions)** - Community Q&A
- **[Issues](https://github.com/context-space/context-space/issues)** - Bug reports & feature requests

---

<div align="center">

**Built with ❤️ by the Context Space community**

> ***"Ultimate Context Engineering Infrastructure."***

*We're building the foundation for AI agents that truly understand context. Start with solid MCPs and integrations, evolve toward intelligent context engineering.*

### ⭐ **Star us on GitHub** • 🍴 **Fork and contribute** • 💬 **Join our [Discord](https://discord.gg/Q74Ta5Xv)**

[![GitHub stars](https://img.shields.io/github/stars/context-space/context-space?style=social)](https://github.com/context-space/context-space)
[![GitHub forks](https://img.shields.io/github/forks/context-space/context-space?style=social)](https://github.com/context-space/context-space)

</div>
