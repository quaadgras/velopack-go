//go:build cgo && (windows || linux || darwin) && (arm64 || amd64)

// Package velopack provides a Go interface to the Velopack library for managing software updates and distribution on desktop.
package velopack

/*
#cgo LDFLAGS: -L${SRCDIR}/../binaries

#include <stdlib.h>
#include "Velopack.h"

#ifdef _WIN32
#define GO_EXPORT __declspec(dllexport)
#else
#define GO_EXPORT
#endif

extern GO_EXPORT char *go_release_feed_callback(uintptr_t user_data, char *psz_releases_name);
extern GO_EXPORT void go_free_release_feed_callback(uintptr_t user_data, char *psz_feed);
extern GO_EXPORT bool go_download_asset_callback(uintptr_t user_data, struct vpkc_asset_t *asset, char *psz_local_path, size_t progress_callback_id);
extern GO_EXPORT void go_progress_callback(uintptr_t user_data, size_t progress);
extern GO_EXPORT void go_after_install_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_before_uninstall_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_before_update_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_after_update_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_first_run_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_restarted_callback(uintptr_t user_data, char *psz_app_version);
extern GO_EXPORT void go_log_callback(uintptr_t user_data, char *level, char *psz_message);

static vpkc_update_source_t *go_vpkc_new_source_custom_callback(uintptr_t user_data) {
    return vpkc_new_source_custom_callback(
        (vpkc_release_feed_delegate_t)go_release_feed_callback,
        (vpkc_free_release_feed_t)go_free_release_feed_callback,
        (vpkc_download_asset_delegate_t)go_download_asset_callback,
        (void *)user_data
    );
}

static bool go_vpkc_download_updates(vpkc_update_manager_t *p_manager,
                                     struct vpkc_update_info_t *p_update,
                                     uintptr_t user_data) {
    return vpkc_download_updates(p_manager, p_update, (vpkc_progress_callback_t)go_progress_callback, (void *)user_data);
}

static void go_set_logger() {
    vpkc_set_logger((vpkc_log_callback_t)go_log_callback, 0);
}

static void go_callbacks() {
    vpkc_app_set_hook_after_install((vpkc_hook_callback_t)go_after_install_callback);
    vpkc_app_set_hook_before_uninstall((vpkc_hook_callback_t)go_before_uninstall_callback);
    vpkc_app_set_hook_before_update((vpkc_hook_callback_t)go_before_update_callback);
    vpkc_app_set_hook_after_update((vpkc_hook_callback_t)go_after_update_callback);
    vpkc_app_set_hook_first_run((vpkc_hook_callback_t)go_first_run_callback);
    vpkc_app_set_hook_restarted((vpkc_hook_callback_t)go_restarted_callback);
}

*/
import "C"
import (
	"errors"
	"runtime"
	"runtime/cgo"
	"sync/atomic"
	"unsafe"
)

//export go_release_feed_callback
func go_release_feed_callback(user_data uintptr, psz_releases_name *C.char) *C.char {
	return C.CString(string(cgo.Handle(user_data).Value().(SourceCustomCallbacks).ReleaseFeedFunc(C.GoString(psz_releases_name))))
}

//export go_free_release_feed_callback
func go_free_release_feed_callback(_ uintptr, psz_feed *C.char) {
	C.free(unsafe.Pointer(psz_feed))
}

//export go_download_asset_callback
func go_download_asset_callback(user_data uintptr, asset *C.vpkc_asset_t, psz_local_path *C.char, progress_callback_id C.size_t) C.bool {
	callbacks := cgo.Handle(user_data).Value().(SourceCustomCallbacks)
	success := callbacks.DownloadAssetFunc(
		toAsset(asset),
		C.GoString(psz_local_path),
		func(progress int16) {
			C.vpkc_source_report_progress(progress_callback_id, C.int16_t(progress))
		},
	)
	return C.bool(success)
}

//export go_progress_callback
func go_progress_callback(user_data uintptr, progress C.size_t) {
	cgo.Handle(user_data).Value().(func(uint))(uint(progress))
}

//export go_after_install_callback
func go_after_install_callback(_ uintptr, psz_app_version *C.char) {
	if app.WindowsHookAfterInstall == nil {
		return
	}
	app.WindowsHookAfterInstall(C.GoString(psz_app_version))
}

//export go_before_uninstall_callback
func go_before_uninstall_callback(_ uintptr, psz_app_version *C.char) {
	if app.WindowsHookBeforeUninstall == nil {
		return
	}
	app.WindowsHookBeforeUninstall(C.GoString(psz_app_version))
}

//export go_before_update_callback
func go_before_update_callback(_ uintptr, psz_app_version *C.char) {
	if app.WindowsHookBeforeUpdate == nil {
		return
	}
	app.WindowsHookBeforeUpdate(C.GoString(psz_app_version))
}

//export go_after_update_callback
func go_after_update_callback(_ uintptr, psz_app_version *C.char) {
	if app.WindowsHookAfterUpdate == nil {
		return
	}
	app.WindowsHookAfterUpdate(C.GoString(psz_app_version))
}

//export go_first_run_callback
func go_first_run_callback(_ uintptr, psz_app_version *C.char) {
	if app.HookFirstRun == nil {
		return
	}
	app.HookFirstRun(C.GoString(psz_app_version))
}

//export go_restarted_callback
func go_restarted_callback(_ uintptr, psz_app_version *C.char) {
	if app.HookRestarted == nil {
		return
	}
	app.HookRestarted(C.GoString(psz_app_version))
}

//export go_log_callback
func go_log_callback(_ uintptr, level, psz_message *C.char) {
	app.Logger(C.GoString(level), C.GoString(psz_message))
}

func toAsset(asset *C.vpkc_asset_t) *Asset {
	if asset == nil {
		return nil
	}
	converted := &Asset{
		handle:        unsafe.Pointer(asset),
		PackageID:     C.GoString(asset.PackageId),
		Version:       C.GoString(asset.Version),
		Type:          AssetType(C.GoString(asset.Type)),
		Filename:      C.GoString(asset.FileName),
		SHA1:          C.GoString(asset.SHA1),
		SHA256:        C.GoString(asset.SHA256),
		Size:          uint64(asset.Size),
		NotesMarkdown: C.GoString(asset.NotesMarkdown),
		NotesHTML:     C.GoString(asset.NotesHtml),
	}
	runtime.AddCleanup(converted, func(handle unsafe.Pointer) {
		C.vpkc_free_asset((*C.vpkc_asset_t)(handle))
	}, converted.handle)
	return converted
}

func (info *UpdateInfo) load(update_info *C.vpkc_update_info_t) *UpdateInfo {
	var deltas []*Asset
	if update_info.DeltasToTarget != nil {
		for ptr := update_info.DeltasToTarget; *ptr != nil; ptr = (**C.vpkc_asset_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*ptr))) {
			deltas = append(deltas, toAsset(*ptr))
		}
	}
	if info.handle != unsafe.Pointer(update_info) {
		runtime.AddCleanup(info, func(handle *C.vpkc_update_info_t) {
			C.vpkc_free_update_info(handle)
		}, update_info)
	}
	*info = UpdateInfo{
		handle:            unsafe.Pointer(update_info),
		TargetFullRelease: toAsset(update_info.TargetFullRelease),
		BaseRelease:       toAsset(update_info.BaseRelease),
		DeltasToTarget:    deltas,
		IsDowngrade:       bool(update_info.IsDowngrade),
	}
	return info
}

func get_last_error() error {
	var buf [512]byte
	n := C.vpkc_get_last_error((*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	return errors.New(string(buf[:min(n, C.size_t(len(buf)))]))
}

// NewSourceFile creates a new FileSource update source for the given local directory path containing updates.
func NewSourceFile(psz_file_path string) (*UpdateSource, error) {
	file_path := C.CString(psz_file_path)
	defer C.free(unsafe.Pointer(file_path))
	source := new(UpdateSource)
	source.handle = C.vpkc_new_source_file(file_path)
	if source.handle == nil {
		return nil, get_last_error()
	}
	runtime.AddCleanup(source, func(handle unsafe.Pointer) {
		C.vpkc_free_source(handle)
	}, source.handle)
	return source, nil
}

// NewSourceHTTP creates a new HttpSource update source for the given HTTP URL of a remote update server.
func NewSourceHTTP(psz_http_url string) (*UpdateSource, error) {
	http_url := C.CString(psz_http_url)
	defer C.free(unsafe.Pointer(http_url))
	source := new(UpdateSource)
	source.handle = C.vpkc_new_source_http_url(http_url)
	if source.handle == nil {
		return nil, get_last_error()
	}
	runtime.AddCleanup(source, func(handle unsafe.Pointer) {
		C.vpkc_free_source(handle)
	}, source.handle)
	return source, nil
}

func optionsLocator(options ...updateOptionsAndLocatorConfig) (*C.vpkc_update_options_t, *C.vpkc_locator_config_t) {
	var p_options *C.vpkc_update_options_t
	var p_locator *C.vpkc_locator_config_t
	if len(options) > 0 {
		for _, option := range options {
			switch option := option.(type) {
			case UpdateOptions:
				ExplicitChannel := C.CString(option.ExplicitChannel)
				defer C.free(unsafe.Pointer(ExplicitChannel))
				p_options = &C.vpkc_update_options_t{
					AllowVersionDowngrade:       C.bool(option.AllowVersionDowngrade),
					ExplicitChannel:             ExplicitChannel,
					MaximumDeltasBeforeFallback: C.int32_t(option.MaximumDeltasBeforeFallback),
				}
			case LocatorConfig:
				RootAppDir := C.CString(option.RootAppDir)
				defer C.free(unsafe.Pointer(RootAppDir))
				UpdateExePath := C.CString(option.UpdateExePath)
				defer C.free(unsafe.Pointer(UpdateExePath))
				PackagesDir := C.CString(option.PackagesDir)
				defer C.free(unsafe.Pointer(PackagesDir))
				ManifestPath := C.CString(option.ManifestPath)
				defer C.free(unsafe.Pointer(ManifestPath))
				CurrentBinaryDir := C.CString(option.CurrentBinaryDir)
				defer C.free(unsafe.Pointer(CurrentBinaryDir))
				p_locator = &C.vpkc_locator_config_t{
					RootAppDir:       RootAppDir,
					UpdateExePath:    UpdateExePath,
					PackagesDir:      PackagesDir,
					ManifestPath:     ManifestPath,
					CurrentBinaryDir: CurrentBinaryDir,
					IsPortable:       C.bool(option.IsPortable),
				}
			}
		}
	}
	return p_options, p_locator
}

// NewUpdateManager creates a new UpdateManager instance from the given update location with optional [UpdateOptions] and [LocatorConfig].
func NewUpdateManager(psz_url_or_path string, options ...updateOptionsAndLocatorConfig) (*UpdateManager, error) {
	p_options, p_locator := optionsLocator(options...)
	url_or_path := C.CString(psz_url_or_path)
	defer C.free(unsafe.Pointer(url_or_path))
	manager := new(UpdateManager)
	if !C.vpkc_new_update_manager(url_or_path, p_options, p_locator, &manager.handle) {
		return nil, get_last_error()
	}
	runtime.AddCleanup(manager, func(handle unsafe.Pointer) {
		C.vpkc_free_update_manager(handle)
	}, manager.handle)
	return manager, nil
}

// NewUpdateManagerFromSource creates a new UpdateManager instance using the given UpdateSource.
func NewUpdateManagerFromSource(source *UpdateSource, options ...updateOptionsAndLocatorConfig) (*UpdateManager, error) {
	p_options, p_locator := optionsLocator(options...)
	manager := new(UpdateManager)
	if !C.vpkc_new_update_manager_with_source(source.handle, p_options, p_locator, &manager.handle) {
		return nil, get_last_error()
	}
	runtime.AddCleanup(manager, func(handle unsafe.Pointer) {
		C.vpkc_free_update_manager(handle)
	}, manager.handle)
	return manager, nil
}

// AppID returns the currently installed app id.
func (up *UpdateManager) AppID() string {
	var len = C.vpkc_get_app_id(up.handle, nil, 0)
	var buf = make([]byte, len+1) // +1 for null terminator
	C.vpkc_get_app_id(up.handle, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len))
	return string(buf[:len]) // return the string without the null terminator
}

// CurrentlyInstalledVersion returns the currently installed version of the app.
func (up *UpdateManager) CurrentlyInstalledVersion() string {
	var len = C.vpkc_get_current_version(up.handle, nil, 0)
	var buf = make([]byte, len+1) // +1 for null terminator
	C.vpkc_get_current_version(up.handle, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len))
	return string(buf[:len]) // return the string without the null terminator
}

// IsPortable returns whether the app is in portable mode. On Windows this can be true or false.
// On MacOS and Linux this will always be true.
func (up *UpdateManager) IsPortable() bool {
	return bool(C.vpkc_is_portable(up.handle))
}

// UpdatePendingRestart returns an asset if there is an update downloaded which still needs to be
// applied. You can pass this asset to [UpdateManager.WaitExitThenApplyUpdates] to apply the update.
func (up *UpdateManager) UpdatePendingRestart() (*Asset, bool) {
	var asset *C.vpkc_asset_t
	if !C.vpkc_update_pending_restart(up.handle, &asset) {
		return nil, false
	}
	return toAsset(asset), true
}

// CheckForUpdates Checks for updates. If there are updates available, this method will return an [UpdateInfo]
// object containing the latest available release, and any delta updates that can be applied if they are available.
func (up *UpdateManager) CheckForUpdates() (*UpdateInfo, UpdateStatus, error) {
	var update_info *C.vpkc_update_info_t
	check_result := C.vpkc_check_for_updates(up.handle, &update_info)
	if update_info == nil {
		if UpdateStatus(check_result) == -1 {
			return nil, 0, get_last_error()
		}
		return nil, UpdateStatus(check_result), nil
	}
	var info UpdateInfo
	info.load(update_info)
	return &info, UpdateStatus(check_result), nil
}

func assetSilentRestart(options ...upto3[Silent, Restart, UnsafeProcessID]) (C.bool, int, **C.char, int, bool) {
	var silent C.bool
	var restart []*C.char
	var pid int
	var has_pid bool
	for _, option := range options {
		switch option := option.(type) {
		case Silent:
			silent = C.bool(option)
		case Restart:
			restart = make([]*C.char, len(option)+1) // +1 for null terminator
			for i, arg := range option {
				arg_cstr := C.CString(arg)
				defer C.free(unsafe.Pointer(arg_cstr))
				restart[i] = arg_cstr
			}
		case UnsafeProcessID:
			pid = int(option)
			has_pid = true
		}
	}
	var restartPtr **C.char
	if len(restart) > 0 {
		restartPtr = (**C.char)(unsafe.Pointer(&restart[0]))
	}
	return silent, len(restart), restartPtr, pid, has_pid
}

// WaitForExitThenApplyUpdates this will launch the Velopack updater and tell it to wait for this program
// to exit gracefully.
//   - You should then clean up any state and exit your app. The updater will apply updates and then
//   - (if [Restart] specified) restart your app. The updater will only wait for 60 seconds before giving up.
func (up *UpdateManager) WaitForExitThenApplyUpdates(update either[*UpdateInfo, *Asset], options ...upto3[Silent, Restart, UnsafeProcessID]) error {
	silent, restart, restartPtr, pid, has_pid := assetSilentRestart(options...)
	p_asset := (*C.vpkc_asset_t)(nil)
	if update != nil {
		switch update := update.(type) {
		case *Asset:
			p_asset = (*C.vpkc_asset_t)(update.handle)
		case *UpdateInfo:
			p_asset = (*C.vpkc_asset_t)(update.TargetFullRelease.handle)
		}
	}
	if has_pid {
		if !C.vpkc_unsafe_apply_updates(up.handle, p_asset, silent, C.uint32_t(pid), restart != 0, restartPtr, C.size_t(restart-1)) {
			return get_last_error()
		}
	}
	if !C.vpkc_wait_exit_then_apply_updates(up.handle, p_asset, silent, restart != 0, restartPtr, C.size_t(restart-1)) {
		return get_last_error()
	}
	return nil
}

/*
Create a new _CUSTOM_ update source with user-provided callbacks to fetch release feeds and download assets.
You can report download progress using [UpdateSource.ReportProgress].
*/
func NewSourceCustomCallback(callbacks SourceCustomCallbacks) (*UpdateSource, error) {
	if callbacks.ReleaseFeedFunc == nil {
		return nil, errors.New("ReleaseFeedFunc must not be nil")
	}
	if callbacks.DownloadAssetFunc == nil {
		return nil, errors.New("DownloadAssetFunc must not be nil")
	}
	var source = new(UpdateSource)
	source.handle = C.go_vpkc_new_source_custom_callback(C.uintptr_t(cgo.NewHandle(callbacks)))
	if source.handle == nil {
		return nil, get_last_error()
	}
	runtime.AddCleanup(source, func(handle unsafe.Pointer) {
		C.vpkc_free_source(handle)
	}, source.handle)
	return source, nil
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
	var progress_callback cgo.Handle = 0
	if progress != nil {
		progress_callback = cgo.NewHandle(progress)
	}
	var info_handle *C.vpkc_update_info_t
	if update_info != nil {
		info_handle = (*C.vpkc_update_info_t)(update_info.handle)
	}
	if !C.go_vpkc_download_updates(up.handle,
		info_handle,
		C.uintptr_t(progress_callback),
	) {
		return get_last_error()
	}
	update_info.load((*C.vpkc_update_info_t)(update_info.handle))
	return nil
}

var once atomic.Bool

// Run helps you to handle app activation events correctly.
// This should be used as early as possible in your application startup code.
// (eg. the beginning of main() or wherever your entry point is).
// This function will not return in some cases. Do not call this function more
// than once in your application.
func Run(a App) {
	if !once.CompareAndSwap(false, true) {
		panic("velopack.Run called more than once")
	}

	app = a
	C.vpkc_app_set_auto_apply_on_startup(C.bool(app.AutoApplyOnStartup))
	if app.Args != nil {
		var args = make([]*C.char, len(app.Args)+1) // +1 for null terminator
		for i, arg := range app.Args {
			arg_cstr := C.CString(arg)
			defer C.free(unsafe.Pointer(arg_cstr))
			args[i] = arg_cstr
		}
		C.vpkc_app_set_args(&args[0], C.size_t(len(app.Args)))
	}
	if app.Locator != nil {
		RootAppDir := C.CString(app.Locator.RootAppDir)
		defer C.free(unsafe.Pointer(RootAppDir))
		UpdateExePath := C.CString(app.Locator.UpdateExePath)
		defer C.free(unsafe.Pointer(UpdateExePath))
		PackagesDir := C.CString(app.Locator.PackagesDir)
		defer C.free(unsafe.Pointer(PackagesDir))
		ManifestPath := C.CString(app.Locator.ManifestPath)
		defer C.free(unsafe.Pointer(ManifestPath))
		CurrentBinaryDir := C.CString(app.Locator.CurrentBinaryDir)
		defer C.free(unsafe.Pointer(CurrentBinaryDir))
		C.vpkc_app_set_locator(&C.vpkc_locator_config_t{
			RootAppDir:       RootAppDir,
			UpdateExePath:    UpdateExePath,
			PackagesDir:      PackagesDir,
			ManifestPath:     ManifestPath,
			CurrentBinaryDir: CurrentBinaryDir,
			IsPortable:       C.bool(app.Locator.IsPortable),
		})
	}
	C.go_callbacks()
	if app.Logger != nil {
		C.go_set_logger()
	}
	C.vpkc_app_run(nil)
}
