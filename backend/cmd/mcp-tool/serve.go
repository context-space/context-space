package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ServeConfig struct {
	Port string
	Host string
}

type CommandRequest struct {
	Command string            `json:"command"`
	Args    map[string]string `json:"args"`
}

type CommandResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	JobID   string `json:"jobId,omitempty"`
}

type JobStatus struct {
	ID        string     `json:"id"`
	Command   string     `json:"command"`
	Status    string     `json:"status"` // running, completed, failed
	Output    []string   `json:"output"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime,omitempty"`
	mutex     sync.RWMutex
}

type WSMessage struct {
	JobID   string `json:"jobId"`
	Line    string `json:"line"`
	IsError bool   `json:"isError"`
}

var (
	jobs           = make(map[string]*JobStatus)
	jobMutex       sync.RWMutex
	wsClients      = make(map[*websocket.Conn]bool)
	wsClientsMutex sync.RWMutex
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
)

func runServeCommand(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Printf(`Usage: mcp-tool serve [options]

Start HTTP server with Web UI for mcp-tool commands

Options:
`)
		fs.PrintDefaults()
		fmt.Printf(`
Examples:
  # Start server on default port
  mcp-tool serve

  # Start server on custom port
  mcp-tool serve -port 8080

  # Start server on custom host and port
  mcp-tool serve -host 0.0.0.0 -port 3000
`)
	}

	var config ServeConfig
	fs.StringVar(&config.Port, "port", "8080", "Port to serve on")
	fs.StringVar(&config.Host, "host", "localhost", "Host to serve on")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	fmt.Printf("Starting mcp-tool web server on http://%s:%s\n", config.Host, config.Port)

	if err := startServer(config); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startServer(config ServeConfig) error {
	// Setup routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/execute", handleExecute)
	http.HandleFunc("/api/jobs", handleJobs)
	http.HandleFunc("/api/jobs/", handleJobDetail)
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/static/", http.FileServer(http.Dir("static")))

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	fmt.Printf("Server starting at http://%s\n", addr)
	fmt.Printf("Open your browser and navigate to the URL above\n")

	return server.ListenAndServe()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Context Space - MCP Adapter Tool Web UI</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f5f5f5;
            line-height: 1.6;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: #2563eb;
            color: white;
            padding: 2rem;
            border-radius: 8px;
            margin-bottom: 2rem;
            text-align: center;
        }
        
        .command-section {
            background: white;
            border-radius: 8px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        
        .command-title {
            font-size: 1.5rem;
            font-weight: 600;
            margin-bottom: 1rem;
            color: #1f2937;
        }
        
        .form-group {
            margin-bottom: 1rem;
        }
        
        .form-row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1rem;
            margin-bottom: 1rem;
        }
        
        label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 500;
            color: #374151;
        }
        
        input, textarea, select {
            width: 100%;
            padding: 0.75rem;
            border: 1px solid #d1d5db;
            border-radius: 4px;
            font-size: 1rem;
        }
        
        input:focus, textarea:focus, select:focus {
            outline: none;
            border-color: #2563eb;
            box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
        }
        
        .btn {
            background: #2563eb;
            color: white;
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: 4px;
            font-size: 1rem;
            font-weight: 500;
            cursor: pointer;
            transition: background 0.2s;
        }
        
        .btn:hover {
            background: #1d4ed8;
        }
        
        .btn:disabled {
            background: #9ca3af;
            cursor: not-allowed;
        }
        
        .output-section {
            background: #1f2937;
            color: #f9fafb;
            padding: 1rem;
            border-radius: 4px;
            font-family: 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace;
            max-height: 400px;
            overflow-y: auto;
            display: none;
        }
        
        .output-line {
            margin-bottom: 0.25rem;
        }
        
        .output-error {
            color: #fca5a5;
        }
        
        .status {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 4px;
            font-size: 0.875rem;
            font-weight: 500;
            margin-left: 1rem;
        }
        
        .status-running {
            background: #fef3c7;
            color: #92400e;
        }
        
        .status-completed {
            background: #d1fae5;
            color: #065f46;
        }
        
        .status-failed {
            background: #fee2e2;
            color: #991b1b;
        }
        
        .required {
            color: #dc2626;
        }
        
        .checkbox-container {
            background: #f8fafc;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
            transition: all 0.2s ease;
        }
        
        .checkbox-container:hover {
            border-color: #2563eb;
            background: #f1f5f9;
        }
        
        .checkbox-wrapper {
            display: flex;
            align-items: center;
            cursor: pointer;
            user-select: none;
        }
        
        .custom-checkbox {
            position: relative;
            margin-right: 12px;
        }
        
        .custom-checkbox input[type="checkbox"] {
            opacity: 0;
            position: absolute;
            width: 0;
            height: 0;
        }
        
        .checkmark {
            display: block;
            width: 20px;
            height: 20px;
            background: white;
            border: 2px solid #d1d5db;
            border-radius: 4px;
            position: relative;
            transition: all 0.2s ease;
        }
        
        .custom-checkbox input[type="checkbox"]:checked + .checkmark {
            background: #2563eb;
            border-color: #2563eb;
        }
        
        .checkmark:after {
            content: "";
            position: absolute;
            display: none;
            left: 6px;
            top: 2px;
            width: 6px;
            height: 10px;
            border: solid white;
            border-width: 0 2px 2px 0;
            transform: rotate(45deg);
        }
        
        .custom-checkbox input[type="checkbox"]:checked + .checkmark:after {
            display: block;
        }
        
        .checkbox-label {
            font-weight: 500;
            color: #374151;
            display: flex;
            flex-direction: column;
        }
        
        .checkbox-title {
            font-size: 1rem;
            margin-bottom: 2px;
        }
        
        .checkbox-description {
            font-size: 0.875rem;
            color: #6b7280;
            font-weight: 400;
        }
        
        .ai-icon {
            display: inline-block;
            margin-right: 6px;
            font-size: 1.1em;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Context Space - MCP Adapter Tool Web UI</h1>
            <p>Visual interface for effective MCP Adapters development</p>
        </div>

        <!-- Generate Command Section -->
        <div class="command-section">
            <h2 class="command-title">Generate Provider Adapter</h2>
            <form id="generateForm">
                <div class="form-row">
                    <div class="form-group">
                        <label for="command">Command <span class="required">*</span></label>
                        <input type="text" id="command" name="command" placeholder="npx" required>
                    </div>
                    <div class="form-group">
                        <label for="identifier">Identifier <span class="required">*</span></label>
                        <input type="text" id="identifier" name="identifier" placeholder="filesystem" required>
                    </div>
                </div>
                
                <div class="form-row">
                    <div class="form-group">
                        <label for="name">Display Name <span class="required">*</span></label>
                        <input type="text" id="name" name="name" placeholder="Filesystem Provider" required>
                    </div>
                    <div class="form-group">
                        <label for="auth">Authentication Type</label>
                        <select id="auth" name="auth">
                            <option value="none">None</option>
                            <option value="apikey">API Key</option>
                        </select>
                    </div>
                </div>

                <div class="form-group">
                    <label for="description">Description <span class="required">*</span></label>
                    <textarea id="description" name="description" rows="2" placeholder="Provider description" required></textarea>
                </div>

                <div class="form-group">
                    <label for="args">Arguments (one per line)</label>
                    <textarea id="args" name="args" rows="3" placeholder="-y&#10;@modelcontextprotocol/server-filesystem&#10;/tmp"></textarea>
                </div>

                <div class="form-row">
                    <div class="form-group">
                        <label for="categories">Categories (comma-separated)</label>
                        <input type="text" id="categories" name="categories" placeholder="storage,files">
                    </div>
                    <div class="form-group">
                        <label for="output">Output Directory</label>
                        <input type="text" id="output" name="output" placeholder="configs/providers" value="configs/providers">
                    </div>
                </div>

                <div class="form-group">
                    <label for="envs">Environment Variables (one per line, KEY=VALUE)</label>
                    <textarea id="envs" name="envs" rows="2" placeholder="API_KEY=secret123&#10;DEBUG=mcp:*"></textarea>
                </div>

                <div class="checkbox-container">
                    <label class="checkbox-wrapper" for="translate">
                        <div class="custom-checkbox">
                            <input type="checkbox" id="translate" name="translate">
                            <span class="checkmark"></span>
                        </div>
                        <div class="checkbox-label">
                            <span class="checkbox-title">
                                <span class="ai-icon">🤖</span>Enable AI Translation
                            </span>
                            <span class="checkbox-description">
                                Automatically generate Chinese translations for provider descriptions (requires OPENAI_API_KEY)
                            </span>
                        </div>
                    </label>
                </div>

                <button type="submit" class="btn" id="generateBtn">Generate Adapter</button>
                <span id="generateStatus"></span>
            </form>

            <div id="generateOutput" class="output-section">
                <div id="generateOutputContent"></div>
            </div>
        </div>

        <!-- Call Command Section -->
        <div class="command-section">
            <h2 class="command-title">Call MCP Operation</h2>
            <form id="callForm">
                <div class="form-row">
                    <div class="form-group">
                        <label for="callCommand">Command <span class="required">*</span></label>
                        <input type="text" id="callCommand" name="command" placeholder="npx" required>
                    </div>
                    <div class="form-group">
                        <label for="operation">Operation <span class="required">*</span></label>
                        <input type="text" id="operation" name="operation" placeholder="list_directory" required>
                    </div>
                </div>

                <div class="form-group">
                    <label for="callArgs">Arguments (one per line)</label>
                    <textarea id="callArgs" name="args" rows="3" placeholder="-y&#10;@modelcontextprotocol/server-filesystem&#10;/tmp"></textarea>
                </div>

                <div class="form-group">
                    <label for="params">Parameters (JSON)</label>
                    <textarea id="params" name="params" rows="3" placeholder='{"path": "/tmp"}'></textarea>
                </div>

                <div class="form-group">
                    <label for="callEnvs">Environment Variables (one per line, KEY=VALUE)</label>
                    <textarea id="callEnvs" name="envs" rows="2" placeholder="API_KEY=secret123"></textarea>
                </div>

                <button type="submit" class="btn" id="callBtn">Call Operation</button>
                <span id="callStatus"></span>
            </form>

            <div id="callOutput" class="output-section">
                <div id="callOutputContent"></div>
            </div>
        </div>
    </div>

    <script>
        let ws = null;
        let currentJobId = null;

        function connectWebSocket() {
            const wsUrl = 'ws://' + window.location.host + '/ws';
            ws = new WebSocket(wsUrl);
            
            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                if (data.jobId === currentJobId) {
                    appendOutput(data.line, data.isError);
                }
            };
            
            ws.onclose = function() {
                setTimeout(connectWebSocket, 3000); // Reconnect after 3 seconds
            };
        }

        function appendOutput(line, isError) {
            const outputContent = document.getElementById(getCurrentOutputContentId());
            const outputSection = document.getElementById(getCurrentOutputId());
            
            const lineDiv = document.createElement('div');
            lineDiv.className = 'output-line' + (isError ? ' output-error' : '');
            lineDiv.textContent = line;
            
            outputContent.appendChild(lineDiv);
            outputSection.style.display = 'block';
            outputSection.scrollTop = outputSection.scrollHeight;
        }

        function getCurrentOutputId() {
            return currentJobId && currentJobId.startsWith('generate') ? 'generateOutput' : 'callOutput';
        }

        function getCurrentOutputContentId() {
            return currentJobId && currentJobId.startsWith('generate') ? 'generateOutputContent' : 'callOutputContent';
        }

        function clearOutput(outputContentId) {
            document.getElementById(outputContentId).innerHTML = '';
        }

        function setStatus(statusElementId, text, className) {
            const statusEl = document.getElementById(statusElementId);
            statusEl.textContent = text;
            statusEl.className = 'status ' + className;
        }

        function submitForm(formId, command, statusElementId, outputContentId) {
            const form = document.getElementById(formId);
            const formData = new FormData(form);
            const btn = form.querySelector('button[type="submit"]');
            
            // Disable button
            btn.disabled = true;
            
            // Clear previous output
            clearOutput(outputContentId);
            
            // Prepare request data
            const args = {};
            for (let [key, value] of formData.entries()) {
                if (value.trim()) {
                    args[key] = value;
                }
            }
            
            const requestData = {
                command: command,
                args: args
            };
            
            setStatus(statusElementId, 'Starting...', 'status-running');
            
            fetch('/api/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    currentJobId = data.jobId;
                    setStatus(statusElementId, 'Running...', 'status-running');
                } else {
                    setStatus(statusElementId, 'Failed: ' + data.message, 'status-failed');
                    btn.disabled = false;
                }
            })
            .catch(error => {
                setStatus(statusElementId, 'Error: ' + error.message, 'status-failed');
                btn.disabled = false;
            });
        }

        // Setup form handlers
        document.getElementById('generateForm').addEventListener('submit', function(e) {
            e.preventDefault();
            submitForm('generateForm', 'generate', 'generateStatus', 'generateOutputContent');
        });

        document.getElementById('callForm').addEventListener('submit', function(e) {
            e.preventDefault();
            submitForm('callForm', 'call', 'callStatus', 'callOutputContent');
        });

        // Connect WebSocket
        connectWebSocket();
        
        // Check job status periodically
        setInterval(function() {
            if (currentJobId) {
                fetch('/api/jobs/' + currentJobId)
                    .then(response => response.json())
                    .then(job => {
                        if (job.status === 'completed') {
                            const statusElementId = currentJobId.startsWith('generate') ? 'generateStatus' : 'callStatus';
                            const btn = document.querySelector(currentJobId.startsWith('generate') ? '#generateBtn' : '#callBtn');
                            setStatus(statusElementId, 'Completed', 'status-completed');
                            btn.disabled = false;
                            currentJobId = null;
                        } else if (job.status === 'failed') {
                            const statusElementId = currentJobId.startsWith('generate') ? 'generateStatus' : 'callStatus';
                            const btn = document.querySelector(currentJobId.startsWith('generate') ? '#generateBtn' : '#callBtn');
                            setStatus(statusElementId, 'Failed', 'status-failed');
                            btn.disabled = false;
                            currentJobId = null;
                        }
                    })
                    .catch(console.error);
            }
        }, 2000);
    </script>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, nil)
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	jobID := fmt.Sprintf("%s_%d", req.Command, time.Now().UnixNano())

	job := &JobStatus{
		ID:        jobID,
		Command:   req.Command,
		Status:    "running",
		Output:    []string{},
		StartTime: time.Now(),
	}

	jobMutex.Lock()
	jobs[jobID] = job
	jobMutex.Unlock()

	// Start command execution in background
	go executeCommand(job, req)

	response := CommandResponse{
		Success: true,
		Message: "Command started",
		JobID:   jobID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleJobs(w http.ResponseWriter, r *http.Request) {
	jobMutex.RLock()
	defer jobMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func handleJobDetail(w http.ResponseWriter, r *http.Request) {
	jobID := strings.TrimPrefix(r.URL.Path, "/api/jobs/")

	jobMutex.RLock()
	job, exists := jobs[jobID]
	jobMutex.RUnlock()

	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	job.mutex.RLock()
	defer job.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	wsClientsMutex.Lock()
	wsClients[conn] = true
	wsClientsMutex.Unlock()

	defer func() {
		wsClientsMutex.Lock()
		delete(wsClients, conn)
		wsClientsMutex.Unlock()
	}()

	// Keep connection alive by reading ping/pong messages
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle ping messages
		if messageType == websocket.PingMessage {
			if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				log.Printf("Failed to send pong: %v", err)
				break
			}
		}
	}
}

func executeCommand(job *JobStatus, req CommandRequest) {
	defer func() {
		job.mutex.Lock()
		endTime := time.Now()
		job.EndTime = &endTime
		job.mutex.Unlock()
	}()

	// Build command arguments
	args := buildCommandArgs(req)

	// Detect how we're running and build the appropriate command
	useGoRun, err := detectRunningMode()
	if err != nil {
		job.mutex.Lock()
		job.Status = "failed"
		job.Output = append(job.Output, fmt.Sprintf("Failed to detect running mode: %v", err))
		job.mutex.Unlock()
		return
	}

	var cmd *exec.Cmd
	if useGoRun {
		// We're running via 'go run'
		// Get current working directory
		wd, _ := os.Getwd()

		// Check if we're in cmd/mcp-tool directory or project root
		var allArgs []string
		if filepath.Base(wd) == "mcp-tool" {
			// We're in cmd/mcp-tool directory, use 'go run .'
			allArgs = append([]string{"run", "."}, args...)
		} else {
			// We're likely in project root, use 'go run ./cmd/mcp-tool'
			allArgs = append([]string{"run", "./cmd/mcp-tool"}, args...)
		}

		cmd = exec.Command("go", allArgs...)
		// Don't set cmd.Dir - keep current working directory
	} else {
		// We're running as a compiled binary
		cmd = exec.Command(os.Args[0], args...)
		// Don't set cmd.Dir - keep current working directory
	}

	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		job.mutex.Lock()
		job.Status = "failed"
		job.Output = append(job.Output, fmt.Sprintf("Failed to create stdout pipe: %v", err))
		job.mutex.Unlock()
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		job.mutex.Lock()
		job.Status = "failed"
		job.Output = append(job.Output, fmt.Sprintf("Failed to create stderr pipe: %v", err))
		job.mutex.Unlock()
		return
	}

	if err := cmd.Start(); err != nil {
		job.mutex.Lock()
		job.Status = "failed"
		job.Output = append(job.Output, fmt.Sprintf("Failed to start command: %v", err))
		job.mutex.Unlock()
		return
	}

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			job.mutex.Lock()
			job.Output = append(job.Output, line)
			job.mutex.Unlock()

			// Broadcast to WebSocket clients
			broadcastOutput(job.ID, line, false)
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			job.mutex.Lock()
			job.Output = append(job.Output, "ERROR: "+line)
			job.mutex.Unlock()

			// Broadcast to WebSocket clients
			broadcastOutput(job.ID, line, true)
		}
	}()

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		job.mutex.Lock()
		job.Status = "failed"
		job.Output = append(job.Output, fmt.Sprintf("Command failed: %v", err))
		job.mutex.Unlock()
	} else {
		job.mutex.Lock()
		job.Status = "completed"
		job.mutex.Unlock()
	}
}

func buildCommandArgs(req CommandRequest) []string {
	args := []string{req.Command}

	// Handle generate command
	if req.Command == "generate" {
		if cmd, ok := req.Args["command"]; ok && cmd != "" {
			args = append(args, "-command", cmd)
		}
		if id, ok := req.Args["identifier"]; ok && id != "" {
			args = append(args, "-identifier", id)
		}
		if name, ok := req.Args["name"]; ok && name != "" {
			args = append(args, "-name", name)
		}
		if desc, ok := req.Args["description"]; ok && desc != "" {
			args = append(args, "-description", desc)
		}
		if auth, ok := req.Args["auth"]; ok && auth != "" {
			args = append(args, "-auth", auth)
		}
		if cats, ok := req.Args["categories"]; ok && cats != "" {
			args = append(args, "-categories", cats)
		}
		if output, ok := req.Args["output"]; ok && output != "" {
			args = append(args, "-output", output)
		}
		if translate, ok := req.Args["translate"]; ok && translate == "on" {
			args = append(args, "-translate")
		}

		// Handle args (multiline)
		if cmdArgs, ok := req.Args["args"]; ok && cmdArgs != "" {
			lines := strings.Split(cmdArgs, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					args = append(args, "-arg", line)
				}
			}
		}

		// Handle envs (multiline)
		if envs, ok := req.Args["envs"]; ok && envs != "" {
			lines := strings.Split(envs, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					args = append(args, "-env", line)
				}
			}
		}
	}

	// Handle call command
	if req.Command == "call" {
		if cmd, ok := req.Args["command"]; ok && cmd != "" {
			args = append(args, "-command", cmd)
		}
		if op, ok := req.Args["operation"]; ok && op != "" {
			args = append(args, "-operation", op)
		}
		if params, ok := req.Args["params"]; ok && params != "" {
			args = append(args, "-params", params)
		}

		// Handle args (multiline)
		if cmdArgs, ok := req.Args["args"]; ok && cmdArgs != "" {
			lines := strings.Split(cmdArgs, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					args = append(args, "-arg", line)
				}
			}
		}

		// Handle envs (multiline)
		if envs, ok := req.Args["envs"]; ok && envs != "" {
			lines := strings.Split(envs, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					args = append(args, "-env", line)
				}
			}
		}
	}

	return args
}

// Simplified broadcast function - in a real implementation you'd manage WebSocket connections
func broadcastOutput(jobID, line string, isError bool) {
	wsClientsMutex.RLock()
	defer wsClientsMutex.RUnlock()

	for conn := range wsClients {
		if err := conn.WriteJSON(WSMessage{JobID: jobID, Line: line, IsError: isError}); err != nil {
			log.Printf("Failed to write to WebSocket connection: %v", err)
			conn.Close() // Close the connection on error
		}
	}
}

// detectRunningMode detects whether the program is running via 'go run' or as a compiled binary
func detectRunningMode() (useGoRun bool, err error) {
	execPath := os.Args[0]

	// Check if we're running from a temporary directory (typical for go run)
	if strings.Contains(execPath, os.TempDir()) || strings.Contains(execPath, "/var/folders/") || strings.Contains(execPath, "\\Temp\\") {
		// We're likely running via 'go run'
		return true, nil
	}

	// We're running as a compiled binary
	return false, nil
}
