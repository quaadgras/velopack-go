//go:build cgo && darwin

package velopack

// #cgo amd64 LDFLAGS: -lvelopack_libc_osx_x64_gnu
// #cgo arm64 LDFLAGS: -lvelopack_libc_osx_arm64_gnu
import "C"
