//go:build cgo && windows

package velopack

// #cgo amd64 LDFLAGS: -lvelopack_libc_windows_x64_msvc
// #cgo arm64 LDFLAGS: -lvelopack_libc_windows_arm64_msvc
import "C"
