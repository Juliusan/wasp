package evmlogger

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/slog"

	hiveLog "github.com/iotaledger/hive.go/log"
)

func Init(hiveLogger hiveLog.Logger) {
	log.SetDefault(log.NewLogger(&hiveLogHandler{hiveLogger}))
}

type hiveLogHandler struct{ hiveLog.Logger }

// Enabled implements slog.Handler.
func (*hiveLogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Handle implements slog.Handler.
func (h *hiveLogHandler) Handle(ctx context.Context, r slog.Record) error {
	switch {
	case r.Level >= slog.LevelError:
		h.Logger.LogError(r.Message)
	case r.Level <= slog.LevelDebug:
		h.Logger.LogDebug(r.Message)
	case r.Level == slog.LevelWarn:
		h.Logger.LogWarn(r.Message)
	default:
		h.Logger.LogInfo(r.Message)
	}
	return nil
}

// WithAttrs implements slog.Handler.
func (h *hiveLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// TODO: unimplemented in hive logger?
	return h
}

// WithGroup implements slog.Handler.
func (h *hiveLogHandler) WithGroup(name string) slog.Handler {
	// TODO: unimplemented in hive logger?
	return h
}
