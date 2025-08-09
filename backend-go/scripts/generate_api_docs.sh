#!/bin/bash

# 生成完整的 API 文档
echo "🚀 生成 TodoIng 完整 API 文档..."

# 运行文档生成器
go run tools/generate_complete_api.go

echo ""
echo "📖 文档访问方式："
echo "   完整 API 文档: http://localhost:5004/api-docs"
echo "   Swagger UI: http://localhost:5004/swagger/"
echo "   文档中心: http://localhost:5004/docs/"
echo ""
echo "📁 生成的文件："
echo "   - docs/api_complete.json (完整 API 定义)"
echo ""
echo "✅ API 文档生成完成！"
