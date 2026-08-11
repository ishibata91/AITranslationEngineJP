//go:build !windows

package api

import "os/exec"

// hideChildProcessWindow はWindows以外で子プロセスの起動設定を変更しない。
func hideChildProcessWindow(_ *exec.Cmd) {}
