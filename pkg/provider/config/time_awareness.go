package config

import (
	"fmt"
	"strings"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
)

type TimeAwarenessSettings struct {
	Enabled    bool   `json:"enabled,omitempty"`
	TimeZone   string `json:"time_zone,omitempty"`
	TimeFormat string `json:"time_format,omitempty"`
}

func (p *TimeAwarenessSettings) Apply(options []llmagent.Option, agentName string) ([]llmagent.Option, error) {
	// Validate TimeAwareness settings before applying
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("invalid TimeAwareness settings for agent %s: %w", agentName, err)
	}

	options = append(options, llmagent.WithAddCurrentTime(p.Enabled))
	if p.TimeZone != "" {
		options = append(options, llmagent.WithTimezone(p.TimeZone))
	}

	if p.TimeFormat != "" {
		options = append(options, llmagent.WithTimeFormat(p.TimeFormat))
	}

	return options, nil

}

// validateTimeAwarenessSettings validates TimeAwareness configuration
func (p *TimeAwarenessSettings) Validate() error {
	if p == nil {
		return nil
	}

	// Validate timezone if provided
	if p.TimeZone != "" {
		// Check if timezone is valid by attempting to load it
		_, err := time.LoadLocation(p.TimeZone)
		if err != nil {
			// Check for common timezone abbreviations that might not work and provide helpful error
			commonTimezones := map[string]string{
				"PST": "America/Los_Angeles",
				"PDT": "America/Los_Angeles",
				"CST": "America/Chicago",
				"CDT": "America/Chicago",
				"MST": "America/Denver",
				"MDT": "America/Denver",
				"EDT": "America/New_York",
			}

			if canonical, exists := commonTimezones[strings.ToUpper(p.TimeZone)]; exists {
				return fmt.Errorf("invalid timezone '%s': use IANA timezone name '%s' instead of abbreviation", p.TimeZone, canonical)
			}

			return fmt.Errorf("invalid timezone '%s': %w. Use IANA timezone names like 'America/New_York', 'Europe/London', or 'UTC'", p.TimeZone, err)
		}
	}

	// Validate time format if provided
	if p.TimeFormat != "" {
		// Test the format by formatting a known time
		testTime := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
		_ = testTime.Format(p.TimeFormat) // Go's time.Format doesn't return an error, it just formats

		// Check for common format mistakes and provide suggestions
		if strings.Contains(p.TimeFormat, "YYYY") || strings.Contains(p.TimeFormat, "MM") || strings.Contains(p.TimeFormat, "DD") {
			return fmt.Errorf("invalid time format '%s': Go uses '2006' for year, '01' for month, '02' for day. Example: '2006-01-02 15:04:05'", p.TimeFormat)
		}
	}

	return nil
}
