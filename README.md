# DDOSKing - DDoS攻击仿真环境

![image-20250302140010433](C:\Users\17128\AppData\Roaming\Typora\typora-user-images\image-20250302140010433.png)

<img src="C:\Users\17128\AppData\Roaming\Typora\typora-user-images\image-20250302163118293.png" alt="image-20250302163118293" style="zoom: 33%;" />

## 项目简介

DDOSKing 是一个基于 Docker 的 DDoS 攻击仿真环境，专为研究和测试各种 DDoS 攻击技术而设计。本项目使用 Seed-emulator 构建仿真网络环境并提供可视化界面，攻击脚本采用 Go 语言开发并部署在 Docker 容器中。

该仿真环境模拟了当前主流的 DDoS 攻击方法和僵尸网络攻击，并加入了针对 DeepSeek 等 AI 服务的攻击仿真。整个仿真环境只需在单台主机上即可模拟完整的互联网和 DDoS 攻击场景。

> **注意**：部署 DeepSeek 1.5B 至少需要 8GB 内存

## 系统需求

为了完整运行仿真环境，建议系统内存大于 24GB。

## 系统组成

本仿真环境模拟了完整的 DDoS 攻击基础设施，包含以下独立构建的组件：

| 组件                | 数量 | 功能描述                          |
| ------------------- | ---- | --------------------------------- |
| C2服务器            | 1台  | 负责向僵尸主机发送攻击指令        |
| 反射放大服务器      | 2个  | 用于 Layer 4 的 DDoS 攻击反射放大 |
| 僵尸主机            | 5台  | 执行各种攻击                      |
| Unbound DNS解析器   | 1台  | 特殊配置用于进行脉冲攻击          |
| DNS权威服务器       | 1台  | 用于脉冲攻击的请求积累和放大      |
| DeepSeek 1.5B服务器 | 1台  | 用于模拟 AI 服务遭受的 HTTP 攻击  |

> **注意**：本仿真环境专注于 DDoS 仿真，暂不模拟僵尸网络的传播、通信机制（如弱口令爆破感染主机、DGA 搜索 C2 等），可能在后续版本中添加。

其余节点由 SeedEmu 自动生成。详细信息可参考 [SeedEmu 官方文档](https://github.com/seed-labs/seed-emulator)。

## 攻击类型

DDOSKing 覆盖了多种 DDoS 攻击类型，主要分为以下几类：

### 1. 链路泛洪攻击 (Layer 4)

利用大量流量占满网络带宽的攻击，包括：

- **直接攻击**：UDP 泛洪
- **反射放大攻击**：DNS、NTP、CLDAP、SSDP 等

### 2. 资源消耗攻击 (Layer 7)

消耗服务器计算资源的攻击，包括：

- **HTTP 泛洪**：为保持环境轻量化，Web 服务仅提供一个页面
- **针对 DeepSeek 构造复杂 prompt**：进行 AI 服务资源消耗
- **SYN 泛洪**：消耗目标半连接队列

### 3. 脉冲攻击

脉冲攻击旨在短时间内发送高带宽的报文（数据量往往较低，因此持续时间较短），使目标队列被占满，造成超时，触发 TCP 的拥塞控制，从而使目标 TCP 服务降级。

- **DNSBomb：**IEEE S&P 24的工作，参考链接：[DNSBomb: A New Practical-and-Powerful Pulsing DoS Attack Exploiting DNS Queries-and-Responses | IEEE Conference Publication | IEEE Xplore](https://ieeexplore.ieee.org/abstract/document/10646654)

- **DNSBoomerang：**DNSBoomerang 基于 DNSBomb 并进行改进，使积累数据量能够大幅增加（攻击积累报文数量随 DNS 反射器增加而增加）。公网实验，攻击者以530kbps速率积累请求，源IP为500台不同的反射服务器（收到响应报文会返回响应报文），积累17s（累计13700条请求），反射带宽达到108Mbps，持续时间约1s，放大204倍。

  <div style=text-align:left>
      <img src="C:\Users\17128\AppData\Roaming\Typora\typora-user-images\image-20250302182435467.png" alt="image-20250302182435467" style="zoom:33%;" />
      <img src="C:\Users\17128\AppData\Roaming\Typora\typora-user-images\image-20250302182256270.png" alt="image-20250302182256270" style="zoom:45%;" />
  </div>

## 环境搭建

建议在Linux环境搭建，windows用户可以使用wsl。

### 安装步骤

#### 1. Docker 安装与配置

- 安装 Docker：[参考官方文档](https://docs.docker.com/engine/install/)
- **注意**：中国大陆用户需要配置 Docker 镜像源，因为默认的 DockerHub 可能无法访问

#### 2. 安装项目依赖

```bash
# 在项目根目录执行
pip3 install -r requirements.txt

# 设置 Python 环境变量
source development.dev
```

#### 3. 准备 Docker 镜像

```bash
# 进入 ddosking 目录并创建 image 目录
cd ddosking
mkdir image
cd image

# 请将下载的镜像文件放在 image 目录中，加载 Docker 镜像
docker load -i ddosking.tar
docker load -i ollama.tar
docker load -i unbound.tar
```

镜像下载地址：链接：https://rec.ustc.edu.cn/share/e961def0-f679-11ef-8767-d33b1f047ce8（密码：1958）

### 启动仿真环境

```bash
# 在 ddosking 目录下执行
python3 ddosking.py

# 构建并启动 Docker 容器
cd output
docker-compose build  # 首次构建大约需要半个小时
docker-compose up

# 关闭仿真环境
docker-compose down
```

为了保证伪造报文发包正常工作，需要清除 Docker 构建的 NAT 规则：

```bash
iptables -t nat -F
```

> **提示**：建议清除前先保存 NAT 规则，以便恢复和调试。

在浏览器中访问以下地址查看网络拓扑图：

```
http://127.0.0.1:8080/map.html
```

## 攻击配置

僵尸网络的 Bot 需要配置 C2 服务器、反射器和 Unbound DNS 解析器的 IP 地址。

### C2服务器设置

```bash
cd /root/c2
go run main.go  # 启动 C2 服务器，开始监听
```

### 僵尸节点设置

```bash
# 自动配置，无需手动操作
cd /root/bot
echo 10.150.0.71 > serverfile/c2.txt  # 添加 C2 服务器 IP 地址
echo -e "10.171.0.71\n10.170.0.71" > serverfile/reflector.txt  # 添加反射器 IP 地址
echo 10.152.0.71 > serverfile/resolver.txt  # 添加 Unbound 服务器 IP 地址
go run main.go  # 启动服务，连接 C2 服务器
```

### 反射器节点设置

```bash
cd /root/reflector
go run main.go  # 启动服务
# 输入 1 开启监听
```

### Unbound服务器设置

```bash
# 用于脉冲攻击
/etc/init.d/unbound start  # 启动服务
```

### DeepSeek节点设置

```bash
cd /usr/local/bin/

# 预装 tmux，可以先输入 tmux 再输入命令启动
tmux
OLLAMA_HOST=0.0.0.0 ollama serve

# 启动后，在另一个终端输入以下命令开启终端会话
ollama run deepseek-r1:1.5b
```

### 攻击参数调整

可以在 `bot/attacker/attack` 目录中修改发包速率和其他攻击参数。

## 注意事项

1. **反射放大攻击流量限制**：反射放大器收到报文后会调用函数构造报文，然后才会转发。受限于 CPU 性能，反射放大攻击的流量相比于 UDP 直接攻击要小不少（攻击停止后，反射器仍然会处理未处理的报文，攻击会延长），可适当调控攻击速率。
2. **网络问题**：使用过程中若遇到网络问题，可尝试清除 iptables 规则。
3. **安全使用**：请在安全的环境中使用本工具，仅用于学习和研究目的。

## 技术栈

- **Docker**: 容器化技术
- **SeedEmu**: 网络仿真框架
- **Go**: 攻击脚本开发语言
- **Python**: 环境配置脚本
- **Ollama**: 部署 DeepSeek 1.5B 模型

## 许可证
本项目基于 [GNU通用公共许可证第3版（GNU GPL v3）](https://www.gnu.org/licenses/gpl-3.0.html) 发布。详细条款请参阅项目根目录中的 `LICENSE` 文件。

## 问题反馈

如遇问题或有改进建议，请通过以下方式联系：

- 提交GitHub Issue：[项目仓库Issue页面](https://github.com/yourusername/ddosking/issues)

## 免责声明

本项目仅供安全研究和教育目的使用，请勿用于任何非法活动。用户需承担因不当使用可能产生的一切法律责任。