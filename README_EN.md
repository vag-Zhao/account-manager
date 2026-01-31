# 🔐 Account Manager

Cross-platform Account Management System | [中文](README.md)

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)
![Wails](https://img.shields.io/badge/Wails-v2.11-DF0000)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

- 📊 **Account Management** - CRUD, AES-256 encryption, CSV import, search/filter/pagination
- 📧 **Email Notifications** - SMTP config, auto expiry reminders, scheduled checks
- 🚀 **Remote Deployment** - SSH deployment, service control, environment detection
- 🌐 **Bilingual UI** - Chinese/English support

## 🛠️ Tech Stack

| Backend | Frontend |
|---------|----------|
| Go + Wails v2 | Vue 3 + TypeScript |
| SQLite + GORM | Naive UI + Pinia |
| AES-256 + SSH | Vite |

## 🚀 Quick Start

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and run
git clone https://github.com/vag-Zhao/account-manager.git
cd account-manager
cd frontend && npm install && cd ..
wails dev
```

## 📸 Screenshots

| Dashboard | Email Settings |
|-----------|----------------|
| ![Dashboard](shortcut/1.png) | ![Email](shortcut/2.png) |

## 📄 License

MIT License - See [LICENSE](LICENSE)

## 📞 Contact

- 📧 zgs3344@hunnu.edu.cn
- 🐛 [Issues](https://github.com/vag-Zhao/account-manager/issues)
