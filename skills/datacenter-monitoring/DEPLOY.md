# Datacenter Monitoring — Deployment Guide

## Quick Start (Mock Sensors)

### 1. Verify mock sensor works

```bash
cd skills/datacenter-monitoring
python3 scripts/mock-sensor.py --all --summary
```

Expected output:

```
🏢 Datacenter Sensor Readings — 2026-02-17 12:00:00 UTC

Rack   Temp (°C)    Humidity (%)   Status
----------------------------------------------
A1     23.8         44.5           normal
A2     24.1         45.2           normal
A3     24.3         42.8           normal
B1     22.2         46.3           normal
B2     23.1         47.1           normal
B3     25.2         41.8           normal

✅ All racks within normal range.
```

### 2. Test anomaly scenarios

```bash
# Temperature spike on rack A1
python3 scripts/mock-sensor.py --rack A1 --spike

# Sensor failure on rack B2
python3 scripts/mock-sensor.py --rack B2 --fail

# Random chaos (10% failures, 15% spikes)
python3 scripts/mock-sensor.py --all --chaos --summary

# HVAC status
python3 scripts/mock-sensor.py --hvac
```

### 3. Configure picoclaw cron

Add this to `~/.picoclaw/config.json` (or use `picoclaw cron add`):

```json
{
  "cron_jobs": [
    {
      "name": "dc_temp_check",
      "schedule": "*/5 * * * *",
      "command": "Read all datacenter rack temperatures by running 'python3 skills/datacenter-monitoring/scripts/mock-sensor.py --all'. Analyze the results. If any rack shows 'warning' or 'critical' status, cross-check with adjacent rack readings. If this is a confirmed hot spot, alert via the configured messaging channel. If only one sensor shows an issue, flag it as a potential sensor malfunction. Also check HVAC status with '--hvac' flag. Report findings using the report_gene tool to solidify the experience."
    }
  ]
}
```

### 4. Start the gateway

```bash
picoclaw gateway
```

### 5. Monitor gene evolution

```bash
# Check gene pool after a few hours
picoclaw gene list
picoclaw gene stats
```

## Production Deployment (Real Sensors)

### Hardware

| Component | Model | Connection | Notes |
|-----------|-------|------------|-------|
| LicheeRV Nano | RISC-V 1GHz | — | Main edge board |
| DS18B20 x3 | 1-Wire temp | GPIO 4 (default) | One per rack, waterproof probe |
| DHT22 | I2C humidity | GPIO 17 | Ambient humidity |
| 4ch Relay | 5V relay module | GPIO 18-21 | Fan/HVAC control |

### Wiring

```
LicheeRV Nano
├── GPIO4  ─── DS18B20 (Rack A) ─── 4.7kΩ pullup to 3.3V
│               DS18B20 (Rack B) ─── (shared 1-Wire bus)
│               DS18B20 (Rack C) ─── (shared 1-Wire bus)
├── GPIO17 ─── DHT22 data pin ─── 10kΩ pullup to 3.3V
├── GPIO18 ─── Relay CH1 (exhaust fan)
├── GPIO19 ─── Relay CH2 (backup cooling)
├── 3.3V   ─── sensor VCC
├── 5V     ─── relay VCC
└── GND    ─── common ground
```

### Software Setup

```bash
# Enable 1-Wire on Linux
echo "w1-gpio" | sudo tee -a /etc/modules
echo "w1-therm" | sudo tee -a /etc/modules

# Read DS18B20
cat /sys/bus/w1/devices/28-*/w1_slave

# Install picoclaw
wget https://github.com/Clawland-AI/picclaw/releases/latest/download/picclaw-linux-arm64
chmod +x picclaw-linux-arm64
sudo mv picclaw-linux-arm64 /usr/local/bin/picoclaw

# Configure
picoclaw onboard
# Edit ~/.picoclaw/config.json with your LLM API key and channel settings
```

Replace `mock-sensor.py` with `read-ds18b20.py` in the cron command for real sensor readings.

## Gene Evolution Timeline

| Period | Expected Behavior |
|--------|-------------------|
| Day 1-3 | Uses seed gene `gene_dc_temp_spike_crosscheck` at 70% confidence |
| Day 3-7 | Confidence rises to 85%+ after successful cross-checks |
| Week 2 | Creates time-of-day pattern genes (if afternoon temps consistently higher) |
| Week 3 | Creates rack-specific genes (if some racks consistently hotter) |
| Month 2+ | High-confidence genes shared to Fleet for other datacenter nodes |

## Cost Comparison

| Item | Traditional | Clawland |
|------|-------------|----------|
| Hardware | $5,000+ monitoring system | $28 sensor kit |
| Software | $2,000/yr DCIM license | Free (open source) |
| Labor | $48,000/yr night shift operator | $0 (automated) |
| **Annual total** | **$55,000+** | **$28 one-time** |
