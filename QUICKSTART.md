# QEMU Monitor - Quick Start Guide

## Installation & Running (3 Simple Steps)

### Option 1: Using the Start Script (Easiest)
```bash
./start.sh
```
This will automatically build and run the application.

### Option 2: Using Make
```bash
make run
```

### Option 3: Manual Build
```bash
go build -o qemu-monitor
./qemu-monitor
```

## Access the Dashboard

Once running, open your browser to:
```
http://localhost:5450
```

## What You'll See

The dashboard shows:
- **Real-time monitoring** of all QEMU instances
- **Auto-refresh** every 5 seconds
- **Color-coded status** (green = running, amber = suspended, blue = snapshot)
- **Detailed information** for each VM:
  - Process ID and resource usage
  - Memory allocation and CPU cores
  - Network interfaces with MAC addresses
  - Disk image paths
  - Instance type (Multipass vs Custom)

## Features

### Filter View
Click the filter buttons to show:
- **All** - Every QEMU instance
- **Running** - Only active VMs
- **Suspended** - Only suspended VMs
- **Multipass** - Only Multipass-managed instances
- **Custom** - Only your custom QEMU VMs (like RDK-B instances)

### API Access
Use the REST API for automation:
```bash
curl http://localhost:5450/api/instances
```

Returns JSON with all instance data.

## Customization

### Change Port
Edit `main.go`, line with:
```go
addr := "0.0.0.0:5450"  // Change 5450 to your preferred port
```

### Change Update Frequency
Edit `main.go`, look for:
```go
ticker := time.NewTicker(5 * time.Second)  // Change duration here
```

### Modify Theme Colors
Edit `html.go`, CSS variables section:
```css
:root {
    --accent-green: #00ff88;  /* Change to your preferred colors */
    --accent-amber: #ffaa00;
    /* ... more variables */
}
```

## Troubleshooting

**Problem**: No instances showing
- **Solution**: Make sure QEMU instances are running with `ps -ef | grep qemu`

**Problem**: Port 5450 already in use
- **Solution**: Change the port in `main.go` or stop the conflicting service

**Problem**: Permission denied
- **Solution**: The app needs to read process info - ensure you can run `ps -ef`

## Integration Ideas

### System Service (macOS)
Create a LaunchAgent to run at startup - see `INTEGRATION.md` for details.

### Alfred/Raycast Workflow
Use the API endpoint in a custom workflow:
```bash
curl -s http://localhost:5450/api/instances | jq '.count'
```

### Terminal Alias
Add to your `.zshrc` or `.bashrc`:
```bash
alias qm='open http://localhost:5450'
```

## What Gets Monitored

The app automatically detects and tracks:
- ✅ Multipass VMs (ubuntu, primary, etc.)
- ✅ Custom QEMU instances (your RDK-B VMs)
- ✅ Memory allocation
- ✅ CPU core count
- ✅ Network configurations
- ✅ Disk images
- ✅ Machine type
- ✅ Runtime status

## Next Steps

1. **Run the app**: `./start.sh`
2. **Open the dashboard**: http://localhost:5450
3. **Explore your instances**: See all your QEMU VMs at a glance
4. **Set up auto-start**: Make it run on boot (optional)

Enjoy monitoring your QEMU instances! 🚀
