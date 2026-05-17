---
# eol-date-jj3j
title: TestFormatRelease is time-flaky
status: todo
type: bug
priority: normal
created_at: 2026-05-17T11:43:57Z
updated_at: 2026-05-17T11:43:57Z
---

\`internal/ui/display_test.go:78\` (\`TestFormatRelease/recent_release\`) fails intermittently with messages like:

\`\`\`
formatRelease().relative = "2m ago", want "3m ago"
\`\`\`

The test compares a relative-time string against a hardcoded expected value, so it fails whenever the test runs across a minute boundary relative to whatever fixed timestamp the input uses. Fails on \`main\` (verified at commit a077e41), independent of any tooling migration.

## Fix candidates

- Inject a clock interface into \`formatRelease\` so the test can pin "now"
- Loosen the assertion to a regex like \`^[0-9]+m ago$\` for the recent_release case
- Use a much larger gap (hours/days) where ±1 minute drift doesn't matter

## Evidence

- File: \`internal/ui/display_test.go:78\`
- Reproduces: \`go test ./internal/ui/ -run TestFormatRelease\` near a minute boundary
