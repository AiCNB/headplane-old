# 裸机模式
[English](Bare-Metal.md) | [简体中文](Bare-Metal.zh-CN.md)

Bare-Metal 模式是部署 Headplane 最灵活的方式。它允许你在无需 Docker 或其他容器运行时的任何系统上运行 Headplane。虽然不推荐，但我理解每个人有不同需求。

> 它可与 **简单** 和 **集成** 部署模式一起使用。参见下方章节了解如何配置集成模式。

## 部署

需求：
- 已部署的 Headscale 0.26 或更高版本
- Node.js 22 LTS 或更高版本
- [PNPM](https://pnpm.io/installation) 10.x
- 一个完成的配置文件（config.yaml）

在安装 Headplane 之前，确保 `/var/lib/headplane` 已存在并且运行 Headplane 服务的用户对其有写权限。你可以使用以下命令创建该目录：

```sh
sudo mkdir -p /var/lib/headplane
# 如果不是 root，请将 headplane:headplane 替换为适当的用户和组。
sudo chown -R headplane:headplane /var/lib/headplane
```

克隆 Headplane 仓库、安装依赖并构建项目：
```sh
git clone https://github.com/tale/headplane
cd headplane
git checkout v0.6.0 # 或任何你想使用的标签
pnpm install
pnpm build
```

## 运行 Headplane
你可以通过 `pnpm start` 或 `node build/headplane/server.js` 来启动 Headplane。运行服务器时需要存在 `build` 目录。该目录的结构非常重要，不要随意修改。

### 集成模式
由于你在裸机环境中运行 Headplane，很可能也在裸机上运行 Headscale。请参阅[**集成模式**](/docs/Integrated-Mode.zh-CN.md)指南，以获取在原生 Linux (/proc) 中设置集成模式的说明。

### 更改管理路径
由于你自行构建 Headplane，可以将管理路径配置为任意值。在运行 `pnpm build` 时，可以传递 `__INTERNAL_PREFIX` 环境变量来更改管理路径。例如：

```sh
__INTERNAL_PREFIX=/admin2 pnpm build
```

请注意，管理路径在运行时不可配置，因此如果想更改它，你需要重新构建项目。另外，`/admin` 以外的任何路径都不受官方支持，未来版本中可能会出现问题。

> 有关如何将 `config.yaml` 文件设置为适当的值，请参阅[**配置**](/docs/Configuration.zh-CN.md)指南。

### Systemd 单元
下面是一个可用于管理 Headplane 服务的 systemd 单元示例：

```ini
[Unit]
Description=Headplane
# 如果在裸机环境下使用集成模式（/proc 集成），取消注释下一行
# PartOf=headscale.service

[Service]
Type=simple
User=headplane
Group=headplane
WorkingDirectory=/path/to/headplane
ExecStart=/usr/bin/node /path/to/headplane/build/server/index.js
Restart=always

[Install]
WantedBy=multi-user.target
```

你需要将 `/path/to/headplane` 替换为系统上 Headplane 仓库的实际路径。将该文件保存为 `/etc/systemd/system/` 中的 `headplane.service`，并运行 `systemctl enable headplane` 以启用服务。

其他字段也可能需要配置，因为该单元假设系统上存在名为 `headplane` 的用户和组。你可以根据系统配置修改这些值。
