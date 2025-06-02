/*
Package open implements an [actions.Runner] that opens a notification using a user config.
*/
package custom

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"text/template"

	"github.com/nobe4/gh-not/internal/colors"
	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

var errTemplateFailure = errors.New("failed to parse template")

type Runner struct {
	Client *gh.Client
	Config *config.Data
}

func (*Runner) Run(n *notifications.Notification, _ []string, w io.Writer) error {
	slog.Debug("open using custom action", "notification", n)

	templ, err := template.New("item").Parse("{{.Subject.Title}}: {{.Subject.URL}}")
	if err != nil {
		return errTemplateFailure
	}

	var buf strings.Builder

	// Execute the template on the notification and write to the buffer
	if err := templ.Execute(&buf, n); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// #nosec G204: We want variable input into this command, and damn the security issues!
	cmd := exec.Command("echo", buf.String())
	cmd.Stdout = w
	cmd.Stderr = w

	fmt.Fprintf(w, "%sEXECUTING: %s\n", colors.Blue("CUSTOM "), cmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}
