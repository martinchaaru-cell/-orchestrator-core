package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
	Status string            `json:"status"`
	Result string            `json:"result"`
	Error  string            `json:"error"`
}

type Agent struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

var (
	orchestratorURL = os.Getenv("ORCHESTRATOR_URL")
	agentName       = os.Getenv("AGENT_NAME")
	agentID         = uuid.New().String()
)

func init() {
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}
	if agentName == "" {
		agentName = "agent-" + agentID[:8]
	}
	log.Printf("🤖 Agent starting: %s (URL: %s)", agentName, orchestratorURL)
}

func requestTask() (*Task, error) {
	agent := Agent{
		ID:       agentID,
		Status:   "idle",
		LastSeen: time.Now(),
	}

	payload, _ := json.Marshal(agent)
	resp, err := http.Post(
		orchestratorURL+"/agent/request",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil, nil
	}

	var task Task
	json.NewDecoder(resp.Body).Decode(&task)
	return &task, nil
}

func executeTask(task *Task) {
	log.Printf("⚙️  Executing task: %s (type: %s)", task.ID, task.Type)

	var cmd *exec.Cmd

	switch task.Type {
	case "hello_world":
		cmd = exec.Command("bash", "/app/tasks/hello_world.sh")

	case "echo_task":
		message := task.Params["message"]
		cmd = exec.Command("bash", "-c", fmt.Sprintf("echo '%s'", message))

	case "nmap_scan":
		target := task.Params["target"]
		cmd = exec.Command("bash", "/app/tasks/nmap_scan.sh", target)

	default:
		task.Status = "failed"
		task.Error = "Unknown task type: " + task.Type
		goto report
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
		task.Result = string(output)
	} else {
		task.Status = "completed"
		task.Result = string(output)
	}

report:
	payload, _ := json.Marshal(task)
	resp, err := http.Post(
		orchestratorURL+"/agent/result",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Printf("❌ Failed to report result: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("✅ Task completed: %s", task.ID)
}

func main() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("🎯 Agent ready, waiting for tasks...")

	for range ticker.C {
		task, err := requestTask()
		if err != nil {
			log.Printf("⚠️  Error requesting task: %v", err)
			continue
		}

		if task == nil {
			continue
		}

		executeTask(task)
	}
}
