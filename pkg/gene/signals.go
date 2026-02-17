package gene

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Signal types for the Clawland edge AI ecosystem.
// Adapted from Evolver's signal types for physical-world monitoring.
const (
	SignalSensorError       = "sensor_error"         // Sensor read failure
	SignalNewPatternDetected = "new_pattern_detected" // LLM identifies novel pattern
	SignalResponseTooSlow   = "response_too_slow"    // Action taken too late
	SignalStrategyProven    = "strategy_proven"       // Same gene succeeds N times
	SignalUnknownSituation  = "unknown_situation"     // No gene matches signals
	SignalThresholdBreach   = "threshold_breach"      // Sensor value exceeds limits
	SignalCrossSensorAnomaly = "cross_sensor_anomaly" // Multiple sensors correlate
	SignalTimePattern       = "time_pattern"          // Recurring time-based issue
)

// SignalExtractor extracts signals from picclaw's existing data sources.
// Ported from Evolver's signals.js, adapted for sensor-world signals.
type SignalExtractor struct {
	workspace string
	memoryDir string
}

// NewSignalExtractor creates a new extractor for the given workspace.
func NewSignalExtractor(workspace string) *SignalExtractor {
	return &SignalExtractor{
		workspace: workspace,
		memoryDir: filepath.Join(workspace, "memory"),
	}
}

// ExtractSignals gathers signals from all available data sources.
func (se *SignalExtractor) ExtractSignals() []string {
	var signals []string

	// Extract from long-term memory
	signals = append(signals, se.extractFromMemory()...)

	// Extract from today's daily note
	signals = append(signals, se.extractFromDailyNotes()...)

	// Deduplicate
	return deduplicate(signals)
}

// ExtractSignalsFromText extracts signals from arbitrary text content.
// Useful for processing tool execution results.
func (se *SignalExtractor) ExtractSignalsFromText(text string) []string {
	var signals []string
	lower := strings.ToLower(text)

	// Sensor error patterns
	sensorErrorPatterns := []string{
		"sensor error", "read failure", "sensor timeout",
		"connection refused", "i2c error", "spi error",
		"no data", "nan reading", "sensor offline",
	}
	for _, p := range sensorErrorPatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalSensorError)
			break
		}
	}

	// Threshold breach patterns
	thresholdPatterns := []string{
		"exceeded", "above threshold", "below threshold",
		"critical", "alarm", "alert", "warning",
		"too high", "too low", "out of range",
	}
	for _, p := range thresholdPatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalThresholdBreach)
			break
		}
	}

	// Response timing patterns
	slowPatterns := []string{
		"too slow", "delayed response", "timeout",
		"late detection", "missed window",
	}
	for _, p := range slowPatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalResponseTooSlow)
			break
		}
	}

	// Cross-sensor correlation patterns
	crossPatterns := []string{
		"correlation", "multiple sensors", "cross-check",
		"coincide", "both sensors", "all sensors",
	}
	for _, p := range crossPatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalCrossSensorAnomaly)
			break
		}
	}

	// Novel pattern detection
	novelPatterns := []string{
		"never seen", "unusual", "anomaly", "unexpected",
		"new pattern", "first time", "novel",
	}
	for _, p := range novelPatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalNewPatternDetected)
			break
		}
	}

	// Time-based pattern detection
	timePatterns := []string{
		"every night", "recurring", "daily pattern",
		"same time", "periodic", "schedule",
		"morning", "evening", "dawn", "dusk",
	}
	for _, p := range timePatterns {
		if strings.Contains(lower, p) {
			signals = append(signals, SignalTimePattern)
			break
		}
	}

	return deduplicate(signals)
}

// extractFromMemory extracts signals from MEMORY.md.
func (se *SignalExtractor) extractFromMemory() []string {
	memFile := filepath.Join(se.memoryDir, "MEMORY.md")
	data, err := os.ReadFile(memFile)
	if err != nil {
		return nil
	}
	return se.ExtractSignalsFromText(string(data))
}

// extractFromDailyNotes extracts signals from recent daily notes (last 3 days).
func (se *SignalExtractor) extractFromDailyNotes() []string {
	var signals []string
	for i := 0; i < 3; i++ {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("20060102")
		monthDir := dateStr[:6]
		filePath := filepath.Join(se.memoryDir, monthDir, dateStr+".md")

		if data, err := os.ReadFile(filePath); err == nil {
			signals = append(signals, se.ExtractSignalsFromText(string(data))...)
		}
	}
	return deduplicate(signals)
}

// ExtractScenarioHints tries to determine the deployment scenario from context.
func (se *SignalExtractor) ExtractScenarioHints() string {
	memFile := filepath.Join(se.memoryDir, "MEMORY.md")
	data, err := os.ReadFile(memFile)
	if err != nil {
		return "generic"
	}
	text := strings.ToLower(string(data))

	scenarioPatterns := map[string]*regexp.Regexp{
		"datacenter": regexp.MustCompile(`data\s*center|server\s*room|rack|ups|cooling`),
		"aquaculture": regexp.MustCompile(`fish|pond|aqua|dissolved\s*oxygen|do\s*sensor|water\s*quality`),
		"greenhouse":  regexp.MustCompile(`greenhouse|ventilat|humidity|soil|plant|grow`),
		"agriculture": regexp.MustCompile(`farm|crop|irrigation|field|soil\s*moisture`),
	}

	for scenario, re := range scenarioPatterns {
		if re.MatchString(text) {
			return scenario
		}
	}
	return "generic"
}

// deduplicate removes duplicate strings from a slice.
func deduplicate(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make(map[string]struct{}, len(s))
	result := make([]string, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
