package sandbox

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBasicRunsAndCapturesOutput(t *testing.T) {
	u, err := Basic{}.Run(context.Background(), Spec{
		Path:   "/bin/sh",
		Args:   []string{"-c", "echo hello"},
		Limits: Limits{WallTime: 2 * time.Second, OutputBytes: 1 << 20},
	})
	if err != nil || u.Status != StatusOK || strings.TrimSpace(string(u.Stdout)) != "hello" {
		t.Fatalf("got status=%s out=%q err=%v", u.Status, u.Stdout, err)
	}
}

func TestBasicEmptyEnvironment(t *testing.T) {
	t.Setenv("SECRET_TOKEN", "leak")
	u, _ := Basic{}.Run(context.Background(), Spec{
		Path:   "/bin/sh",
		Args:   []string{"-c", "echo ${SECRET_TOKEN:-none}"},
		Limits: Limits{WallTime: 2 * time.Second, OutputBytes: 1 << 20},
	})
	if strings.TrimSpace(string(u.Stdout)) != "none" {
		t.Fatalf("environment leaked into solution: %q", u.Stdout)
	}
}

func TestBasicTimeLimitKillsChildren(t *testing.T) {
	u, _ := Basic{}.Run(context.Background(), Spec{
		Path:   "/bin/sh",
		Args:   []string{"-c", "sleep 10 & sleep 10"},
		Limits: Limits{WallTime: 300 * time.Millisecond, OutputBytes: 1 << 20},
	})
	if u.Status != StatusTimeLimit || u.WallTime > 3*time.Second {
		t.Fatalf("expected fast time limit, got status=%s wall=%v", u.Status, u.WallTime)
	}
}

func TestBasicOutputLimit(t *testing.T) {
	u, _ := Basic{}.Run(context.Background(), Spec{
		Path:   "/bin/sh",
		Args:   []string{"-c", "head -c 100000 /dev/zero"},
		Limits: Limits{WallTime: 2 * time.Second, OutputBytes: 1000},
	})
	if u.Status != StatusOutputLimit || len(u.Stdout) != 1000 {
		t.Fatalf("expected output limit, got status=%s len=%d", u.Status, len(u.Stdout))
	}
}
