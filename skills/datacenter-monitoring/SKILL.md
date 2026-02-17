---
name: datacenter-monitoring
description: Monitor datacenter rack temperatures with cross-sensor validation, anomaly detection, and HVAC status checking. Demonstrates Gene Evolution (CGEP) in action.
homepage: https://github.com/Clawland-AI/picclaw
metadata: {"nanobot":{"emoji":"🏢","requires":{"bins":["python3"]}}}
---

# Datacenter Monitoring

End-to-end datacenter temperature monitoring skill for PicoClaw. This is the reference implementation demonstrating the Gene Evolution Protocol (CGEP) in a real-world scenario.

## Overview

This skill monitors server rack temperatures using a mock sensor script (or real sensors via I2C/GPIO), cross-validates readings across adjacent racks, detects anomalies, and evolves monitoring strategies over time through the Gene Evolution system.

## Architecture

```
  Rack Sensor A ─┐
  Rack Sensor B ─┤── picclaw agent ── Gene Evolution ── Fleet (optional)
  Rack Sensor C ─┘       │
                     HVAC status
                          │
                     Alert channels
                     (Telegram/Discord/etc)
```

## How It Works

### 1. Periodic Temperature Check

Set up a cron job to read rack temperatures every 5 minutes:

```bash
# Using the mock sensor (development/testing)
python3 scripts/mock-sensor.py --rack A1

# Using real DS18B20 sensors (production on Linux)
python3 scripts/read-ds18b20.py --device /sys/bus/w1/devices/28-*/w1_slave
```

### 2. Cross-Sensor Validation (Gene: `gene_dc_temp_spike_crosscheck`)

When a temperature exceeds the threshold (default: 35°C):

1. The agent reads the alerting rack sensor
2. Reads 2 adjacent rack sensors for cross-verification
3. If all 3 sensors agree: confirmed hot spot → escalate alert
4. If only 1 sensor is elevated: suspect sensor fault → flag for inspection
5. Checks HVAC status and cooling logs

### 3. Gene Evolution in Action

Over time, the agent learns:

- **Typical patterns**: "Rack A1 always runs 2°C hotter after 2 PM due to sun exposure" → adjusts threshold
- **False positives**: "Sensor B3 drifts high every Tuesday during UPS testing" → adds time-based exception
- **Novel situations**: "All three racks rose together during power blip" → creates new compound alert gene

Each learned pattern becomes a Gene that can be shared with other picclaw instances monitoring different facilities.

## Setup with Cron

Add to picoclaw's cron configuration:

```json
{
  "cron_jobs": [
    {
      "name": "dc_temp_check",
      "schedule": "*/5 * * * *",
      "command": "Read datacenter temperatures from all rack sensors using scripts/mock-sensor.py. Compare readings against 35°C threshold. If any sensor exceeds threshold, cross-check with adjacent sensors. Report findings and take appropriate action based on Gene Evolution strategies."
    }
  ]
}
```

## Mock Sensor Script

The `scripts/mock-sensor.py` script simulates datacenter conditions:

```bash
# Normal reading
python3 scripts/mock-sensor.py --rack A1
# Output: {"rack":"A1","temperature":24.3,"humidity":45.2,"timestamp":"...","status":"normal"}

# Simulate spike
python3 scripts/mock-sensor.py --rack A1 --spike
# Output: {"rack":"A1","temperature":38.7,"humidity":52.1,"timestamp":"...","status":"warning"}

# Simulate sensor failure
python3 scripts/mock-sensor.py --rack A1 --fail
# Output: {"rack":"A1","error":"sensor read failure: I2C timeout","timestamp":"...","status":"error"}

# All racks at once
python3 scripts/mock-sensor.py --all
# Output: [{"rack":"A1",...},{"rack":"A2",...},{"rack":"B1",...}]
```

## Expected Gene Evolution Behavior

After running for several days, you should see:

1. **Day 1-3**: Agent uses the seed gene `gene_dc_temp_spike_crosscheck` (confidence: 70%)
2. **Day 3-7**: After successful cross-check validations, confidence rises to 85%+
3. **Week 2+**: Agent creates new genes for learned time-of-day patterns
4. **Week 3+**: New genes for specific rack behaviors (if some racks consistently run hotter)

Check evolution progress:

```bash
picoclaw gene list    # See all genes and confidence scores
picoclaw gene stats   # See drift intensity and capsule count
```

## Hardware Requirements (Production)

For real deployment:

| Component | Model | Cost |
|-----------|-------|------|
| Edge board | LicheeRV Nano | $10 |
| Temp sensor | DS18B20 (x3) | $6 |
| Humidity sensor | DHT22 | $5 |
| Relay module | 4-channel | $3 |
| USB power | 5V 2A | $4 |
| **Total** | | **$28** |

For mock/development: just Python 3.6+.
