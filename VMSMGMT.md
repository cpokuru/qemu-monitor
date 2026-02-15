# QEMU Monitor - VM Management Setup Guide

## New Features

The QEMU Monitor now includes full VM lifecycle management:

✅ **Start VMs** - Launch configured VMs with one click
✅ **Stop VMs** - Gracefully shutdown running instances  
✅ **Shell Access** - Get SSH/HTTP connection info
✅ **Auto-Discovery** - See both running and available VMs

## Quick Setup

### 1. Create your VMs configuration file

Copy the example and edit with your VM details:

```bash
cp vms.json.example vms.json
nano vms.json  # or use your preferred editor
```

### 2. Configure Your VMs

Edit `vms.json` with your actual VM configurations. Here's the structure:

```json
{
  "vms": [
    {
      "name": "rdk.snapshot",
      "disk": "rdk.snapshot.qcow2",
      "memory": "6144",
      "cpus": "6",
      "bios": "u-boot.bin",
      "snapshot": false,
      "networks": [
        {
          "type": "vmnet-shared",
          "id": "testwan",
          "mac": "52:54:00:2d:6e:99"
        }
      ],
      "ssh_port": null,
      "http_port": null,
      "working_dir": "/Users/your-name/path/to/qemu-vms"
    }
  ]
}
```

### 3. Set Your Working Directory

**IMPORTANT**: Update `working_dir` to the actual path where your QEMU disk images and u-boot.bin are located:

```json
"working_dir": "/Users/yourname/Documents/qemu-vms"
```

The application will `cd` to this directory before running QEMU commands.

### 4. Configure SSH/HTTP Access (Optional)

If your VM has port forwarding for SSH or HTTP, add the ports:

```json
{
  "name": "RDK-B-Digital-Twin",
  "networks": [
    {
      "type": "user",
      "id": "testlan2",
      "mac": "52:54:00:2d:6e:9c",
      "port_forwards": [
        {"host": 2223, "guest": 22},
        {"host": 8081, "guest": 80}
      ]
    }
  ],
  "ssh_port": 2223,
  "http_port": 8081,
  ...
}
```

## Configuration Examples

### Example 1: RDK-B Snapshot (6GB RAM, 6 CPUs)

```json
{
  "name": "rdk.snapshot",
  "disk": "rdk.snapshot.qcow2",
  "memory": "6144",
  "cpus": "6",
  "bios": "u-boot.bin",
  "networks": [
    {
      "type": "vmnet-shared",
      "id": "testwan",
      "mac": "52:54:00:2d:6e:99"
    },
    {
      "type": "vmnet-host",
      "id": "testlan",
      "mac": "52:54:00:2d:6e:98"
    },
    {
      "type": "user",
      "id": "testlan2",
      "mac": "52:54:00:2d:6e:97"
    }
  ],
  "working_dir": "/Users/chandra/qemu-vms"
}
```

### Example 2: RDK-B Digital Twin (with SSH/HTTP, snapshot mode)

```json
{
  "name": "RDK-B-Digital-Twin",
  "disk": "rdk.twin.qcow2",
  "memory": "1024",
  "cpus": "2",
  "bios": "u-boot.bin",
  "snapshot": true,
  "networks": [
    {
      "type": "vmnet-shared",
      "id": "testwan",
      "mac": "52:54:00:2d:6e:9a"
    },
    {
      "type": "vmnet-host",
      "id": "testlan",
      "mac": "52:54:00:2d:6e:9b"
    },
    {
      "type": "user",
      "id": "testlan2",
      "mac": "52:54:00:2d:6e:9c",
      "port_forwards": [
        {"host": 2223, "guest": 22},
        {"host": 8081, "guest": 80}
      ]
    }
  ],
  "ssh_port": 2223,
  "http_port": 8081,
  "working_dir": "/Users/chandra/qemu-vms"
}
```

## Using the Web Interface

### Starting a VM

1. Configured but stopped VMs appear in the "Available VMs" section
2. Click the **Start** button
3. Confirm the action
4. VM will appear in "Running Instances" after a few seconds

### Stopping a VM

1. Find the running VM
2. Click the **Stop** button
3. Confirm the action
4. VM will gracefully shutdown (SIGTERM)

### Shell Access

1. Click the **Shell** button on any VM (running or stopped)
2. A modal will show:
   - SSH command (if configured): `ssh -p 2223 root@localhost`
   - HTTP URL (if configured): `http://localhost:8081`
3. Copy and paste the command in your terminal

## API Endpoints

The monitor exposes these REST APIs:

### Get All Instances
```bash
curl http://localhost:5450/api/instances
```

### Start a VM
```bash
curl -X POST http://localhost:5450/api/start \
  -H "Content-Type: application/json" \
  -d '{"name": "rdk.snapshot"}'
```

### Stop a VM
```bash
curl -X POST http://localhost:5450/api/stop \
  -H "Content-Type: application/json" \
  -d '{"pid": "12345"}'
```

### Get Shell Info
```bash
curl -X POST http://localhost:5450/api/shell \
  -H "Content-Type: application/json" \
  -d '{"name": "RDK-B-Digital-Twin"}'
```

### Get VM Configurations
```bash
curl http://localhost:5450/api/vms
```

## Permissions

Starting and stopping VMs requires `sudo` permissions. The app will execute:

**Starting:**
```bash
sudo qemu-system-aarch64 [arguments...]
```

**Stopping:**
```bash
sudo kill -TERM <pid>
```

You'll be prompted for your sudo password when performing these operations (unless you've configured passwordless sudo).

### Optional: Configure Passwordless Sudo for QEMU

Add these lines to `/etc/sudoers` (using `sudo visudo`):

```
yourusername ALL=(ALL) NOPASSWD: /usr/local/bin/qemu-system-aarch64
yourusername ALL=(ALL) NOPASSWD: /bin/kill
```

Replace `yourusername` with your actual username and adjust paths if needed.

## Troubleshooting

### VMs don't appear in "Available VMs"

- Check that `vms.json` exists in the same directory as the `qemu-monitor` binary
- Verify JSON syntax is valid: `python3 -m json.tool vms.json`

### Can't start VM

- Check that `working_dir` is correct and contains your disk images
- Verify `bios` file (u-boot.bin) exists in the working directory
- Check disk image file exists: `ls /path/to/working_dir/rdk.snapshot.qcow2`
- Ensure you have sudo permissions

### Shell button shows "No shell access configured"

- Add `ssh_port` to your VM configuration
- Make sure the VM has port forwarding configured in the `networks` section

### VM starts but doesn't appear in the list

- Wait a few seconds - the monitor refreshes every 5 seconds
- Check if the process is actually running: `ps aux | grep qemu-system-aarch64`
- Look at the application logs for errors

## Advanced: Multiple Network Configurations

You can configure complex networking setups:

```json
"networks": [
  {
    "type": "vmnet-shared",
    "id": "wan",
    "mac": "52:54:00:aa:bb:cc"
  },
  {
    "type": "vmnet-host", 
    "id": "lan1",
    "mac": "52:54:00:dd:ee:ff"
  },
  {
    "type": "user",
    "id": "mgmt",
    "mac": "52:54:00:11:22:33",
    "port_forwards": [
      {"host": 2222, "guest": 22},
      {"host": 8080, "guest": 80},
      {"host": 8443, "guest": 443}
    ]
  }
]
```

## Next Steps

1. Edit `vms.json` with your actual VM configurations
2. Set correct `working_dir` paths
3. Build and run: `./start.sh`
4. Open http://localhost:5450
5. Start managing your QEMU VMs! 🚀
