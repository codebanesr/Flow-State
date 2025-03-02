<div align="center">

<pre style="color: #3cb371">   
   ▄████████  ▄█        ▄██████▄   ▄█     █▄     ▄████████     ███        ▄████████     ███        ▄████████ 
  ███    ███ ███       ███    ███ ███     ███   ███    ███ ▀█████████▄   ███    ███ ▀█████████▄   ███    ███ 
  ███    █▀  ███       ███    ███ ███     ███   ███    █▀     ▀███▀▀██   ███    ███    ▀███▀▀██   ███    █▀  
 ▄███▄▄▄     ███       ███    ███ ███     ███   ███            ███   ▀   ███    ███     ███   ▀  ▄███▄▄▄     
▀▀███▀▀▀     ███       ███    ███ ███     ███ ▀███████████     ███     ▀███████████     ███     ▀▀███▀▀▀     
  ███        ███       ███    ███ ███     ███          ███     ███       ███    ███     ███       ███    █▄  
  ███        ███▌    ▄ ███    ███ ███ ▄█▄ ███    ▄█    ███     ███       ███    ███     ███       ███    ███ 
  ███        █████▄▄██  ▀██████▀   ▀███▀███▀   ▄████████▀     ▄████▀     ███    █▀     ▄████▀     ██████████ 
             ▀                                                                                               
</pre>

# CloudStar

### Enterprise-Grade Cloud Desktops for AI Agents & Collaborative Computing

[![Docker Pulls](https://img.shields.io/docker/pulls/shanurcsenitap/vnc_chrome_debug?style=flat-square&color=3cb371&labelColor=333333)](https://hub.docker.com/r/shanurcsenitap/vnc_chrome_debug)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-3cb371.svg?style=flat-square)](https://www.gnu.org/licenses/agpl-3.0)
[![CI/CD Pipeline](https://img.shields.io/github/actions/workflow/status/codebanesr/orchestrator/build.yml?style=flat-square&color=3cb371&labelColor=333333)](https://github.com/codebanesr/orchestrator/actions)
[![Documentation](https://img.shields.io/badge/docs-latest-3cb371.svg?style=flat-square&labelColor=333333)](https://github.com/codebanesr/orchestrator/wiki)
[![Discord](https://img.shields.io/discord/1234567890?style=flat-square&color=3cb371&labelColor=333333&label=discord)](https://discord.gg/cloudstar)

</div>

## 🚀 Overview

CloudStar is a powerful, scalable platform for virtual Linux containers running full-featured desktops and browsers with complete external control. Our solution enables AI agents to interact with real desktop environments just like humans, while also offering collaborative capabilities for teams and individuals.

> **Featured On**: [Awesome-Containers List](https://github.com/awesome-containers) | **Demo**: [live.orchestrator.dev](https://live.orchestrator.dev)

## 🎯 Key Use Cases

### AI Desktop Agents
- **Autonomous Web Research**: AI agents with full browsing capabilities can research, summarize, and analyze web content
- **Visual Decision Making**: Agents can interact with visual interfaces, analyzing charts, graphs, and UI elements
- **Multi-step Workflows**: Agents can execute complex workflows across multiple applications
- **Automation Testing**: Build and test applications with AI agents that mimic real user behavior

### Collaborative Computing
- **Pair Programming**: Code together in real-time with shared terminal and IDE access
- **Interactive Presentations**: Give live demos with audience participation
- **Virtual Classroom**: Teach programming, design, or any digital skill with shared desktops
- **Co-browsing Customer Support**: Assist customers through shared browser sessions

### Entertainment & Community
- **Synchronized Movie Watching**: Watch films together with perfect synchronization
- **Collaborative Gaming**: Play browser games together or watch friends play
- **Virtual Watch Parties**: Create themed viewing rooms for TV shows, sports events, or conferences
- **Creative Collaboration**: Work on digital art, music production, or video editing together

### Enterprise Solutions
- **Secure Remote Work**: Provide isolated desktop environments that never store sensitive data locally
- **Compliance & Audit**: Record desktop sessions for regulatory compliance or training
- **Legacy Application Access**: Run legacy software in containers accessible from any device
- **Cross-platform Testing**: Test applications across multiple browser and OS configurations

## 🏗️ Architecture

```mermaid
flowchart TD
    classDef client fill:#99FF99,stroke:#00CC00,color:#1a5c1a
    classDef desktop fill:#99FF99,stroke:#00CC00,color:#1a5c1a
    classDef core fill:#99FF99,stroke:#00CC00,color:#1a5c1a
    
    Client[("Client / AI Agent")]:::client
    
    subgraph VirtualDesktops["Virtual Desktop Containers"]
        direction TB
        D1["Desktop Container 1"]:::desktop
        D2["Desktop Container 2"]:::desktop
        DN["Desktop Container N..."]:::desktop
    end
    
    subgraph ControlInfrastructure["Control Infrastructure"]
        direction TB
        CH["Control Hub Service<br/><i>port: 8090</i>"]:::core
        LB["Load Balancer<br/>(Fabio LB)<br/><i>ports: 9999, 9998</i>"]:::core
        SD["Service Discovery<br/>(Consul)<br/><i>port: 8500</i>"]:::core
    end
    
    Client -->|"HTTP Request"| LB
    LB -->|"Load Balance"| CH
    CH -->|"Manage Containers"| D1
    CH -->|"Manage Containers"| D2
    CH -->|"..."| DN
    CH -.->|"Service Registration"| SD
    LB -.->|"Health Checks"| SD
```

## 💡 Features

<div align="center">

| Isolation & Control | Scalability | Monitoring & Security |
|---------------------|-------------|-----------------------|
| <img src="https://img.icons8.com/?size=100&id=RjmW1_uvskWO&format=png&color=228B22" width=50> Fully Isolated Containers | <img src="https://img.icons8.com/3d-fluency/50/factory.png" width=50> Auto-scaling Clusters | <img src="https://img.icons8.com/3d-fluency/50/visible.png" width=50> Real-time Metrics |
| **External Control** | **Enterprise Ready** | **Zero Trust Security** |
| <img src="https://img.icons8.com/3d-fluency/50/remote-desktop.png" width=50> Remotely Managed Desktops & Browsers | <img src="https://img.icons8.com/3d-fluency/50/network.png" width=50> Multi-Node Support | <img src="https://img.icons8.com/3d-fluency/50/lock.png" width=50> Mutual TLS & RBAC |

</div>

- **Ephemeral Environments**: Spin up disposable desktops in seconds, then destroy them without a trace
- **Full Desktop Access**: Complete Linux desktop environment with audio, video, and peripheral support
- **Collaborative Tools**: Built-in text, voice, and video chat for human-to-human or AI-to-human collaboration
- **Programmable Interface**: Control desktops via API or WebSocket for automation and integration
- **Resource Efficient**: Optimized containers use minimal resources while providing full desktop functionality
- **Cross-Platform**: Access from any device with a modern web browser - no client software needed

## 🚤 Quick Start

### Prerequisites
- Docker 20.10+
- Go 1.22+
- 4GB RAM (8GB recommended)
- Linux kernel >5.10

### Installation & Setup
```bash
# Clone with depth for a faster download
git clone --depth=1 https://github.com/codebanesr/orchestrator.git
cd orchestrator

# Install dependencies
go mod download

# Build and run using make commands
make swagger  # Generate Swagger documentation
make build    # Build the application
make run      # Run the application (includes swagger generation)
```

## ⚙️ Configuration

### 🔧 Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `WIDTH` | Width of the virtual desktop/browser window | `1024` |
| `HEIGHT` | Height of the virtual desktop/browser window | `768` |
| `LOG_LEVEL` | Logging verbosity level (debug, info, warn, error) | `debug` |
| `PORT` | Port on which the service will run | `3000` |
| `BEHIND_PROXY` | Whether the service is running behind a proxy | `false` |
| `CONSUL_HTTP_ADDR` | Address of the Consul service | `localhost:8500` |

### 🔌 Exposed Ports & Endpoints

| Service | Port | Access URL | Description |
|---------|------|------------|-------------|
| Orchestrator API | 8090 | `http://localhost:8090/` | Main control API for container management |
| Fabio Load Balancer | 9999 | `http://localhost:9999/` | HTTP traffic to containers |
| Fabio UI | 9998 | `http://localhost:9998/` | Fabio management interface |
| Consul | 8500 | `http://localhost:8500/` | Service discovery dashboard |
| Container noVNC | 6901 | `/{containerID}/novnc/` | Web-based VNC client |
| Container VNC | 5901 | `/{containerID}/vnc/` | Direct VNC server access |

## 🔒 Security
- **Zero Trust Architecture**: Mutual TLS between all components
- **Isolated Environments**: Each desktop runs in its own container with resource limits
- **Ephemeral Sessions**: No data persists after container termination
- **Secure Communication**: Encrypted WebRTC for audio/video and HTTPS for all traffic
- **Granular Permissions**: Role-based access control for all container operations

## 📚 Documentation
The API documentation is available via Swagger UI:
- When running behind proxy: `http://localhost:9999/orchestrator/swagger/`
- Direct access: `http://localhost:8090/swagger/`

## 🌐 Community & Support

- 📚 [Documentation](https://github.com/codebanesr/orchestrator/wiki)
- 💬 [Discord Community](https://discord.gg/cloudstar)
- 🐦 [Twitter](https://twitter.com/cloudstar)
- 📧 Enterprise Support: contact@orchestrator.dev

## 🤝 Contributing
We welcome contributions! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Add tests for new functionality
4. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

## 📜 License
This project is licensed under the AGPL v3 License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <sub>Built with ❤️ by Codebanesr | Inspired by <a href="https://github.com/m1k1o/neko">Neko</a></sub>
</div>