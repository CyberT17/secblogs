package main

import (
	"testing"
	"time"
)

func TestIsWithinDateRange(t *testing.T) {
	now := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	pastCutoff := now.AddDate(0, 0, -7)
	futureCutoff := now.AddDate(0, 0, 7)

	tests := []struct {
		name      string
		published *time.Time
		want      bool
	}{
		{
			name:      "nil publication date",
			published: nil,
			want:      false,
		},
		{
			name: "published 8 days ago (older than past cutoff)",
			published: func() *time.Time {
				d := now.AddDate(0, 0, -8)
				return &d
			}(),
			want: false,
		},
		{
			name: "published exactly at past cutoff",
			published: func() *time.Time {
				d := pastCutoff
				return &d
			}(),
			want: false,
		},
		{
			name: "published 3 days ago",
			published: func() *time.Time {
				d := now.AddDate(0, 0, -3)
				return &d
			}(),
			want: true,
		},
		{
			name: "published right now",
			published: func() *time.Time {
				d := now
				return &d
			}(),
			want: true,
		},
		{
			name: "published 3 days in future (within 1 week)",
			published: func() *time.Time {
				d := now.AddDate(0, 0, 3)
				return &d
			}(),
			want: true,
		},
		{
			name: "published exactly at future cutoff (7 days in future)",
			published: func() *time.Time {
				d := futureCutoff
				return &d
			}(),
			want: true,
		},
		{
			name: "published 7 days and 1 second in future (beyond a week in future)",
			published: func() *time.Time {
				d := futureCutoff.Add(1 * time.Second)
				return &d
			}(),
			want: false,
		},
		{
			name: "published 8 days in future (beyond a week in future)",
			published: func() *time.Time {
				d := now.AddDate(0, 0, 8)
				return &d
			}(),
			want: false,
		},
		{
			name: "published far in future (e.g. 1 year ahead)",
			published: func() *time.Time {
				d := now.AddDate(1, 0, 0)
				return &d
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isWithinDateRange(tt.published, pastCutoff, futureCutoff)
			if got != tt.want {
				t.Errorf("isWithinDateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
