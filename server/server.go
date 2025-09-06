package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CycleZero/mc-yggdrasil-go/models"
	"github.com/CycleZero/mc-yggdrasil-go/service"
)

// YggdrasilServer 表示Yggdrasil认证服务器

type YggdrasilServer struct {
	Port    int
	Service service.YggdrasilService
	server  *http.Server
}

// NewYggdrasilServer 创建一个新的Yggdrasil服务器
func NewYggdrasilServer(port int, service service.YggdrasilService) *YggdrasilServer {
	return &YggdrasilServer{
		Port:    port,
		Service: service,
	}
}

// Start 启动Yggdrasil服务器
func (s *YggdrasilServer) Start() error {
	// 注册路由
	r := http.NewServeMux()
	r.HandleFunc("/authserver/authenticate", s.handleAuthenticate)
	r.HandleFunc("/authserver/refresh", s.handleRefresh)
	r.HandleFunc("/authserver/validate", s.handleValidate)
	r.HandleFunc("/authserver/invalidate", s.handleInvalidate)
	r.HandleFunc("/authserver/signout", s.handleSignout)
	r.HandleFunc("/", s.handleRoot)

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: r,
	}

	// 启动服务器
	log.Printf("Yggdrasil服务器启动在端口 %d...\n", s.Port)
	return s.server.ListenAndServe()
}

// Stop 停止Yggdrasil服务器（使用优雅关闭）
func (s *YggdrasilServer) Stop(ctx context.Context) error {
	if s.server != nil {
		log.Println("正在停止Yggdrasil服务器...")
		return s.server.Shutdown(ctx)
	}
	return nil
}

// handleAuthenticate 处理认证请求
// POST /authserver/authenticate
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E8%AE%A4%E8%AF%81
func (s *YggdrasilServer) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "IllegalArgumentException", err.Error())
		return
	}
	defer r.Body.Close()

	// 调用服务处理认证
	resp, err := s.Service.Auth(req)
	if err != nil {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", err.Error())
		return
	}

	// 写入响应
	s.writeJSONResponse(w, http.StatusOK, resp)
}

// handleRefresh 处理刷新令牌请求
// POST /authserver/refresh
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E5%88%B7%E6%96%B0
func (s *YggdrasilServer) handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "IllegalArgumentException", err.Error())
		return
	}
	defer r.Body.Close()

	// 调用服务处理刷新
	resp, err := s.Service.Refresh(req)
	if err != nil {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", err.Error())
		return
	}

	// 写入响应
	s.writeJSONResponse(w, http.StatusOK, resp)
}

// handleValidate 处理验证令牌请求
// POST /authserver/validate
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E9%AA%8C%E8%AF%81
func (s *YggdrasilServer) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "IllegalArgumentException", err.Error())
		return
	}
	defer r.Body.Close()

	// 调用服务处理验证
	valid, err := s.Service.Validate(req)
	if err != nil {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", err.Error())
		return
	}

	// 根据技术规范，验证成功返回204 No Content
	if valid {
		w.WriteHeader(http.StatusNoContent)
	} else {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", "Invalid token.")
	}
}

// handleInvalidate 处理使令牌失效请求
// POST /authserver/invalidate
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E5%A4%B1%E6%95%88
func (s *YggdrasilServer) handleInvalidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.InvalidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "IllegalArgumentException", err.Error())
		return
	}
	defer r.Body.Close()

	// 调用服务处理失效
	err := s.Service.Invalidate(req)
	if err != nil {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", err.Error())
		return
	}

	// 根据技术规范，成功返回204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// handleSignout 处理登出请求
// POST /authserver/signout
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E7%99%BB%E5%87%BA
func (s *YggdrasilServer) handleSignout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var req models.SignoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, "IllegalArgumentException", err.Error())
		return
	}
	defer r.Body.Close()

	// 调用服务处理登出
	err := s.Service.Signout(req)
	if err != nil {
		s.writeErrorResponse(w, http.StatusForbidden, "ForbiddenOperationException", err.Error())
		return
	}

	// 根据技术规范，成功返回204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// handleRoot 处理根路径请求
func (s *YggdrasilServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Yggdrasil认证服务器正在运行",
	})
}

// writeJSONResponse 写入JSON响应
func (s *YggdrasilServer) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("写入响应失败: %v\n", err)
	}
}

// writeErrorResponse 写入错误响应
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E9%94%99%E8%AF%AF%E4%BF%A1%E6%81%AF%E6%A0%BC%E5%BC%8F
func (s *YggdrasilServer) writeErrorResponse(w http.ResponseWriter, statusCode int, errorType, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	errResp := models.ErrorResponse{
		Error:        errorType,
		ErrorMessage: errorMessage,
	}
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		log.Printf("写入错误响应失败: %v\n", err)
	}
}
