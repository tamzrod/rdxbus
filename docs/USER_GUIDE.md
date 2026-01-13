# RDXBus User Guide

RDXBus is a Modbus TCP client and protocol workbench designed for interactive testing, debugging, and stress testing of Modbus devices and servers.

---

## Quick Start

### Easy Mode (Interactive)

The easiest way to get started is **Easy Mode**, which provides an interactive menu:

```bash
./rdxbus easy
```

This mode guides you through common tasks with simple prompts.

### Expert Mode (CLI Flags)

For advanced users and automation:

```bash
./rdxbus -target 127.0.0.1:502 -fc 3 -address 0 -quantity 10
```

---

## Easy Mode: Complete Feature Guide

When you run `./rdxbus easy`, you'll see an interactive menu with three main sections:

### 1. Read Once

**Purpose:** Perform a single read operation from a Modbus register.

**Steps:**
1. Enter target device address (default: `127.0.0.1:502`)
2. Enter Unit ID (default: `1`)
3. Select "Read once"
4. Choose function code:
   - `1` - Read Coil Status (FC 01)
   - `2` - Read Discrete Inputs (FC 02)
   - `3` - Read Holding Registers (FC 03)
   - `4` - Read Input Registers (FC 04)
5. Enter starting register address
6. Enter number of registers to read
7. See the register values returned

**Example:**
```
Target address [127.0.0.1:502]: 192.168.1.100:502
Unit ID [1]: 1
Selection [1]: 1  (Read once)
Function Code [3]: 3
Starting address [0]: 100
Number of registers [10]: 5
```

### 2. Poll Continuously

**Purpose:** Repeatedly read from the same register to monitor changes over time.

**Steps:**
1. Enter target device address
2. Enter Unit ID
3. Select "Poll continuously"
4. Configure polling parameters:
   - Function code (1-4)
   - Starting address
   - Number of registers
   - Poll interval (in seconds)
5. Polling continues until you press Ctrl+C

**Example:**
```
Target address [127.0.0.1:502]: 127.0.0.1:502
Unit ID [1]: 1
Selection [1]: 2  (Poll continuously)
Function Code [3]: 3
Starting address [0]: 0
Number of registers [10]: 2
Poll interval in seconds [1]: 2
```

The tool will display:
- Timestamp of each poll
- Register values
- Any errors encountered

### 3. Scan Helpers

Advanced discovery tools for identifying responsive devices and registers.

#### 3.1 Find Unit ID

**Purpose:** Automatically discover which Unit IDs are responding on the target device.

**How it works:**
1. Enter target address
2. Tool scans Unit IDs 1-247 (by default, in steps of 50)
3. Reports the first responding Unit ID

**Example:**
```
Target address [127.0.0.1:502]: 192.168.1.100:502
Selection: 1  (Find Unit ID)
[Scanning...]
Found Unit ID: 32
```

#### 3.2 Scan Address Range

**Purpose:** Discover which register addresses respond to read requests.

**How it works:**
1. Enter target address
2. Enter Unit ID
3. Tool scans address range (0-1000 by default, in steps of 10)
4. Reports the first responding address
5. Performs refined search around the found address

**Example:**
```
Target address [127.0.0.1:502]: 192.168.1.100:502
Unit ID [1]: 5
Selection: 2  (Scan address range)
[Scanning...]
First responding address: 450
```

---

## Expert Mode: CLI Flags Reference

For users who prefer command-line automation and scripting.

### Connection Settings

| Flag | Default | Description |
|------|---------|-------------|
| `-target` | `127.0.0.1:502` | Modbus TCP endpoint (host:port) |
| `-timeout` | `100ms` | Socket read/write timeout |

### Modbus Parameters

| Flag | Default | Description |
|------|---------|-------------|
| `-unit` | `1` | Modbus Unit ID (0-247) |
| `-fc` | `3` | Function Code (1=coils, 2=discrete inputs, 3=holding registers, 4=input registers) |
| `-address` | `0` | Starting register address |
| `-quantity` | `10` | Number of registers to read |

### Concurrency & Load Testing

| Flag | Default | Description |
|------|---------|-------------|
| `-workers` | `10` | Number of concurrent workers |
| `-rate` | `0` (unlimited) | Requests per second |
| `-duration` | `10s` | Test duration |

### Ramp Testing

Create multi-step load tests:

| Flag | Default | Description |
|------|---------|-------------|
| `-ramp` | *(none)* | Comma-separated request rates, e.g., `"100,500,1000"` |
| `-step-duration` | `5s` | Duration per ramp step |

### Advanced Options

| Flag | Default | Description |
|------|---------|-------------|
| `-strict` | `false` | Enable strict Modbus TCP framing validation |
| `-quiet` | `false` | Suppress detailed output |

---

## Expert Mode: Usage Examples

### Single Read

Read 5 holding registers starting at address 100:

```bash
./rdxbus -target 192.168.1.100:502 -unit 1 -fc 3 -address 100 -quantity 5
```

### Stress Test

Run 50 concurrent workers at 1000 req/s for 30 seconds:

```bash
./rdxbus -target 192.168.1.100:502 -workers 50 -rate 1000 -duration 30s
```

### Ramp Test

Start at 100 req/s, step up to 500, then 1000 (5 seconds each):

```bash
./rdxbus -target 192.168.1.100:502 -ramp "100,500,1000" -step-duration 5s
```

### Function Code Examples

Read coil statuses (FC 01):
```bash
./rdxbus -target 192.168.1.100:502 -fc 1 -address 0 -quantity 16
```

Read discrete inputs (FC 02):
```bash
./rdxbus -target 192.168.1.100:502 -fc 2 -address 0 -quantity 16
```

Read holding registers (FC 03):
```bash
./rdxbus -target 192.168.1.100:502 -fc 3 -address 0 -quantity 10
```

Read input registers (FC 04):
```bash
./rdxbus -target 192.168.1.100:502 -fc 4 -address 0 -quantity 10
```

---

## Common Workflows

### Discover a Modbus Device

1. Start Easy Mode: `./rdxbus easy`
2. Enter the device's IP and port
3. Select "Scan helpers" → "Find Unit ID"
4. Note the responding Unit ID
5. Select "Scan helpers" → "Scan address range"
6. Note the first responding address

### Monitor a Register Over Time

1. Start Easy Mode: `./rdxbus easy`
2. Select "Poll continuously"
3. Enter function code, address, and quantity
4. Set a reasonable poll interval (e.g., 2-5 seconds)
5. Observe values as they change

### Stress Test a Device

```bash
# Gradually increase load
./rdxbus -target 192.168.1.100:502 \
  -ramp "10,50,100,500" \
  -step-duration 10s \
  -workers 100
```

### Validate Protocol Compliance

```bash
# Use strict mode for strict Modbus TCP validation
./rdxbus -target 192.168.1.100:502 \
  -fc 3 \
  -address 0 \
  -quantity 10 \
  -strict
```

---

## Output Interpretation

### Successful Read

```
read successful
latency: 12.345ms
values: [100 200 300 400 500]
```

### Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| `read error: connection refused` | Device not responding at target address | Check IP/port, verify device is online |
| `read error: i/o timeout` | Device took too long to respond | Increase `-timeout`, check network latency |
| `modbus exception fc=3 code=2` | Modbus exception (code 2 = Illegal Data Address) | Verify address is valid for the device |
| `modbus exception fc=3 code=3` | Modbus exception (code 3 = Illegal Data Value) | Check quantity doesn't exceed device limits |
| `function code mismatch` | Device returned unexpected function code | May indicate protocol issue or device malfunction |

---

## Modbus Background

RDXBus supports Modbus TCP, which uses:

- **Standard Port:** 502 (or 502 + offset for virtual instances)
- **Function Codes:**
  - FC 1: Read Coil Status (bits)
  - FC 2: Read Discrete Inputs (bits)
  - FC 3: Read Holding Registers (16-bit words)
  - FC 4: Read Input Registers (16-bit words)
- **Unit ID:** Device identifier on a Modbus network (0-247, default 1)
- **Transaction ID:** Automatically managed by RDXBus
- **Protocol ID:** Always 0 (Modbus TCP standard)

---

## Tips & Troubleshooting

### Connection Issues

**Problem:** Cannot connect to device
- Verify the IP address and port are correct
- Check network connectivity: `ping 192.168.1.100`
- Confirm the device is running and Modbus TCP is enabled
- Check firewall rules on both ends

### Timeout Issues

**Problem:** Frequent timeouts
- Increase `-timeout` (default is 100ms)
- Check network latency with `ping` or `mtr`
- Verify the device is not overloaded
- Try reducing concurrency (`-workers`)

### Modbus Exceptions

**Problem:** Receiving Modbus exception errors
- Verify the Unit ID matches the device configuration
- Check that addresses are within the device's register range
- Ensure the quantity requested fits in available registers
- Confirm function code matches register type (coils vs. registers)

### Scanning Tips

- Scanning discovers *responsive* devices, not all configured ones
- Unit ID scanning scans 1-247 in steps of 50 by default
- Address scanning finds the first responding address then refines
- Scanning can be slow on high-latency networks

---

## Security Considerations

- RDXBus sends unencrypted Modbus TCP frames
- No authentication is performed (standard Modbus TCP)
- Use network segmentation to protect your Modbus devices
- Modbus TCP should only be used on trusted networks
- For security-critical applications, use VPNs or network isolation

---

## Support & Reporting Issues

RDXBus is a protocol workbench designed for testing and debugging.

- **For protocol issues:** Verify device Modbus TCP compliance
- **For connection problems:** Check network configuration
- **For bugs:** Describe the exact command used and error message

---

## Version & Help

Display command-line help:

```bash
./rdxbus -h
```

Check version (if available):

```bash
./rdxbus -version
```

---

**Happy testing!** 🎯
