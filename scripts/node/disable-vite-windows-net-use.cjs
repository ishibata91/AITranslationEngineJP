const childProcess = require("node:child_process");

const originalExec = childProcess.exec;

function createNoopChild() {
  return {
    pid: 0,
    kill() {
      return true;
    },
    on() {
      return this;
    },
    once() {
      return this;
    },
    stdout: null,
    stderr: null,
    stdin: null
  };
}

childProcess.exec = function patchedExec(command, ...args) {
  if (process.platform === "win32" && typeof command === "string" && command.trim().toLowerCase() === "net use") {
    const callback = args.find((arg) => typeof arg === "function");

    queueMicrotask(() => {
      if (callback) {
        callback(new Error("disabled by preload"), "", "");
      }
    });

    return createNoopChild();
  }

  return originalExec.call(this, command, ...args);
};
