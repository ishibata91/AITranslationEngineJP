//go:build windows

package api

import (
	"os/exec"
	"testing"
)

func TestHideChildProcessWindow(t *testing.T) {
	cmd := exec.Command("dotnet")
	hideChildProcessWindow(cmd)

	if cmd.SysProcAttr == nil || !cmd.SysProcAttr.HideWindow {
		t.Fatal("子プロセスのウィンドウを非表示に設定しなかった")
	}
}
