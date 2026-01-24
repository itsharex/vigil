# Vigil

> **Proactive, lightweight server monitoring.**

![Build Status](https://github.com/pineappledr/vigil/actions/workflows/build.yml/badge.svg)
![License](https://img.shields.io/github/license/pineappledr/vigil)
![Go Version](https://img.shields.io/github/go-mod/go-version/pineappledr/vigil)
![SQLite Version](https://img.shields.io/badge/SQLite-v1.44.3-003B57?logo=sqlite&logoColor=white)

**Vigil** is a next-generation monitoring system built for speed and simplicity. It provides instant visibility into your infrastructure with a modern web interface and predictive health analysis, ensuring you never miss a critical hardware failure.

Works on **any Linux system** (Ubuntu, Debian, Proxmox, Unraid, Fedora, etc.).

---

## 🚀 Features

- **🔥 Lightweight Agent:** Single Go binary with zero dependencies. Deploy it on any server in seconds.
- **🐳 Docker Server:** The central hub is containerized for easy deployment via Docker or Compose.
- **💻 Responsive Web Dashboard:** Beautiful Flutter-based web interface that works perfectly on Desktop and Mobile browsers.
- **🔍 Predictive Health Check:** Advanced analysis to determine if a drive is *actually* dying or just old.
- **🤖 Telegram Alerts:** Get instant notifications via a Telegram Bot when a drive fails.

---

## 📋 Requirements

Vigil is lightweight, but it relies on standard system tools to talk to your hardware.

**Essential:**
- **Linux OS:** (64-bit recommended)
- **Root/Sudo Access:** Required to read physical disk health.
- **smartmontools:** The core engine for reading HDD, SSD, and NVMe health data.

**Recommended:**
- **nvme-cli:** Provides enhanced detail for NVMe drives.

**Install Requirements:**
```bash
# Ubuntu / Debian / Proxmox
sudo apt update && sudo apt install smartmontools nvme-cli -y

# Fedora / CentOS / RHEL
sudo dnf install smartmontools nvme-cli

# Arch Linux
sudo pacman -S smartmontools nvme-cli