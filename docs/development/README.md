# 开发文档

本目录包含 todoIng 项目的开发相关文档。

## 文档列表

- [开发计划概览](development-plan.md) - 总体开发计划
- [核心模块开发计划](core-module-plan.md) - 用户认证和基础任务管理功能
- [历史追踪模块开发计划](history-tracking-plan.md) - 任务变更历史记录和追踪功能
- [团队管理模块开发计划](team-management-plan.md) - 团队协作和权限管理功能
- [高级功能模块开发计划](advanced-features-plan.md) - 生命历程、事项总结、成就系统等
- [优化和完善阶段开发计划](optimization-plan.md) - 性能优化、测试完善和部署准备

## 目录结构

```
development/
├── README.md                    # 本文件
├── development-plan.md          # 开发计划概览
├── core-module-plan.md          # 核心模块开发计划
├── history-tracking-plan.md     # 历史追踪模块开发计划
├── team-management-plan.md      # 团队管理模块开发计划
├── advanced-features-plan.md    # 高级功能模块开发计划
└── optimization-plan.md         # 优化和完善阶段开发计划
```

## 开发规范

在进行开发前，请阅读以下文档以了解项目的技术架构和设计原则：

- [技术设计](../technical-design.md)
- [数据库设计](../database-design.md)
- [API 设计](../api-design.md)
- [UI/UX 设计](../ui-ux-design.md)
- [实现技术方案](../implementation-plan.md)

## 开发环境要求

请确保您的开发环境满足以下要求：

- Node.js >= 18.x (推荐使用 nvm 管理多个 Node.js 版本)
- MongoDB >= 5.0 (建议使用 MongoDB Atlas 作为云数据库解决方案)
- npm >= 8.x (或使用 yarn >= 1.22 作为替代包管理器)
- Git (版本控制工具)
- Docker (可选，用于容器化部署)
- VS Code (推荐开发编辑器)

## 贡献指南

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 代码规范

- 遵循项目中已有的代码风格
- 使用 TypeScript 进行类型检查
- 编写单元测试以确保代码质量
- 提交前运行 lint 和测试命令