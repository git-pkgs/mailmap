// Package mailmap parses Git .mailmap files and resolves author identities.
//
// Parse accepts the four mailmap forms for replacing names, email addresses,
// or both. Names and email addresses match without regard to case. Entries
// matching both name and email take priority over entries matching only email;
// within either group, the last matching entry wins.
//
// Blank lines, comments starting with #, and malformed entries are skipped.
// Parse returns a partially populated Mailmap together with any reader or
// scanner error. It uses bufio.Scanner's default token limit.
//
// An empty, zero-value, or nil Mailmap returns identities unchanged.
package mailmap
