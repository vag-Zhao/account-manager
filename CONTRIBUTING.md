# 贡献指南 | Contributing Guide

感谢你考虑为 Account Manager 做出贡献！

## 开发环境设置

### 前置要求
- Go 1.24+
- Node.js 16+
- Wails CLI v2.11.0+

### 安装步骤
```bash
# 克隆仓库
git clone https://github.com/vag-Zhao/account-manager.git
cd account-manager

# 安装前端依赖
cd frontend
npm install
cd ..

# 安装 Go 依赖
go mod download

# 运行开发模式
wails dev
```

## 开发流程

### 1. Fork 仓库
点击右上角的 "Fork" 按钮

### 2. 创建分支
```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/your-bug-fix
```

### 3. 提交代码
```bash
git add .
git commit -m "feat: add some feature"
```

提交信息格式：
- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `style:` 代码格式化
- `refactor:` 代码重构
- `test:` 测试相关
- `chore:` 构建/工具相关

### 4. 推送到 Fork
```bash
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request
在 GitHub 上创建 Pull Request

## 代码规范

### Go 代码
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释

### TypeScript/Vue 代码
- 遵循 ESLint 规则
- 使用 TypeScript 类型注解
- 组件命名使用 PascalCase

## 测试

```bash
# 运行 Go 测试
go test ./...

# 构建项目
wails build
```

## 报告问题

使用 GitHub Issues 报告问题，请包含：
- 清晰的问题描述
- 复现步骤
- 期望行为和实际行为
- 环境信息（OS、版本等）
- 相关截图或日志

## 行为准则

- 尊重所有贡献者
- 保持友好和专业
- 接受建设性的批评
- 关注对项目最有利的事情

## 许可证

通过贡献代码，你同意你的贡献将在 MIT 许可证下发布。
