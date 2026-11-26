package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

// Basic — слой 1: таймаут, пустое окружение, лимит вывода, убийство всей группы процессов.
// Решение при этом видит сеть и файловую систему: только для разработки!
type Basic struct{}

func (Basic) Run(ctx context.Context, spec Spec) (Usage, error) {
	if spec.Limits.WallTime > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, spec.Limits.WallTime)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, spec.Path, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = []string{} // никаких DATABASE_URL и прочих секретов
	cmd.Stdin = spec.Stdin

	stdout := &limitedBuffer{limit: spec.Limits.OutputBytes}
	stderr := &limitedBuffer{limit: 64 << 10}
	cmd.Stdout, cmd.Stderr = stdout, stderr

	// Своя группа процессов, чтобы при таймауте убить и потомков (fork-бомбу).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = time.Second

	start := time.Now()
	err := cmd.Run()
	u := Usage{WallTime: time.Since(start), Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	if ps := cmd.ProcessState; ps != nil {
		u.CPUTime = ps.UserTime() + ps.SystemTime()
		u.ExitCode = ps.ExitCode()
	}

	var exitErr *exec.ExitError
	switch {
	case ctx.Err() == context.DeadlineExceeded:
		u.Status = StatusTimeLimit
	case stdout.overflow:
		u.Status = StatusOutputLimit
	case err == nil:
		u.Status = StatusOK
	case errors.As(err, &exitErr):
		u.Status = StatusRuntimeError
	default:
		u.Status = StatusInternalError
		return u, err
	}
	return u, nil
}

// limitedBuffer хранит не больше limit байт и запоминает, что вывод был обрезан.
// bytes.Buffer лежит в поле, а не встроен: иначе io.Copy нашёл бы у него ReadFrom
// и записал бы всё мимо нашего Write, в обход лимита.
type limitedBuffer struct {
	buf      bytes.Buffer
	limit    int64
	overflow bool
}

func (b *limitedBuffer) Bytes() []byte { return b.buf.Bytes() }

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.limit > 0 {
		room := b.limit - int64(b.buf.Len())
		if int64(len(p)) > room {
			b.overflow = true
			if room > 0 {
				b.buf.Write(p[:room])
			}
			return len(p), nil // молча выбрасываем лишнее, чтобы решение не упало на записи
		}
	}
	return b.buf.Write(p)
}
