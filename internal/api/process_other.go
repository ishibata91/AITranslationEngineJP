//go:build !windows

package api

import "os/exec"

// hideChildProcessWindow は非Windowsでは子プロセスの表示設定を変更しない。
func hideChildProcessWindow(_ *exec.Cmd) {}
