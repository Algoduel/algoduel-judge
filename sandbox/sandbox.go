// Package sandbox запускает решение игрока в изолированном окружении.
//
// Песочница строится по слоям (см. judge/README.md). Сейчас реализован
// только слой 1 (Basic): таймаут, пустое окружение, лимит на вывод.
// Basic НЕ изолирует решение и годится только для разработки.
package sandbox

import (
	"context"
	"io"
	"time"
)

// Limits — ограничения на один запуск.
type Limits struct {
	WallTime    time.Duration // таймаут по реальному времени (страховка от зависаний)
	MemoryBytes int64         // применяется начиная со слоя cgroups
	MaxPids     int64         // применяется начиная со слоя cgroups
	OutputBytes int64         // максимальный размер stdout
}

// Spec — что и как запустить.
type Spec struct {
	Path   string   // путь к исполняемому файлу (или интерпретатору)
	Args   []string // аргументы
	Dir    string   // рабочая директория
	Stdin  io.Reader
	Limits Limits
}

// Status — чем закончился запуск.
type Status string

const (
	StatusOK            Status = "ok"
	StatusTimeLimit     Status = "time_limit"
	StatusMemoryLimit   Status = "memory_limit"
	StatusOutputLimit   Status = "output_limit"
	StatusRuntimeError  Status = "runtime_error"
	StatusInternalError Status = "internal_error"
)

// Usage — что получилось и сколько ресурсов ушло.
type Usage struct {
	Status       Status        `json:"status"`
	ExitCode     int           `json:"exit_code"`
	WallTime     time.Duration `json:"wall_time"`
	CPUTime      time.Duration `json:"cpu_time"`
	Instructions uint64        `json:"instructions"` // 0, пока не реализован слой perf
	MemoryPeak   int64         `json:"memory_peak"`  // 0, пока не реализован слой cgroups
	Stdout       []byte        `json:"-"`
	Stderr       []byte        `json:"-"`
}

// Runner — общий интерфейс всех реализаций песочницы.
// Следующие слои (namespaces, cgroups, seccomp, perf) реализуют его же,
// поэтому остальной код judge не меняется.
type Runner interface {
	Run(ctx context.Context, spec Spec) (Usage, error)
}
