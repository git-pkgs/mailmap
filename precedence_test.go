package mailmap_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/git-pkgs/mailmap"
)

func TestDuplicatePrecedence(t *testing.T) {
	cases := []struct {
		title, input        string
		name, email         string
		wantName, wantEmail string
	}{
		{
			title: "name replacement",
			input: "First Name <old@example.com>\nLater Name <old@example.com>\n",
			name:  "jane", email: "old@example.com",
			wantName: "Later Name", wantEmail: "old@example.com",
		},
		{
			title: "email replacement",
			input: "<first@example.com> <OLD@EXAMPLE.COM>\n<later@example.com> <old@example.com>\n",
			name:  "jane", email: "Old@Example.com",
			wantName: "jane", wantEmail: "later@example.com",
		},
		{
			title: "name and email replacement",
			input: "First Name <first@example.com> <OLD@EXAMPLE.COM>\nLater Name <later@example.com> <old@example.com>\n",
			name:  "jane", email: "Old@Example.com",
			wantName: "Later Name", wantEmail: "later@example.com",
		},
		{
			title: "specific replacement before later generic entry",
			input: `First Name <first@example.com> JANE <OLD@EXAMPLE.COM>
Later Name <later@example.com> jane <old@example.com>
Generic Name <generic@example.com> <old@example.com>
`,
			name: "Jane", email: "Old@Example.com",
			wantName: "Later Name", wantEmail: "later@example.com",
		},
		{
			title: "unrelated entries do not override a match",
			input: "Jane Doe <jane@example.com> <old@example.com>\nOther Name <other@example.com> <other-old@example.com>\n",
			name:  "jane", email: "old@example.com",
			wantName: "Jane Doe", wantEmail: "jane@example.com",
		},
	}
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			m, err := mailmap.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			name, email := m.Resolve(tc.name, tc.email)
			if name != tc.wantName || email != tc.wantEmail {
				t.Errorf("Resolve(%q, %q) = (%q, %q), want (%q, %q)",
					tc.name, tc.email, name, email, tc.wantName, tc.wantEmail)
			}
			t.Run("git", func(t *testing.T) {
				gitPath, err := exec.LookPath("git")
				if err != nil {
					t.Skip("git is not installed")
				}
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, ".mailmap"), []byte(tc.input), 0o600); err != nil {
					t.Fatal(err)
				}
				cmd := exec.Command(gitPath, "init", "--quiet", dir)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("git init: %v\n%s", err, out)
				}
				cmd = exec.Command(gitPath, "check-mailmap", fmt.Sprintf("%s <%s>", tc.name, tc.email))
				cmd.Dir = dir
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("git check-mailmap: %v\n%s", err, out)
				}
				want := fmt.Sprintf("%s <%s>", tc.wantName, tc.wantEmail)
				if got := strings.TrimSpace(string(out)); got != want {
					t.Errorf("git check-mailmap = %q, want %q", got, want)
				}
			})
		})
	}
}
