//go:build cgo && windows

package velopack

// #cgo amd64 LDFLAGS: -lvelopack_libc_win_x64_gnu -lws2_32 -lgcc_s -lbcrypt
// #cgo arm64 LDFLAGS: -lvelopack_libc_win_arm64_gnu
import "C"
