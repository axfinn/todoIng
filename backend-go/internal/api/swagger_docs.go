package api

// 这个文件专门用于聚合所有swagger注释，确保swag能够扫描到所有接口

// TaskHandler接口文档聚合

// CreateTask 创建新任务
// @Summary 创建新任务
// @Description 创建一个新的任务项，支持设置标题、描述、状态、优先级等
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param task body object true "任务信息" example({"title":"新任务","description":"任务描述","status":"To Do","priority":"Medium"})
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks [post]

// ListTasks 获取任务列表
// @Summary 获取用户的所有任务
// @Description 获取当前用户创建的所有任务列表，按创建时间倒序排列
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {array} map[string]interface{} "任务列表"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks [get]

// GetTask 获取单个任务详情
// @Summary 获取任务详情
// @Description 根据任务ID获取任务的详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]interface{} "任务详情"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [get]

// UpdateTask 更新任务
// @Summary 更新任务信息
// @Description 根据任务ID更新任务的详细信息
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "任务ID"
// @Param task body object true "更新的任务信息" example({"title":"更新任务","status":"In Progress"})
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [put]

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 根据任务ID删除指定的任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/{id} [delete]

// ExportTasks 导出任务
// @Summary 导出所有任务
// @Description 将用户的所有任务导出为JSON格式
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {array} map[string]interface{} "导出的任务数据"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/export [get]

// ImportTasks 导入任务
// @Summary 导入任务数据
// @Description 从JSON数据批量导入任务
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param tasks body array true "要导入的任务数组"
// @Success 200 {object} map[string]interface{} "导入成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/tasks/import [post]

// ReportHandler接口文档聚合

// ListReports 获取报表列表
// @Summary 获取用户的所有报表
// @Description 获取当前用户创建的所有报表列表，按创建时间倒序排列
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {array} map[string]interface{} "报表列表"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports [get]

// GetReport 获取报表详情
// @Summary 获取单个报表详情
// @Description 根据报表ID获取报表的详细信息，包含关联的任务数据
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "报表ID"
// @Success 200 {object} map[string]interface{} "报表详情"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "报表不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports/{id} [get]

// CreateReport 创建报表
// @Summary 创建新报表
// @Description 根据指定的任务创建新的报表
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param report body object true "报表信息" example({"title":"周报","tasks":["task_id_1","task_id_2"]})
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports [post]

// UpdateReport 更新报表
// @Summary 更新报表信息
// @Description 根据报表ID更新报表的详细信息
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "报表ID"
// @Param report body object true "更新的报表信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "报表不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports/{id} [put]

// DeleteReport 删除报表
// @Summary 删除报表
// @Description 根据报表ID删除指定的报表
// @Tags 报表管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Param id path string true "报表ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 401 {object} map[string]string "未授权"
// @Failure 404 {object} map[string]string "报表不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/reports/{id} [delete]

// CaptchaHandler接口文档聚合

// GenerateCaptcha 生成验证码
// @Summary 生成验证码图片
// @Description 生成一个新的验证码图片，返回base64编码的SVG图片和验证码ID
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "验证码图片和ID" example({"image":"data:image/svg+xml;base64,...","id":"captcha_id"})
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/auth/captcha [get]

// VerifyCaptcha 验证验证码
// @Summary 验证验证码
// @Description 验证用户输入的验证码是否正确
// @Tags 验证码
// @Accept json
// @Produce json
// @Param body body object true "验证码内容和ID" example({"captcha":"ABCD","captchaId":"captcha_id"})
// @Success 200 {object} map[string]string "验证成功"
// @Failure 400 {object} map[string]string "验证失败或请求参数错误"
// @Router /api/auth/verify-captcha [post]

// 假的函数，仅用于让swag扫描到上述注释
func swaggerDocumentationAggregator() {}
