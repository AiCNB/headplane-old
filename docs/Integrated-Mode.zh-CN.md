# 集成模式
[English](Integrated-Mode.md) | [简体中文](Integrated-Mode.zh-CN.md)

<picture>
    <source
        media="(prefers-color-scheme: dark)"
        srcset="../assets/dns-dark.png"
    >
    <source
        media="(prefers-color-scheme: light)"
        srcset="../assets/dns-light.png"
    >
    <img
        alt="集成预览"
        src="../assets/dns-dark.png"
    >
</picture>

集成模式是一种部署方式，允许你在自动管理 DNS 和 Headplane 设置的情况下部署 Headplane。对于大多数用户来说，这是推荐的部署方式，因为它提供了更完整的功能体验。

## 部署
> 如果你不打算使用 Docker 部署，请查看 [**裸机模式**](/docs/Bare-Metal.zh-CN.md) 部署指南。底部的“集成模式”部分包含一些注意事项。

需求：
- Docker 与 Docker Compose
- Headscale 0.26 或更高版本
- 一个完成的配置文件（config.yaml）

这里是一个示例的 Docker Compose 部署：
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
      # 该路径应与 config.yaml 中的 headscale.config_path 相匹配
      - './headscale-config/config.yaml:/etc/headscale/config.yaml'

      # 如果在 Headscale 中使用 dns.extra_records（推荐），
      # 该路径应与 config.yaml 中的 headscale.dns_records_path 相匹配
      - './headscale-config/dns_records.json:/etc/headscale/dns_records.json'

      # Headplane 会将数据存储在此目录
      - './headplane-data:/var/lib/headplane'

      # 如果使用 Docker 集成，挂载 Docker 套接字
      - '/var/run/docker.sock:/var/run/docker.sock:ro'
  headscale:
    image: headscale/headscale:0.26.0
    container_name: headscale
    restart: unless-stopped
    command: serve
    labels:
      # 使 Headplane 能找到并向其发送信号
      me.tale.headplane.target: headscale
    ports:
      - '8080:8080'
    volumes:
      - './headscale-data:/var/lib/headscale'
      - './headscale-config:/etc/headscale'
```

这样部署后，Headplane 的界面将在你部署服务器的 `/admin` 路径下可用。除非你自行构建容器或在裸机模式下运行 Headplane，否则 `/admin` 路径目前无法配置。

> 有关如何设置 `config.yaml` 的更多信息，请参阅[**配置**](/docs/Configuration.zh-CN.md)指南。

## Docker 集成
Docker 集成最容易设置，因为只需要将 Docker 套接字挂载到容器并进行一些基础配置。Headplane 使用 Docker 标签发现 Headscale 容器。只要 Headplane 能访问 Docker 套接字并通过标签或名称识别 Headscale 容器，它就会自动将配置和 DNS 更改传播到 Headscale，无需额外设置。或者，你也可以直接指定容器名称，而不是使用标签动态确定。

## 原生 Linux (/proc) 集成
当在非 Docker 环境中运行 Headscale 和 Headplane 时，使用 `proc` 集成。Headplane 会尝试通过 `/proc` 文件系统定位 Headscale 进程的 PID 并直接与其通信。要使其工作，Headplane 进程必须具备以下权限：

- 读取 `/proc` 文件系统
- 向 Headscale 进程发送信号（`SIGTERM`）

最佳方式是以与 Headscale 相同的用户运行 Headplane（或都以 `root` 运行）。由于当前集成方式的限制，如果 Headscale 的 PID 发生变化，Headplane 不会重新检查。这意味着如果手动重启 Headscale，你也需要重启 Headplane。

## Kubernetes 集成
Kubernetes 集成最复杂，因为需要创建具有适当权限的服务账号。该服务账号必须具备以下权限，示例：
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: headplane-agent
  namespace: default # 根据需要调整命名空间
rules:
- apiGroups: ['']
  resources: ['pods']
  verbs: ['get', 'list']
- apiGroups: ['apps']
  resources: ['deployments']
  verbs: ['get', 'list']
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: headplane-agent
  namespace: default # 根据需要调整命名空间
roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: Role
    name: headplane-agent
subjects:
- kind: ServiceAccount
  name: default # 如果使用其他服务账号请修改
  namespace: default # 根据需要调整命名空间
```

要在 Kubernetes 中成功部署 Headplane，你需要在同一个 Pod 中运行 Headplane 和 Headscale 容器。这是因为 Headplane 需要访问 Headscale 的 PID 以与其通信。以下是示例，注意 **`shareProcessNamespace: true`** 字段：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: headplane
  namespace: default # 根据需要调整命名空间
  labels:
    app: headplane
spec:
  replicas: 1
  selector:
    matchLabels:
      app: headplane
  template:
    metadata:
      labels:
        app: headplane
    spec:
      shareProcessNamespace: true
      serviceAccountName: default
      containers:
      - name: headplane
        image: ghcr.io/tale/headplane:0.6.0
        env:
        # 如果 Headscale 的 Pod 名称不是固定的，设置以下变量
        # 我们将使用 downward API 获取 Pod 名称
        - name: HEADPLANE_LOAD_ENV_OVERRIDES
          value: 'true'
        - name: 'HEADPLANE_INTEGRATION__KUBERNETES__POD_NAME'
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        volumeMounts:
        - name: headscale-config
          mountPath: /etc/headscale
        - name: headplane-data
          mountPath: /var/lib/headplane

      - name: headscale
        image: headscale/headscale:0.26.0
        command: ['serve']
        volumeMounts:
        - name: headscale-data
          mountPath: /var/lib/headscale
        - name: headscale-config
          mountPath: /etc/headscale

      volumes:
        - name: headscale-data
          persistentVolumeClaim:
            claimName: headscale-data
        - name: headplane-data
          persistentVolumeClaim:
            claimName: headplane-data
        - name: headscale-config
          persistentVolumeClaim:
            claimName: headscale-config
```
