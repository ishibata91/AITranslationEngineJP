//go:build !windows

package api

import (
	"os/exec"
	"testing"
)

func TestHideChildProcessWindowDoesNothingOutsideWindows(t *testing.T) {
	cmd := exec.Command("echo", "test") //nolint:gosec // テスト用の固定コマンド。
	hideChildProcessWindow(cmd)
	if cmd.SysProcAttr != nil {
		t.Fatalf("非Windowsで子プロセス設定を変更した: %+v", cmd.SysProcAttr)
	}
}
