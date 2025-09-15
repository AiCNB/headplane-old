# Nix
[English](Nix.md) | [简体中文](Nix.zh-CN.md)

提供的 [flake.nix](../flake.nix)：
```
$ nix flake show . --all-systems
git+file:///home/igor/personal/headplane?ref=refs/heads/nix&rev=2d78a95a0648a3778e114fb246ea436e96475d62
├───devShell
│   ├───aarch64-darwin: development environment 'headplane'
│   ├───x86_64-darwin: development environment 'headplane'
│   └───x86_64-linux: development environment 'headplane'
├───formatter
│   ├───aarch64-darwin: package 'alejandra-3.1.0'
│   ├───x86_64-darwin: package 'alejandra-3.1.0'
│   └───x86_64-linux: package 'alejandra-3.1.0'
├───nixosModules
│   └───headplane: NixOS module
├───overlays
│   └───default: Nixpkgs overlay
└───packages
    ├───aarch64-darwin
    │   ├───headplane: package 'headplane-0.5.3-SNAPSHOT'
    │   └───headplane-agent: package 'hp_agent-0.5.3-SNAPSHOT'
    ├───x86_64-darwin
    │   ├───headplane: package 'headplane-0.5.3-SNAPSHOT'
    │   └───headplane-agent: package 'hp_agent-0.5.3-SNAPSHOT'
    └───x86_64-linux
        ├───headplane: package 'headplane-0.5.3-SNAPSHOT'
        └───headplane-agent: package 'hp_agent-0.5.3-SNAPSHOT'
```

## NixOS 模块选项
定义为 `services.headplane.*`，详情见 `./nix/` 目录。

## 使用方法

1. 添加 `github:tale/headplane` flake 输入。
2. 导入默认 overlay 以添加 `pkgs.headplane` 和 `pkgs.headplane-agent`。
3. 导入 NixOS 模块以获得 `services.headplane.*`。

```nix
# Your flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    headplane = {
      url = "github:igor-ramazanov/headplane/nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    nixpkgs,
    headplane,
    ...
  }: {
    nixosConfigurations.MY_MACHINE = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        # 提供 `services.headplane.*` NixOS 选项。
        headplane.nixosModules.headplane
        {
          # 提供 `pkgs.headplane` 和 `pkgs.headplane-agent`。
          nixpkgs.overlays = [ headplane.overlays.default ];
        }
        {
          {config, pkgs, ...}:
          let
            format = pkgs.formats.yaml {};

            # 生成一个 Headplane 接受的有效 Headscale 配置（当 `config_strict == true` 时）。
            settings = lib.recursiveUpdate config.services.headscale.settings {
              acme_email = "/dev/null";
              tls_cert_path = "/dev/null";
              tls_key_path = "/dev/null";
              policy.path = "/dev/null";
              oidc.client_secret_path = "/dev/null";
            };

            headscaleConfig = format.generate "headscale.yml" settings;
          in {
            services.headplane = {
              enable = true;
              agent = {
                # 仅作为示例。
                # 撰写本文档时 Headplane Agent 尚未准备好。
                enable = true;
                settings = {
                  HEADPLANE_AGENT_DEBUG = true;
                  HEADPLANE_AGENT_HOSTNAME = "localhost";
                  HEADPLANE_AGENT_TS_SERVER = "https://example.com";
                  HEADPLANE_AGENT_TS_AUTHKEY = "xxxxxxxxxxxxxx";
                  HEADPLANE_AGENT_HP_SERVER = "https://example.com/admin/dns";
                  HEADPLANE_AGENT_HP_AUTHKEY = "xxxxxxxxxxxxxx";
                };
              };
              settings = {
                server = {
                  host = "127.0.0.1";
                  port = 3000;
                  cookie_secret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
                  cookie_secure = true;
                };
                headscale = {
                  url = "https://example.com";
                  config_path = "${headscaleConfig}";
                  config_strict = true;
                };
                integration.proc.enabled = true;
                oidc = {
                  issuer = "https://oidc.example.com";
                  client_id = "headplane";
                  client_secret = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
                  disable_api_key_login = true;
                  # 与 Authelia 集成时可能需要。
                  token_endpoint_auth_method = "client_secret_basic";
                  headscale_api_key = "xxxxxxxxxx.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
                  redirect_uri = "https://oidc.example.com/admin/oidc/callback";
                };
              };
            };
          }
        }
      ];
    };
  };
}
```
