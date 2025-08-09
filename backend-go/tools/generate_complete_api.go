package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// 生成合并的 API 文档
	combinedDoc := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "TodoIng Complete API",
			"description": "TodoIng 项目完整 API 文档 - 包含 REST 和 gRPC 接口",
			"version":     "1.0.0",
			"contact": map[string]interface{}{
				"name":  "TodoIng Team",
				"email": "support@todoing.com",
			},
		},
		"host":     "localhost:5004",
		"basePath": "/api",
		"schemes":  []string{"http", "https"},
		"consumes": []string{"application/json"},
		"produces": []string{"application/json"},
		"tags": []interface{}{
			map[string]interface{}{
				"name":        "REST-认证",
				"description": "REST API 用户认证接口",
			},
			map[string]interface{}{
				"name":        "REST-任务",
				"description": "REST API 任务管理接口",
			},
			map[string]interface{}{
				"name":        "REST-报表",
				"description": "REST API 报表管理接口",
			},
			map[string]interface{}{
				"name":        "gRPC-认证",
				"description": "gRPC API 用户认证服务",
			},
			map[string]interface{}{
				"name":        "gRPC-任务",
				"description": "gRPC API 任务管理服务",
			},
			map[string]interface{}{
				"name":        "gRPC-报表",
				"description": "gRPC API 报表管理服务",
			},
			map[string]interface{}{
				"name":        "gRPC-验证码",
				"description": "gRPC API 验证码服务",
			},
		},
		"securityDefinitions": map[string]interface{}{
			"BearerAuth": map[string]interface{}{
				"type":        "apiKey",
				"name":        "Authorization",
				"in":          "header",
				"description": "JWT token (Bearer <token>)",
			},
		},
		"paths":       map[string]interface{}{},
		"definitions": map[string]interface{}{},
	}

	// 添加 REST API 路径
	restPaths := map[string]interface{}{
		// 认证相关接口
		"/auth/register": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "用户注册",
				"description": "注册新用户账号",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "注册信息",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/RegisterRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "注册成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/RegisterResponse",
						},
					},
					"400": map[string]interface{}{
						"description": "请求参数错误",
					},
					"409": map[string]interface{}{
						"description": "用户已存在",
					},
				},
			},
		},
		"/auth/login": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "用户登录",
				"description": "通过邮箱和密码登录，支持邮箱验证码登录",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "登录信息",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/LoginRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "登录成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/LoginResponse",
						},
					},
					"401": map[string]interface{}{
						"description": "认证失败",
					},
				},
			},
		},
		"/auth/me": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "获取当前用户信息",
				"description": "获取当前登录用户的详细信息",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "用户信息",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/User",
						},
					},
					"401": map[string]interface{}{
						"description": "未授权",
					},
				},
			},
		},
		"/auth/send-email-code": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "发送注册邮箱验证码",
				"description": "向指定邮箱发送注册验证码",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "邮箱地址",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"email": map[string]interface{}{
									"type":    "string",
									"example": "user@example.com",
								},
							},
							"required": []string{"email"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "发送成功",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"message": map[string]interface{}{
									"type": "string",
								},
								"emailCodeId": map[string]interface{}{
									"type": "string",
								},
							},
						},
					},
				},
			},
		},
		"/auth/send-login-email-code": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "发送登录邮箱验证码",
				"description": "向指定邮箱发送登录验证码",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "邮箱地址",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"email": map[string]interface{}{
									"type":    "string",
									"example": "user@example.com",
								},
							},
							"required": []string{"email"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "发送成功",
					},
				},
			},
		},
		"/auth/captcha": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "获取验证码",
				"description": "生成图形验证码",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "验证码生成成功",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{
									"type":        "string",
									"description": "验证码ID",
								},
								"image": map[string]interface{}{
									"type":        "string",
									"description": "验证码图片 (Base64)",
								},
							},
						},
					},
				},
			},
		},
		"/auth/verify-captcha": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-认证"},
				"summary":     "验证验证码",
				"description": "验证用户输入的图形验证码",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "验证码信息",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{
									"type": "string",
								},
								"code": map[string]interface{}{
									"type": "string",
								},
							},
							"required": []string{"id", "code"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "验证成功",
					},
				},
			},
		},
		// 任务管理接口
		"/tasks": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "获取任务列表",
				"description": "获取当前用户的所有任务",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "任务列表",
						"schema": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/Task",
							},
						},
					},
					"401": map[string]interface{}{
						"description": "未授权",
					},
				},
			},
			"post": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "创建任务",
				"description": "创建新的任务",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "任务信息",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/CreateTaskRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "创建成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Task",
						},
					},
					"400": map[string]interface{}{
						"description": "请求参数错误",
					},
					"401": map[string]interface{}{
						"description": "未授权",
					},
				},
			},
		},
		"/tasks/export/all": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "导出所有任务",
				"description": "导出当前用户的所有任务数据",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "导出成功",
						"schema": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		"/tasks/import": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "导入任务",
				"description": "批量导入任务数据",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "任务数据",
						"schema": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/CreateTaskRequest",
							},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "导入成功",
					},
				},
			},
		},
		"/tasks/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "获取任务详情",
				"description": "根据ID获取任务详情",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "任务ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "任务详情",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Task",
						},
					},
					"404": map[string]interface{}{
						"description": "任务不存在",
					},
				},
			},
			"put": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "更新任务",
				"description": "更新任务信息",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "任务ID",
					},
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "更新的任务信息",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/UpdateTaskRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "更新成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Task",
						},
					},
					"404": map[string]interface{}{
						"description": "任务不存在",
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":        []string{"REST-任务"},
				"summary":     "删除任务",
				"description": "删除指定任务",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "任务ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "删除成功",
					},
					"404": map[string]interface{}{
						"description": "任务不存在",
					},
				},
			},
		},
		// 报表管理接口
		"/reports": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "获取报表列表",
				"description": "获取用户的所有报表",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "报表列表",
						"schema": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/Report",
							},
						},
					},
				},
			},
		},
		"/reports/generate": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "生成报表",
				"description": "生成新的报表",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "报表生成参数",
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"type": map[string]interface{}{
									"type":        "string",
									"description": "报表类型",
									"enum":        []string{"daily", "weekly", "monthly"},
								},
								"dateRange": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"start": map[string]interface{}{
											"type":   "string",
											"format": "date",
										},
										"end": map[string]interface{}{
											"type":   "string",
											"format": "date",
										},
									},
								},
							},
							"required": []string{"type"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "生成成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Report",
						},
					},
				},
			},
		},
		"/reports/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "获取报表详情",
				"description": "根据ID获取报表详情",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "报表ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "报表详情",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Report",
						},
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "删除报表",
				"description": "删除指定报表",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "报表ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "删除成功",
					},
				},
			},
		},
		"/reports/{id}/polish": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "润色报表",
				"description": "对报表内容进行AI润色",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "报表ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "润色成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/Report",
						},
					},
				},
			},
		},
		"/reports/{id}/export/{format}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"REST-报表"},
				"summary":     "导出报表",
				"description": "以指定格式导出报表",
				"security":    []interface{}{map[string]interface{}{"BearerAuth": []string{}}},
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "报表ID",
					},
					map[string]interface{}{
						"name":        "format",
						"in":          "path",
						"required":    true,
						"type":        "string",
						"description": "导出格式",
						"enum":        []string{"pdf", "excel", "word"},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "导出成功",
						"schema": map[string]interface{}{
							"type":   "string",
							"format": "binary",
						},
					},
				},
			},
		},
	}

	// 添加 gRPC API 路径（模拟为 HTTP 接口）
	grpcPaths := map[string]interface{}{
		"/grpc/auth/register": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"gRPC-认证"},
				"summary":     "gRPC 用户注册",
				"description": "通过 gRPC 接口注册用户",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "gRPC 注册请求",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcRegisterRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "注册成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcRegisterResponse",
						},
					},
				},
			},
		},
		"/grpc/auth/login": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"gRPC-认证"},
				"summary":     "gRPC 用户登录",
				"description": "通过 gRPC 接口登录",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "gRPC 登录请求",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcLoginRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "登录成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcLoginResponse",
						},
					},
				},
			},
		},
		"/grpc/tasks/create": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"gRPC-任务"},
				"summary":     "gRPC 创建任务",
				"description": "通过 gRPC 接口创建任务",
				"parameters": []interface{}{
					map[string]interface{}{
						"name":        "body",
						"in":          "body",
						"required":    true,
						"description": "gRPC 创建任务请求",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcCreateTaskRequest",
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "创建成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcCreateTaskResponse",
						},
					},
				},
			},
		},
		"/grpc/tasks/list": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"gRPC-任务"},
				"summary":     "gRPC 获取任务列表",
				"description": "通过 gRPC 接口获取任务列表",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "任务列表",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcListTasksResponse",
						},
					},
				},
			},
		},
		"/grpc/captcha/generate": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"gRPC-验证码"},
				"summary":     "gRPC 生成验证码",
				"description": "通过 gRPC 接口生成验证码",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "验证码生成成功",
						"schema": map[string]interface{}{
							"$ref": "#/definitions/GrpcGenerateCaptchaResponse",
						},
					},
				},
			},
		},
	}

	// 合并路径
	for path, methods := range restPaths {
		combinedDoc["paths"].(map[string]interface{})[path] = methods
	}
	for path, methods := range grpcPaths {
		combinedDoc["paths"].(map[string]interface{})[path] = methods
	}

	// 添加数据模型定义
	definitions := map[string]interface{}{
		// 用户相关模型
		"User": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
				"email": map[string]interface{}{
					"type":        "string",
					"description": "邮箱地址",
					"example":     "user@example.com",
				},
				"nickname": map[string]interface{}{
					"type":        "string",
					"description": "用户昵称",
					"example":     "用户昵称",
				},
				"avatar": map[string]interface{}{
					"type":        "string",
					"description": "头像URL",
					"example":     "https://example.com/avatar.jpg",
				},
				"createdAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "创建时间",
					"example":     "2023-12-01T10:00:00Z",
				},
				"updatedAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "更新时间",
					"example":     "2023-12-01T10:00:00Z",
				},
			},
			"required": []string{"id", "email"},
		},
		// 认证相关模型
		"RegisterRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"email": map[string]interface{}{
					"type":        "string",
					"description": "邮箱地址",
					"example":     "user@example.com",
				},
				"password": map[string]interface{}{
					"type":        "string",
					"description": "密码",
					"minLength":   6,
					"example":     "password123",
				},
				"nickname": map[string]interface{}{
					"type":        "string",
					"description": "昵称",
					"example":     "用户昵称",
				},
				"emailCode": map[string]interface{}{
					"type":        "string",
					"description": "邮箱验证码",
					"example":     "123456",
				},
				"emailCodeId": map[string]interface{}{
					"type":        "string",
					"description": "邮箱验证码ID",
					"example":     "code-id-123",
				},
			},
			"required": []string{"email", "password", "emailCode", "emailCodeId"},
		},
		"RegisterResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user": map[string]interface{}{
					"$ref": "#/definitions/User",
				},
				"token": map[string]interface{}{
					"type":        "string",
					"description": "访问令牌",
					"example":     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
			},
		},
		"LoginRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"email": map[string]interface{}{
					"type":        "string",
					"description": "邮箱地址",
					"example":     "user@example.com",
				},
				"password": map[string]interface{}{
					"type":        "string",
					"description": "密码",
					"example":     "password123",
				},
				"emailCode": map[string]interface{}{
					"type":        "string",
					"description": "邮箱验证码（可选，用于邮箱验证码登录）",
					"example":     "123456",
				},
			},
			"required": []string{"email"},
		},
		"LoginResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user": map[string]interface{}{
					"$ref": "#/definitions/User",
				},
				"token": map[string]interface{}{
					"type":        "string",
					"description": "访问令牌",
					"example":     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
			},
		},
		// 任务相关模型
		"Task": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "任务ID",
					"example":     "60d5ecb74eb3b8001f8b4568",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "任务标题",
					"example":     "完成项目文档",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "任务描述",
					"example":     "编写API文档和用户手册",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "任务状态",
					"enum":        []string{"todo", "doing", "done"},
					"example":     "todo",
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": "优先级",
					"enum":        []string{"low", "medium", "high"},
					"example":     "medium",
				},
				"dueDate": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "截止时间",
					"example":     "2023-12-31T23:59:59Z",
				},
				"userId": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
				"createdAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "创建时间",
					"example":     "2023-12-01T10:00:00Z",
				},
				"updatedAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "更新时间",
					"example":     "2023-12-01T10:00:00Z",
				},
			},
			"required": []string{"id", "title", "status", "userId"},
		},
		"CreateTaskRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "任务标题",
					"example":     "完成项目文档",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "任务描述",
					"example":     "编写API文档和用户手册",
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": "优先级",
					"enum":        []string{"low", "medium", "high"},
					"default":     "medium",
					"example":     "medium",
				},
				"dueDate": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "截止时间",
					"example":     "2023-12-31T23:59:59Z",
				},
			},
			"required": []string{"title"},
		},
		"UpdateTaskRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "任务标题",
					"example":     "更新后的任务标题",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "任务描述",
					"example":     "更新后的任务描述",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "任务状态",
					"enum":        []string{"todo", "doing", "done"},
					"example":     "doing",
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": "优先级",
					"enum":        []string{"low", "medium", "high"},
					"example":     "high",
				},
				"dueDate": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "截止时间",
					"example":     "2023-12-31T23:59:59Z",
				},
			},
		},
		// 报表相关模型
		"Report": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "报表ID",
					"example":     "report-123",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "报表标题",
					"example":     "2023年12月工作报表",
				},
				"content": map[string]interface{}{
					"type":        "string",
					"description": "报表内容",
					"example":     "本月完成了10个任务...",
				},
				"type": map[string]interface{}{
					"type":        "string",
					"description": "报表类型",
					"enum":        []string{"daily", "weekly", "monthly"},
					"example":     "monthly",
				},
				"period": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"start": map[string]interface{}{
							"type":    "string",
							"format":  "date",
							"example": "2023-12-01",
						},
						"end": map[string]interface{}{
							"type":    "string",
							"format":  "date",
							"example": "2023-12-31",
						},
					},
					"description": "统计周期",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "报表状态",
					"enum":        []string{"generating", "completed", "failed"},
					"example":     "completed",
				},
				"userId": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
				"createdAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "创建时间",
					"example":     "2023-12-01T10:00:00Z",
				},
				"updatedAt": map[string]interface{}{
					"type":        "string",
					"format":      "date-time",
					"description": "更新时间",
					"example":     "2023-12-01T10:00:00Z",
				},
			},
			"required": []string{"id", "title", "type", "userId"},
		},
		// gRPC 相关模型
		"AuthLoginRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"email": map[string]interface{}{
					"type":        "string",
					"description": "邮箱地址",
					"example":     "user@example.com",
				},
				"password": map[string]interface{}{
					"type":        "string",
					"description": "密码",
					"example":     "password123",
				},
			},
			"required": []string{"email", "password"},
		},
		"AuthLoginResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"token": map[string]interface{}{
					"type":        "string",
					"description": "访问令牌",
					"example":     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				},
				"user": map[string]interface{}{
					"$ref": "#/definitions/User",
				},
			},
		},
		"AuthRegisterRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"email": map[string]interface{}{
					"type":        "string",
					"description": "邮箱地址",
					"example":     "user@example.com",
				},
				"password": map[string]interface{}{
					"type":        "string",
					"description": "密码",
					"example":     "password123",
				},
				"nickname": map[string]interface{}{
					"type":        "string",
					"description": "昵称",
					"example":     "用户昵称",
				},
			},
			"required": []string{"email", "password"},
		},
		"AuthRegisterResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user": map[string]interface{}{
					"$ref": "#/definitions/User",
				},
			},
		},
		"GetUserInfoRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
			},
			"required": []string{"id"},
		},
		"GetUserInfoResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user": map[string]interface{}{
					"$ref": "#/definitions/User",
				},
			},
		},
		"CreateTaskGrpcRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "任务标题",
					"example":     "完成项目文档",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "任务描述",
					"example":     "编写API文档和用户手册",
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": "优先级",
					"example":     "medium",
				},
				"dueDate": map[string]interface{}{
					"type":        "string",
					"description": "截止时间",
					"example":     "2023-12-31T23:59:59Z",
				},
				"userId": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
			},
			"required": []string{"title", "userId"},
		},
		"CreateTaskGrpcResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task": map[string]interface{}{
					"$ref": "#/definitions/Task",
				},
			},
		},
		"ListTasksRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"userId": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "任务状态过滤",
					"example":     "todo",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"example":     1,
				},
				"pageSize": map[string]interface{}{
					"type":        "integer",
					"description": "每页大小",
					"example":     20,
				},
			},
			"required": []string{"userId"},
		},
		"ListTasksResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"tasks": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/Task",
					},
					"description": "任务列表",
				},
				"total": map[string]interface{}{
					"type":        "integer",
					"description": "总数量",
					"example":     100,
				},
			},
		},
		"GenerateReportRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"userId": map[string]interface{}{
					"type":        "string",
					"description": "用户ID",
					"example":     "60d5ecb74eb3b8001f8b4567",
				},
				"type": map[string]interface{}{
					"type":        "string",
					"description": "报表类型",
					"example":     "monthly",
				},
				"startDate": map[string]interface{}{
					"type":        "string",
					"description": "开始日期",
					"example":     "2023-12-01",
				},
				"endDate": map[string]interface{}{
					"type":        "string",
					"description": "结束日期",
					"example":     "2023-12-31",
				},
			},
			"required": []string{"userId", "type"},
		},
		"GenerateReportResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"report": map[string]interface{}{
					"$ref": "#/definitions/Report",
				},
			},
		},
		"GenerateCaptchaRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"width": map[string]interface{}{
					"type":        "integer",
					"description": "验证码图片宽度",
					"example":     120,
				},
				"height": map[string]interface{}{
					"type":        "integer",
					"description": "验证码图片高度",
					"example":     40,
				},
			},
		},
		"GenerateCaptchaResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "验证码ID",
					"example":     "captcha-123",
				},
				"image": map[string]interface{}{
					"type":        "string",
					"description": "验证码图片 (Base64)",
					"example":     "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
				},
			},
		},
		"VerifyCaptchaRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "验证码ID",
					"example":     "captcha-123",
				},
				"code": map[string]interface{}{
					"type":        "string",
					"description": "验证码",
					"example":     "ABC123",
				},
			},
			"required": []string{"id", "code"},
		},
		"VerifyCaptchaResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"valid": map[string]interface{}{
					"type":        "boolean",
					"description": "验证结果",
					"example":     true,
				},
			},
		},
		"Error": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"code": map[string]interface{}{
					"type":        "integer",
					"description": "错误代码",
					"example":     400,
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "错误信息",
					"example":     "请求参数错误",
				},
			},
		},
	}

	combinedDoc["definitions"] = definitions

	// 生成 JSON 文件
	jsonData, err := json.MarshalIndent(combinedDoc, "", "  ")
	if err != nil {
		log.Fatal("生成 JSON 失败:", err)
	}

	if err := ioutil.WriteFile("docs/api_complete.json", jsonData, 0644); err != nil {
		log.Fatal("写入文件失败:", err)
	}

	fmt.Println("✅ 完整 API 文档生成成功!")
	fmt.Printf("📄 文件位置: docs/api_complete.json\n")
	fmt.Printf("📊 包含接口: %d 个 REST API + %d 个 gRPC API\n", len(restPaths), len(grpcPaths))
	fmt.Printf("📋 数据模型: %d 个\n", len(definitions))
	fmt.Printf("🔗 访问地址: http://localhost:5004/swagger/ (加载 api_complete.json)\n")
}
