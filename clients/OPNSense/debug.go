package OPNSense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

// ANSI color codes
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

// highlightJSON adds color to JSON strings using regex.
func highlightJSON(input string) string {
	replacements := []struct {
		pattern *regexp.Regexp
		color   string
	}{
		{regexp.MustCompile(`"([^"]+)"\s*:`), cyan},   // Keys
		{regexp.MustCompile(`:\s*"([^"]*)"`), green},  // Strings
		{regexp.MustCompile(`:\s*([0-9.]+)`), yellow}, // Numbers
		{regexp.MustCompile(`:\s*(true|false)`), red}, // Booleans
		{regexp.MustCompile(`:\s*(null)`), red},       // null
	}

	colored := input
	for _, r := range replacements {
		colored = r.pattern.ReplaceAllStringFunc(colored, func(s string) string {
			return r.pattern.ReplaceAllString(s, r.color+"$0"+reset)
		})
	}

	return colored
}

func prettyJson(input string) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, []byte(input), "", "  "); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	formatted := out.String()
	formatted = highlightJSON(formatted)

	return formatted, nil
}

func PrintJson(input string) error {
	formatted, err := prettyJson(input)
	if err != nil {
		return err
	}

	fmt.Printf(formatted)
	return nil
}
