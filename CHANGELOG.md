# TodoIng 变更日志

## [1.9.0] - 2025-08-10

### 🐳 Docker生态系统重大升级

- **完成双后端架构部署方案**
  - 新增Golang版本完整Docker Compose配置（docker-compose.golang.yml）
  - 完善Node.js版本Docker Compose配置（docker-compose.nodejs.yml）
  - 两套方案并行，支持不同技术栈需求

### 🔧 构建与部署优化

- **Golang后端Docker化改进**
  - 修复Dockerfile多阶段构建配置，添加production和grpc目标
  - 更新端口配置从5001到5004，避免与Node.js版本冲突
  - 添加健康检查支持（curl）和安全配置
  
- **前端配置专业化**
  - 创建nginx-golang.conf专门用于Golang后端部署
  - 创建nginx-nodejs.conf专门用于Node.js后端部署
  - 修复服务间通信配置，确保正确的后端服务发现

### 🐛 环境变量与配置修复

- **邮件服务配置优化**
  - 修复EMAIL_FROM环境变量特殊字符问题
  - 完善邮件验证码功能的SMTP配置
  - 解决"501 mail from address must be same as authorization user"错误

- **Docker Compose配置完善**
  - 移除无效的target配置避免构建失败
  - 添加缺失的EMAIL_FROM环境变量映射
  - 修复前端构建参数传递问题

### 🚀 部署验证

- **功能验证完成**
  - Node.js版本：所有服务正常运行，邮件验证码功能正常
  - Golang版本：所有服务正常运行，API健康检查通过
  - 验证码显示问题已解决，环境变量正确生效

### 📝 文档与规范

- 完善各种部署场景的Docker配置
- 优化服务间依赖关系和健康检查
- 统一容器命名规范，避免冲突

### 🔄 技术债务清理

- 清理无效的Docker配置项
- 优化构建上下文和文件映射
- 统一环境变量命名和默认值

---

*本版本完成了TodoIng项目的完整Docker生态系统升级，实现了Node.js和Golang双后端架构的容器化部署，为项目的扩展性和维护性奠定了坚实基础。*
