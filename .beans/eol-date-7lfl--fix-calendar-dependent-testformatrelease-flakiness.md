---
# eol-date-7lfl
title: Fix calendar-dependent TestFormatRelease flakiness
status: completed
type: bug
priority: high
created_at: 2026-05-17T16:02:14Z
updated_at: 2026-05-17T16:19:35Z
---

## Problem

CI run https://github.com/oliverandrich/eol-date/actions/runs/25993844856 failed in `internal/ui/display_test.go:78`:

```
--- FAIL: TestFormatRelease/recent_release
    display_test.go:78: formatRelease().relative = "2m ago", want "3m ago"
```

## Root Cause

`display_test.go:62` builds the release time with `now.AddDate(0, -3, 0)`. From 2026-05-17 that lands on 2026-02-17, which is **89 days** away (Feb 28 + Mar 31 + Apr 30). `formatDuration` in `display.go:35` truncates `89 days / 30 = 2 months`, so it returns `"2m"` instead of the expected `"3m"`.

The bug surfaces whenever the 3-month-back window contains a 28-day February — roughly when the test runs in March–May.

## Plan

- [x] Replace `AddDate(...)` with absolute `time.Duration` offsets (`now.Add(-90*24h)` etc.) in `TestFormatRelease` so the duration math is stable across calendar months.
- [x] Run `mise run test`, `mise run lint`, `mise run build` locally.
- [x] Run `/simplify` before commit.

## Acceptance

- `go test ./internal/ui/...` passes regardless of current date.
- CI green.

## Summary of Changes

- `internal/ui/display_test.go`: `TestFormatRelease` now builds the release-time offsets with `time.Duration` arithmetic (`now.Add(-90*24h)` for the recent case, `now.Add(-30*30*24h)` for the old case) instead of `AddDate`. The 30-day-per-month rounding inside `formatDuration` now matches the test expectations across all calendar months.
- No production code change; the bug was an incorrect assertion against calendar-aware `AddDate` while `formatDuration` works in 30-day buckets.
