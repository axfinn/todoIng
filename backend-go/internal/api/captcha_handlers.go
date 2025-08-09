package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/axfinn/todoIng/backend-go/internal/captcha"
	"github.com/gorilla/mux"
)

type CaptchaDeps struct{ Store *captcha.Store }

// Generate 生成验证码
// @Summary 生成验证码图片
// @Description 生成一个新的验证码图片，返回base64编码的SVG图片和验证码ID
// @Tags 验证码
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "验证码图片和ID"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/auth/captcha [get]
func (d *CaptchaDeps) Generate(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_CAPTCHA") != "true" {
		// 兼容前端：返回一个透明的 1x1 SVG 占位图片 + msg，前端拿到非错误结构即可继续
		placeholder := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(`<svg xmlns='http://www.w3.org/2000/svg' width='80' height='30'></svg>`))
		JSON(w, 200, map[string]string{"image": placeholder, "id": "disabled", "msg": "captcha disabled"})
		return
	}
	id, text := d.Store.Generate(6)
	// simple SVG
	svg := fmt.Sprintf(`<svg xmlns='http://www.w3.org/2000/svg' width='150' height='50'><rect width='100%%' height='100%%' fill='#f0f0f0'/>`)
	for i, c := range text {
		x := 15 + i*20
		y := 30
		color := fmt.Sprintf("#%06X", time.Now().UnixNano()%0xFFFFFF)
		svg += fmt.Sprintf("<text x='%d' y='%d' font-size='24' fill='%s'>%c</text>", x, y, color, c)
	}
	svg += "</svg>"
	JSON(w, 200, map[string]string{"image": "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg)), "id": id})
}

// Verify 验证验证码
// @Summary 验证验证码
// @Description 验证用户输入的验证码是否正确
// @Tags 验证码
// @Accept json
// @Produce json
// @Param body body object{captcha=string,captchaId=string} true "验证码内容和ID"
// @Success 200 {object} map[string]string "验证成功"
// @Failure 400 {object} map[string]string "验证失败或请求参数错误"
// @Router /api/auth/verify-captcha [post]
func (d *CaptchaDeps) Verify(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("ENABLE_CAPTCHA") != "true" {
		JSON(w, 200, map[string]string{"msg": "Captcha bypassed"})
		return
	}
	var body struct {
		Captcha   string `json:"captcha"`
		CaptchaId string `json:"captchaId"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSON(w, 400, map[string]string{"msg": "Invalid body"})
		return
	}
	if body.Captcha == "" || body.CaptchaId == "" {
		JSON(w, 400, map[string]string{"msg": "Captcha and CaptchaId required"})
		return
	}
	if d.Store.Verify(body.CaptchaId, body.Captcha) {
		JSON(w, 200, map[string]string{"msg": "Captcha verified successfully"})
	} else {
		JSON(w, 400, map[string]string{"msg": "Invalid or expired captcha"})
	}
}

func SetupCaptchaRoutes(r *mux.Router, deps *CaptchaDeps) {
	r.HandleFunc("/api/auth/captcha", deps.Generate).Methods(http.MethodGet)
	r.HandleFunc("/api/auth/verify-captcha", deps.Verify).Methods(http.MethodPost)
}
