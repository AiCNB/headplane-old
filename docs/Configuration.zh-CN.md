# 配置
[English](Configuration.md) | [简体中文](Configuration.zh-CN.md)
> 早期版本的 Headplane 使用环境变量而没有配置文件。自 0.5 起，你需要手动迁移配置到新的格式。

Headplane 使用配置文件来管理其设置（[**config.example.yaml**](../config.example.yaml)）。默认情况下，Headplane 会在 `/etc/headplane/config.yaml` 查找该文件。可以使用 **`HEADPLANE_CONFIG_PATH`** 环境变量将其指向其他位置。

Headplane 默认还会在 `/var/lib/headplane` 目录中存储数据。可以在配置文件中按部分配置此路径，但必须确保该目录是持久的并且 Headplane 具有写权限。

## 环境变量
也可以使用环境变量覆盖配置文件。这些更改在加载配置文件之后合并，因此具有更高优先级。环境变量遵循以下模式：**`HEADPLANE_<SECTION>__<KEY_NAME>`**。例如，要覆盖 `oidc.client_secret`，可以设置 `HEADPLANE_OIDC__CLIENT_SECRET` 为所需的值。

以下是更多示例：

- `HEADPLANE_HEADSCALE__URL`：`headscale.url`
- `HEADPLANE_SERVER__PORT`：`server.port`

**此功能默认未启用！**
要启用它，请设置环境变量 **`HEADPLANE_LOAD_ENV_OVERRIDES=true`**。设置后 Headplane 还会将相对的 `.env` 文件加载到环境中。
> 另外请注意，这仅用于配置覆盖，而不是通用环境变量，因此不能指定如 `HEADPLANE_DEBUG_LOG=true` 或 `HEADPLANE_CONFIG_PATH=/etc/headplane/config.yaml` 等变量。

## 调试
要启用调试日志，请设置 **`HEADPLANE_DEBUG_LOG=true`** 环境变量。这会启用所有调试日志，可能会很快占满日志空间。在生产环境中不推荐。

## 反向代理
在部署 Web 应用时，反向代理非常常见。Headscale 和 Headplane 在这方面很相似。你可以使用任何熟悉的反向代理配置。以下是使用 Traefik 的示例：

> 这里重要的一点是 CORS 中间件。前端需要它来与后端通信。如果你使用其他反向代理，请确保添加必要的头部，以允许前端与后端通信。

```yaml
http:
  routers:
    headscale:
      rule: 'Host(`headscale.tale.me`)'
      service: 'headscale'
      middlewares:
        - 'cors'

    rewrite:
      rule: 'Host(`headscale.tale.me`) && Path(`/`)'
      service: 'headscale'
      middlewares:
        - 'rewrite'

    headplane:
      rule: 'Host(`headscale.tale.me`) && PathPrefix(`/admin`)'
      service: 'headplane'

  services:
    headscale:
      loadBalancer:
        servers:
          - url: 'http://headscale:8080'

    headplane:
      loadBalancer:
        servers:
          - url: 'http://headplane:3000'

  middlewares:
    rewrite:
      addPrefix:
        prefix: '/admin'
    cors:
      headers:
        accessControlAllowHeaders: '*'
        accessControlAllowMethods:
          - 'GET'
          - 'POST'
          - 'PUT'
        accessControlAllowOriginList:
          - 'https://headscale.tale.me'
        accessControlMaxAge: 100
        addVaryHeader: true
```
