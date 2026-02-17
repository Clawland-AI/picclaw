#!/usr/bin/env python3
"""
Mock datacenter sensor script for PicoClaw datacenter-monitoring skill.

Simulates temperature/humidity readings from server rack sensors.
Use this for development, testing, and demonstrating Gene Evolution.

Usage:
    python3 mock-sensor.py --rack A1              # Normal reading
    python3 mock-sensor.py --rack A1 --spike      # Simulate temp spike
    python3 mock-sensor.py --rack A1 --fail       # Simulate sensor failure
    python3 mock-sensor.py --all                   # All racks
    python3 mock-sensor.py --all --chaos           # Random anomalies
"""

import argparse
import json
import random
import sys
from datetime import datetime, timezone

RACKS = ["A1", "A2", "A3", "B1", "B2", "B3"]

# Baseline temperatures per rack (simulates real-world variation)
RACK_BASELINES = {
    "A1": {"temp": 23.5, "humidity": 44.0},
    "A2": {"temp": 24.0, "humidity": 45.0},
    "A3": {"temp": 24.5, "humidity": 43.0},
    "B1": {"temp": 22.0, "humidity": 46.0},
    "B2": {"temp": 23.0, "humidity": 47.0},
    "B3": {"temp": 25.0, "humidity": 42.0},
}

THRESHOLD_TEMP = 35.0
THRESHOLD_HUMIDITY = 70.0


def read_sensor(rack_id, spike=False, fail=False):
    """Simulate a sensor reading for a given rack."""
    now = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

    if fail:
        return {
            "rack": rack_id,
            "error": "sensor read failure: I2C timeout after 3 retries",
            "timestamp": now,
            "status": "error",
        }

    baseline = RACK_BASELINES.get(rack_id, {"temp": 24.0, "humidity": 45.0})

    # Normal jitter: +/- 1.5°C temp, +/- 3% humidity
    temp = baseline["temp"] + random.gauss(0, 0.5)
    humidity = baseline["humidity"] + random.gauss(0, 1.0)

    if spike:
        # Simulate a real temperature spike: +10-18°C above baseline
        temp = baseline["temp"] + random.uniform(10.0, 18.0)
        humidity = baseline["humidity"] + random.uniform(5.0, 15.0)

    temp = round(temp, 1)
    humidity = round(max(0, min(100, humidity)), 1)

    status = "normal"
    if temp > THRESHOLD_TEMP:
        status = "warning"
    if temp > THRESHOLD_TEMP + 5:
        status = "critical"

    return {
        "rack": rack_id,
        "temperature": temp,
        "humidity": humidity,
        "threshold_temp": THRESHOLD_TEMP,
        "threshold_humidity": THRESHOLD_HUMIDITY,
        "timestamp": now,
        "status": status,
    }


def read_all_racks(chaos=False):
    """Read all racks, optionally injecting random anomalies."""
    readings = []
    for rack_id in RACKS:
        spike = False
        fail = False

        if chaos:
            roll = random.random()
            if roll < 0.1:
                fail = True
            elif roll < 0.25:
                spike = True

        readings.append(read_sensor(rack_id, spike=spike, fail=fail))
    return readings


def check_hvac():
    """Simulate HVAC status check."""
    status = random.choice(["running", "running", "running", "running", "standby", "offline"])
    return {
        "hvac_status": status,
        "cooling_mode": "auto" if status == "running" else "manual",
        "set_point": 22.0,
        "timestamp": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
    }


def main():
    parser = argparse.ArgumentParser(description="Mock datacenter sensor for PicoClaw")
    parser.add_argument("--rack", type=str, help="Rack ID (e.g., A1, B2)")
    parser.add_argument("--all", action="store_true", help="Read all racks")
    parser.add_argument("--spike", action="store_true", help="Simulate temperature spike")
    parser.add_argument("--fail", action="store_true", help="Simulate sensor failure")
    parser.add_argument("--chaos", action="store_true", help="Random anomalies across all racks")
    parser.add_argument("--hvac", action="store_true", help="Check HVAC status")
    parser.add_argument("--summary", action="store_true", help="Print human-readable summary")

    args = parser.parse_args()

    if args.hvac:
        result = check_hvac()
        print(json.dumps(result, indent=2))
        return

    if args.all:
        readings = read_all_racks(chaos=args.chaos)
        if args.summary:
            print_summary(readings)
        else:
            print(json.dumps(readings, indent=2))
        return

    if args.rack:
        rack_id = args.rack.upper()
        if rack_id not in RACK_BASELINES:
            print(json.dumps({"error": f"Unknown rack: {rack_id}", "available": RACKS}))
            sys.exit(1)
        result = read_sensor(rack_id, spike=args.spike, fail=args.fail)
        print(json.dumps(result, indent=2))
        return

    # Default: read all racks with summary
    readings = read_all_racks()
    print_summary(readings)


def print_summary(readings):
    """Print a human-readable summary of all rack readings."""
    now_str = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    print(f"\n[DC] Datacenter Sensor Readings -- {now_str} UTC\n")
    print(f"{'Rack':<6} {'Temp (C)':<12} {'Humidity (%)':<14} {'Status'}")
    print("-" * 46)

    warnings = []
    errors = []

    for r in readings:
        rack = r["rack"]
        if "error" in r:
            print(f"{rack:<6} {'ERROR':<12} {'---':<14} {r['status']}")
            errors.append(r)
        else:
            temp = r["temperature"]
            hum = r["humidity"]
            status = r["status"]
            marker = ""
            if status == "warning":
                marker = " [!]"
                warnings.append(r)
            elif status == "critical":
                marker = " [!!!]"
                warnings.append(r)
            print(f"{rack:<6} {temp:<12} {hum:<14} {status}{marker}")

    print()
    if warnings:
        print(f"[WARNING] {len(warnings)} rack(s) above threshold ({THRESHOLD_TEMP}C):")
        for w in warnings:
            print(f"  - {w['rack']}: {w['temperature']}C")
    if errors:
        print(f"[ERROR] {len(errors)} sensor error(s):")
        for e in errors:
            print(f"  - {e['rack']}: {e['error']}")
    if not warnings and not errors:
        print("[OK] All racks within normal range.")


if __name__ == "__main__":
    main()
