# todoIng

todoIng 是一个创新的任务管理系统，它不仅提供基本的任务管理功能，还为每个任务提供完整的生命周期追踪和变更历史记录，就像 Git 管理代码变更一样管理任务。

## 核心特性

- **任务生命周期管理** - 完整的任务状态跟踪（创建、进行中、暂停、完成、取消）
- **变更历史追踪** - 记录任务的每一次变更，支持随时回溯任务状态
- **可视化历史展示** - 以图形化方式展示任务变更历史，类似 Git 提交历史
- **Web 图形化界面** - 直观的用户界面，便于操作和查看任务状态
- **实时同步** - 多设备间任务状态实时同步
- **数据统计** - 提供任务完成情况的统计和分析

## 文档

详细的项目设计文档请查看 [docs](./docs) 目录：

- [项目概述](docs/project-overview.md)
- [技术设计](docs/technical-design.md)
- [UI/UX 设计](docs/ui-ux-design.md)
- [API 设计](docs/api-design.md)
- [数据库设计](docs/database-design.md)
- [实现技术方案](docs/implementation-plan.md)

开发相关文档请查看 [docs/development](./docs/development) 目录：

- [开发计划](docs/development/development-plan.md)

## 技术栈

### 前端
- React 18
- TypeScript
- Redux Toolkit
- React Router
- D3.js (用于历史数据可视化)

### 后端
- Node.js
- Express.js
- MongoDB
- Mongoose
- JWT 认证

## 安装和运行

```bash
# 克隆项目
git clone <repository-url>

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 许可证

[MIT](LICENSE)