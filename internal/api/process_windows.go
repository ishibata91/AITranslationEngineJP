//go:build windows

package api

import (
	"os/exec"
	"syscall"
)

// hideChildProcessWindow は GUI アプリから起動する console 子プロセスのウィンドウを表示しない。
func hideChildProcessWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
