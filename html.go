package main

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>QEMU Instance Monitor</title>
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Orbitron:wght@700;900&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-primary: #0a0e14;
            --bg-secondary: #131920;
            --bg-card: #1a2027;
            --border-color: #2a3441;
            --accent-green: #00ff88;
            --accent-amber: #ffaa00;
            --accent-red: #ff4444;
            --accent-blue: #00aaff;
            --text-primary: #e6edf3;
            --text-secondary: #8b949e;
            --text-dim: #525960;
            --glow-green: rgba(0, 255, 136, 0.3);
            --glow-amber: rgba(255, 170, 0, 0.3);
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'JetBrains Mono', monospace;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
            background-image: 
                repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(0, 255, 136, 0.03) 2px, rgba(0, 255, 136, 0.03) 4px),
                radial-gradient(circle at 20% 80%, rgba(0, 170, 255, 0.05) 0%, transparent 50%),
                radial-gradient(circle at 80% 20%, rgba(255, 170, 0, 0.05) 0%, transparent 50%);
        }

        .header {
            background: var(--bg-secondary);
            border-bottom: 2px solid var(--accent-green);
            padding: 2rem;
            box-shadow: 0 4px 20px rgba(0, 255, 136, 0.1);
            position: sticky;
            top: 0;
            z-index: 100;
            backdrop-filter: blur(10px);
        }

        .header-content {
            max-width: 1400px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 1rem;
        }

        .logo {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .logo-icon {
            width: 48px;
            height: 48px;
            background: linear-gradient(135deg, var(--accent-green), var(--accent-blue));
            border-radius: 8px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 24px;
            font-weight: 900;
            font-family: 'Orbitron', sans-serif;
            box-shadow: 0 4px 16px var(--glow-green);
            animation: pulse 3s ease-in-out infinite;
        }

        @keyframes pulse {
            0%, 100% { box-shadow: 0 4px 16px var(--glow-green); }
            50% { box-shadow: 0 4px 24px var(--glow-green), 0 0 40px var(--glow-green); }
        }

        h1 {
            font-family: 'Orbitron', sans-serif;
            font-size: 2rem;
            font-weight: 900;
            letter-spacing: 1px;
            background: linear-gradient(135deg, var(--accent-green), var(--accent-blue));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }

        .stats {
            display: flex;
            gap: 2rem;
            align-items: center;
            font-size: 0.9rem;
        }

        .stat-item {
            display: flex;
            flex-direction: column;
            align-items: flex-end;
        }

        .stat-label {
            color: var(--text-dim);
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .stat-value {
            color: var(--accent-green);
            font-size: 1.5rem;
            font-weight: 700;
            text-shadow: 0 0 10px var(--glow-green);
        }

        .last-updated {
            color: var(--text-secondary);
            font-size: 0.85rem;
        }

        .last-updated .time {
            color: var(--accent-amber);
            font-weight: 500;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
        }

        .filter-bar {
            display: flex;
            gap: 1rem;
            margin-bottom: 2rem;
            flex-wrap: wrap;
        }

        .filter-btn {
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            color: var(--text-secondary);
            padding: 0.6rem 1.2rem;
            border-radius: 6px;
            cursor: pointer;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.85rem;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
        }

        .filter-btn::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, transparent, var(--accent-green), transparent);
            opacity: 0.1;
            transition: left 0.5s ease;
        }

        .filter-btn:hover::before {
            left: 100%;
        }

        .filter-btn:hover {
            border-color: var(--accent-green);
            color: var(--accent-green);
            transform: translateY(-2px);
        }

        .filter-btn.active {
            background: var(--accent-green);
            color: var(--bg-primary);
            border-color: var(--accent-green);
            font-weight: 700;
            box-shadow: 0 4px 16px var(--glow-green);
        }

        .instances-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
            gap: 1.5rem;
            animation: fadeIn 0.5s ease;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .instance-card {
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 1.5rem;
            transition: all 0.3s ease;
            position: relative;
            overflow: hidden;
            animation: slideIn 0.5s ease backwards;
        }

        .instance-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 3px;
            background: linear-gradient(90deg, var(--accent-green), var(--accent-blue));
            opacity: 0;
            transition: opacity 0.3s ease;
        }

        .instance-card:hover {
            border-color: var(--accent-green);
            transform: translateY(-4px);
            box-shadow: 0 8px 32px rgba(0, 255, 136, 0.15);
        }

        .instance-card:hover::before {
            opacity: 1;
        }

        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateX(-20px);
            }
            to {
                opacity: 1;
                transform: translateX(0);
            }
        }

        .instance-header {
            display: flex;
            justify-content: space-between;
            align-items: start;
            margin-bottom: 1rem;
            gap: 1rem;
        }

        .instance-name {
            font-size: 1.1rem;
            font-weight: 700;
            color: var(--text-primary);
            word-break: break-word;
        }

        .status-badge {
            padding: 0.3rem 0.8rem;
            border-radius: 20px;
            font-size: 0.7rem;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 1px;
            white-space: nowrap;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
        }

        .status-running {
            background: var(--accent-green);
            color: var(--bg-primary);
            box-shadow: 0 2px 12px var(--glow-green);
        }

        .status-suspended {
            background: var(--accent-amber);
            color: var(--bg-primary);
            box-shadow: 0 2px 12px var(--glow-amber);
        }

        .status-snapshot {
            background: var(--accent-blue);
            color: var(--bg-primary);
        }

        .instance-type {
            display: inline-block;
            padding: 0.2rem 0.6rem;
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 4px;
            font-size: 0.7rem;
            color: var(--text-dim);
            text-transform: uppercase;
            letter-spacing: 1px;
            margin-bottom: 1rem;
        }

        .instance-details {
            display: grid;
            gap: 0.8rem;
        }

        .detail-row {
            display: grid;
            grid-template-columns: 120px 1fr;
            gap: 1rem;
            align-items: start;
        }

        .detail-label {
            color: var(--text-dim);
            font-size: 0.75rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .detail-value {
            color: var(--text-primary);
            font-size: 0.85rem;
            word-break: break-all;
        }

        .detail-value.mono {
            font-family: 'JetBrains Mono', monospace;
            color: var(--accent-green);
        }

        .networks {
            display: flex;
            flex-direction: column;
            gap: 0.4rem;
        }

        .network-item {
            background: var(--bg-secondary);
            padding: 0.5rem;
            border-radius: 4px;
            border-left: 2px solid var(--accent-blue);
            font-size: 0.75rem;
        }

        .network-type {
            color: var(--accent-blue);
            font-weight: 700;
            text-transform: uppercase;
            font-size: 0.7rem;
            margin-bottom: 0.2rem;
        }

        .network-mac {
            color: var(--text-secondary);
            font-family: 'JetBrains Mono', monospace;
        }

        .actions {
            display: flex;
            gap: 0.5rem;
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid var(--border-color);
        }

        .action-btn {
            flex: 1;
            padding: 0.6rem 1rem;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            background: var(--bg-secondary);
            color: var(--text-primary);
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
            cursor: pointer;
            transition: all 0.3s ease;
            letter-spacing: 0.5px;
        }

        .action-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
        }

        .action-btn.start {
            border-color: var(--accent-green);
            color: var(--accent-green);
        }

        .action-btn.start:hover {
            background: var(--accent-green);
            color: var(--bg-primary);
        }

        .action-btn.stop {
            border-color: var(--accent-red);
            color: var(--accent-red);
        }

        .action-btn.stop:hover {
            background: var(--accent-red);
            color: var(--bg-primary);
        }

        .action-btn.force-stop {
            border-color: var(--accent-red);
            color: var(--accent-red);
            background: rgba(255, 68, 68, 0.1);
        }

        .action-btn.force-stop:hover {
            background: var(--accent-red);
            color: var(--bg-primary);
        }

        .action-btn.shell {
            border-color: var(--accent-blue);
            color: var(--accent-blue);
        }

        .action-btn.shell:hover {
            background: var(--accent-blue);
            color: var(--bg-primary);
        }

        .action-btn:disabled {
            opacity: 0.3;
            cursor: not-allowed;
            transform: none;
        }

        .action-btn:disabled:hover {
            transform: none;
            box-shadow: none;
            background: var(--bg-secondary);
        }

        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.8);
            z-index: 1000;
            align-items: center;
            justify-content: center;
        }

        .modal.show {
            display: flex;
        }

        .modal-content {
            background: var(--bg-card);
            border: 2px solid var(--accent-green);
            border-radius: 12px;
            padding: 2rem;
            max-width: 500px;
            width: 90%;
            box-shadow: 0 8px 32px rgba(0, 255, 136, 0.3);
        }

        .modal-header {
            font-size: 1.3rem;
            font-weight: 700;
            color: var(--accent-green);
            margin-bottom: 1rem;
        }

        .modal-body {
            margin-bottom: 1.5rem;
            line-height: 1.8;
        }

        .modal-code {
            background: var(--bg-secondary);
            padding: 1rem;
            border-radius: 6px;
            border-left: 3px solid var(--accent-green);
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9rem;
            margin: 1rem 0;
            overflow-x: auto;
            color: var(--accent-green);
        }

        .modal-actions {
            display: flex;
            gap: 1rem;
            justify-content: flex-end;
        }

        .modal-btn {
            padding: 0.7rem 1.5rem;
            border: 1px solid var(--border-color);
            border-radius: 6px;
            background: var(--bg-secondary);
            color: var(--text-primary);
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.85rem;
            font-weight: 700;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .modal-btn.primary {
            background: var(--accent-green);
            border-color: var(--accent-green);
            color: var(--bg-primary);
        }

        .modal-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
        }

        .no-instances {
            text-align: center;
            padding: 4rem 2rem;
            color: var(--text-dim);
            font-size: 1.2rem;
        }

        .loading {
            text-align: center;
            padding: 4rem 2rem;
            color: var(--accent-green);
            font-size: 1.2rem;
            animation: pulse-text 1.5s ease-in-out infinite;
        }

        @keyframes pulse-text {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        @media (max-width: 768px) {
            .instances-grid {
                grid-template-columns: 1fr;
            }
            
            .stats {
                flex-direction: column;
                align-items: flex-start;
                gap: 1rem;
            }
            
            .detail-row {
                grid-template-columns: 1fr;
                gap: 0.3rem;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="header-content">
            <div class="logo">
                <div class="logo-icon">Q</div>
                <h1>QEMU Monitor</h1>
            </div>
            <div class="stats">
                <div class="stat-item">
                    <span class="stat-label">Instances</span>
                    <span class="stat-value" id="instance-count">0</span>
                </div>
                <div class="stat-item">
                    <span class="last-updated">Updated: <span class="time" id="last-updated">--</span></span>
                </div>
            </div>
        </div>
    </div>

    <div class="container">
        <div class="filter-bar">
            <button class="filter-btn active" data-filter="all">All</button>
            <button class="filter-btn" data-filter="running">Running</button>
            <button class="filter-btn" data-filter="suspended">Suspended</button>
            <button class="filter-btn" data-filter="multipass">Multipass</button>
            <button class="filter-btn" data-filter="custom">Custom</button>
        </div>

        <div id="instances-container">
            <div class="loading">⟳ Loading instances...</div>
        </div>
    </div>

    <div id="shell-modal" class="modal">
        <div class="modal-content">
            <div class="modal-header">Shell Access</div>
            <div class="modal-body" id="shell-modal-body"></div>
            <div class="modal-actions">
                <button class="modal-btn" onclick="closeShellModal()">Close</button>
            </div>
        </div>
    </div>

    <script>
        let currentFilter = 'all';
        let allInstances = [];
        let vmsConfig = { vms: [] };

        function formatUptime(cpuTime) {
            if (!cpuTime) return 'N/A';
            return cpuTime;
        }

        async function loadVMsConfig() {
            try {
                const response = await fetch('/api/vms');
                vmsConfig = await response.json();
            } catch (error) {
                console.error('Failed to load VMs config:', error);
            }
        }

        function findVMConfig(name) {
            return vmsConfig.vms.find(function(vm) { return vm.name === name; });
        }

        function isVMRunning(name) {
            return allInstances.some(function(inst) { return inst.name === name; });
        }

        async function startVM(name) {
            if (!confirm('Start VM: ' + name + '?')) return;
            
            try {
                const response = await fetch('/api/start', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: name })
                });
                const data = await response.json();
                
                if (data.error) {
                    alert('Error: ' + data.error);
                } else {
                    alert('VM started successfully! Refreshing...');
                    setTimeout(fetchInstances, 2000);
                }
            } catch (error) {
                alert('Failed to start VM: ' + error.message);
            }
        }

        async function stopVM(pid, name, force) {
            if (!force) {
                if (!confirm('Stop VM: ' + name + ' (PID: ' + pid + ')?\n\nThis will send SIGTERM for graceful shutdown.')) {
                    return;
                }
            }
            
            try {
                const response = await fetch('/api/stop', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ pid: pid, force: force || false })
                });
                const data = await response.json();
                
                if (data.error) {
                    if (data.error.includes('operation not permitted') || data.error.includes('permission denied')) {
                        alert('Permission denied. You may need to enter your sudo password in the terminal running the server.');
                    } else {
                        alert('Error: ' + data.error);
                    }
                } else {
                    alert('VM ' + data.status + ' successfully! Refreshing...');
                    setTimeout(fetchInstances, 2000);
                }
            } catch (error) {
                alert('Failed to stop VM: ' + error.message);
            }
        }

        async function forceStopVM(pid, name) {
            if (!confirm('⚠️ FORCE STOP VM: ' + name + ' (PID: ' + pid + ')?\n\nThis will send SIGKILL and immediately terminate the VM.\nUse this only if graceful shutdown failed.\n\nContinue?')) {
                return;
            }
            await stopVM(pid, name, true);
        }

        async function showShell(name) {
            try {
                const response = await fetch('/api/shell', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: name })
                });
                const data = await response.json();
                
                if (data.error) {
                    alert('Error: ' + data.error);
                    return;
                }

                let html = '<div><strong>VM:</strong> ' + name + '</div>';
                
                if (!data.running) {
                    html += '<div style="color: var(--accent-amber); margin-top: 1rem;">VM is not currently running</div>';
                }
                
                if (data.ssh_command) {
                    html += '<div class="modal-code">' + data.ssh_command + '</div>';
                }
                
                if (data.http_url) {
                    html += '<div style="margin-top: 1rem;"><strong>HTTP:</strong></div>';
                    html += '<div class="modal-code">' + data.http_url + '</div>';
                }

                if (!data.ssh_command && !data.http_url) {
                    html += '<div style="color: var(--text-dim); margin-top: 1rem;">No shell access configured for this VM</div>';
                }

                document.getElementById('shell-modal-body').innerHTML = html;
                document.getElementById('shell-modal').classList.add('show');
            } catch (error) {
                alert('Failed to get shell info: ' + error.message);
            }
        }

        function closeShellModal() {
            document.getElementById('shell-modal').classList.remove('show');
        }

        function createInstanceCard(instance, index) {
            const statusClass = instance.status.toLowerCase();
            const animationDelay = index * 0.05;
            
            let networksHtml = '';
            if (instance.networks && instance.networks.length > 0) {
                networksHtml = instance.networks.map(function(net) {
                    return '<div class="network-item">' +
                           '<div class="network-type">' + (net.type || 'unknown') + '</div>' +
                           '<div class="network-mac">' + net.mac + '</div>' +
                           '</div>';
                }).join('');
            } else {
                networksHtml = '<div style="color: var(--text-dim); font-size: 0.75rem;">No networks</div>';
            }

            return '<div class="instance-card" style="animation-delay: ' + animationDelay + 's" data-type="' + instance.type + '" data-status="' + instance.status + '">' +
                   '<div class="instance-header">' +
                   '<div class="instance-name">' + (instance.name || 'Unnamed Instance') + '</div>' +
                   '<div class="status-badge status-' + statusClass + '">' + instance.status + '</div>' +
                   '</div>' +
                   '<div class="instance-type">' + instance.type + '</div>' +
                   '<div class="instance-details">' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">PID</div>' +
                   '<div class="detail-value mono">' + instance.pid + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">Memory</div>' +
                   '<div class="detail-value mono">' + (instance.memory || 'N/A') + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">CPU Cores</div>' +
                   '<div class="detail-value mono">' + (instance.cpu_count || 'N/A') + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">CPU Time</div>' +
                   '<div class="detail-value mono">' + formatUptime(instance.uptime) + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">Machine</div>' +
                   '<div class="detail-value">' + (instance.machine || 'N/A') + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">Disk Image</div>' +
                   '<div class="detail-value mono">' + (instance.disk_image || 'N/A') + '</div>' +
                   '</div>' +
                   '<div class="detail-row">' +
                   '<div class="detail-label">Networks</div>' +
                   '<div class="detail-value">' +
                   '<div class="networks">' + networksHtml + '</div>' +
                   '</div>' +
                   '</div>' +
                   '</div>' +
                   '<div class="actions">' +
                   '<button class="action-btn stop" onclick="stopVM(\'' + instance.pid + '\', \'' + (instance.name || 'VM') + '\', false)" title="Graceful shutdown (SIGTERM)">Stop</button>' +
                   '<button class="action-btn force-stop" onclick="forceStopVM(\'' + instance.pid + '\', \'' + (instance.name || 'VM') + '\')" title="Force kill (SIGKILL)">Kill</button>' +
                   '<button class="action-btn shell" onclick="showShell(\'' + (instance.name || '') + '\')">Shell</button>' +
                   '</div>' +
                   '</div>';
        }

        function renderInstances(instances) {
            allInstances = instances;
            const container = document.getElementById('instances-container');
            
            let filtered = instances;
            if (currentFilter !== 'all') {
                if (currentFilter === 'running' || currentFilter === 'suspended') {
                    filtered = instances.filter(function(i) { return i.status === currentFilter; });
                } else {
                    filtered = instances.filter(function(i) { return i.type === currentFilter; });
                }
            }
            
            let html = '';
            
            // Show available VMs that can be started
            const runningNames = instances.map(function(i) { return i.name; });
            const availableVMs = vmsConfig.vms.filter(function(vm) {
                return !runningNames.includes(vm.name);
            });

            if (availableVMs.length > 0 && currentFilter === 'all') {
                html += '<div style="margin-bottom: 2rem;">';
                html += '<h2 style="color: var(--text-dim); font-size: 1rem; margin-bottom: 1rem; text-transform: uppercase; letter-spacing: 1px;">Available VMs</h2>';
                html += '<div class="instances-grid">';
                availableVMs.forEach(function(vm, index) {
                    const animationDelay = index * 0.05;
                    html += '<div class="instance-card" style="animation-delay: ' + animationDelay + 's; opacity: 0.6;">';
                    html += '<div class="instance-header">';
                    html += '<div class="instance-name">' + vm.name + '</div>';
                    html += '<div class="status-badge" style="background: var(--text-dim);">STOPPED</div>';
                    html += '</div>';
                    html += '<div class="instance-type">configured</div>';
                    html += '<div class="instance-details">';
                    html += '<div class="detail-row"><div class="detail-label">Memory</div><div class="detail-value mono">' + vm.memory + 'M</div></div>';
                    html += '<div class="detail-row"><div class="detail-label">CPU Cores</div><div class="detail-value mono">' + vm.cpus + '</div></div>';
                    html += '<div class="detail-row"><div class="detail-label">Disk</div><div class="detail-value mono">' + vm.disk + '</div></div>';
                    html += '</div>';
                    html += '<div class="actions">';
                    html += '<button class="action-btn start" onclick="startVM(\'' + vm.name + '\')">Start</button>';
                    html += '<button class="action-btn shell" onclick="showShell(\'' + vm.name + '\')">Shell Info</button>';
                    html += '</div>';
                    html += '</div>';
                });
                html += '</div></div>';
            }

            if (filtered.length === 0 && availableVMs.length === 0) {
                container.innerHTML = '<div class="no-instances">No instances found</div>';
                return;
            }

            if (filtered.length > 0) {
                if (availableVMs.length > 0 && currentFilter === 'all') {
                    html += '<h2 style="color: var(--text-dim); font-size: 1rem; margin-bottom: 1rem; text-transform: uppercase; letter-spacing: 1px;">Running Instances</h2>';
                }
                html += '<div class="instances-grid">';
                html += filtered.map(createInstanceCard).join('');
                html += '</div>';
            }
            
            container.innerHTML = html;
        }

        async function fetchInstances() {
            try {
                const response = await fetch('/api/instances');
                const data = await response.json();
                
                document.getElementById('instance-count').textContent = data.count;
                document.getElementById('last-updated').textContent = data.last_updated;
                
                renderInstances(data.instances);
            } catch (error) {
                console.error('Error fetching instances:', error);
                document.getElementById('instances-container').innerHTML = 
                    '<div class="no-instances">Error loading instances</div>';
            }
        }

        // Close modal on outside click
        document.getElementById('shell-modal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeShellModal();
            }
        });

        // Filter functionality
        document.querySelectorAll('.filter-btn').forEach(function(btn) {
            btn.addEventListener('click', function() {
                document.querySelectorAll('.filter-btn').forEach(function(b) {
                    b.classList.remove('active');
                });
                btn.classList.add('active');
                currentFilter = btn.dataset.filter;
                renderInstances(allInstances);
            });
        });

        // Initialize
        loadVMsConfig().then(function() {
            fetchInstances();
            setInterval(fetchInstances, 5000);
        });
    </script>
</body>
</html>
`
