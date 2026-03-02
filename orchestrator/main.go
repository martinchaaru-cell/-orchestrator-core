package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Params    map[string]string `json:"params"`
	Status    string            `json:"status"`
	Result    string            `json:"result"`
	Error     string            `json:"error"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Agent struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func init() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Invalid Redis URL: %v", err)
	}

	redisClient = redis.NewClient(opt)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("✅ Redis connected")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/task", createTask)
	router.GET("/task/:id", getTask)
	router.GET("/tasks", listTasks)
	router.POST("/agent/request", requestTask)
	router.POST("/agent/result", reportResult)
	router.GET("/agents", listAgents)
	router.GET("/health", healthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Orchestrator listening on :%s", port)
	router.Run(":" + port)
}

func createTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.ID = uuid.New().String()
	task.Status = "pending"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	taskJSON, _ := json.Marshal(task)
	if err := redisClient.LPush(ctx, "task_queue", string(taskJSON)).Err(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	redisClient.Set(ctx, "task:"+task.ID, string(taskJSON), 24*time.Hour)

	c.JSON(201, task)
	log.Printf("📋 Task submitted: %s (type: %s)", task.ID, task.Type)
}

func getTask(c *gin.Context) {
	taskID := c.Param("id")
	taskJSON, err := redisClient.Get(ctx, "task:"+taskID).Result()
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	var task Task
	json.Unmarshal([]byte(taskJSON), &task)
	c.JSON(200, task)
}

func listTasks(c *gin.Context) {
	var tasks []Task
	keys, _ := redisClient.Keys(ctx, "task:*").Result()

	for _, key := range keys {
		taskJSON, _ := redisClient.Get(ctx, key).Result()
		var task Task
		json.Unmarshal([]byte(taskJSON), &task)
		tasks = append(tasks, task)
	}

	c.JSON(200, gin.H{"tasks": tasks, "count": len(tasks)})
}

func requestTask(c *gin.Context) {
	var agent Agent
	if err := c.BindJSON(&agent); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	agentJSON, _ := json.Marshal(agent)
	redisClient.Set(ctx, "agent:"+agent.ID, string(agentJSON), 5*time.Minute)

	taskJSON, err := redisClient.RPop(ctx, "task_queue").Result()
	if err == redis.Nil {
		c.JSON(204, nil)
		return
	}

	var task Task
	json.Unmarshal([]byte(taskJSON), &task)
	task.Status = "running"
	task.UpdatedAt = time.Now()
	redisClient.Set(ctx, "task:"+task.ID, string(json.RawMessage(taskJSON)), 24*time.Hour)

	c.JSON(200, task)
	log.Printf("⚙️  Task assigned to agent %s: %s", agent.ID, task.ID)
}

func reportResult(c *gin.Context) {
	var result Task
	if err := c.BindJSON(&result); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result.UpdatedAt = time.Now()
	taskJSON, _ := json.Marshal(result)
	redisClient.Set(ctx, "task:"+result.ID, string(taskJSON), 24*time.Hour)

	c.JSON(200, result)
	log.Printf("✅ Result received for task %s (status: %s)", result.ID, result.Status)
}

func listAgents(c *gin.Context) {
	var agents []Agent
	keys, _ := redisClient.Keys(ctx, "agent:*").Result()

	for _, key := range keys {
		agentJSON, _ := redisClient.Get(ctx, key).Result()
		var agent Agent
		json.Unmarshal([]byte(agentJSON), &agent)
		agents = append(agents, agent)
	}

	c.JSON(200, gin.H{"agents": agents, "count": len(agents)})
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
}
