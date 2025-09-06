package service

import (
	"errors"
	"sync"
	"time"

	"github.com/CycleZero/mc-yggdrasil-go/models"
	"github.com/CycleZero/mc-yggdrasil-go/utils"
)

// YggdrasilService 定义Yggdrasil服务的核心接口
// 将客户端层和逻辑处理层分开，便于其他项目按需调用

type YggdrasilService interface {
	// Auth 执行认证请求
	Auth(req models.AuthRequest) (*models.AuthResponse, error)
	
	// Refresh 刷新访问令牌
	Refresh(req models.RefreshRequest) (*models.AuthResponse, error)
	
	// Validate 验证访问令牌是否有效
	Validate(req models.ValidateRequest) (bool, error)
	
	// Invalidate 使访问令牌失效
	Invalidate(req models.InvalidateRequest) error
	
	// Signout 使用用户名和密码登出
	Signout(req models.SignoutRequest) error
}

// MemoryYggdrasilService 是YggdrasilService的内存实现
// 用于演示和测试，实际项目中可能需要持久化存储

type MemoryYggdrasilService struct {
	// 用户数据存储
	users map[string]UserCredentials // 用户名 -> 用户凭证
	
	// 令牌存储
	accessTokens map[string]AccessTokenInfo // 访问令牌 -> 令牌信息
	clientTokens map[string]string // 客户端令牌 -> 访问令牌
	
	// 角色存储
	profiles map[string]*models.Profile // 用户ID -> 角色
	
	// 锁，用于并发控制
	mu sync.RWMutex
}

// UserCredentials 表示用户凭证

type UserCredentials struct {
	ID       string
	Password string
}

// AccessTokenInfo 表示访问令牌信息

type AccessTokenInfo struct {
	UserID      string
	ClientToken string
	ProfileID   string
	CreatedAt   time.Time
}

// NewMemoryYggdrasilService 创建一个新的内存实现的Yggdrasil服务
func NewMemoryYggdrasilService() *MemoryYggdrasilService {
	return &MemoryYggdrasilService{
		users:        make(map[string]UserCredentials),
		accessTokens: make(map[string]AccessTokenInfo),
		clientTokens: make(map[string]string),
		profiles:     make(map[string]*models.Profile),
	}
}

// Auth 实现认证请求
func (s *MemoryYggdrasilService) Auth(req models.AuthRequest) (*models.AuthResponse, error) {
	s.mu.RLock()
	userCreds, exists := s.users[req.Username]
	s.mu.RUnlock()
	
	// 检查用户是否存在且密码正确
	if !exists || userCreds.Password != req.Password {
		return nil, errors.New("Invalid credentials. Invalid username or password.")
	}
	
	s.mu.RLock()
	profile, profileExists := s.profiles[userCreds.ID]
	s.mu.RUnlock()
	
	// 检查角色是否存在
	if !profileExists {
		return nil, errors.New("No profile found for user")
	}
	
	// 生成访问令牌和客户端令牌
	accessToken := utils.GenerateUUID()
	clientToken := req.ClientToken
	if clientToken == "" {
		clientToken = utils.GenerateUUID()
	}
	
	// 存储令牌信息
	s.mu.Lock()
	s.accessTokens[accessToken] = AccessTokenInfo{
		UserID:      userCreds.ID,
		ClientToken: clientToken,
		ProfileID:   profile.ID,
		CreatedAt:   time.Now(),
	}
	s.clientTokens[clientToken] = accessToken
	s.mu.Unlock()
	
	// 构建响应
	resp := &models.AuthResponse{
		AccessToken:  accessToken,
		ClientToken:  clientToken,
		SelectedProfile: profile,
	}
	
	// 如果请求了用户信息，添加用户信息
	if req.RequestUser {
		resp.User = &models.User{
			ID: userCreds.ID,
			Properties: []models.Property{
				{
					Name:  "preferredLanguage",
					Value: "en",
				},
			},
		}
	}
	
	return resp, nil
}

// Refresh 实现刷新访问令牌
func (s *MemoryYggdrasilService) Refresh(req models.RefreshRequest) (*models.AuthResponse, error) {
	s.mu.RLock()
	tokenInfo, exists := s.accessTokens[req.AccessToken]
	s.mu.RUnlock()
	
	// 检查令牌是否存在
	if !exists {
		return nil, errors.New("Invalid token.")
	}
	
	// 检查客户端令牌是否匹配
	if req.ClientToken != "" && req.ClientToken != tokenInfo.ClientToken {
		return nil, errors.New("Invalid token.")
	}
	
	s.mu.RLock()
	profile, profileExists := s.profiles[tokenInfo.UserID]
	s.mu.RUnlock()
	
	// 检查角色是否存在
	if !profileExists {
		return nil, errors.New("No profile found for user")
	}
	
	// 生成新的访问令牌
	newAccessToken := utils.GenerateUUID()
	
	// 更新令牌信息
	s.mu.Lock()
	// 删除旧的访问令牌
	delete(s.accessTokens, req.AccessToken)
	// 删除旧的客户端令牌映射
	if oldAccessToken, exists := s.clientTokens[tokenInfo.ClientToken]; exists && oldAccessToken == req.AccessToken {
		delete(s.clientTokens, tokenInfo.ClientToken)
	}
	// 存储新的访问令牌
	s.accessTokens[newAccessToken] = AccessTokenInfo{
		UserID:      tokenInfo.UserID,
		ClientToken: tokenInfo.ClientToken,
		ProfileID:   tokenInfo.ProfileID,
		CreatedAt:   time.Now(),
	}
	s.clientTokens[tokenInfo.ClientToken] = newAccessToken
	s.mu.Unlock()
	
	// 构建响应
	resp := &models.AuthResponse{
		AccessToken:     newAccessToken,
		ClientToken:     tokenInfo.ClientToken,
		SelectedProfile: profile,
	}
	
	// 如果请求了用户信息，添加用户信息
	if req.RequestUser {
		resp.User = &models.User{
			ID: tokenInfo.UserID,
			Properties: []models.Property{
				{
					Name:  "preferredLanguage",
					Value: "en",
				},
			},
		}
	}
	
	return resp, nil
}

// Validate 实现验证访问令牌
func (s *MemoryYggdrasilService) Validate(req models.ValidateRequest) (bool, error) {
	s.mu.RLock()
	tokenInfo, exists := s.accessTokens[req.AccessToken]
	s.mu.RUnlock()
	
	// 检查令牌是否存在
	if !exists {
		return false, nil
	}
	
	// 检查客户端令牌是否匹配
	if req.ClientToken != "" && req.ClientToken != tokenInfo.ClientToken {
		return false, nil
	}
	
	return true, nil
}

// Invalidate 实现使访问令牌失效
func (s *MemoryYggdrasilService) Invalidate(req models.InvalidateRequest) error {
	s.mu.RLock()
	tokenInfo, exists := s.accessTokens[req.AccessToken]
	s.mu.RUnlock()
	
	// 检查令牌是否存在
	if !exists {
		return errors.New("Invalid token.")
	}
	
	// 检查客户端令牌是否匹配
	if req.ClientToken != tokenInfo.ClientToken {
		return errors.New("Invalid token.")
	}
	
	// 删除令牌
	s.mu.Lock()
	delete(s.accessTokens, req.AccessToken)
	if oldAccessToken, exists := s.clientTokens[tokenInfo.ClientToken]; exists && oldAccessToken == req.AccessToken {
		delete(s.clientTokens, tokenInfo.ClientToken)
	}
	s.mu.Unlock()
	
	return nil
}

// Signout 实现使用用户名和密码登出
func (s *MemoryYggdrasilService) Signout(req models.SignoutRequest) error {
	s.mu.RLock()
	userCreds, exists := s.users[req.Username]
	s.mu.RUnlock()
	
	// 检查用户是否存在且密码正确
	if !exists || userCreds.Password != req.Password {
		return errors.New("Invalid credentials. Invalid username or password.")
	}
	
	// 找出并删除该用户的所有访问令牌
	s.mu.Lock()
	for accessToken, tokenInfo := range s.accessTokens {
		if tokenInfo.UserID == userCreds.ID {
			delete(s.accessTokens, accessToken)
			if oldAccessToken, exists := s.clientTokens[tokenInfo.ClientToken]; exists && oldAccessToken == accessToken {
				delete(s.clientTokens, tokenInfo.ClientToken)
			}
		}
	}
	s.mu.Unlock()
	
	return nil
}

// 添加方法用于管理用户和角色（仅用于演示和测试）

// AddUser 添加一个用户
func (s *MemoryYggdrasilService) AddUser(username, password string) (string, error) {
	userID := utils.GenerateUUID()
	
	s.mu.Lock()
	s.users[username] = UserCredentials{
		ID:       userID,
		Password: password,
	}
	s.mu.Unlock()
	
	return userID, nil
}

// AddProfile 添加一个角色
func (s *MemoryYggdrasilService) AddProfile(userID, name string) (*models.Profile, error) {
	// 生成与离线验证系统兼容的UUID
	profileID, err := utils.GenerateOfflinePlayerUUID(name)
	if err != nil {
		return nil, err
	}
	
	profile := &models.Profile{
		ID:   profileID,
		Name: name,
	}
	
	s.mu.Lock()
	s.profiles[userID] = profile
	s.mu.Unlock()
	
	return profile, nil
}