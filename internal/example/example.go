// SPDX-License-Identifier: EUPL-1.2
// Copyright (c) 2025 Oliver Andrich

// Package example provides a simple example of an internal package.
package example

import "fmt"

// Greet returns a greeting message for the given name.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
