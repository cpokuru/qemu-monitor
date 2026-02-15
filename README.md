# QEMU Instance Monitor

A real-time web-based monitoring dashboard for tracking QEMU virtual machine instances on macOS (and Linux).

## Features

- 🔍 **Auto-Discovery**: Automatically detects all running QEMU instances
- 📊 **Real-Time Monitoring**: Updates every 5 seconds
- 🎨 **Modern UI**: Industrial/terminal-inspired aesthetic with dark theme
- 🔧 **Detailed Info**: Shows memory, CPU, networks, disk images, and status
- 🏷️ **Instance Types**: Distinguishes between Multipass and custom instances
- 🎯 **Smart Filtering**: Filter by status (running/suspended) or type
- 📱 **Responsive**: Works on desktop and mobile

## Prerequisites

- Go 1.21+ installed
- macOS with QEMU instances (or Linux)
- Permissions to run `ps -ef` command

## Installation

1. Build the application:
```bash
go build -o qemu-monitor
```

2. Run the application:
```bash
./qemu-monitor
```

3. Open your browser and navigate to:
```
http://localhost:5450
```

## Usage

The dashboard will automatically detect and display:

- **Multipass VMs**: Ubuntu instances managed by Multipass
- **Custom QEMU instances**: Your RDK-B and other custom VMs

### Extracted Information

For each instance, the monitor displays:

- **Name**: Instance/VM name
- **Status**: Running, Suspended, or Snapshot mode
- **PID**: Process ID
- **Memory**: Allocated RAM
- **CPU Cores**: Number of vCPUs
- **CPU Time**: Total CPU time used
- **Machine Type**: QEMU machine type (e.g., virt)
- **Disk Image**: Path to the disk image file
- **Networks**: All network interfaces with MAC addresses

### Filtering

Use the filter buttons at the top to:
- View all instances
- Show only running instances
- Show only suspended instances
- Filter by Multipass instances
- Filter by custom instances

## API Endpoints

The application exposes a REST API:

### GET /api/instances

Returns JSON with all QEMU instances:

```json
{
  "instances": [
    {
      "pid": "46682",
      "ppid": "46594",
      "user": "root",
      "cpu_time": "26:07.05",
      "memory": "6144M",
      "cpu_count": "6",
      "disk_image": "rdk.snapshot.qcow2",
      "name": "rdk.snapshot",
      "machine": "virt",
      "networks": [
        {"type": "vmnet-shared", "mac": "52:54:00:2d:6e:99"},
        {"type": "vmnet-host", "mac": "52:54:00:2d:6e:98"}
      ],
      "type": "custom",
      "status": "running"
    }
  ],
  "count": 1,
  "last_updated": "2026-02-15 12:30:45"
}
```

## Configuration

The application runs on `0.0.0.0:5450` by default. To change the port, modify the `addr` variable in `main.go`:

```go
addr := "0.0.0.0:YOUR_PORT"
```

## Architecture

- **Backend**: Go HTTP server with periodic polling
- **Frontend**: Vanilla JavaScript with modern CSS
- **Update Interval**: 5 seconds (configurable in code)
- **Data Source**: `ps -ef` command output

## Customization

### Change Update Interval

In `main.go`, modify the ticker duration:

```go
ticker := time.NewTicker(5 * time.Second) // Change to your preference
```

### Modify UI Theme

Edit the CSS variables in `html.go`:

```css
:root {
    --bg-primary: #0a0e14;
    --accent-green: #00ff88;
    /* ... more variables */
}
```

## Troubleshooting

### No instances detected

- Ensure QEMU instances are running
- Check that you have permission to run `ps -ef`
- Verify instances show up in: `ps -ef | grep qemu-system-aarch64`

### Port already in use

Change the port in `main.go` or stop the conflicting service.

### Page not loading

- Check if the server started successfully
- Look for errors in the console output
- Ensure firewall allows connections to port 5450

## Development

### Project Structure

```
qemu-monitor/
├── main.go      # Backend logic and HTTP server
├── html.go      # Embedded HTML/CSS/JS UI
├── go.mod       # Go module definition
└── README.md    # This file
```

### Adding New Features

1. **Backend**: Modify parsing logic in `parseQEMUProcess()`
2. **API**: Add new fields to `QEMUInstance` struct
3. **Frontend**: Update `createInstanceCard()` to display new fields

## License

MIT License - Feel free to use and modify as needed.

## Contributing

This is a standalone tool for personal use. Feel free to fork and customize for your needs!

## Credits

Built for RDK-B development workflow monitoring.
