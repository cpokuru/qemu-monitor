package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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

type VMNetwork struct {
	Type         string        `json:"type"`
	ID           string        `json:"id"`
	MAC          string        `json:"mac"`
	PortForwards []PortForward `json:"port_forwards,omitempty"`
}

type PortForward struct {
	Host  int `json:"host"`
	Guest int `json:"guest"`
}

type VMConfig struct {
	Name       string      `json:"name"`
	Disk       string      `json:"disk"`
	Memory     string      `json:"memory"`
	CPUs       string      `json:"cpus"`
	BIOS       string      `json:"bios"`
	Snapshot   bool        `json:"snapshot"`
	Networks   []VMNetwork `json:"networks"`
	SSHPort    *int        `json:"ssh_port"`
	HTTPPort   *int        `json:"http_port"`
	WorkingDir string      `json:"working_dir"`
}

type VMsConfig struct {
	VMs []VMConfig `json:"vms"`
}

var (
	cachedInstances []QEMUInstance
	lastUpdate      time.Time
	vmsConfig       VMsConfig
	configPath      = "vms.json"
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

func loadVMsConfig() error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Warning: %s not found, VM management disabled", configPath)
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &vmsConfig)
}

func findVMConfig(name string) *VMConfig {
	for i := range vmsConfig.VMs {
		if vmsConfig.VMs[i].Name == name {
			return &vmsConfig.VMs[i]
		}
	}
	return nil
}

func buildQEMUCommand(vm *VMConfig) *exec.Cmd {
	args := []string{
		"qemu-system-aarch64",
		"-nographic",
		"-accel", "hvf",
		"-cpu", "cortex-a72",
		"-machine", "virt",
		"-bios", vm.BIOS,
		"-smp", vm.CPUs,
		"-m", vm.Memory,
		"-device", "virtio-rng-pci",
		"-drive", fmt.Sprintf("file=%s,format=qcow2,if=virtio", vm.Disk),
	}

	// Add networks
	for _, net := range vm.Networks {
		netdevArg := net.Type + ",id=" + net.ID
		if net.Type == "user" && len(net.PortForwards) > 0 {
			for _, pf := range net.PortForwards {
				netdevArg += fmt.Sprintf(",hostfwd=tcp::%d-:%d", pf.Host, pf.Guest)
			}
		}
		args = append(args, "-netdev", netdevArg)
		args = append(args, "-device", fmt.Sprintf("virtio-net-pci,netdev=%s,mac=%s", net.ID, net.MAC))
	}

	// Add name
	args = append(args, "-name", vm.Name)

	// Add snapshot mode if enabled
	if vm.Snapshot {
		args = append(args, "-snapshot")
	}

	// Add serial
	args = append(args, "-serial", "mon:stdio")

	cmd := exec.Command("sudo", args...)
	if vm.WorkingDir != "" {
		cmd.Dir = vm.WorkingDir
	}
	return cmd
}

func startVM(name string) error {
	vm := findVMConfig(name)
	if vm == nil {
		return fmt.Errorf("VM configuration not found: %s", name)
	}

	// Check if already running
	for _, inst := range cachedInstances {
		if inst.Name == name {
			return fmt.Errorf("VM already running with PID %s", inst.PID)
		}
	}

	cmd := buildQEMUCommand(vm)
	
	// Start in background
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VM: %v", err)
	}

	log.Printf("Started VM %s with PID %d", name, cmd.Process.Pid)
	return nil
}

func stopVM(pid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return fmt.Errorf("invalid PID: %s", pid)
	}

	// First try graceful shutdown with SIGTERM
	cmd := exec.Command("sudo", "kill", "-TERM", pid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop VM (PID %d): %v - %s", pidInt, err, string(output))
	}

	log.Printf("Sent SIGTERM to PID %d", pidInt)
	return nil
}

func forceStopVM(pid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return fmt.Errorf("invalid PID: %s", pid)
	}

	// Force kill with SIGKILL
	cmd := exec.Command("sudo", "kill", "-9", pid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to force stop VM (PID %d): %v - %s", pidInt, err, string(output))
	}

	log.Printf("Sent SIGKILL to PID %d", pidInt)
	return nil
}

func getShellInfo(name string) (map[string]interface{}, error) {
	vm := findVMConfig(name)
	if vm == nil {
		return nil, fmt.Errorf("VM configuration not found: %s", name)
	}

	// Check if VM is running
	var running bool
	for _, inst := range cachedInstances {
		if inst.Name == name {
			running = true
			break
		}
	}

	info := map[string]interface{}{
		"name":    name,
		"running": running,
	}

	if vm.SSHPort != nil {
		info["ssh_command"] = fmt.Sprintf("ssh -p %d root@localhost", *vm.SSHPort)
		info["ssh_port"] = *vm.SSHPort
	}

	if vm.HTTPPort != nil {
		info["http_url"] = fmt.Sprintf("http://localhost:%d", *vm.HTTPPort)
		info["http_port"] = *vm.HTTPPort
	}

	return info, nil
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

func handleStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if err := startVM(req.Name); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "started", "name": req.Name})
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PID   string `json:"pid"`
		Force bool   `json:"force"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	var err error
	if req.Force {
		err = forceStopVM(req.PID)
	} else {
		err = stopVM(req.PID)
	}

	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	action := "stopped"
	if req.Force {
		action = "force stopped"
	}
	json.NewEncoder(w).Encode(map[string]string{"status": action, "pid": req.PID})
}

func handleShell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	info, err := getShellInfo(req.Name)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(info)
}

func handleVMsConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(w).Encode(vmsConfig)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexHTML)
}

func main() {
	// Load VM configuration
	if err := loadVMsConfig(); err != nil {
		log.Printf("Warning: Failed to load VMs config: %v", err)
	} else {
		log.Printf("Loaded configuration for %d VMs", len(vmsConfig.VMs))
	}

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
	http.HandleFunc("/api/start", handleStart)
	http.HandleFunc("/api/stop", handleStop)
	http.HandleFunc("/api/shell", handleShell)
	http.HandleFunc("/api/vms", handleVMsConfig)

	addr := "0.0.0.0:5450"
	log.Printf("QEMU Instance Tracker starting on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
