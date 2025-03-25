# 🛡️ DDoSKING - Automated DDoS Attack Simulation Tool

<div align="center">

### *A DDoS Attack Simulation Environment for Cybersecurity Research and Testing*
  
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?logo=docker)](https://docs.docker.com/engine/install/)
[![Memory](https://img.shields.io/badge/Memory-Recommended24GB-red)](https://github.com/seed-labs/seed-emulator)
[![Python](https://img.shields.io/badge/Python-3.x-yellow?logo=python)](https://www.python.org/)
[![Go](https://img.shields.io/badge/Go-Driver-00ADD8?logo=go)](https://golang.org/)
  
</div>

## 📋 Project Overview

DDoSKING is a cutting-edge, Docker-based DDoS attack simulation environment designed for researching and testing various DDoS attack techniques. Built using the Seed-emulator framework, it provides a visualized network environment. Attack scripts are developed in Go and deployed within Docker containers.

This comprehensive simulation environment replicates mainstream DDoS attack methods and botnet attacks, including simulations targeting AI services like DeepSeek. It enables the simulation of complete internet and DDoS attack scenarios on a single host, making it an advanced all-in-one DDoS testing platform for security researchers and professionals.

> **⚠️ Note**: Deploying DeepSeek 1.5B requires at least 8GB of memory.

## ✨ Why Choose DDoSKING?

DDoSKING stands out as the preferred DDoS simulation tool with the following key advantages:

- **🌟 All-in-One Solution**: Simulate attack infrastructure in a single environment
- **🔄 Comprehensive Attack Coverage**: Includes major DDoS attack types and new pulse attacks like DNSBomb
- **🧠 AI Service Attack Simulation**: Simulates attacks on AI models like DeepSeek
- **🖼️ Visualized Network Topology**: Interactive visualization for better understanding of attack paths
- **🔬 Research-Grade Testing**: Suitable for academic research and security testing
- **🛠️ Highly Customizable**: Easily adjustable parameters for tailored attack scenarios
- **🌐 Fully Customizable Network Topology**: Build custom botnet infrastructures for specific research scenarios

## 🎯 What Can You Do with DDoSKING?

- **Security Research**: Understand modern DDoS attack mechanisms in a controlled environment
- **Custom Botnet Development**: Design and build your own botnet architecture for testing specific attack scenarios
- **Defense Testing**: Develop and test DDoS mitigation strategies without impacting production systems
- **Education**: A perfect teaching tool for cybersecurity courses and training programs
- **Security Auditing**: Assess network infrastructure resilience against various attack types
- **AI Service Hardening**: Test and improve AI services' robustness against targeted attacks
- **Performance Benchmarking**: Measure how different systems handle various traffic loads
- **Security Product Testing**: Validate the effectiveness of DDoS protection products

## 🌐 Network Customization Features

DDoSKING offers unparalleled flexibility in network design, enabling researchers to:

- **Create Custom Botnet Topologies**: Design network structures that match real-world scenarios or theoretical models
- **Scale Your Botnet**: Add as many zombie nodes as hardware allows to test large-scale attacks
- **Configure Node Relationships**: Define C2 server hierarchies and zombie communication patterns
- **Simulate Geographic Distribution**: Create network segments simulating geographically distributed botnets
- **Simulate Different Network Conditions**: Test scenarios with varying bandwidth, latency, and packet loss
- **Integrate Custom Attack Scripts**: Add your own attack methods implemented in Go
- **Combine Attack Types**: Create complex multi-vector attack scenarios

The SeedEmu framework supporting DDoSKING allows you to simulate nearly any network configuration, giving you complete freedom to build the precise botnet infrastructure needed for your research or testing.

## 📊 Performance and Features

DDoSKING has been rigorously tested to provide:

- Simulated network environments with hundreds of nodes
- Support for multiple attack vectors simultaneously
- Real-time monitoring of attack effects
- Highly realistic network behavior modeling
- Containerized approach for easy deployment and scalability

## 🏆 Comparison with Other Solutions

| Feature | DDoSKING | Traditional DDoS Tools | Network Simulators |
|:-------:|:--------:|:----------------------:|:------------------:|
| Complete DDoS Infrastructure | ✅ | ❌ | ❌ |
| AI Service Attack Simulation | ✅ | ❌ | ❌ |
| Pulse Attack Support | ✅ | ❌ | ❌ |
| Visualized Topology | ✅ | ❌ | ✅ |
| Single-Host Deployment | ✅ | ✅ | ❌ |
| Educational Value | ✅ | ⚠️ | ✅ |
| Research Applications | ✅ | ⚠️ | ✅ |

## 💻 System Requirements

To fully run the simulation environment, it is recommended to have more than 24GB of memory.

## 🧩 System Components

This simulation environment replicates a complete DDoS attack infrastructure, consisting of the following independently built components:

| Component | Quantity | Description |
|:--------------------:|:--------:|:-------------------------------------:|
| 🎮 C2 Server | 1 | Sends attack commands to zombie machines |
| 🔄 Reflection Amplification Servers | 5 | Used for Layer 4 DDoS reflection amplification attacks |
| 🤖 Zombie Machines | 2 | Execute various attacks |
| 🔍 Unbound DNS Resolver | 1 | Configured for pulse attacks |
| 🌐 DNS Authoritative Server | 1 | Accumulates and amplifies pulse attack requests |
| 🧠 DeepSeek 1.5B Server | 1 | Simulates HTTP attacks targeting AI services |

> **ℹ️ Note**: This simulation environment focuses on DDoS simulation and does not simulate botnet propagation or communication mechanisms (e.g., brute-forcing weak passwords, DGA searching for C2, etc.). These features may be added in future versions.

Other nodes are automatically generated by SeedEmu. For more details, refer to the [SeedEmu official documentation](https://github.com/seed-labs/seed-emulator).

## 💥 Attack Types

DDoSKING covers various types of DDoS attacks, primarily categorized as follows:

### 1. 🌊 Link Flooding Attacks (Layer 4)

Attacks that saturate network bandwidth with high traffic, including:

- **Direct Attacks**: UDP Flood
- **Reflection Amplification Attacks**: DNS, NTP, CLDAP, SSDP, etc.

### 2. 🔋 Resource Exhaustion Attacks (Layer 7)

Attacks that deplete server computational resources, including:

- **HTTP Flood**: Lightweight web service with a single page
- **Complex Prompts for DeepSeek**: Consumes AI service resources
- **SYN Flood**: Exhausts the target's half-open connection queue

### 3. ⚡ Pulse Attacks

Pulse attacks send high-bandwidth packets in short bursts, filling target queues, causing timeouts, and triggering TCP congestion control, degrading the target's TCP services.

- **DNSBomb**: A practical and powerful pulsing DoS attack exploiting DNS queries and responses. Reference: [DNSBomb: A New Practical-and-Powerful Pulsing DoS Attack Exploiting DNS Queries-and-Responses | IEEE Conference Publication | IEEE Xplore](https://ieeexplore.ieee.org/abstract/document/10646654)

- **DNSBoomerang**: An enhancement of DNSBomb that significantly increases cumulative data (the number of cumulative packets grows with the number of DNS reflectors). In public experiments, attackers accumulated requests at 530kbps from 500 different reflectors. After 17 seconds (13,700 requests), the reflected bandwidth reached 108Mbps for about 1 second, achieving a 204x amplification.
  <div style=text-align:left>
    <img src="pictures/dnsboomerang2.png" width="40%" />
    <img src="pictures/dnsboomerang1.png" width="40%" />
  </div>

## 🚀 Environment Setup

It is recommended to set up the environment on Linux. Windows users can use WSL.

### 📦 Installation Steps

#### 1. Install and Configure Docker

- Install Docker: [Refer to the official documentation](https://docs.docker.com/engine/install/)
- **Note**: Users in mainland China may need to configure Docker mirror sources due to potential access issues with DockerHub.

#### 2. Install Project Dependencies

```bash
# Run in the project root directory
pip3 install -r requirements.txt

# Set Python environment variables
source development.dev
```

### 🏃 Start the Simulation Environment

```bash
# Run in the root directory
python3 ddosking.py

# Build and start Docker containers
cd output
docker-compose build  # Initial build may take about 30 minutes
docker-compose up

# Shut down the simulation environment
docker-compose down
```

To ensure proper operation of spoofed packet sending, you need to clear the NAT rules built by Docker:

```bash
iptables -t nat -F
```

> **💡 Tip**: It is recommended to save the NAT rules before clearing them so they can be restored during debugging.

Access the following URL in your browser to view the network topology:

```
http://127.0.0.1:8080/map.html
```

## ⚙️ Attack Configuration

Zombie machines in the botnet need to be configured with the IP addresses of the C2 server, reflectors, and Unbound DNS resolver.

### C2 Server Setup

```bash
cd /root/c2
go run .  # Start the C2 server to begin listening
```

### Zombie Node Setup

```bash
# Automatically configured, no manual operation required
cd /root/bot
go run .  # Start the service and connect to the C2 server
```

### Reflector Node Setup

```bash
cd /root/reflector
go run .  # Start the service
# Enter 1 to begin listening
```

### Unbound Server Setup

```bash
# Used for pulse attacks
service unbound start  # Start the service
```

### DeepSeek Node Setup

```bash
# Pre-install tmux, then enter tmux and run the command to start
tmux
OLLAMA_HOST=0.0.0.0 ollama serve

# After starting, press ctrl b+d to exit, then start the terminal session in another terminal
ollama run deepseek-r1:1.5b
```

### Attack Parameter Adjustment

You can modify the packet sending rate and other attack parameters in `bot/attacker/attack/config.go`.

## 📝 Notes

1. **Reflection Amplification Attack Traffic Limitation**: After receiving packets, reflectors construct and forward packets. Due to CPU performance limitations, reflection amplification attack traffic is much smaller than direct UDP attacks (after the attack stops, reflectors continue processing unhandled packets, extending the attack duration). Adjust the attack rate accordingly.
2. **Network Issues**: If network issues occur during use, try clearing the iptables rules.
3. **Safe Usage**: Use this tool in a secure environment for learning and research purposes only.

## 🔧 Tech Stack

- **🐳 Docker**: Containerization technology
- **🌐 SeedEmu**: Network simulation framework
- **🚀 Go**: Attack script development language
- **🐍 Python**: Environment setup scripts
- **🧠 Ollama**: Deploying the DeepSeek 1.5B model

## 🔮 Future Development

We are continuously improving DDoSKING and plan to add:

- More AI service attack simulations
- Botnet propagation mechanisms
- Additional attack vectors
- Enhanced visualization and analysis
- Integration with common security tools
- Performance optimizations

## 📜 License
This project is licensed under the [GNU General Public License v3.0 (GNU GPL v3)](https://www.gnu.org/licenses/gpl-3.0.html). For detailed terms, refer to the `LICENSE` file in the project root directory.

## ⚠️ Disclaimer

This project is intended for security research and educational purposes only. Do not use it for any illegal activities. Users are fully responsible for any consequences arising from improper use.
