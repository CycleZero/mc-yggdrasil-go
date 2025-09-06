package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/CycleZero/mc-yggdrasil-go/client"
	"github.com/CycleZero/mc-yggdrasil-go/models"
	"github.com/CycleZero/mc-yggdrasil-go/server"
	"github.com/CycleZero/mc-yggdrasil-go/service"
	"github.com/CycleZero/mc-yggdrasil-go/utils"
)

func main() {
	// 示例1: 基本UUID工具使用
	fmt.Println("=== 示例1: 基本UUID工具使用 ===")
	demonstrateUUIDTools()

	// 示例2: 使用本地服务实现（不依赖服务器）
	fmt.Println("\n=== 示例2: 使用本地服务实现 ===")
	demonstrateLocalService()

	// 示例3: 启动Yggdrasil服务器并使用客户端连接
	fmt.Println("\n=== 示例3: 启动Yggdrasil服务器 ===")
	demonstrateServer()
}

// demonstrateUUIDTools 演示UUID相关工具的使用
func demonstrateUUIDTools() {
	playerName := "Steve"
	offlineUUID, err := utils.GenerateOfflinePlayerUUID(playerName)
	if err != nil {
		fmt.Printf("生成离线UUID失败: %v\n", err)
	} else {
		fmt.Printf("玩家 '%s' 的离线UUID: %s\n", playerName, offlineUUID)
	}

	randomUUID := utils.GenerateUUID()
	fmt.Printf("随机生成的UUID: %s\n", randomUUID)

	// 格式化UUID示例
	formattedUUID, err := utils.FormatUUID(randomUUID)
	if err != nil {
		fmt.Printf("格式化UUID失败: %v\n", err)
	} else {
		fmt.Printf("格式化后的UUID: %s\n", formattedUUID)
	}
}

// demonstrateLocalService 演示使用本地服务实现（不依赖服务器）
func demonstrateLocalService() {
	// 创建内存实现的服务
	memoryService := service.NewMemoryYggdrasilService()

	// 添加测试用户
	fmt.Println("创建测试用户...")
	userID, err := memoryService.AddUser("test@example.com", "password123")
	if err != nil {
		fmt.Printf("添加用户失败: %v\n", err)
		return
	}
	fmt.Printf("用户已添加，用户ID: %s\n", userID)

	// 添加测试角色
	profile, err := memoryService.AddProfile(userID, "TestPlayer")
	if err != nil {
		fmt.Printf("添加角色失败: %v\n", err)
		return
	}
	fmt.Printf("角色已添加，角色名称: %s, UUID: %s\n", profile.Name, profile.ID)

	// 创建使用本地服务的客户端
	localClient := client.NewLocalClient(memoryService)

	// 执行认证
	fmt.Println("执行认证...")
	authReq := models.AuthRequest{
		Agent: models.Agent{
			Name:    "Minecraft",
			Version: 1,
		},
		Username:    "test@example.com",
		Password:    "password123",
		RequestUser: true,
	}

	authResp, err := localClient.Auth(authReq)
	if err != nil {
		fmt.Printf("认证失败: %v\n", err)
		return
	}
	fmt.Printf("认证成功! 访问令牌: %s\n", authResp.AccessToken)
	fmt.Printf("选定角色: %s (UUID: %s)\n",
		authResp.SelectedProfile.Name,
		authResp.SelectedProfile.ID)

	// 验证令牌
	fmt.Println("验证令牌...")
	validateReq := models.ValidateRequest{
		AccessToken: authResp.AccessToken,
		ClientToken: authResp.ClientToken,
	}

	isValid, err := localClient.Validate(validateReq)
	if err != nil {
		fmt.Printf("验证令牌失败: %v\n", err)
	} else if isValid {
		fmt.Println("令牌验证成功，令牌有效!")
	} else {
		fmt.Println("令牌验证失败，令牌无效!")
	}

	// 刷新令牌
	fmt.Println("刷新令牌...")
	refreshReq := models.RefreshRequest{
		AccessToken: authResp.AccessToken,
		ClientToken: authResp.ClientToken,
		RequestUser: true,
	}

	refreshResp, err := localClient.Refresh(refreshReq)
	if err != nil {
		fmt.Printf("刷新令牌失败: %v\n", err)
	} else {
		fmt.Printf("刷新令牌成功! 新的访问令牌: %s\n", refreshResp.AccessToken)
	}

	// 使令牌失效
	fmt.Println("使令牌失效...")
	invalidateReq := models.InvalidateRequest{
		AccessToken: refreshResp.AccessToken,
		ClientToken: refreshResp.ClientToken,
	}

	if err := localClient.Invalidate(invalidateReq); err != nil {
		fmt.Printf("使令牌失效失败: %v\n", err)
	} else {
		fmt.Println("令牌已成功失效!")
	}
}

// demonstrateServer 演示启动Yggdrasil服务器并使用客户端连接
func demonstrateServer() {
	// 创建内存实现的服务
	memoryService := service.NewMemoryYggdrasilService()

	// 添加测试用户和角色
	userID, _ := memoryService.AddUser("server@example.com", "server123")
	memoryService.AddProfile(userID, "ServerPlayer")

	// 创建并启动服务器
	yggServer := server.NewYggdrasilServer(8080, memoryService)

	// 使用goroutine启动服务器
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := yggServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v\n", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(1 * time.Second)

	// 创建客户端连接到服务器
	yggClient := client.NewYggdrasilClient("http://localhost:8080")

	// 执行认证
	fmt.Println("通过HTTP客户端执行认证...")
	authReq := models.AuthRequest{
		Agent: models.Agent{
			Name:    "Minecraft",
			Version: 1,
		},
		Username:    "server@example.com",
		Password:    "server123",
		RequestUser: true,
	}

	authResp, err := yggClient.Auth(authReq)
	if err != nil {
		fmt.Printf("认证失败: %v\n", err)
	} else {
		fmt.Printf("认证成功! 访问令牌: %s\n", authResp.AccessToken)
		fmt.Printf("选定角色: %s (UUID: %s)\n",
			authResp.SelectedProfile.Name,
			authResp.SelectedProfile.ID)
	}

	// 设置信号处理，优雅关闭服务器
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// 等待信号或超时（5秒后自动停止）
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("\n服务器演示结束，正在停止服务器...")
		stop <- os.Interrupt
	}()

	<-stop

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := yggServer.Stop(ctx); err != nil {
		log.Printf("服务器关闭失败: %v\n", err)
	}

	wg.Wait()
	fmt.Println("服务器已停止")
}
