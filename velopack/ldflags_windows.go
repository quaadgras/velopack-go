//go:build cgo && windows

package velopack

// #cgo amd64 LDFLAGS: -l:velopack_libc_windows_x64_msvc.lib
// #cgo arm64 LDFLAGS: -l:velopack_libc_windows_arm64_msvc.lib
import "C"
