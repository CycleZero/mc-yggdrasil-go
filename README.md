# mc-yggdrasil-go

一个基于Go语言实现的Yggdrasil依赖库，用于Minecraft身份验证和会话管理。该库提供了完整的Yggdrasil API客户端实现、服务层实现以及服务器实现，支持完整的身份验证流程。

## 功能特性

- **完整的Yggdrasil API实现**
  - 身份验证 (Authenticate)
  - 刷新令牌 (Refresh)
  - 验证令牌 (Validate)
  - 使令牌失效 (Invalidate)
  - 登出 (Signout)
- **分层架构设计**
  - 客户端层：处理HTTP请求和响应
  - 服务层：实现核心业务逻辑
  - 数据模型层：定义数据结构
  - 工具层：提供UUID生成等工具函数
- **可独立运行的服务器**：符合Yggdrasil技术规范的HTTP服务器实现
- **内存存储实现**：支持开发和测试环境
- **UUID生成与处理工具**
  - 生成离线玩家UUID
  - 格式化UUID
  - 解析无连字符UUID
  - 生成随机UUID
- **简洁易用的API接口**
- **符合Yggdrasil技术规范的请求和响应格式**

## 项目结构

```
mc-yggdrasil-go/
├── client/        # Yggdrasil客户端实现
├── service/       # Yggdrasil服务层实现
├── server/        # Yggdrasil服务器实现
├── models/        # 数据模型定义
├── utils/         # 工具函数
├── main.go        # 使用示例
├── README.md      # 项目文档
└── go.mod         # Go模块定义
```

## 安装

将该库添加到您的Go项目中：

```bash
go get -u github.com/yourusername/mc-yggdrasil-go
```

## 使用示例

### 导入包

```go
import (
	"mc-yggdrasil-go/client"
	"mc-yggdrasil-go/models"
	"mc-yggdrasil-go/service"
	"mc-yggdrasil-go/server"
	"mc-yggdrasil-go/utils"
)
```

### 1. UUID工具使用示例

```go
// 生成离线玩家UUID
playerName := "Steve"
offlineUUID, err := utils.GenerateOfflinePlayerUUID(playerName)
if err != nil {
	fmt.Printf("生成离线UUID失败: %v\n", err)
} else {
	fmt.Printf("玩家 '%s' 的离线UUID: %s\n", playerName, offlineUUID)
}

// 生成随机UUID
randomUUID := utils.GenerateUUID()
fmt.Printf("随机生成的UUID: %s\n", randomUUID)

// 格式化UUID	uuid, err := utils.FormatUUID(randomUUID)
if err != nil {
	fmt.Printf("格式化UUID失败: %v\n", err)
} else {
	fmt.Printf("格式化后的UUID: %s\n", uuid)
}
```

### 2. 本地服务使用示例（不依赖HTTP）

```go
// 创建内存实现的服务
memoryService := service.NewMemoryYggdrasilService()

// 添加测试用户
userID, err := memoryService.AddUser("test@example.com", "password123")
if err != nil {
	fmt.Printf("添加用户失败: %v\n", err)
}

// 添加测试角色
profile, err := memoryService.AddProfile(userID, "TestPlayer")
if err != nil {
	fmt.Printf("添加角色失败: %v\n", err)
}

// 创建使用本地服务的客户端
localClient := client.NewLocalClient(memoryService)

// 执行认证	authReq := models.AuthRequest{
	Agent: models.Agent{
		Name:    "Minecraft",
		Version: 1,
	},
	Username:   "test@example.com",
	Password:   "password123",
	RequestUser: true,
}
	authResp, err := localClient.Auth(authReq)
if err != nil {
	fmt.Printf("认证失败: %v\n", err)
} else {
	fmt.Printf("认证成功! 访问令牌: %s\n", authResp.AccessToken)
}
```

### 3. HTTP客户端使用示例

```go
// 创建连接到远程Yggdrasil服务器的客户端	yggClient := client.NewYggdrasilClient("https://authserver.mojang.com")

// 准备认证请求	authReq := models.AuthRequest{
	Agent: models.Agent{
		Name:    "Minecraft",
		Version: 1,
	},
	Username:   "your-email@example.com",
	Password:   "your-password",
	RequestUser: true,
}

// 执行认证	authResp, err := yggClient.Auth(authReq)
if err != nil {
	fmt.Printf("认证失败: %v\n", err)
} else {
	fmt.Printf("认证成功! 访问令牌: %s\n", authResp.AccessToken)
}
```

### 4. 启动Yggdrasil服务器示例

```go
// 创建内存实现的服务
memoryService := service.NewMemoryYggdrasilService()

// 添加测试用户和角色
userID, _ := memoryService.AddUser("server@example.com", "server123")
memoryService.AddProfile(userID, "ServerPlayer")

// 创建并启动服务器
server := server.NewYggdrasilServer(8080, memoryService)

// 使用goroutine启动服务器
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	if err := server.Start(); err != nil {
		log.Fatalf("启动服务器失败: %v\n", err)
	}
}()

// ... 服务器运行中 ...

// 优雅关闭服务器
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
server.Stop(ctx)
wg.Wait()
```

## API参考

### 客户端层 (client)

#### NewYggdrasilClient(baseURL string) *YggdrasilClient
创建一个连接到远程Yggdrasil服务器的客户端实例。

- **参数**:
  - baseURL: Yggdrasil服务器的基础URL
- **返回值**:
  - *YggdrasilClient: Yggdrasil客户端实例

#### NewLocalClient(service service.YggdrasilService) *YggdrasilClient
创建一个使用本地服务实现的客户端实例（不依赖HTTP）。

- **参数**:
  - service: Yggdrasil服务层实现
- **返回值**:
  - *YggdrasilClient: Yggdrasil客户端实例

#### (c *YggdrasilClient) Auth(req models.AuthRequest) (*models.AuthResponse, error)
执行身份验证请求。

- **参数**:
  - req: 认证请求对象
- **返回值**:
  - *models.AuthResponse: 认证响应对象
  - error: 错误信息

#### (c *YggdrasilClient) Refresh(req models.RefreshRequest) (*models.RefreshResponse, error)
刷新访问令牌。

- **参数**:
  - req: 刷新请求对象
- **返回值**:
  - *models.RefreshResponse: 刷新响应对象
  - error: 错误信息

#### (c *YggdrasilClient) Validate(req models.ValidateRequest) (bool, error)
验证访问令牌是否有效。

- **参数**:
  - req: 验证请求对象
- **返回值**:
  - bool: 令牌是否有效
  - error: 错误信息

#### (c *YggdrasilClient) Invalidate(req models.InvalidateRequest) error
使访问令牌失效。

- **参数**:
  - req: 失效请求对象
- **返回值**:
  - error: 错误信息

#### (c *YggdrasilClient) Signout(req models.SignoutRequest) error
使用用户名和密码登出所有会话。

- **参数**:
  - req: 登出请求对象
- **返回值**:
  - error: 错误信息

### 服务层 (service)

#### NewMemoryYggdrasilService() *MemoryYggdrasilService
创建一个基于内存实现的Yggdrasil服务。

- **返回值**:
  - *MemoryYggdrasilService: 内存实现的Yggdrasil服务实例

#### (s *MemoryYggdrasilService) Auth(req models.AuthRequest) (*models.AuthResponse, error)
执行身份验证逻辑。

- **参数**:
  - req: 认证请求对象
- **返回值**:
  - *models.AuthResponse: 认证响应对象
  - error: 错误信息

#### (s *MemoryYggdrasilService) Refresh(req models.RefreshRequest) (*models.RefreshResponse, error)
执行令牌刷新逻辑。

- **参数**:
  - req: 刷新请求对象
- **返回值**:
  - *models.RefreshResponse: 刷新响应对象
  - error: 错误信息

#### (s *MemoryYggdrasilService) Validate(req models.ValidateRequest) (bool, error)
执行令牌验证逻辑。

- **参数**:
  - req: 验证请求对象
- **返回值**:
  - bool: 令牌是否有效
  - error: 错误信息

#### (s *MemoryYggdrasilService) Invalidate(req models.InvalidateRequest) error
执行令牌失效逻辑。

- **参数**:
  - req: 失效请求对象
- **返回值**:
  - error: 错误信息

#### (s *MemoryYggdrasilService) Signout(req models.SignoutRequest) error
执行用户登出逻辑。

- **参数**:
  - req: 登出请求对象
- **返回值**:
  - error: 错误信息

#### (s *MemoryYggdrasilService) AddUser(email, password string) (string, error)
添加新用户到内存存储中。

- **参数**:
  - email: 用户邮箱
  - password: 用户密码
- **返回值**:
  - string: 用户ID
  - error: 错误信息

#### (s *MemoryYggdrasilService) AddProfile(userID, name string) (*models.Profile, error)
为用户添加新角色。

- **参数**:
  - userID: 用户ID
  - name: 角色名称
- **返回值**:
  - *models.Profile: 角色信息
  - error: 错误信息

### 服务器层 (server)

#### NewYggdrasilServer(port int, service service.YggdrasilService) *YggdrasilServer
创建一个新的Yggdrasil HTTP服务器。

- **参数**:
  - port: 服务器监听端口
  - service: Yggdrasil服务层实现
- **返回值**:
  - *YggdrasilServer: Yggdrasil服务器实例

#### (s *YggdrasilServer) Start() error
启动Yggdrasil服务器。

- **返回值**:
  - error: 错误信息（如果启动失败）

#### (s *YggdrasilServer) Stop(ctx context.Context) error
优雅关闭Yggdrasil服务器。

- **参数**:
  - ctx: 上下文，用于控制关闭超时
- **返回值**:
  - error: 错误信息（如果关闭失败）

### 工具函数 (utils)

#### GenerateOfflinePlayerUUID(username string) (string, error)
生成与Minecraft离线验证系统兼容的UUID。

- **参数**:
  - username: 玩家用户名
- **返回值**:
  - string: 生成的UUID字符串
  - error: 错误信息

#### FormatUUID(uuid string) (string, error)
格式化UUID为标准格式（带连字符）。

- **参数**:
  - uuid: 无连字符的UUID字符串
- **返回值**:
  - string: 格式化后的UUID字符串
  - error: 错误信息

#### ParseUndashedUUID(uuid string) (string, error)
解析无连字符的UUID字符串。

- **参数**:
  - uuid: 无连字符的UUID字符串
- **返回值**:
  - string: 解析后的UUID字符串
  - error: 错误信息

#### GenerateUUID() string
生成一个随机的UUID字符串。

- **返回值**:
  - string: 生成的UUID字符串

#### ValidateUndashedUUID(uuid string) bool
验证无连字符的UUID字符串是否有效。

- **参数**:
  - uuid: 无连字符的UUID字符串
- **返回值**:
  - bool: UUID是否有效

## 架构说明

该项目采用分层架构设计，具体如下：

1. **数据模型层 (models)**：定义所有请求和响应的数据结构，确保与Yggdrasil技术规范一致。

2. **工具层 (utils)**：提供通用功能，如UUID生成和处理，不依赖于其他层。

3. **服务层 (service)**：实现核心业务逻辑，提供接口定义和内存实现，可以独立使用或被其他层调用。

4. **客户端层 (client)**：提供HTTP客户端功能，可连接远程服务器或使用本地服务实现。

5. **服务器层 (server)**：提供HTTP服务器实现，暴露符合Yggdrasil技术规范的API端点。

这种分层设计使得代码结构清晰，各层职责明确，方便测试和扩展。您可以根据需要只使用其中的某一层，例如只使用服务层进行本地身份验证，或只使用客户端连接到现有的Yggdrasil服务器。

## 注意事项

- 该库需要Go 1.24或更高版本
- 内存存储实现（MemoryYggdrasilService）主要用于开发和测试环境，生产环境建议实现自定义的持久化服务层
- 使用该库访问官方Mojang认证服务时，请遵守Mojang的使用条款
- 对于自定义Yggdrasil服务器，请确保其API兼容官方Yggdrasil技术规范

## 许可证

MIT License

## 致谢

该项目参考了Yggdrasil认证服务的技术规范，特别感谢authlib-injector项目提供的详细文档。