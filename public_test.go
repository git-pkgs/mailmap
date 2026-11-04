package mailmap_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/git-pkgs/mailmap"
)

func ExampleParse() {
	m, err := mailmap.Parse(strings.NewReader(`Jane Doe <jane@example.com> <jane@laptop.(none)>
`))
	if err != nil {
		panic(err)
	}
	name, email := m.Resolve("jane", "jane@laptop.(none)")
	fmt.Printf("%s <%s>\n", name, email)
	// Output: Jane Doe <jane@example.com>
}

func TestParseRepositoryMailmap(t *testing.T) {
	input := `
  # Contributors with old addresses and a shared build account.
Jane Doe <jane@example.com> <jane@laptop.(none)>
Build Team <build@example.com> <builder@workstation.(none)>
Jane Doe <jane@example.com> jane <builder@workstation.(none)>
Later Name <later@example.com> JANE <builder@workstation.(none)>
Later Team <later-team@example.com> <builder@workstation.(none)>
<new@example.com> <old@example.com>
Sam Developer <sam@example.com>
not a mailmap entry
Broken <unclosed@example.com
Too Many <one@example.com> <two@example.com> <three@example.com>
`
	m, err := mailmap.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, email         string
		wantName, wantEmail string
	}{
		{"jane", "jane@laptop.(none)", "Jane Doe", "jane@example.com"},
		{"JANE", "BUILDER@WORKSTATION.(NONE)", "Later Name", "later@example.com"},
		{"Other", "builder@workstation.(none)", "Later Team", "later-team@example.com"},
		{"Old Name", "OLD@EXAMPLE.COM", "Old Name", "new@example.com"},
		{"sam", "sam@example.com", "Sam Developer", "sam@example.com"},
		{"New Contributor", "new-contributor@example.com", "New Contributor", "new-contributor@example.com"},
		{"Broken", "unclosed@example.com", "Broken", "unclosed@example.com"},
		{"Too Many", "three@example.com", "Too Many", "three@example.com"},
	}
	for _, tc := range cases {
		t.Run(tc.name+"/"+tc.email, func(t *testing.T) {
			name, email := m.Resolve(tc.name, tc.email)
			if name != tc.wantName || email != tc.wantEmail {
				t.Fatalf("Resolve(%q, %q) = (%q, %q), want (%q, %q)",
					tc.name, tc.email, name, email, tc.wantName, tc.wantEmail)
			}
		})
	}
}

func TestParseReaderError(t *testing.T) {
	wantErr := errors.New("mailmap read failed")
	r := io.MultiReader(strings.NewReader("Jane Doe <jane@example.com>\n"), iotest.ErrReader(wantErr))
	m, err := mailmap.Parse(r)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Parse error = %v, want %v", err, wantErr)
	}
	name, email := m.Resolve("jane", "jane@example.com")
	if name != "Jane Doe" || email != "jane@example.com" {
		t.Fatalf("partial map resolved to (%q, %q)", name, email)
	}
}

func TestParseScannerLimit(t *testing.T) {
	input := "Jane Doe <jane@example.com>\n" + strings.Repeat("x", bufio.MaxScanTokenSize)
	m, err := mailmap.Parse(strings.NewReader(input))
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("Parse error = %v, want %v", err, bufio.ErrTooLong)
	}
	name, email := m.Resolve("jane", "jane@example.com")
	if name != "Jane Doe" || email != "jane@example.com" {
		t.Fatalf("partial map resolved to (%q, %q)", name, email)
	}
}

func TestResolveUninitialized(t *testing.T) {
	var zero mailmap.Mailmap
	for _, m := range []*mailmap.Mailmap{nil, &zero, mailmap.New()} {
		name, email := m.Resolve("Jane", "JANE@EXAMPLE.COM")
		if name != "Jane" || email != "JANE@EXAMPLE.COM" {
			t.Errorf("empty map resolved to (%q, %q)", name, email)
		}
	}
}
