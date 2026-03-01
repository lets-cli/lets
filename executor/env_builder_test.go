package executor

import (
	"testing"
)

func TestEnvBuilder(t *testing.T) {
	t.Run("should build environment with multiple layers", func(t *testing.T) {
		builder := NewEnvBuilder().
			Add("layer1", map[string]string{"A": "1", "B": "2"}).
			Add("layer2", map[string]string{"C": "3"}).
			Add("layer3", map[string]string{"A": "override"})

		result := builder.Build([]string{"BASE=value"})

		// Check base is preserved
		found := false
		for _, v := range result {
			if v == "BASE=value" {
				found = true
				break
			}
		}
		if !found {
			t.Error("base environment should be preserved")
		}

		// Check all layers are present
		hasA := false
		hasB := false
		hasC := false
		aOverridden := false

		for _, v := range result {
			switch v {
			case "A=1":
				hasA = true
			case "A=override":
				aOverridden = true
			case "B=2":
				hasB = true
			case "C=3":
				hasC = true
			}
		}

		if !hasA && !aOverridden {
			t.Error("layer1 A should be present")
		}
		if !hasB {
			t.Error("layer1 B should be present")
		}
		if !hasC {
			t.Error("layer2 C should be present")
		}
		if !aOverridden {
			t.Error("layer3 should override A")
		}
	})

	t.Run("should handle nil values gracefully", func(t *testing.T) {
		builder := NewEnvBuilder().
			Add("layer1", nil).
			Add("layer2", map[string]string{"A": "1"})

		result := builder.Build(nil)

		if len(result) != 1 {
			t.Errorf("expected 1 entry, got %d", len(result))
		}
	})

	t.Run("should return layers for debugging", func(t *testing.T) {
		builder := NewEnvBuilder().
			Add("first", map[string]string{"A": "1"}).
			Add("second", map[string]string{"B": "2"})

		layers := builder.Layers()

		if len(layers) != 2 {
			t.Errorf("expected 2 layers, got %d", len(layers))
		}
		if layers[0].Name != "first" {
			t.Errorf("expected first layer name 'first', got '%s'", layers[0].Name)
		}
		if layers[1].Name != "second" {
			t.Errorf("expected second layer name 'second', got '%s'", layers[1].Name)
		}
	})
}

func TestBuildChecksumEnv(t *testing.T) {
	t.Run("should create checksum env vars", func(t *testing.T) {
		checksums := map[string]string{
			"__default_checksum__": "abc123",
			"my-source":            "def456",
		}

		result := BuildChecksumEnv(checksums)

		if result["LETS_CHECKSUM"] != "abc123" {
			t.Errorf("expected LETS_CHECKSUM=abc123, got %s", result["LETS_CHECKSUM"])
		}
		if result["LETS_CHECKSUM_MY_SOURCE"] != "def456" {
			t.Errorf("expected LETS_CHECKSUM_MY_SOURCE=def456, got %s", result["LETS_CHECKSUM_MY_SOURCE"])
		}
	})
}

func TestBuildChecksumChangedEnv(t *testing.T) {
	t.Run("should detect changed checksums", func(t *testing.T) {
		current := map[string]string{
			"__default_checksum__": "new-hash",
			"source1":              "unchanged",
		}
		persisted := map[string]string{
			"__default_checksum__": "old-hash",
			"source1":              "unchanged",
		}

		result := BuildChecksumChangedEnv(current, persisted)

		if result["LETS_CHECKSUM_CHANGED"] != "true" {
			t.Errorf("expected LETS_CHECKSUM_CHANGED=true, got %s", result["LETS_CHECKSUM_CHANGED"])
		}
		if result["LETS_CHECKSUM_SOURCE1_CHANGED"] != "false" {
			t.Errorf("expected LETS_CHECKSUM_SOURCE1_CHANGED=false, got %s", result["LETS_CHECKSUM_SOURCE1_CHANGED"])
		}
	})

	t.Run("should mark new checksums as changed", func(t *testing.T) {
		current := map[string]string{
			"new-source": "hash",
		}
		persisted := map[string]string{}

		result := BuildChecksumChangedEnv(current, persisted)

		if result["LETS_CHECKSUM_NEW_SOURCE_CHANGED"] != "true" {
			t.Errorf("expected new checksum to be marked as changed")
		}
	})
}

func TestFormatEnvForDebug(t *testing.T) {
	t.Run("should format empty env", func(t *testing.T) {
		result := FormatEnvForDebug(nil)
		if result != "(empty)" {
			t.Errorf("expected '(empty)', got '%s'", result)
		}
	})

	t.Run("should format env entries", func(t *testing.T) {
		env := []string{"A=1", "B=2"}
		result := FormatEnvForDebug(env)

		if len(result) == 0 {
			t.Error("expected non-empty result")
		}
	})
}
