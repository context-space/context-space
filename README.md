![Context Space](https://r2.context.space/resources/20250724-235344_1753372441182.jpg)


<div align="center">

### Context Space: The First Context Engineering Infrastructure to 10× Your Productivity

     
[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![Docker](https://img.shields.io/badge/docker-supported-blue.svg)](https://docker.com)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![API Docs](https://img.shields.io/badge/API-documented-green.svg)](http://api.context.space/v1/docs)
[![Contributors](https://img.shields.io/badge/contributors-welcome-orange.svg)]()

English | [中文](docs/README.zh-CN.md)


</div>

**Context Space** offers unified MCP tools, secure & verified integrations, and a 5-minute setup — perfect for AI agents, automation workflows, and developer tools. As the **first context engineering infrastructure**, it turns theory into practice by delivering better context and enabling agents to interact effectively with the real world.

## Our Vision

Today's Al agents excel at reasoning but terrible at acting in the real world. They’re cut off from live data and tools, trapped behind scattered APIs, inconsistent sources, and complex authentication..

Context Space changes this. It packages core agent capabilities like task orchestration and memory into standardized, callable tools. With built-in tool discovery and recommendation, it gives agents a clear, controllable, and interpretable path to invoke real-world context.

Context Space makes AI agents truly usable. By combining enterprise-grade security with zero-config simplicity, we are building tool-first context engineering infrastructure that enable agents to seamlessly and securely interact with any service or data source.

## Start Context Engineering with Context Space

Context engineering is the foundation for building reliable AI agents. It goes beyond prompt engineering by managing not only what users say to the model, but also the broader context that shapes its behavior, such as tools, memory, and data.

MCP defines a standard path for agents to securely access real-world services. Context Space brings that vision to life by turning MCP into production-ready infrastructure.

Today, Context Space delivers a secure integration layer with persistent credential management. Guided by MCP principles, it’s evolving into a complete context engineering platform for the next generation of AI.

## 🚀 One-Click AI Integration

Transform your AI assistant into a powerful agent **in seconds**.

**Cursor IDE** - One-click install via `cursor://` deep links. Click "Add to Cursor" and instantly give Claude access to GitHub, Slack, Notion, and 38+ services without editing any JSON files.

**Claude Code** - Simple CLI integration:
```bash
claude mcp add "context-space" https://api.context.space/api/mcp --header "Authorization: Bearer YOUR_API_KEY"
```


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


## Roadmap: From Foundation to Frontier

Our development is structured in clear phases, evolving from the robust production foundation available today to the intelligent context engine of tomorrow.

### 1️⃣ Phase 1: Production-Ready Foundation (Available Now)
The initial phase solves the most critical challenges of using context protocols in production environments, delivering a stable, secure, and scalable infrastructure.

| Challenge in Production | The Context Space Solution |
| :----- | :----- |
| Manual, Insecure Credential Handling | **One-Click OAuth & Vault Security:**<br>Connect to 14+ services with secure OAuth flows, backed by HashiCorp Vault for enterprise-grade credential management. |
| Inconsistent and Complex APIs | **A Single, Unified RESTful API:**<br>Interact with all services through one clean, consistent, and reliable API that you'll actually enjoy using. |
| Complex Deployment & Scattered MCP Servers | **Unified Context Plane with Tool Aggregation:**<br>Connect once, and access everything. Manage all capabilities from a single mcp server endpoint. |


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
[![Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/BsNjUyxQYF)


### Resources

- **[API Documentation](https://api.context.space/v1/docs)** - Complete API reference
- **[Discord Community](https://discord.gg/BsNjUyxQYF)** - Real-time chat and collaboration
- **[GitHub Discussions](https://github.com/context-space/context-space/discussions)** - Community Q&A
- **[Issues](https://github.com/context-space/context-space/issues)** - Bug reports & feature requests

---

<div align="center">

**🌟 Star & Share the Project**

Starring the repository increases our visibility and helps other developers discover the project. If you like Context Space, don't hesitate to share it on Twitter, Reddit, or with your colleagues.

[![GitHub stars](https://img.shields.io/github/stars/context-space/context-space?style=social)](https://github.com/context-space/context-space)
[![GitHub forks](https://img.shields.io/github/forks/context-space/context-space?style=social)](https://github.com/context-space/context-space)

</div>
