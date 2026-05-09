package runtime

import (
	"io"
	"log/slog"
	"strings"
)

const defaultDiagnosticLogSource = "backend"

// InstallDiagnosticLogger installs a lightweight structured slog logger.
func InstallDiagnosticLogger(writer io.Writer, source string, level slog.Leveler) {
	if writer == nil {
		return
	}
	normalizedSource := strings.TrimSpace(source)
	if normalizedSource == "" {
		normalizedSource = defaultDiagnosticLogSource
	}
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level}).WithAttrs([]slog.Attr{
		slog.String("source", normalizedSource),
	})
	slog.SetDefault(slog.New(handler))
}
