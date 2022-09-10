package heredoc

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
)

func Doc(s string, padding int) string {
	s = heredoc.Doc(s)
	if padding == 0 {
		return s
	}

	b := strings.Builder{}
	p := strings.Repeat(" ", padding)

	for _, line := range strings.Split(s, "\n") {
		b.WriteString(p)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}
