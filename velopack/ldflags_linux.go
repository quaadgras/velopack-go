//go:build cgo && linux

package velopack

// #cgo amd64 LDFLAGS: -l:velopack_libc_linux_x64_gnu.a
// #cgo arm64 LDFLAGS: -l:velopack_libc_linux_arm64_gnu.a
import "C"
