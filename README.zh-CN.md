# Headplane
[English](README.md) | [简体中文](README.zh-CN.md)
> 一个功能完整的 [Headscale](https://headscale.net) Web 界面

<picture>
    <source
        media="(prefers-color-scheme: dark)"
        srcset="./assets/preview-dark.png"
    >
    <source
        media="(prefers-color-scheme: light)"
        srcset="./assets/preview-light.png"
    >
    <img
        alt="预览"
        src="./assets/preview-dark.png"
    >
</picture>

Headscale 是 Tailscale 的事实上的自托管版本，这是一款流行的基于 Wireguard 的 VPN 服务。默认情况下它不提供 Web 界面，这就是 Headplane 登场的地方。Headplane 是一个功能完整的 Headscale Web UI，使你能够轻松管理节点、网络和 ACL。

Headplane 旨在复刻官方 Tailscale 产品和控制台所提供的功能，是目前最为完整的 Headscale UI 之一。Headplane 提供的功能包括：

- 机器管理，包括到期时间、网络路由、名称和所有者管理
- 访问控制列表（ACL）和标签配置以执行 ACL
- 支持 OpenID Connect（OIDC）作为登录提供方
- 能够编辑 DNS 设置并自动配置 Headscale
- Headscale 设置的可配置性

## 部署
Headplane 作为一个基于服务器的 Web 应用运行，这意味着你需要一台服务器来运行它。它可以通过 Docker 镜像（推荐）或手动安装来使用。部署 Headplane 有两种方式：

- ### [集成模式（推荐）](/docs/Integrated-Mode.zh-CN.md)
  集成模式解锁了 Headplane 的全部功能，是最完整的部署方式。它直接与 Headscale 通信。

- ### [简单模式](/docs/Simple-Mode.zh-CN.md)
  简单模式不包含对 DNS 和 Headplane 设置的自动管理，在进行更改时需要手动编辑和重新加载。

### 版本控制
Headplane 使用 [语义化版本](https://semver.org/)（自 v0.6.0 起）进行发布。预发布版本在 `next` 标签下提供，并在新的发布 PR 打开且处于测试时更新。

### 贡献
Headplane 是一个开源项目，欢迎贡献！如果你有任何建议、错误报告或功能请求，请提交 issue。更多信息请参阅[贡献指南](./docs/CONTRIBUTING.zh-CN.md)。

---

<picture>
    <source
        media="(prefers-color-scheme: dark)"
        srcset="./assets/acls-dark.png"
    >
    <source
        media="(prefers-color-scheme: light)"
        srcset="./assets/acls-light.png"
    >
    <img
        alt="ACLs"
        src="./assets/acls-dark.png"
    >
</picture>

<picture>
    <source
        media="(prefers-color-scheme: dark)"
        srcset="./assets/machine-dark.png"
    >
    <source
        media="(prefers-color-scheme: light)"
        srcset="./assets/machine-light.png"
    >
    <img
        alt="设备管理"
        src="./assets/machine-dark.png"
    >
</picture>

> Copyright (c) 2025 Aarnav Tale
