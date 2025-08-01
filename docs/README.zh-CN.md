![Context Space](https://r2.context.space/resources/20250724-235344_1753372441182.jpg)

<div align="center">

### Context Space：首个提升10倍生产力的上下文工程基础设施

[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org/dl/)
[![Docker](https://img.shields.io/badge/docker-supported-blue.svg)](https://docker.com)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![API Docs](https://img.shields.io/badge/API-documented-green.svg)](http://api.context.space/v1/docs)
[![Contributors](https://img.shields.io/badge/contributors-welcome-orange.svg)]()

中文 | [English](../README.md)

</div>

**Context Space** 提供统一的MCP工具、安全可信的集成服务，以及5分钟快速设置，可广泛应用于AI智能体、自动化工作流和跨平台工具的开发。作为**首个上下文工程基础设施**，Context Space将理论转化为实践，提供更好的上下文和外部数据管理功能，使智能体能够有效地与真实世界交互。

## 愿景

AI智能体擅长推理，但在真实世界中的执行能力很差，分散的API、混乱的数据源和复杂的认证机制阻碍了它们与外部数据和工具的连接。

Context Space正在改变这一切。通过将任务编排和记忆等核心能力封装为标准化的可调用工具，并内置工具发现和推荐功能，为智能体提供清晰、可控、可解释的上下文调用路径。

Context Space让AI智能体真正可用。结合企业级安全与零配置的易用性，我们正在构建工具优先的上下文工程基础设施，使AI智能体能够无缝且安全地与任何外部服务或数据源交互。

## 使用Context Space开始上下文工程

上下文工程是构建可靠AI智能体的基础。上下文工程在提示词工程的基础上进一步拓展，不仅关注用户输入的内容，还对AI模型推理过程中依赖的所有信息进行系统性设计与管理，例如工具、记忆和数据等。

MCP为 AI 工具提供安全访问外部数据和服务的标准化方式，Context Space将MCP的愿景转化为生产可用的基础设施。

目前Context Space打造了一个具备持久凭据管理能力的安全集成层，我们正将这一基础设施拓展为面向下一代 AI应用的完整上下文工程平台。


## 一键AI集成

**几秒钟**即可将您的AI助手转变为强大的智能体。

**Cursor IDE** - 通过`cursor://`一键安装。点击"添加到Cursor"，即可让Cursor瞬间访问GitHub、Slack、Notion和38+服务，无需编辑任何JSON文件。

**Claude Code** - 简单的CLI集成：
```bash
claude mcp add "context-space" https://api.context.space/api/mcp --header "Authorization: Bearer YOUR_API_KEY"
```


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
