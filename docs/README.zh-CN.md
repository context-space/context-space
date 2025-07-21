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


中文 | [English](../README.md)  

**Context Space 是一个面向 AI Agent 构建的工具优先（Tool-first）上下文工程基础设施**，将任务编排、记忆等 Agent 的核心能力统一封装为标准化可调用工具，凭借强大的工具发现和推荐能力，为 Agent 提供清晰、可控、可解释的上下文调用路径，成为构建复杂智能行为的坚实底座。

**🔗 即插即用** • **🔐 企业级安全** • **🚀 生产可用** • **🤖 上下文工程**

[Demo演示](#demo-演示) • [技术路线图](#技术路线图) • [API 文档](http://api.context.space/v1/docs)
</div>

![Homepage Screenshot](https://r2.context.space/resources/readme-homepage-screenshot.jpg)

---

## 使用Context Space开始上下文工程

上下文工程是构建可靠 AI 智能体的基础。不仅管理用户对模型说了什么，更进一步掌控塑造模型行为的广义上下文，例如工具、记忆和数据，从而超越了传统的提示词工程。上下文工程在提示词工程的基础上进一步拓展，不仅关注用户输入的内容，还对AI模型推理过程中依赖的所有信息进行系统性设计与管理，例如工具、记忆和数据等。

当 [MCP（模型上下文协议）](https://modelcontextprotocol.io/)在 2024 年底出现时，所描绘的愿景令人振奋：为 AI 工具提供安全访问外部数据和服务的标准化方式。MCP的广泛应用具有突破性意义，真正定义了AI智能体如何与世界进行交互。

我们基于MCP的愿景构建Context Space，将这项前瞻设想真正落地为生产可用的基础设施。

目前Context Space打造了一个具备持久凭据管理能力的安全集成层，我们正将这一基础设施拓展为面向下一代 AI应用的完整上下文工程平台。

### Demo 演示

#### 1️⃣ OAuth 流程

*简化 OAuth 设置 - 无需编辑配置文件*

![OAuth 演示](https://r2.context.space/resources/readme-demo-oauth-flow-github-v2.gif)


#### 2️⃣ 一键给GitHub点星

*GitHub 集成 - 用自然语言为Github仓库点 Star*

![GitHub 点星演示](https://r2.context.space/resources/readme-demo-github-star-repo.gif)

#### 3️⃣ 网络搜索
*实时网络搜索 - 立即获取最新信息*

![网络搜索演示](https://r2.context.space/resources/readme-demo-web-search.gif)

**在线试用**：[https://context.space/integrations](https://context.space/integrations)

---

## 技术路线图

我们规划了清晰完善的开发节奏，目前已构建稳定可用的基础设施，将在12个月内逐步成为上下文工程的重要平台。

### 第一阶段：生产可用的基础设施（已上线）

这一阶段聚焦于解决MCP实际部署过程中的关键难题，打造稳定、安全且可扩展的基础设施。

| **MCP部署难点** | **Context Space 解决方案** |
|------------|----------------------|
| 凭证手动管理、易泄漏 | ✅ **一键式 OAuth & Vault 安全机制**：支持 14+ 服务的安全 OAuth 流程，结合 HashiCorp Vault，实现企业级凭证管理。 |
| API复杂、多种类 | ✅ **统一 RESTful 接口**：与所有服务交互只需一个简洁一致的 API，开发体验显著提升。 |
| MCP Server分散、部署繁琐 | ✅ **统一上下文平面与工具聚合**：一次连接，全面访问。所有能力统一由一个 MCP 端点管理。 |

### 第二阶段：智能上下文层（开发中）

在第一阶段基础上，我们将拓展AI能力，持续构建更智能的上下文工程。

**路线图时间线：**

| 时间线 | 关键功能 | MCP 集成 |
|-------|---------|---------|
| 未来 6 个月 | 原生 MCP 支持、上下文记忆、智能聚合 | 完整 MCP 协议兼容性 |
| 6-12 个月 | 语义检索、上下文优化、实时更新 | 增强的 MCP 工具能力 |
| 12+ 个月 | 上下文合成、预测加载、AI 上下文推理 | 高级 MCP 生态系统功能 |

---

## 支持的服务和上下文源

### 生产环境集成

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

**✅ 14+ 集成服务即插即用 • 每周持续更新**

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

**我们邀请你一同构建上下文工程的未来！**

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

- | 利益相关者 | 今天 | 明天 |
|----------|------|------|
| **👥 用户** | 免费生产访问 | 更广泛的生态系统兼容性 |
| **👨‍💻 贡献者** | 免受剥削保护 | 最大社区影响力 |

---

## 社区与支持

Context Space 是一个由社区驱动的项目。我们相信只有开放的环境才能构建最好的基础设施，我们希望汇聚来自全球开发者的创意与专业能力，一起推动上下文工程的可能性边界。

### 加入我们的社区一起成长

[![Twitter](https://img.shields.io/twitter/follow/hi_contextspace?style=social)](https://twitter.com/hi_contextspace)
[![Discord](https://img.shields.io/badge/Discord-Join%20Us-5865F2?logo=discord&logoColor=white)](https://discord.gg/Q74Ta5Xv)

### 资源

- **[API 文档](https://api.context.space/v1/docs)** - 完整 API 参考
- **[Discord 社区](https://discord.gg/Q74Ta5Xv)** - 实时聊天和协作
- **[GitHub 讨论](https://github.com/context-space/context-space/discussions)** - 社区问答
- **[Issues](https://github.com/context-space/context-space/issues)** - Bug 报告和功能请求

---


<div align="center">

**⭐star & 分享我们的项目**

为项目点星🌟可以提升项目曝光度，帮助更多开发者发现 Context Space。如果你喜欢这个项目，欢迎把 Context Space 分享给你的朋友和同事，让更多人加入我们。

[![GitHub stars](https://img.shields.io/github/stars/context-space/context-space?style=social)](https://github.com/context-space/context-space)
[![GitHub forks](https://img.shields.io/github/forks/context-space/context-space?style=social)](https://github.com/context-space/context-space)

</div>
