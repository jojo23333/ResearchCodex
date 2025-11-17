package timeutil

import "time"

// TimestampSlug returns YYYYMMDD_HHMMSS formatted time in UTC.
func TimestampSlug(now time.Time) string {
	return now.UTC().Format("20060102_150405")
}

// ISO8601 returns an RFC3339 timestamp in UTC.
func ISO8601(now time.Time) string {
	return now.UTC().Format(time.RFC3339)
}

// NowTimestamps returns both slug and iso formatted timestamps.
func NowTimestamps() (slug string, iso string, now time.Time) {
	now = time.Now().UTC()
	return TimestampSlug(now), ISO8601(now), now
}
