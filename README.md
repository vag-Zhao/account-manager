# 🔐 Account Manager

一个基于 Wails v2 构建的跨平台账户管理系统，提供安全的凭证管理、自动化邮件提醒和远程部署功能。

A cross-platform account management system built with Wails v2, featuring secure credential management, automated email reminders, and remote deployment capabilities.

---

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![Vue Version](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript)
![Wails](https://img.shields.io/badge/Wails-v2.11-DF0000?logo=wails)
![License](https://img.shields.io/badge/License-MIT-green.svg)

---

## ✨ 核心特性 | Core Features

### 📊 账户管理 | Account Management
- ✅ 完整的 CRUD 操作（创建、读取、更新、删除）
- 🔒 AES-256 密码加密存储
- 🏷️ 账户类型分类（PLUS、BUSINESS、FREE）
- 📈 状态跟踪（已售/未售）
- 📅 到期日期管理
- 📥 CSV 批量导入
- 🔍 搜索、过滤、分页（50条/页）

### 🛡️ 安全特性 | Security Features
- 🔒 AES-256 密码加密存储
- 📝 审计日志记录
- 🔒 SSH 主机密钥验证

### 📧 邮件通知 | Email Notifications
- 📬 SMTP 配置支持（QQ、163、Gmail、Outlook 等）
- ⏰ 自动到期提醒
- 📝 可自定义邮件模板
- 📊 邮件发送日志
- 🕐 每小时定时检查

### 🚀 远程部署 | Remote Deployment
- 🖥️ 基于 SSH 的 Linux 服务器部署
- ⚙️ 服务启动/停止控制
- 🔍 服务器环境检测
- 📊 部署状态跟踪

### 🌍 其他功能 | Additional Features
- 🌐 双语支持（中文/英文）
- 📊 统计仪表板
- 🎨 现代化 UI 设计
- 💾 SQLite 本地数据库
- 🔄 实时数据同步

---

## 🛠️ 技术栈 | Tech Stack

### 后端 | Backend
- 🐹 **Go 1.24** - 高性能后端语言
- 🎯 **Wails v2.11.0** - 跨平台桌面应用框架
- 💾 **SQLite + GORM** - 轻量级数据库与 ORM
- 🔒 **AES-256 Encryption** - 企业级加密标准
- 📬 **SMTP Email Support** - 邮件发送功能
- ⏰ **Cron Scheduling** - 定时任务调度
- 🔑 **SSH Client** - 远程服务器连接

### 前端 | Frontend
- ⚡ **Vue 3** - 渐进式 JavaScript 框架
- 📘 **TypeScript 5.x** - 类型安全的 JavaScript
- 🎨 **Naive UI** - 现代化 Vue 3 组件库
- 📦 **Pinia** - Vue 3 状态管理
- 🌐 **Vue Router** - 单页应用路由
- 🔧 **Vite** - 下一代前端构建工具

---

## 🚀 快速开始 | Quick Start

### 前置要求 | Prerequisites

- **Go** 1.24 或更高版本
- **Node.js** 16+ 和 npm/pnpm
- **Wails CLI** v2.11.0+

安装 Wails CLI:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 安装步骤 | Installation

1. **克隆仓库 | Clone Repository**
```bash
git clone https://github.com/yourusername/account-manager.git
cd account-manager
```

2. **安装依赖 | Install Dependencies**
```bash
# 安装前端依赖
cd frontend
npm install
cd ..

# 安装 Go 依赖
go mod download
```

3. **开发模式 | Development Mode**
```bash
wails dev
```

应用将在开发模式下启动，支持热重载。前端开发服务器运行在 `http://localhost:34115`。

4. **构建生产版本 | Build for Production**
```bash
# Windows
wails build

# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64
```

构建产物位于 `build/bin/` 目录。

---

## 📁 项目结构 | Project Structure

```
account-manager/
├── frontend/               # Vue 3 前端应用
│   ├── src/
│   │   ├── components/    # Vue 组件
│   │   ├── views/         # 页面视图
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── router/        # 路由配置
│   │   └── assets/        # 静态资源
│   └── package.json
├── internal/              # Go 后端代码
│   ├── models/           # 数据模型
│   ├── database/         # 数据库操作
│   ├── services/         # 业务逻辑
│   │   ├── account.go    # 账户服务
│   │   ├── email.go      # 邮件服务
│   │   └── deploy.go     # 部署服务
│   └── utils/            # 工具函数
├── app.go                # Wails 应用入口
├── main.go               # 主程序入口
├── wails.json            # Wails 配置
├── go.mod                # Go 依赖
└── README.md
```

---

## ⚙️ 配置说明 | Configuration

### 数据库配置 | Database Configuration

应用首次启动时会自动创建 SQLite 数据库文件：
- **Windows**: `%APPDATA%/account-manager/accounts.db`
- **macOS**: `~/Library/Application Support/account-manager/accounts.db`
- **Linux**: `~/.config/account-manager/accounts.db`

### SMTP 邮件配置 | SMTP Email Configuration

在应用设置中配置 SMTP 服务器：

```yaml
SMTP 服务器: smtp.example.com
端口: 587
用户名: your-email@example.com
密码: your-app-password
发件人: your-email@example.com
启用 TLS: true
```

**常用 SMTP 配置示例：**

| 邮箱服务商 | SMTP 服务器 | 端口 | TLS |
|-----------|------------|------|-----|
| QQ 邮箱 | smtp.qq.com | 587 | ✅ |
| 163 邮箱 | smtp.163.com | 465 | ✅ |
| Gmail | smtp.gmail.com | 587 | ✅ |
| Outlook | smtp.office365.com | 587 | ✅ |

### SSH 部署配置 | SSH Deployment Configuration

配置远程服务器信息：
- **主机地址**: 服务器 IP 或域名
- **端口**: SSH 端口（默认 22）
- **用户名**: SSH 登录用户
- **密码/密钥**: 认证方式
- **部署路径**: 应用部署目录

---

## 🔒 安全说明 | Security

### 加密机制 | Encryption

- **密码加密**: 使用 AES-256 加密账户密码
- **本地存储**: 数据存储在本地 SQLite 数据库中

### 最佳实践 | Best Practices

1. ✅ 定期备份数据库文件
2. ✅ 启用审计日志以追踪操作记录
3. ✅ 不要在不安全的网络环境下使用

### 审计日志 | Audit Logging

所有关键操作都会记录审计日志：
- 账户创建/修改/删除
- 邮件发送记录
- 远程部署操作

---

## 📸 功能截图 | Screenshots

### 仪表板 | Dashboard
![Dashboard](shortcut/1.png)

### 邮件设置 |  EmailSetting
![Account List](shortcut/2.png)

---

## 🤝 贡献指南 | Contributing

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 开源协议 | License

本项目采用 MIT 协议开源 - 详见 [LICENSE](LICENSE) 文件。

---

## 📞 联系方式 | Contact

如有问题或建议，请通过以下方式联系：

- 📧 Email: zgs3344@hunnu.edu.cn
- 🐛 Issues: [GitHub Issues](https://github.com/vag-Zhao/account-manager/issues)

---

## 🙏 致谢 | Acknowledgments

- [Wails](https://wails.io/) - 优秀的 Go + Web 桌面应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Naive UI](https://www.naiveui.com/) - 现代化 Vue 3 组件库
- [GORM](https://gorm.io/) - Go 语言 ORM 库

---

<div align="center">
Made with ❤️ using Wails v2
</div>
