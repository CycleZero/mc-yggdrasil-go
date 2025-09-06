package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CycleZero/mc-yggdrasil-go/models"
	"github.com/CycleZero/mc-yggdrasil-go/service"
)

// YggdrasilClient 表示与Yggdrasil服务器通信的客户端
// 客户端层负责与服务器通信，不处理业务逻辑

type YggdrasilClient struct {
	BaseURL    string       // Yggdrasil服务器的基础URL
	HTTPClient *http.Client // HTTP客户端
	// 可选：本地服务实现，用于测试或离线模式
	LocalService service.YggdrasilService
}

// NewYggdrasilClient 创建一个新的Yggdrasil客户端
func NewYggdrasilClient(baseURL string) *YggdrasilClient {
	return &YggdrasilClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// NewLocalClient 创建一个使用本地服务实现的客户端
// 用于测试或离线模式，不依赖外部服务器
func NewLocalClient(service service.YggdrasilService) *YggdrasilClient {
	return &YggdrasilClient{
		LocalService: service,
	}
}

// Auth 执行认证请求
// 如果配置了本地服务，则使用本地服务；否则发送HTTP请求
func (c *YggdrasilClient) Auth(req models.AuthRequest) (*models.AuthResponse, error) {
	if c.LocalService != nil {
		return c.LocalService.Auth(req)
	}

	url := c.BaseURL + "/authserver/authenticate"
	resp, err := c.doPostRequest(url, req)
	if err != nil {
		return nil, err
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

// Refresh 刷新访问令牌
// 如果配置了本地服务，则使用本地服务；否则发送HTTP请求
func (c *YggdrasilClient) Refresh(req models.RefreshRequest) (*models.AuthResponse, error) {
	if c.LocalService != nil {
		return c.LocalService.Refresh(req)
	}

	url := c.BaseURL + "/authserver/refresh"
	resp, err := c.doPostRequest(url, req)
	if err != nil {
		return nil, err
	}

	var refreshResp models.AuthResponse
	if err := json.Unmarshal(resp, &refreshResp); err != nil {
		return nil, err
	}

	return &refreshResp, nil
}

// Validate 验证访问令牌是否有效
// 如果配置了本地服务，则使用本地服务；否则发送HTTP请求
func (c *YggdrasilClient) Validate(req models.ValidateRequest) (bool, error) {
	if c.LocalService != nil {
		return c.LocalService.Validate(req)
	}

	url := c.BaseURL + "/authserver/validate"
	respData, err := c.doPostRequest(url, req)

	// 成功验证时返回204 No Content，失败时返回错误
	if err != nil {
		var apiErr models.ErrorResponse
		if jsonErr := json.Unmarshal(respData, &apiErr); jsonErr == nil {
			return false, errors.New(apiErr.ErrorMessage)
		}
		return false, err
	}

	return true, nil
}

// Invalidate 使访问令牌失效（登出）
// 如果配置了本地服务，则使用本地服务；否则发送HTTP请求
func (c *YggdrasilClient) Invalidate(req models.InvalidateRequest) error {
	if c.LocalService != nil {
		return c.LocalService.Invalidate(req)
	}

	url := c.BaseURL + "/authserver/invalidate"
	_, err := c.doPostRequest(url, req)
	return err
}

// Signout 使用用户名和密码登出
// 如果配置了本地服务，则使用本地服务；否则发送HTTP请求
func (c *YggdrasilClient) Signout(req models.SignoutRequest) error {
	if c.LocalService != nil {
		return c.LocalService.Signout(req)
	}

	url := c.BaseURL + "/authserver/signout"
	_, err := c.doPostRequest(url, req)
	return err
}

// doPostRequest 执行HTTP POST请求并返回响应内容
func (c *YggdrasilClient) doPostRequest(url string, body interface{}) ([]byte, error) {
	// 序列化请求体
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 发送请求
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody := make([]byte, 1024*1024) // 限制1MB
	n, err := resp.Body.Read(respBody)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	respBody = respBody[:n]

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 尝试解析错误响应
		var apiErr models.ErrorResponse
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil {
			return respBody, errors.New(apiErr.ErrorMessage)
		}
		return respBody, errors.New("request failed with status code: " + resp.Status)
	}

	return respBody, nil
}
