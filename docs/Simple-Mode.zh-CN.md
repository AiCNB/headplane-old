# 简单模式
[English](Simple-Mode.md) | [简体中文](Simple-Mode.zh-CN.md)

简单模式可以让你快速部署 Headplane，适用于测试或简单环境。它不包含对 DNS 和 Headplane 设置的自动管理，在进行更改时需要手动编辑并重新加载。如果你需要功能更完整的部署方式，请查看[**集成模式**](/docs/Integrated-Mode.zh-CN.md)。

## 部署
> 如果你不打算使用 Docker 部署，请查看 [**裸机模式**](/docs/Bare-Metal.zh-CN.md) 部署指南。

需求：
- Docker 与 Docker Compose
- 已部署的 Headscale 0.26 或更高版本
- 一个完成的配置文件（config.yaml）

以下是一个示例的 Docker Compose 部署：
```yaml
services:
  headplane:
    # 建议固定到特定版本
    image: ghcr.io/tale/headplane:0.6.0
    container_name: headplane
    restart: unless-stopped
    ports:
      - '3000:3000'
    volumes:
      - './config.yaml:/etc/headplane/config.yaml'
      - './headplane-data:/var/lib/headplane'
```

这样部署后，Headplane 的界面将在你部署服务器的 `/admin` 路径下可用。除非你自行构建容器或在裸机模式下运行 Headplane，否则 `/admin` 路径目前无法配置。

> 有关如何设置 `config.yaml` 的更多信息，请参阅[**配置**](/docs/Configuration.zh-CN.md)指南。
