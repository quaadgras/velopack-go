//go:build !(cgo && (windows || linux || darwin) && (arm64 || amd64))

package velopack

import (
	"runtime"
)

func newError() error {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" || runtime.GOARCH == "amd64" {
			return ErrDisabledCGO
		}
	}
	return ErrPlatformNotSupported
}

// NewSourceFile creates a new FileSource update source for the given local directory path containing updates.
func NewSourceFile(psz_file_path string) (*UpdateSource, error) { return nil, newError() }

// NewSourceHTTP creates a new HttpSource update source for the given HTTP URL of a remote update server.
func NewSourceHTTP(psz_http_url string) (*UpdateSource, error) { return nil, newError() }

// NewUpdateManager creates a new UpdateManager instance from the given update location with optional [UpdateOptions] and [LocatorConfig].
func NewUpdateManager(psz_url_or_path string, options ...updateOptionsAndLocatorConfig) (*UpdateManager, error) {
	return nil, newError()
}

// NewUpdateManagerFromSource creates a new UpdateManager instance using the given UpdateSource.
func NewUpdateManagerFromSource(source *UpdateSource, options ...updateOptionsAndLocatorConfig) (*UpdateManager, error) {
	return nil, newError()
}

// AppID returns the currently installed app id.
func (up *UpdateManager) AppID() string { panic(newError()) }

// CurrentlyInstalledVersion returns the currently installed version of the app.
func (up *UpdateManager) CurrentlyInstalledVersion() string { panic(newError()) }

// IsPortable returns whether the app is in portable mode. On Windows this can be true or false.
// On MacOS and Linux this will always be true.
func (up *UpdateManager) IsPortable() bool { return false }

// UpdatePendingRestart returns an asset if there is an update downloaded which still needs to be
// applied. You can pass this asset to [UpdateManager.WaitExitThenApplyUpdates] to apply the update.
func (up *UpdateManager) UpdatePendingRestart() (*Asset, bool) { return nil, false }

// CheckForUpdates Checks for updates. If there are updates available, this method will return an [UpdateInfo]
// object containing the latest available release, and any delta updates that can be applied if they are available.
func (up *UpdateManager) CheckForUpdates() (*UpdateInfo, UpdateStatus, error) {
	return nil, UpdateError, newError()
}

// WaitForExitThenApplyUpdates this will launch the Velopack updater and tell it to wait for this program
// to exit gracefully.
//   - You should then clean up any state and exit your app. The updater will apply updates and then
//   - (if [Restart] specified) restart your app. The updater will only wait for 60 seconds before giving up.
func (up *UpdateManager) WaitForExitThenApplyUpdates(update either[*UpdateInfo, *Asset], options ...upto3[Silent, Restart, UnsafeProcessID]) error {
	return newError()
}

/*
Create a new _CUSTOM_ update source with user-provided callbacks to fetch release feeds and download assets.
You can report download progress using [UpdateSource.ReportProgress].
*/
func NewSourceCustomCallback(callbacks SourceCustomCallbacks) (*UpdateSource, error) {
	return nil, newError()
}

// DownloadUpdates downloads the specified updates to the local app packages directory.
// Progress is reported back to the caller via an optional callback. This function will
// acquire a global update lock so may fail if there is already another update operation
// in progress.
//   - If the update contains delta packages and the delta feature is enabled
//   - this method will attempt to unpack and prepare them.
//   - If there is no delta update available, or there is an error preparing delta
//     packages, this method will fall back to downloading the full version of the update.
func (up *UpdateManager) DownloadUpdates(update_info *UpdateInfo, progress func(progress uint)) error {
	return newError()
}

// Run helps you to handle app activation events correctly.
// This should be used as early as possible in your application startup code.
// (eg. the beginning of main() or wherever your entry point is).
// This function will not return in some cases.
func Run(a App) { return }
