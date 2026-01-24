# Vigil 👁️

> **Modern, lightweight server monitoring for the 2026 era.**

![Build Status](https://github.com/KutonnoZer0/vigil/actions/workflows/build.yml/badge.svg)
![License](https://img.shields.io/github/license/KutonnoZer0/vigil)
![Go Version](https://img.shields.io/github/go-mod/go-version/KutonnoZer0/vigil)

**Vigil** is a next-generation monitoring system designed to replace legacy tools like Scrutiny. It focuses on simplicity, mobile-first design, and AI-powered health analysis for your home lab.

---

## 🚀 Features

- **🔥 Single Binary Architecture:** No complex databases or multi-container setups. Just one file.
- **📱 Mobile-First:** Native iOS & Android app (Flutter) for monitoring on the go.
- **🧠 AI Analysis:** Integrated LLM checks to determine if a drive is *actually* dying or just old.
- **⚡ Real-time S.M.A.R.T. Tracking:** Monitors temperature, reallocated sectors, and power-on hours.
- **🔔 Push Notifications:** Get alerted instantly on your phone when a drive fails.

---

## 🛠️ Architecture

Vigil follows a clean **Hub & Spoke** model:

1.  **Vigil Agent (Go):** A lightweight binary that runs on your servers (Proxmox, Ubuntu, Unraid). It wraps `smartctl` to read raw disk health.
2.  **Vigil Server (Go):** The central hub that receives data, stores it in SQLite, and serves the API.
3.  **Vigil UI (Flutter):** A beautiful, responsive interface that runs as a Web App *and* a Native Mobile App.

---

## 📦 Installation

### 1. The Agent (Proxmox / Linux)
The agent is a single static binary. You can download it from the [Releases](https://github.com/KutonnoZer0/vigil/releases) page or build it yourself.

**One-Liner Install (Coming Soon):**
```bash
curl -sL https://vigil.sh/install-agent | sudo bash
```

**Manual Build:**
```bash
git clone https://github.com/KutonnoZer0/vigil.git
cd vigil
go build -o vigil-agent ./cmd/agent
sudo ./vigil-agent
```

### 2. The Server (Docker)
```bash
docker run -d   -p 8090:8090   -v vigil-data:/data   --name vigil   ghcr.io/kutonnozer0/vigil:latest
```

---

## 🗺️ Roadmap

- [x] **Phase 1: The Agent** - Build a Go binary to read local SMART data.
- [ ] **Phase 2: The Server** - Create the API to receive agent data.
- [ ] **Phase 3: The UI** - Build the Flutter Web Dashboard.
- [ ] **Phase 4: Mobile App** - Compile Flutter for iOS/Android.
- [ ] **Phase 5: AI Integration** - Connect to Ollama for drive health analysis.

---

## 🤝 Contributing

Contributions are welcome! Please read the [CONTRIBUTING.md](CONTRIBUTING.md) details (coming soon) for details on our code of conduct, and the process for submitting pull requests.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
