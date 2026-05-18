#!/usr/bin/env python3
"""
Mock pond sensor and aerator script for the aquaculture-monitoring skill.

Usage:
    python3 mock-pond.py --pond pond-01
    python3 mock-pond.py --pond pond-01 --low-do
    python3 mock-pond.py --pond pond-01 --bad-ph
    python3 mock-pond.py --pond pond-01 --hot
    python3 mock-pond.py --all --summary
    python3 mock-pond.py --aerator on --pond pond-01
    python3 mock-pond.py --daily-report
"""

import argparse
import json
import random
import sys
from datetime import datetime, timedelta, timezone

PONDS = ["pond-01", "pond-02", "pond-03"]

POND_BASELINES = {
    "pond-01": {"do": 5.9, "ph": 7.4, "temp": 27.0},
    "pond-02": {"do": 6.4, "ph": 7.8, "temp": 26.2},
    "pond-03": {"do": 5.5, "ph": 7.1, "temp": 28.1},
}

THRESHOLDS = {
    "do_warning": 5.0,
    "do_critical": 4.0,
    "do_emergency": 3.0,
    "ph_warning_low": 6.5,
    "ph_warning_high": 8.5,
    "ph_critical_low": 6.0,
    "ph_critical_high": 9.0,
    "temp_warning_low": 18.0,
    "temp_warning_high": 32.0,
    "temp_critical_low": 12.0,
    "temp_critical_high": 35.0,
}


def utc_now():
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def classify(reading):
    reasons = []
    status = "normal"

    dissolved_oxygen = reading["dissolved_oxygen_mg_l"]
    ph = reading["ph"]
    temperature = reading["water_temp_c"]

    if dissolved_oxygen < THRESHOLDS["do_emergency"]:
        status = "emergency"
        reasons.append("dissolved oxygen below emergency threshold")
    elif dissolved_oxygen < THRESHOLDS["do_critical"]:
        status = "critical"
        reasons.append("dissolved oxygen below critical threshold")
    elif dissolved_oxygen < THRESHOLDS["do_warning"]:
        status = "warning"
        reasons.append("dissolved oxygen below warning threshold")

    if ph < THRESHOLDS["ph_critical_low"] or ph > THRESHOLDS["ph_critical_high"]:
        status = max_status(status, "critical")
        reasons.append("pH outside critical range")
    elif ph < THRESHOLDS["ph_warning_low"] or ph > THRESHOLDS["ph_warning_high"]:
        status = max_status(status, "warning")
        reasons.append("pH outside warning range")

    if temperature < THRESHOLDS["temp_critical_low"] or temperature > THRESHOLDS["temp_critical_high"]:
        status = max_status(status, "critical")
        reasons.append("water temperature outside critical range")
    elif temperature < THRESHOLDS["temp_warning_low"] or temperature > THRESHOLDS["temp_warning_high"]:
        status = max_status(status, "warning")
        reasons.append("water temperature outside warning range")

    reading["status"] = status
    reading["reasons"] = reasons
    reading["recommended_action"] = recommend_action(reading)
    return reading


def max_status(left, right):
    order = {"normal": 0, "warning": 1, "critical": 2, "emergency": 3}
    return left if order[left] >= order[right] else right


def recommend_action(reading):
    dissolved_oxygen = reading["dissolved_oxygen_mg_l"]
    if dissolved_oxygen < THRESHOLDS["do_emergency"]:
        return "turn aerator on immediately, alert all channels, re-check every 10 minutes"
    if dissolved_oxygen < THRESHOLDS["do_critical"]:
        return "turn aerator on and notify pond manager"
    if dissolved_oxygen < THRESHOLDS["do_warning"]:
        return "watch DO trend and consider preventive aeration before dawn"
    if reading["status"] in ("warning", "critical"):
        return "inspect water quality and confirm with a second reading"
    return "append reading to Memory"


def read_pond(pond_id, low_do=False, bad_ph=False, hot=False, chaos=False):
    if pond_id not in POND_BASELINES:
        raise ValueError(f"unknown pond: {pond_id}")

    baseline = POND_BASELINES[pond_id]
    dissolved_oxygen = baseline["do"] + random.gauss(0, 0.25)
    ph = baseline["ph"] + random.gauss(0, 0.08)
    temperature = baseline["temp"] + random.gauss(0, 0.35)

    if chaos:
        roll = random.random()
        if roll < 0.20:
            low_do = True
        elif roll < 0.30:
            bad_ph = True
        elif roll < 0.40:
            hot = True

    if low_do:
        dissolved_oxygen = random.uniform(2.6, 4.3)
    if bad_ph:
        ph = random.choice([random.uniform(5.6, 6.2), random.uniform(8.8, 9.4)])
    if hot:
        temperature = random.uniform(32.5, 36.0)

    reading = {
        "pond": pond_id,
        "dissolved_oxygen_mg_l": round(dissolved_oxygen, 2),
        "ph": round(ph, 2),
        "water_temp_c": round(temperature, 1),
        "aerator": "off",
        "thresholds": THRESHOLDS,
        "timestamp": utc_now(),
    }
    return classify(reading)


def read_all_ponds(chaos=False):
    return [read_pond(pond_id, chaos=chaos) for pond_id in PONDS]


def set_aerator(pond_id, action):
    if pond_id not in POND_BASELINES:
        raise ValueError(f"unknown pond: {pond_id}")
    if action not in {"on", "off", "status"}:
        raise ValueError("aerator action must be on, off, or status")

    state = "on" if action == "on" else "off"
    if action == "status":
        state = random.choice(["on", "off", "off"])

    return {
        "pond": pond_id,
        "aerator": state,
        "action": action,
        "timestamp": utc_now(),
        "note": "mock relay result; replace with GPIO/contactor command in production",
    }


def daily_report():
    end = datetime.now(timezone.utc)
    start = end - timedelta(hours=24)
    rows = []
    for pond_id, baseline in POND_BASELINES.items():
        min_do = round(baseline["do"] - random.uniform(0.9, 1.7), 2)
        avg_do = round(baseline["do"] - random.uniform(0.1, 0.4), 2)
        max_temp = round(baseline["temp"] + random.uniform(1.5, 4.2), 1)
        aerator_minutes = random.choice([0, 20, 35, 50, 75])
        rows.append({
            "pond": pond_id,
            "min_do_mg_l": min_do,
            "avg_do_mg_l": avg_do,
            "max_water_temp_c": max_temp,
            "aerator_runtime_minutes": aerator_minutes,
            "recommendation": "start preventive aeration at 04:15" if min_do < 4.8 else "no schedule change",
        })

    return {
        "report_type": "aquaculture_daily_report",
        "window_start": start.isoformat().replace("+00:00", "Z"),
        "window_end": end.isoformat().replace("+00:00", "Z"),
        "ponds": rows,
    }


def print_summary(readings):
    print(f"\n[AQUA] Pond Readings -- {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M:%S')} UTC\n")
    print(f"{'Pond':<9} {'DO mg/L':<9} {'pH':<6} {'Temp C':<8} {'Status':<10} Action")
    print("-" * 78)
    for reading in readings:
        marker = ""
        if reading["status"] == "warning":
            marker = " [!]"
        elif reading["status"] == "critical":
            marker = " [!!]"
        elif reading["status"] == "emergency":
            marker = " [!!!]"
        print(
            f"{reading['pond']:<9} "
            f"{reading['dissolved_oxygen_mg_l']:<9} "
            f"{reading['ph']:<6} "
            f"{reading['water_temp_c']:<8} "
            f"{reading['status'] + marker:<10} "
            f"{reading['recommended_action']}"
        )


def main():
    parser = argparse.ArgumentParser(description="Mock aquaculture pond sensor for PicClaw")
    parser.add_argument("--pond", default="pond-01", help="Pond ID")
    parser.add_argument("--all", action="store_true", help="Read all ponds")
    parser.add_argument("--low-do", action="store_true", help="Simulate low dissolved oxygen")
    parser.add_argument("--bad-ph", action="store_true", help="Simulate unsafe pH")
    parser.add_argument("--hot", action="store_true", help="Simulate high water temperature")
    parser.add_argument("--chaos", action="store_true", help="Inject random anomalies")
    parser.add_argument("--summary", action="store_true", help="Print human-readable summary")
    parser.add_argument("--aerator", choices=["on", "off", "status"], help="Mock aerator relay action")
    parser.add_argument("--daily-report", action="store_true", help="Generate mock daily report")
    args = parser.parse_args()

    try:
        if args.daily_report:
            print(json.dumps(daily_report(), indent=2))
            return

        if args.aerator:
            print(json.dumps(set_aerator(args.pond, args.aerator), indent=2))
            return

        if args.all:
            readings = read_all_ponds(chaos=args.chaos)
            if args.summary:
                print_summary(readings)
            else:
                print(json.dumps(readings, indent=2))
            return

        reading = read_pond(args.pond, low_do=args.low_do, bad_ph=args.bad_ph, hot=args.hot)
        print(json.dumps(reading, indent=2))
    except ValueError as error:
        print(json.dumps({"error": str(error), "available_ponds": PONDS}))
        sys.exit(1)


if __name__ == "__main__":
    main()
