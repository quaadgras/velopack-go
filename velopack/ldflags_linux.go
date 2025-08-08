//go:build cgo && linux

package velopack

// #cgo amd64 LDFLAGS: -lvelopack_libc_linux_x64_gnu
// #cgo arm64 LDFLAGS: -lvelopack_libc_linux_arm64_gnu
import "C"
