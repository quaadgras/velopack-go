package velopack

import (
	"log/slog"
	"os"
)

// DownloadUpdatesInTheBackground should be called in a separate goroutine to download updates automatically,
// if there is an update, it will be applied on the next application startup. Uses slog for leveled logging.
// Calls [Run]. Only checks for updates once.
func DownloadUpdatesInTheBackground(source string) {
	Run(App{
		AutoApplyOnStartup: true,
		Logger: func(level, message string) {
			switch level {
			case "error":
				slog.Error(message)
			case "trace":
				slog.Debug(message)
			case "info":
				slog.Info(message)
			}
		},
	})
	manager, err := NewUpdateManager(source)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	latest, status, err := manager.CheckForUpdates()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	if status == UpdateAvailable {
		if err := manager.DownloadUpdates(latest, nil); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

// ApplyUpdatesAndRestart will exit your app immediately, apply updates, and then optionally relaunch the app using the specified
// restart [os.Args]. If you need to save state or clean up, you should do that before calling this method. The user may be
// prompted during the update, if the update requires additional frameworks to be installed etc. You can check if there are
// pending updates by checking UpdatePendingRestart.
func (up *UpdateManager) ApplyUpdatesAndRestart(update either[*UpdateInfo, *Asset], args ...string) error {
	if err := up.WaitForExitThenApplyUpdates(update, Restart(args)); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
