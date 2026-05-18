---
name: aquaculture-monitoring
description: Monitor pond dissolved oxygen, pH, and water temperature; control aerators; generate daily reports; and learn pre-dawn oxygen trends through Memory and CGEP genes.
homepage: https://github.com/Clawland-AI/picclaw
metadata: {"nanobot":{"emoji":"🐟","requires":{"bins":["python3"]}}}
---

# Aquaculture Monitoring

Pond monitoring skill for PicoClaw/PicClaw edge nodes. Use it to watch dissolved oxygen (DO), pH, and water temperature, trigger aerator control through relay scripts, send alerts, and record daily Memory so the agent can learn recurring pond patterns.

## Monitoring Contract

Default thresholds:

```yaml
pond:
  check_interval: "*/5 * * * *"
  metrics:
    dissolved_oxygen_mg_l:
      target: "5.5-8.0"
      warning_below: 5.0
      critical_below: 4.0
      emergency_below: 3.0
    ph:
      target: "6.8-8.4"
      warning_low: 6.5
      warning_high: 8.5
      critical_low: 6.0
      critical_high: 9.0
    water_temp_c:
      target: "22-30"
      warning_low: 18.0
      warning_high: 32.0
      critical_low: 12.0
      critical_high: 35.0
  actuators:
    aerator:
      start_when:
        - "dissolved_oxygen_mg_l < 4.0"
        - "dissolved_oxygen_mg_l is falling faster than 0.25 mg/L/hour before dawn"
      stop_when: "dissolved_oxygen_mg_l >= 5.5 for two consecutive checks"
      safety_max_runtime_minutes: 90
```

Operational rules:

- `EMERGENCY`: DO below 3.0 mg/L. Turn aerator on immediately, send all-channel alert, re-check every 10 minutes, and escalate if DO does not recover above 5.0 mg/L within 30 minutes.
- `CRITICAL`: DO below 4.0 mg/L, pH outside 6.0-9.0, or water temp outside 12-35 C. Turn aerator on for low DO, alert the pond manager, and write the event to Memory.
- `WARNING`: DO below 5.0 mg/L, pH outside 6.5-8.5, or water temp outside 18-32 C. Watch trend, recommend inspection, and start preventive aeration if the trend is still falling near dawn.
- `NORMAL`: Append readings to daily Memory and keep the aerator state unchanged.

## Sensor Commands

Development and demos use the bundled mock reader:

```bash
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --all --summary
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --pond pond-01 --low-do
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --pond pond-01 --bad-ph
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --pond pond-01 --hot
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --aerator on --pond pond-01
python3 skills/aquaculture-monitoring/scripts/mock-pond.py --daily-report
```

Production readers should return the same JSON shape:

```json
{
  "pond": "pond-01",
  "dissolved_oxygen_mg_l": 5.8,
  "ph": 7.4,
  "water_temp_c": 27.2,
  "aerator": "off",
  "status": "normal",
  "timestamp": "2026-05-18T04:30:00Z"
}
```

## Cron Setup

Add a 5-minute water quality check:

```json
{
  "cron_jobs": [
    {
      "name": "aquaculture_water_check",
      "schedule": "*/5 * * * *",
      "command": "Run 'python3 skills/aquaculture-monitoring/scripts/mock-pond.py --all' or the production pond sensor reader. Evaluate dissolved oxygen, pH, and water temperature against the aquaculture-monitoring thresholds. If DO is critical or falling before dawn, turn on the aerator with the relay script, alert the configured channel, append the event to Memory, and report the experience using report_gene."
    },
    {
      "name": "aquaculture_daily_report",
      "schedule": "0 7 * * *",
      "command": "Generate the aquaculture daily report from the last 24 hours of pond Memory. Include min/avg/max DO, pH, and water temperature, aerator runtime, alerts, recovery outcomes, and any new pre-dawn trend recommendations."
    }
  ]
}
```

## Memory Trend Learning

Write compact observations to `MEMORY.md` or daily notes:

```markdown
- 2026-05-18 04:40 pond-01: DO fell from 5.2 to 4.1 mg/L between 03:30 and 04:30. Aerator started at 04:35. DO recovered to 5.6 by 05:05. Outcome: preventive aeration at 04:15 likely avoids critical alert.
```

When similar entries repeat, use seed gene `gene_pond_do_night_emergency` and the generic threshold optimization gene:

1. Detect pre-dawn DO drop from recent Memory.
2. Start aerator before DO crosses the critical threshold.
3. Re-check every 10 minutes until recovery.
4. Record runtime, recovery time, and whether the action prevented an alert.
5. Report the successful strategy with `report_gene` so it can be reused on similar ponds.

## Production Hardware

| Component | Example | Purpose |
|-----------|---------|---------|
| Edge board | LicheeRV Nano or Raspberry Pi Zero 2W | Runs PicClaw |
| Dissolved oxygen probe | RS485 optical DO probe | Primary life-safety metric |
| pH probe | Industrial pH electrode + ADC | Water quality guardrail |
| Temperature probe | DS18B20 waterproof | Temperature compensation and heat alerts |
| Relay or contactor | 5V relay plus AC contactor | Aerator control |
| Network | Ethernet, Wi-Fi, or 4G modem | Alerting and fleet sync |
