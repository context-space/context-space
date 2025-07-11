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

---

**🌐 语言** | **中文** | [English](../README.md) | 

---

**Context Space 致力于成为 AI 智能体上下文工程的基础设施。**
当行业关注提示工程时，我们相信下一个前沿领域是上下文工程——在正确的时间、以正确的格式为 AI 系统提供正确的信息。
加入我们的 [Discord](https://discord.gg/Q74Ta5Xv ) 并帮助我们一起构建！有问题？[提交 issue](https://github.com/context-space/context-space/issues) 或加入我们的讨论。

**🔗 集成服务** • **🔐 凭据安全** • **🚀 生产就绪** • **🤖 上下文工程**

[快速开始](#-快速开始) • [在线演示](#-在线演示) • [API 文档](http://api.context.space/v1/docs)

</div>

![首页截图](https://r2.context.space/resources/readme-homepage-screenshot.jpg)

---

## 什么是上下文工程？

**上下文工程** 是对 AI 模型推理过程中所有周围信息的系统性设计和管理。虽然[提示工程](https://blog.langchain.com/context-engineering-for-agents/)优化的是您对模型*说什么*，上下文工程管理的是*模型在生成响应时知道什么*。

**核心组件：**
- **系统指令** - 指导模型行为的规则和示例
- **动态记忆** - 对话历史和持久知识
- **检索信息** - 来自文档、API 和数据库的实时数据
- **可用工具** - 模型可以使用的函数（搜索、发送邮件等）
- **用户状态** - 偏好设置、上下文和会话信息

---

## 为什么上下文工程很重要

### 从提示到上下文的演进

| **提示工程** | **上下文工程** |
|------------|-------------|
| 制作有效的指令 | 编排信息生态系统 |
| 优化单次交互 | 管理持久的多轮工作流 |
| 即时、定向的响应 | 动态知识集成与记忆 |

两者都是必要且互补的技能。上下文工程建立在提示工程基础之上并对其进行扩展。

### 上下文工程解决的关键问题

- **幻觉现象** → 通过基于真实外部数据的基础事实减少
- **无状态性** → 通过记忆缓冲区和用户建模替代
- **过时知识** → 通过检索流水线和动态更新解决
- **弱个性化** → 通过用户状态跟踪解决

正如 [Andrej Karpathy 最近强调的](https://x.com/karpathy/status/1937902205765607626)，上下文工程正在成为可在复杂真实世界环境中运行的可靠 AI 系统的基础。

---

## 我们的故事：基于 MCP 的上下文工程愿景构建

### 为什么我们喜爱 MCP

当 [MCP（模型上下文协议）](https://modelcontextprotocol.io/) 在 2024 年末出现时，我们*真正感到兴奋*。这个愿景非常精准：为 AI 工具安全访问外部数据和服务提供标准化方式。MCP 代表了在思考 AI 智能体应如何与世界交互方面的突破。

**MCP 的正确之处：**
- 工具-AI 通信的**清晰协议设计**
- 外部集成的**安全优先方法**
- 可与生态系统共同成长的**架构**
- 团结社区的**开放标准**

### 我们使用 MCP 的开发之旅

当我们深入使用 MCP 构建生产系统时，遇到了许多早期采用者都面临的实施挑战：

- 企业环境的**复杂部署设置**
- 生产级安全要求中的 **OAuth 流程缺口**
- 需要持久上下文管理的**无状态特性**
- 每个服务连接的**集成开销**

这些不是 MCP 设计中的缺陷——而是构建 MCP 使之成为可能的缺失基础设施层的机会。

### 我们的解决方案：Context Space

我们意识到 MCP 为我们提供了构建更大项目的完美基础。因此我们花费数月构建了 **Context Space** —— 将 MCP 愿景扩展为完整上下文工程的生产就绪基础设施。

**今天**：具有安全凭据处理和持久记忆的生产就绪上下文提供者集成层。
**明天**：基于 MCP 原则构建的完整上下文工程基础设施。

#### 🎬 在线演示

##### OAuth 流程实操
![OAuth 演示](https://r2.context.space/resources/readme-demo-oauth-flow-github-v2.gif)
*简单的 OAuth 设置 - 无需再编辑配置文件*

##### 为 GitHub 仓库点星
![GitHub 点星演示](https://r2.context.space/resources/readme-demo-github-star-repo.gif)
*GitHub 集成 - 用自然语言为仓库点 Star*

##### 网络搜索

![网络搜索演示](https://r2.context.space/resources/readme-demo-web-search.gif)
*实时网络搜索 - 即时获取最新信息*

**在线试用**：[https://context.space/integrations](https://context.space/integrations)

### 为什么我们选择开源

当 Karpathy 最近强调"上下文工程"是下一个前沿领域时，我们知道 MCP 是这种演进的完美起点。我们开源 Context Space 是因为：

1. **MCP 和集成是上下文工程的绝佳起点**
2. **我们仍在构建走向成熟的上下文工程基础设施**
3. **社区可以帮助我们一起构建令人惊叹的东西**

我们将自己视为 MCP 布道者，构建使 MCP 愿景为每个人所用的基础设施层。

---

## 🎯 现状 vs. 未来

### 当前现实：生产就绪的集成基础设施

**我们已经解决了 MCP 的痛点：**
- **14+ 服务集成**：GitHub、Slack、Airtable、HubSpot、Notion、Figma、Spotify、Stripe 等等。
- **简化的安全 OAuth**：不再需要手动编辑 MCP 配置文件。内置 OAuth 流程可正确处理。
- **企业级安全**：HashiCorp Vault 集成、自动令牌轮换、加密凭据存储。
- **生产基础设施**：Docker、Kubernetes、PostgreSQL、Redis、全面监控。
- **RESTful API**：真正可靠工作的清洁 HTTP 端点。

**Context Space 如何增强 MCP 的生产使用：**

| **MCP 基础** | **Context Space 增强** |
|------------|----------------------|
| 清洁的协议设计 | ✅ + 企业部署基础设施 |
| 安全感知架构 | ✅ + OAuth UI 流程 & Vault 集成 |
| 可扩展工具接口 | ✅ + 持久凭据管理 |
| 开放生态系统愿景 | ✅ + 14+ 即用集成 |

### 🔍 愿景：完整的上下文工程基础设施

我们相信 MCP 为上下文工程基础设施提供了*完美的基础*。以下是我们如何扩展 MCP 的愿景：

**路线图时间线：**

| 阶段 | 时间线 | 关键功能 | MCP 集成 |
|-----|-------|---------|---------|
| **1** | 未来 6 个月 | 原生 MCP 支持、上下文记忆、智能聚合 | 完整 MCP 协议兼容性 |
| **2** | 6-12 个月 | 语义检索、上下文优化、实时更新 | 增强的 MCP 工具能力 |
| **3** | 12+ 个月 | 上下文合成、预测加载、AI 上下文推理 | 高级 MCP 生态系统功能 |

---

## 支持的服务和上下文源

### 生产就绪集成

| 服务 | 类别 | 认证 | 上下文能力 | 状态 |
|------|------|------|----------|------|
| **GitHub** | 开发 | OAuth | 代码仓库、issues、PRs、提交历史 | 就绪 |
| **Slack** | 通信 | OAuth | 团队对话、频道、工作流 | 就绪 |
| **Airtable** | 数据管理 | OAuth | 结构化业务数据、CRM 记录 | 就绪 |
| **HubSpot** | CRM | OAuth | 客户数据、销售管道、交互 | 就绪 |
| **Notion** | 知识 | OAuth | 文档、项目计划、Wiki | 就绪 |
| **Spotify** | 个人 | OAuth | 音乐偏好、收听模式 | 就绪 |
| **Stripe** | 金融 | API Key | 支付数据、客户行为 | 就绪 |
| **更多...** | 各种 | 各种 | 5+ 其他集成 | 就绪 |

**✅ 14+ 集成即用 • 每周增加更多**

**[查看所有集成 →](https://context.space/integrations)**

---

## 📖 API 文档

### 快速 API 示例

#### 🔐 身份验证
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     https://api.context.space/v1/users/me
```

#### 🔗 创建 OAuth 授权 URL
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     -X POST \
     https://api.context.space/v1/credentials/auth/oauth/github/auth-url
```

#### ⚡ 执行操作
```bash
curl -H "Authorization: Bearer <jwt-token>" \
     -X POST \
     https://api.context.space/v1/invocations/github/list_repositories
```

**完整 API 文档**：http://api.context.space/v1/docs

---

## 贡献

**我们需要您的帮助来构建上下文工程的未来！**

**加入我们的 [Discord](https://discord.gg/Q74Ta5Xv) 与其他贡献者联系并获得 PR 帮助！**

[![Contributors](https://contrib.rocks/image?repo=context-space/context-space&anon=1)](https://github.com/context-space/context-space/graphs/contributors)

### 快速贡献指南

1. **签署 [CLA](CLA.md)**：在您的第一个 PR 上评论"I have read the CLA Document and I hereby sign the CLA"
2. **Fork & 分支**：`git checkout -b feat/amazing-feature`
3. **遵循标准**：使用 `make lint` 并包含测试
4. **提交 PR**：附上清晰的描述

**完整贡献指南**：[CONTRIBUTING.md](CONTRIBUTING.md)

### 适合新手的 Issues

| 类型 | 难度 | 示例 |
|------|------|------|
| **Bug 修复** | 简单 | 修复 API 响应格式 |
| **文档** | 简单 | 改进 API 示例 |
| **新集成** | 中等 | 添加 Discord/Twitter 支持 |
| **上下文功能** | 困难 | 实现语义搜索 |

**[查看所有 Issues →](https://github.com/context-space/context-space/issues)**

---

## 许可证

### 当前许可证：AGPL v3 → Apache 2.0 过渡

**为什么采用这种方法？**
- **现在**：AGPL v3 在我们的初创阶段提供保护
- **未来**：Apache 2.0 过渡（随着社区增长）以最大化采用
- **CLA**：贡献者签署我们的 CLA 以实现这种过渡

| 利益相关者 | 今天 | 明天 |
|----------|------|------|
| **👥 用户** | 免费生产访问 | 更广泛的生态系统兼容性 |
| **👨‍💻 贡献者** | 免受剥削保护 | 最大社区影响力 |

---

## 社区与支持

### 加入我们不断成长的社区

[![Twitter](https://img.shields.io/twitter/follow/hi_contextspace?style=social)](https://twitter.com/hi_contextspace)
[![Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/Q74Ta5Xv)

### 资源

- **[API 文档](https://api.context.space/v1/docs)** - 完整 API 参考
- **[Discord 社区](https://discord.gg/Q74Ta5Xv)** - 实时聊天和协作
- **[GitHub 讨论](https://github.com/context-space/context-space/discussions)** - 社区问答
- **[Issues](https://github.com/context-space/context-space/issues)** - Bug 报告和功能请求

---

<div align="center">

**由 Context Space 社区用 ❤️ 构建**

> ***"终极上下文工程基础设施。"***

*我们正在为真正理解上下文的 AI 智能体构建基础。从稳固的 MCP 和集成开始，演进为智能上下文工程。*

### ⭐ **在 GitHub 上为我们点星** • 🍴 **Fork 并贡献** • 💬 **加入我们的 [Discord](https://discord.gg/Q74Ta5Xv)**

[![GitHub stars](https://img.shields.io/github/stars/context-space/context-space?style=social)](https://github.com/context-space/context-space)
[![GitHub forks](https://img.shields.io/github/forks/context-space/context-space?style=social)](https://github.com/context-space/context-space)

</div>
