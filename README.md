<div align="center">
<a href="https://context.space">
<h1 align="center"> Context Space </h1>
</a>

[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![Docker](https://img.shields.io/badge/docker-supported-blue.svg)](https://docker.com)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![API Docs](https://img.shields.io/badge/API-documented-green.svg)](http://api.context.space/v1/docs)
[![Contributors](https://img.shields.io/badge/contributors-welcome-orange.svg)]()

English | [中文](docs/README.zh-CN.md)


**Context Space is a production-ready infrastructure layer for Context Engineering, making it easy to seamlessly adopt MCP servers and connect AI agents to external data and services.** While the industry focuses on prompt engineering, we believe the next frontier lies in Context Engineering that provides AI systems with the right information, at the right time, in the right format.

**🔗 Zero-Config Integrations** • **🔐 Enterprise-Grade Security** • **🚀 Production Ready** • **🤖 Context Engineering**

[Quick Start](#-quick-start) • [Live Demo](#-live-demo) • [API Documentation](http://api.context.space/v1/docs)

</div>

![Homepage Screenshot](https://r2.context.space/resources/readme-homepage-screenshot.jpg)


---

## What is Context Engineering?

As Andrej Karpathy recently [noted](https://x.com/karpathy/status/1937902205765607626), context engineering is becoming the foundation for reliable AI systems that can operate in complex, real-world environments.

**Context engineering** is the systematic design and management of all information surrounding an AI model during inference. Context engineering builds upon and extends prompt engineering. While [prompt engineering](https://blog.langchain.com/context-engineering-for-agents/) optimizes *what* you say to the model, context engineering governs what the model knows when it generates a response.

**Key Components:**
- System Instructions - Rules and examples that guide model behavior
- Dynamic Memory - Conversation history and persistent knowledge
- Retrieved Information - Real-time data from documents, APIs, and databases
- Available Tools - Functions the model can use (search, send_email, etc.)
- User State - Preferences, context, and session information

## Start Context Engineering with Context Space

When [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) appeared in late 2024, the vision was spot-on: a standardized way for AI tools to securely access external data and services. MCP represented a breakthrough in thinking about how AI agents should interact with the world.

Recognizing MCP as the perfect foundation, we built Context Space to extend its vision into production-ready infrastructure. 

Today, Context Space delivers a secure integration layer with persistent credential management. Guided by MCP’s principles, we are expanding this foundation into a complete context engineering platform for the next generation of AI.

### Live Demo

#### 1️⃣ OAuth Flow in Action
*Simple OAuth setup - no more config file editing*

![OAuth Demo](https://r2.context.space/resources/readme-demo-oauth-flow-github-v2.gif)

#### 2️⃣ Star a GitHub Repository
*GitHub integration - Star repositories with natural language*

![GitHub Star Demo](https://r2.context.space/resources/readme-demo-github-star-repo.gif)

#### 3️⃣ Web Search
*Real-time web search - get the latest information instantly*

![Web Search Demo](https://r2.context.space/resources/readme-demo-web-search.gif)

**Try Live**: [https://context.space/integrations](https://context.space/integrations)

---


## 🎯 Roadmap: From Foundation to Frontier

Our development is structured in clear phases, evolving from the robust production foundation available today to the intelligent context engine of tomorrow.

### 1️⃣ Phase 1: Production-Ready Foundation (Available Now)
The initial phase solves the most critical challenges of using context protocols in production environments, delivering a stable, secure, and scalable infrastructure.

| Challenge in Production | The Context Space Solution |
| :----- | :----- |
| Manual, Insecure Credential Handling | **One-Click OAuth & Vault Security:**<br>Connect to 14+ services with secure OAuth flows, backed by HashiCorp Vault for enterprise-grade credential management. |
| Inconsistent and Complex APIs | **A Single, Unified RESTful API:**<br>Interact with all services through one clean, consistent, and reliable API that you'll actually enjoy using. |
| Difficult Deployment & Scaling | **Battle-Tested Infrastructure:**<br>Deploy with confidence using our official support for Docker, Kubernetes, and PostgreSQL. |


### 2️⃣ Phase 2: The Intelligent Context Layer ( In Development)

Building on this foundation, our future work focuses on enabling more advanced AI capabilities. 

**Roadmap Timeline:**

|Timeline | Key Features | MCP Integration |
|----------|--------------|------------------|
| Next 6 months | Native MCP Support, Context Memory, Smart Aggregation | Full MCP protocol compatibility |
| 6-12 months | Semantic Retrieval, Context Optimization, Real-time Updates | Enhanced MCP tool capabilities |
| 12+ months | Context Synthesis, Predictive Loading, AI Context Reasoning | Advanced MCP ecosystem features |

---

## Supported Services & Context Sources

### Production-Ready Integrations

| Service | Category | Auth | Context Capabilities | Status |
|---------|----------|------|---------------------|--------|
| **GitHub** | Development | OAuth | Code repos, issues, PRs, commit history | Ready |
| **Slack** | Communication | OAuth | Team conversations, channels, workflows | Ready |
| **Airtable** | Data Management | OAuth | Structured business data, CRM records | Ready |
| **HubSpot** | CRM | OAuth | Customer data, sales pipeline, interactions | Ready |
| **Notion** | Knowledge | OAuth | Documentation, project plans, wikis | Ready |
| **Spotify** | Personal | OAuth | Music preferences, listening patterns | Ready |
| **Stripe** | Financial | API Key | Payment data, customer behavior | Ready |
| **More...** | Various | Various | 5+ additional integrations | Ready |

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

**You are invited to help shape the future of context engineering.** 

[![Contributors](https://contrib.rocks/image?repo=context-space/context-space&anon=1)](https://github.com/context-space/context-space/graphs/contributors)

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

Context Space is a community-driven project. We believe the best infrastructure is built in the open, with developers from all over the world contributing their ideas and expertise. Every contribution, big or small, helps us push the boundaries of what's possible.

### Join Our Growing Community

[![Twitter](https://img.shields.io/twitter/follow/hi_contextspace?style=social)](https://twitter.com/hi_contextspace)
[![Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/Q74Ta5Xv)


### Resources

- **[API Documentation](https://api.context.space/v1/docs)** - Complete API reference
- **[Discord Community](https://discord.gg/Q74Ta5Xv)** - Real-time chat and collaboration
- **[GitHub Discussions](https://github.com/context-space/context-space/discussions)** - Community Q&A
- **[Issues](https://github.com/context-space/context-space/issues)** - Bug reports & feature requests

---

<div align="center">

**🌟 Star & Share the Project**

Starring the repository increases our visibility and helps other developers discover the project. If you like Context Space, don't hesitate to share it on Twitter, Reddit, or with your colleagues.

[![GitHub stars](https://img.shields.io/github/stars/context-space/context-space?style=social)](https://github.com/context-space/context-space)
[![GitHub forks](https://img.shields.io/github/forks/context-space/context-space?style=social)](https://github.com/context-space/context-space)

</div>
