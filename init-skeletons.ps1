# Initialize project skeletons for microclaw and clawland-fleet via GitHub API
# Run from project root: .\init-skeletons.ps1

$ErrorActionPreference = "Stop"

function Create-File {
    param([string]$Repo, [string]$Path, [string]$Content, [string]$Msg)
    $b64 = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($Content))
    $body = @{message=$Msg; content=$b64} | ConvertTo-Json
    $body | gh api -X PUT "/repos/Clawland-AI/$Repo/contents/$Path" --input -
}

# --- microclaw ---
$pio = @'
[env:esp32]
platform = espressif32
board = esp32dev
framework = arduino
monitor_speed = 115200
lib_deps =
    knolleary/PubSubClient@^2.8
    adafruit/DHT sensor library@^1.4
    adafruit/Adafruit Unified Sensor@^1.1
'@
Create-File -Repo microclaw -Path platformio.ini -Content $pio -Msg "chore: initialize project structure"
Create-File -Repo microclaw -Path src/main.cpp -Content (@'
/**
 * MicroClaw — Sensor-level micro AI Agent
 * Runs on ESP32 with <1MB RAM, $2-5 hardware cost.
 * Part of the Clawland edge AI agent network.
 */
#include <Arduino.h>
#include <WiFi.h>

const char* AGENT_NAME = "microclaw";
const char* VERSION = "0.1.0";

void setup() {
    Serial.begin(115200);
    Serial.println("🦀 MicroClaw Agent v0.1.0");
    Serial.println("   MCU-level sensor agent starting...");
    Serial.println("   Waiting for sensor configuration...");
}

void loop() {
    delay(1000);
}
'@) -Msg "chore: initialize project structure"

$mcg = @'
.pio/
.vscode/
lib/
include/
'@
Create-File -Repo microclaw -Path .gitignore -Content $mcg -Msg "chore: initialize project structure"
Write-Host "microclaw: skeleton created"

# --- clawland-fleet ---
$gomod = @'
module github.com/Clawland-AI/clawland-fleet

go 1.22
'@
Create-File -Repo clawland-fleet -Path go.mod -Content $gomod -Msg "chore: initialize project structure"

$maingo = @'
// Package main is the entry point for the Clawland Fleet Manager.
// Fleet Manager handles Cloud-Edge orchestration: node registration,
// heartbeat monitoring, event collection, and command dispatch.
package main

import (
	"fmt"
	"log"
	"os"
)

const version = "0.1.0"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🦀 Clawland Fleet Manager v%s\n", version)
	fmt.Printf("   Cloud-Edge orchestration starting on :%s...\n", port)
	fmt.Println("   Waiting for edge agent registrations...")

	log.Printf("Fleet Manager listening on :%s", port)
}
'@
Create-File -Repo clawland-fleet -Path cmd/fleet/main.go -Content $maingo -Msg "chore: initialize project structure"

$registry = @'
// Package fleet provides the core Fleet Manager functionality.
package fleet

import (
	"sync"
	"time"
)

// Node represents a registered edge agent.
type Node struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"` // picclaw, nanoclaw, microclaw
	Capabilities []string          `json:"capabilities"`
	Location     string            `json:"location,omitempty"`
	LastSeen     time.Time         `json:"last_seen"`
	Status       string            `json:"status"` // online, offline, degraded
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Registry manages registered edge nodes.
type Registry struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

// NewRegistry creates a new node registry.
func NewRegistry() *Registry {
	return &Registry{nodes: make(map[string]*Node)}
}

// Register adds or updates a node in the registry.
func (r *Registry) Register(node *Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	node.LastSeen = time.Now()
	node.Status = "online"
	r.nodes[node.ID] = node
}

// Heartbeat updates the last seen time for a node.
func (r *Registry) Heartbeat(nodeID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n, ok := r.nodes[nodeID]; ok {
		n.LastSeen = time.Now()
		n.Status = "online"
		return true
	}
	return false
}

// List returns all registered nodes.
func (r *Registry) List() []*Node {
	r.mu.RLock()
	defer r.mu.RUnlock()
	nodes := make([]*Node, 0, len(r.nodes))
	for _, n := range r.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}
'@
Create-File -Repo clawland-fleet -Path pkg/fleet/registry.go -Content $registry -Msg "chore: initialize project structure"

$fleetg = @'
/fleet
*.exe
*.test
*.out
.env
vendor/
'@
Create-File -Repo clawland-fleet -Path .gitignore -Content $fleetg -Msg "chore: initialize project structure"
Write-Host "clawland-fleet: skeleton created"
