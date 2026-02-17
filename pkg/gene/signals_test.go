package gene

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractSignalsFromText_SensorError(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("The DHT22 sensor reported a read failure at 3am")

	found := false
	for _, s := range signals {
		if s == SignalSensorError {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected sensor_error signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_ThresholdBreach(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("Temperature exceeded critical threshold")

	found := false
	for _, s := range signals {
		if s == SignalThresholdBreach {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected threshold_breach signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_ResponseTooSlow(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("The alert was too slow, missed the window")

	found := false
	for _, s := range signals {
		if s == SignalResponseTooSlow {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected response_too_slow signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_CrossSensorAnomaly(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("Detected correlation between temperature and humidity sensors")

	found := false
	for _, s := range signals {
		if s == SignalCrossSensorAnomaly {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected cross_sensor_anomaly signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_NewPattern(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("This is an unusual anomaly we haven't seen before")

	found := false
	for _, s := range signals {
		if s == SignalNewPatternDetected {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected new_pattern_detected signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_TimePattern(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("This recurring issue appears every night")

	found := false
	for _, s := range signals {
		if s == SignalTimePattern {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected time_pattern signal, got %v", signals)
	}
}

func TestExtractSignalsFromText_MultipleSignals(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("Sensor error detected and temperature exceeded critical alert threshold")

	if len(signals) < 2 {
		t.Errorf("expected at least 2 signals, got %d: %v", len(signals), signals)
	}
}

func TestExtractSignalsFromText_NoSignals(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("Everything is normal and operating as expected")

	if len(signals) != 0 {
		t.Errorf("expected 0 signals for normal text, got %d: %v", len(signals), signals)
	}
}

func TestExtractSignalsFromText_CaseInsensitive(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("SENSOR ERROR DETECTED")

	found := false
	for _, s := range signals {
		if s == SignalSensorError {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("signal extraction should be case-insensitive, got %v", signals)
	}
}

func TestExtractSignalsFromText_Deduplicated(t *testing.T) {
	se := NewSignalExtractor(t.TempDir())
	signals := se.ExtractSignalsFromText("sensor error and another sensor timeout happened")

	count := 0
	for _, s := range signals {
		if s == SignalSensorError {
			count++
		}
	}
	if count > 1 {
		t.Errorf("expected deduplicated signals, got %d occurrences of sensor_error", count)
	}
}

func TestExtractScenarioHints_Datacenter(t *testing.T) {
	tmpDir := t.TempDir()
	memDir := filepath.Join(tmpDir, "memory")
	os.MkdirAll(memDir, 0755)
	os.WriteFile(filepath.Join(memDir, "MEMORY.md"), []byte("Monitoring the data center server room cooling systems"), 0644)

	se := NewSignalExtractor(tmpDir)
	scenario := se.ExtractScenarioHints()

	if scenario != "datacenter" {
		t.Errorf("expected 'datacenter' scenario, got %s", scenario)
	}
}

func TestExtractScenarioHints_Aquaculture(t *testing.T) {
	tmpDir := t.TempDir()
	memDir := filepath.Join(tmpDir, "memory")
	os.MkdirAll(memDir, 0755)
	os.WriteFile(filepath.Join(memDir, "MEMORY.md"), []byte("Monitoring dissolved oxygen levels in the fish pond"), 0644)

	se := NewSignalExtractor(tmpDir)
	scenario := se.ExtractScenarioHints()

	if scenario != "aquaculture" {
		t.Errorf("expected 'aquaculture' scenario, got %s", scenario)
	}
}

func TestExtractScenarioHints_Generic(t *testing.T) {
	tmpDir := t.TempDir()
	se := NewSignalExtractor(tmpDir)
	scenario := se.ExtractScenarioHints()

	if scenario != "generic" {
		t.Errorf("expected 'generic' scenario when no MEMORY.md, got %s", scenario)
	}
}

func TestDeduplicate(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	result := deduplicate(input)

	if len(result) != 3 {
		t.Errorf("expected 3 unique items, got %d: %v", len(result), result)
	}
}

func TestDeduplicate_Empty(t *testing.T) {
	result := deduplicate([]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestDeduplicate_Nil(t *testing.T) {
	result := deduplicate(nil)
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
