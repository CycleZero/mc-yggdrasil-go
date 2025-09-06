package models

// ErrorResponse 表示API返回的错误信息
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E9%94%99%E8%AF%AF%E4%BF%A1%E6%81%AF%E6%A0%BC%E5%BC%8F

type ErrorResponse struct {
	Error        string `json:"error,omitempty"`        // 错误的简要描述（机器可读）
	ErrorMessage string `json:"errorMessage,omitempty"` // 错误的详细信息（人类可读）
	Cause        string `json:"cause,omitempty"`        // 该错误的原因（可选）
}

// Property 表示用户或角色的属性
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E7%94%A8%E6%88%B7

type Property struct {
	Name      string `json:"name"`                // 属性的名称
	Value     string `json:"value"`               // 属性的值
	Signature string `json:"signature,omitempty"` // 属性值的数字签名（仅在特定情况下需要包含）
}

// User 表示Yggdrasil系统中的用户
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E7%94%A8%E6%88%B7

type User struct {
	ID         string     `json:"id"`         // 用户的ID（无符号UUID）
	Properties []Property `json:"properties"` // 用户的属性
}

// Profile 表示Minecraft中的角色
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E8%A7%92%E8%89%B2profile

type Profile struct {
	ID         string     `json:"id"`                   // 角色UUID（无符号）
	Name       string     `json:"name"`                 // 角色名称
	Properties []Property `json:"properties,omitempty"` // 角色的属性（仅在特定情况下需要包含）
}

// TextureType 表示材质类型

type TextureType string

const (
	TextureSkin TextureType = "SKIN"
	TextureCape TextureType = "CAPE"
)

// TextureModel 表示材质模型
// https://github.com/yushijinhun/authlib-injector/wiki/Yggdrasil-%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8A%80%E6%9C%AF%E8%A7%84%E8%8C%83#%E8%A7%92%E8%89%B2profile

type TextureModel string

const (
	TextureModelDefault TextureModel = "default" // 正常手臂宽度（4px）的皮肤
	TextureModelSlim    TextureModel = "slim"    // 细手臂（3px）的皮肤
)

// AuthRequest 表示认证请求

type AuthRequest struct {
	Agent       Agent  `json:"agent"`                 // 客户端代理信息
	Username    string `json:"username"`              // 用户名（邮箱）
	Password    string `json:"password"`              // 密码
	ClientToken string `json:"clientToken,omitempty"` // 客户端令牌（可选）
	RequestUser bool   `json:"requestUser,omitempty"` // 是否请求用户信息（可选）
}

// Agent 表示客户端代理信息

type Agent struct {
	Name    string `json:"name"`    // 代理名称，如"Minecraft"
	Version int    `json:"version"` // 代理版本
}

// AuthResponse 表示认证响应

type AuthResponse struct {
	AccessToken       string    `json:"accessToken"`                 // 访问令牌
	ClientToken       string    `json:"clientToken"`                 // 客户端令牌
	AvailableProfiles []Profile `json:"availableProfiles,omitempty"` // 可用的角色列表（可选）
	SelectedProfile   *Profile  `json:"selectedProfile,omitempty"`   // 选定的角色（可选）
	User              *User     `json:"user,omitempty"`              // 用户信息（可选，只有当requestUser为true时才返回）
}

// RefreshRequest 表示刷新令牌请求

type RefreshRequest struct {
	AccessToken     string   `json:"accessToken"`               // 当前的访问令牌
	ClientToken     string   `json:"clientToken"`               // 客户端令牌
	RequestUser     bool     `json:"requestUser,omitempty"`     // 是否请求用户信息（可选）
	SelectedProfile *Profile `json:"selectedProfile,omitempty"` // 选定的角色（可选）
}

// ValidateRequest 表示验证令牌请求

type ValidateRequest struct {
	AccessToken string `json:"accessToken"`           // 要验证的访问令牌
	ClientToken string `json:"clientToken,omitempty"` // 客户端令牌（可选）
}

// InvalidateRequest 表示使令牌失效请求（登出）

type InvalidateRequest struct {
	AccessToken string `json:"accessToken"` // 要失效的访问令牌
	ClientToken string `json:"clientToken"` // 客户端令牌
}

// SignoutRequest 表示登出请求（使用用户名和密码）

type SignoutRequest struct {
	Username string `json:"username"` // 用户名（邮箱）
	Password string `json:"password"` // 密码
}
