//go:build tools

package tools

import (
	// This import points to the main gqlgen package, which includes the CLI tool's dependencies.
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)

// This file ensures Go modules see the above tools as dependencies
// even though they are only used during development.
