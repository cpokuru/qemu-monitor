package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type QEMUInstance struct {
	PID          string    `json:"pid"`
	PPID         string    `json:"ppid"`
	User         string    `json:"user"`
	CPUTime      string    `json:"cpu_time"`
	StartTime    string    `json:"start_time"`
	Memory       string    `json:"memory"`
	CPUCount     string    `json:"cpu_count"`
	DiskImage    string    `json:"disk_image"`
	Name         string    `json:"name"`
	Machine      string    `json:"machine"`
	Networks     []Network `json:"networks"`
	Type         string    `json:"type"` // "multipass" or "custom"
	Status       string    `json:"status"`
	Uptime       string    `json:"uptime"`
}

type Network struct {
	Type string `json:"type"`
	MAC  string `json:"mac"`
}

type Response struct {
	Instances   []QEMUInstance `json:"instances"`
	Count       int            `json:"count"`
	LastUpdated string         `json:"last_updated"`
}

var (
	cachedInstances []QEMUInstance
	lastUpdate      time.Time
)

func parseQEMUProcess(fields []string, cmdline string) QEMUInstance {
	instance := QEMUInstance{
		User:      fields[0],
		PID:       fields[1],
		PPID:      fields[2],
		CPUTime:   fields[6],
		StartTime: fields[4],
		Networks:  []Network{},
	}

	// Determine type
	if strings.Contains(cmdline, "multipass") {
		instance.Type = "multipass"
	} else {
		instance.Type = "custom"
	}

	// Extract memory
	memRegex := regexp.MustCompile(`-m\s+(\d+[MG]?)`)
	if match := memRegex.FindStringSubmatch(cmdline); len(match) > 1 {
		instance.Memory = match[1]
	}

	// Extract CPU count
	cpuRegex := regexp.MustCompile(`-smp\s+(\d+)`)
	if match := cpuRegex.FindStringSubmatch(cmdline); len(match) > 1 {
		instance.CPUCount = match[1]
	}

	// Extract disk image
	diskRegex := regexp.MustCompile(`file=([^,]+\.(?:qcow2|img))`)
	if match := diskRegex.FindStringSubmatch(cmdline); len(match) > 1 {
		diskPath := match[1]
		parts := strings.Split(diskPath, "/")
		instance.DiskImage = parts[len(parts)-1]
	}

	// Extract name
	nameRegex := regexp.MustCompile(`-name\s+([^\s]+)`)
	if match := nameRegex.FindStringSubmatch(cmdline); len(match) > 1 {
		instance.Name = match[1]
	} else if strings.Contains(cmdline, "multipass") {
		// Extract from path for multipass
		pathRegex := regexp.MustCompile(`instances/([^/]+)/`)
		if match := pathRegex.FindStringSubmatch(cmdline); len(match) > 1 {
			instance.Name = match[1]
		}
	} else if instance.DiskImage != "" {
		// Use disk image name as fallback
		instance.Name = strings.TrimSuffix(instance.DiskImage, ".qcow2")
	}

	// Extract machine type
	machineRegex := regexp.MustCompile(`-machine\s+([^,\s]+)`)
	if match := machineRegex.FindStringSubmatch(cmdline); len(match) > 1 {
		instance.Machine = match[1]
	}

	// Extract networks
	macRegex := regexp.MustCompile(`mac=([0-9a-f:]+)`)
	netTypeRegex := regexp.MustCompile(`-netdev\s+([^,]+)`)
	
	macMatches := macRegex.FindAllStringSubmatch(cmdline, -1)
	netTypeMatches := netTypeRegex.FindAllStringSubmatch(cmdline, -1)
	
	for i, match := range macMatches {
		net := Network{MAC: match[1]}
		if i < len(netTypeMatches) {
			net.Type = netTypeMatches[i][1]
		}
		instance.Networks = append(instance.Networks, net)
	}

	// Determine status
	if strings.Contains(cmdline, "-loadvm suspend") {
		instance.Status = "suspended"
	} else if strings.Contains(cmdline, "-snapshot") {
		instance.Status = "snapshot"
	} else {
		instance.Status = "running"
	}

	// Calculate uptime from CPU time
	instance.Uptime = instance.CPUTime

	return instance
}

func getQEMUInstances() ([]QEMUInstance, error) {
	cmd := exec.Command("ps", "-ef")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	instances := []QEMUInstance{}

	for _, line := range lines {
		if !strings.Contains(line, "qemu-system-aarch64") || strings.Contains(line, "grep") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		// Get full command line
		cmdline := strings.Join(fields[7:], " ")
		
		// Skip sudo wrapper processes - we only want the actual qemu process
		if strings.HasPrefix(cmdline, "sudo ") {
			continue
		}
		
		instance := parseQEMUProcess(fields, cmdline)
		instances = append(instances, instance)
	}

	return instances, nil
}

func updateInstances() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		instances, err := getQEMUInstances()
		if err != nil {
			log.Printf("Error getting instances: %v", err)
		} else {
			cachedInstances = instances
			lastUpdate = time.Now()
		}
		<-ticker.C
	}
}

func handleInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := Response{
		Instances:   cachedInstances,
		Count:       len(cachedInstances),
		LastUpdated: lastUpdate.Format("2006-01-02 15:04:05"),
	}

	json.NewEncoder(w).Encode(response)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexHTML)
}

func main() {
	// Initial load
	instances, err := getQEMUInstances()
	if err != nil {
		log.Fatal(err)
	}
	cachedInstances = instances
	lastUpdate = time.Now()

	// Start background updater
	go updateInstances()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/instances", handleInstances)

	addr := "0.0.0.0:5450"
	log.Printf("QEMU Instance Tracker starting on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
