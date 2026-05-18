// Package source resolves log input sources for logslice.
//
// It supports three modes of operation:
//
//  1. Stdin — when no file paths are provided, the process's standard input
//     (or a configured substitute) is used as the sole source.
//
//  2. Explicit file paths — one or more file paths are opened in order.
//
//  3. Glob patterns — shell-style wildcards (e.g. /var/log/*.log) are
//     expanded via [ExpandGlobs] before being passed to [Resolver.Resolve].
//
// Typical usage:
//
//	paths, err := source.ExpandGlobs(rawArgs)
//	if err != nil { ... }
//	rs := source.New()
//	readers, err := rs.Resolve(paths)
//	if err != nil { ... }
//	defer func() {
//		for _, r := range readers { r.Close() }
//	}()
package source
