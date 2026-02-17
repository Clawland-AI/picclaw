package gene

import "time"

// DefaultSeedGenes returns the 6 seed genes for initializing a new Gene store.
// Three generic genes (repair/optimize/innovate) + three scenario-specific.
func DefaultSeedGenes() []Gene {
	now := time.Now().Format(time.RFC3339)

	return []Gene{
		// === Generic Genes ===
		{
			ID:       "gene_sensor_repair_from_errors",
			Category: "repair",
			Scenario: "generic",
			SignalsMatch: []string{
				SignalSensorError,
			},
			Strategy: []string{
				"Identify which sensor is reporting errors from recent logs",
				"Attempt to re-read the sensor with a 5-second delay",
				"If retry fails, check if adjacent or backup sensors are available",
				"If a backup sensor exists, switch to reading from it",
				"Log the failure and note the sensor ID for maintenance",
				"If no backup available, report the gap in monitoring coverage",
			},
			Constraints: GeneConstraints{
				MaxRetries:     3,
				TimeoutSeconds: 60,
				Requires:       []string{"exec"},
				Forbids:        []string{"reboot"},
			},
			Validation: []ValidationStep{
				{Description: "Sensor read returns valid data after retry", Expected: "non-empty numeric value"},
			},
			Confidence: 0.7,
			VerifiedBy: 0,
			CreatedAt:  now,
		},
		{
			ID:       "gene_threshold_optimize_from_history",
			Category: "optimize",
			Scenario: "generic",
			SignalsMatch: []string{
				SignalThresholdBreach,
				SignalStrategyProven,
			},
			Strategy: []string{
				"Review the last 7 days of sensor readings from daily notes",
				"Calculate mean and standard deviation for the breached metric",
				"If current threshold causes >3 false alarms per day, widen by 1 sigma",
				"If current threshold missed real events, tighten by 0.5 sigma",
				"Update the threshold in MEMORY.md with rationale and date",
				"Monitor for 24 hours and evaluate false-positive rate",
			},
			Constraints: GeneConstraints{
				MaxRetries:     1,
				TimeoutSeconds: 120,
			},
			Validation: []ValidationStep{
				{Description: "New threshold reduces false alarms", Expected: "fewer alerts in next 24h"},
			},
			Confidence: 0.6,
			VerifiedBy: 0,
			CreatedAt:  now,
		},
		{
			ID:       "gene_cross_sensor_innovate",
			Category: "innovate",
			Scenario: "generic",
			SignalsMatch: []string{
				SignalCrossSensorAnomaly,
				SignalNewPatternDetected,
			},
			Strategy: []string{
				"Identify which sensors show correlated anomalies",
				"Check if the correlation is temporal (same time window) or causal",
				"Look for environmental factors that could link the sensors (e.g. shared power, shared location)",
				"If a causal pattern is found, document it as a new monitoring rule",
				"Create a compound alert that triggers only when multiple sensors agree",
				"Log the new pattern in MEMORY.md for future reference",
			},
			Constraints: GeneConstraints{
				MaxRetries:     2,
				TimeoutSeconds: 300,
			},
			Validation: []ValidationStep{
				{Description: "Compound alert reduces false positives", Expected: "correlation confirmed by 2+ events"},
			},
			Confidence: 0.5,
			VerifiedBy: 0,
			CreatedAt:  now,
		},

		// === Scenario-Specific Genes ===
		{
			ID:       "gene_dc_temp_spike_crosscheck",
			Category: "repair",
			Scenario: "datacenter",
			SignalsMatch: []string{
				SignalThresholdBreach,
				SignalCrossSensorAnomaly,
			},
			Strategy: []string{
				"Read temperature from the alerting rack sensor",
				"Read temperature from 2 adjacent rack sensors for cross-verification",
				"If all 3 sensors show elevated temp, confirm real hot spot — escalate alert",
				"If only 1 sensor is elevated, suspect sensor malfunction — flag for inspection",
				"Check HVAC status and cooling system logs",
				"If HVAC is offline, trigger emergency cooling protocol notification",
			},
			Constraints: GeneConstraints{
				MaxRetries:     2,
				TimeoutSeconds: 120,
				Requires:       []string{"exec"},
			},
			Validation: []ValidationStep{
				{Description: "Cross-check reduces false temperature alarms", Expected: "3 sensors agree on reading"},
			},
			Confidence: 0.7,
			VerifiedBy: 0,
			CreatedAt:  now,
		},
		{
			ID:       "gene_pond_do_night_emergency",
			Category: "repair",
			Scenario: "aquaculture",
			SignalsMatch: []string{
				SignalThresholdBreach,
				SignalTimePattern,
			},
			Strategy: []string{
				"Read dissolved oxygen (DO) level from pond sensor",
				"If DO < 3 mg/L and current time is between 22:00-06:00, trigger emergency",
				"Activate aerator immediately via relay control command",
				"Send alert to pond manager via configured channel",
				"Re-check DO level every 10 minutes",
				"If DO recovers above 5 mg/L, log recovery and reduce aerator to normal",
				"If DO does not recover within 30 minutes, escalate to urgent",
			},
			Constraints: GeneConstraints{
				MaxRetries:     5,
				TimeoutSeconds: 1800,
				Requires:       []string{"exec", "message"},
			},
			Validation: []ValidationStep{
				{Description: "DO level recovers after aerator activation", Expected: "DO > 5 mg/L within 30 minutes"},
			},
			Confidence: 0.8,
			VerifiedBy: 0,
			CreatedAt:  now,
		},
		{
			ID:       "gene_greenhouse_ventilation_schedule",
			Category: "optimize",
			Scenario: "greenhouse",
			SignalsMatch: []string{
				SignalTimePattern,
				SignalStrategyProven,
			},
			Strategy: []string{
				"Read temperature and humidity from greenhouse sensors",
				"Compare current readings with the optimal range from MEMORY.md",
				"If temperature > optimal range and time is daytime (06:00-18:00), open vents",
				"If humidity > 80%, activate exhaust fans",
				"Track ventilation events and correlate with temperature drop rate",
				"Optimize vent open/close schedule based on 7-day rolling data",
				"Update schedule in MEMORY.md with new timing recommendations",
			},
			Constraints: GeneConstraints{
				MaxRetries:     2,
				TimeoutSeconds: 300,
			},
			Validation: []ValidationStep{
				{Description: "Temperature stays within optimal range", Expected: "within 2°C of target for 80% of daylight hours"},
			},
			Confidence: 0.6,
			VerifiedBy: 0,
			CreatedAt:  now,
		},
	}
}
