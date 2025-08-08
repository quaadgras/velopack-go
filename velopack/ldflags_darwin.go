//go:build cgo && darwin

package velopack

// #cgo darwin amd64 LDFLAGS: -l:velopack_libc_osx_x64_gnu.a
// #cgo darwin arm64 LDFLAGS: -l:velopack_libc_osx_arm64_gnu.a
import "C"
