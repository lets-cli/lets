package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"

	"github.com/lets-cli/lets/internal/checksum"
	"github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/env"
)

func TestMain(m *testing.M) {
	// executor.runCmd calls env.DebugLevel() which panics if SetDebugLevel hasn't
	// been called first (production callers go through main.go which always calls it).
	env.SetDebugLevel(0)
	os.Exit(m.Run())
}

// invocation records a single ScriptRunner call.
type invocation struct {
	script  string
	command *config.Command
}

// RecordingRunner is a test double for ScriptRunner. It records every invocation
// in order and can be configured to return controlled errors on specific calls.
type RecordingRunner struct {
	mu     sync.Mutex
	calls  []invocation
	errors map[int]error // keyed by 0-based call index
}

func newRecordingRunner() *RecordingRunner {
	return &RecordingRunner{errors: make(map[int]error)}
}

// failOn configures the runner to return err on the Nth call (0-based).
// Must be called before Execute() — not safe for concurrent use.
func (r *RecordingRunner) failOn(callIndex int, err error) {
	r.errors[callIndex] = err
}

// run is the ScriptRunner implementation. Thread-safe for parallel command dispatch.
func (r *RecordingRunner) run(command *config.Command, script string) error {
	r.mu.Lock()
	idx := len(r.calls)
	r.calls = append(r.calls, invocation{script: script, command: command})
	r.mu.Unlock()

	// r.errors is written only during setup, before any goroutines start — no lock needed.
	return r.errors[idx]
}

func (r *RecordingRunner) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.calls)
}

func (r *RecordingRunner) allScripts() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.calls))
	for i, c := range r.calls {
		out[i] = c.script
	}
	return out
}

// newTestCfg returns a minimal *config.Config backed by a temp directory.
// It pre-creates DotLetsDir so that Config.CreateChecksumsDir() can use os.Mkdir
// (not os.MkdirAll) to create the checksums sub-directory inside it.
func newTestCfg(t *testing.T) *config.Config {
	t.Helper()
	dir := t.TempDir()
	dotLets := filepath.Join(dir, ".lets")
	if err := os.Mkdir(dotLets, 0o755); err != nil {
		t.Fatalf("newTestCfg: create .lets dir: %v", err)
	}
	return &config.Config{
		WorkDir:      dir,
		FilePath:     filepath.Join(dir, "lets.yaml"),
		Shell:        "sh",
		Commands:     config.Commands{},
		DotLetsDir:   dotLets,
		ChecksumsDir: filepath.Join(dotLets, "checksums"),
	}
}

// newCmd builds a *config.Command with one or more sequential scripts.
// SkipDocopts is always true — tests don't need docopt parsing.
func newCmd(name string, scripts ...string) *config.Command {
	cmds := make([]*config.Cmd, len(scripts))
	for i, s := range scripts {
		cmds[i] = &config.Cmd{Script: s}
	}
	return &config.Command{
		Name:        name,
		SkipDocopts: true,
		Cmds:        config.Cmds{Commands: cmds},
	}
}

// newParallelCmd builds a *config.Command whose scripts are dispatched in parallel.
func newParallelCmd(name string, scripts ...string) *config.Command {
	cmds := make([]*config.Cmd, len(scripts))
	for i, s := range scripts {
		cmds[i] = &config.Cmd{Name: fmt.Sprintf("%s_%d", name, i), Script: s}
	}
	return &config.Command{
		Name:        name,
		SkipDocopts: true,
		Cmds:        config.Cmds{Commands: cmds, Parallel: true},
	}
}

// execCtx wraps a command in a fresh top-level executor Context.
func execCtx(command *config.Command) *Context {
	return NewExecutorCtx(context.Background(), command)
}

// TestInitRunsOnce verifies that cfg.Init is executed on the first Execute() call
// and skipped on all subsequent calls to the same Executor.
func TestInitRunsOnce(t *testing.T) {
	cfg := newTestCfg(t)
	cfg.Init = "init-script"
	cmd := newCmd("foo", "foo-script")
	cfg.Commands["foo"] = cmd

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("first Execute: %v", err)
	}
	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("second Execute: %v", err)
	}

	scripts := r.allScripts()
	if len(scripts) < 2 {
		t.Fatalf("expected at least 2 calls, got %d: %v", len(scripts), scripts)
	}
	if scripts[0] != "init-script" {
		t.Errorf("first call must be init-script, got %q", scripts[0])
	}

	initCount := 0
	for _, s := range scripts {
		if s == "init-script" {
			initCount++
		}
	}
	if initCount != 1 {
		t.Errorf("init-script must run exactly once across both Execute() calls, ran %d times", initCount)
	}
}

// TestDependsRunInDeclarationOrder verifies that Depends are executed before the
// command's own scripts and in the order they were declared.
func TestDependsRunInDeclarationOrder(t *testing.T) {
	cfg := newTestCfg(t)

	dep1 := newCmd("dep1", "dep1-script")
	dep2 := newCmd("dep2", "dep2-script")
	cfg.Commands["dep1"] = dep1
	cfg.Commands["dep2"] = dep2

	deps := &config.Deps{}
	deps.Set("dep1", config.Dep{Name: "dep1"})
	deps.Set("dep2", config.Dep{Name: "dep2"})

	main := newCmd("main", "main-script")
	main.Depends = deps

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(main)); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	scripts := r.allScripts()
	want := []string{"dep1-script", "dep2-script", "main-script"}
	if len(scripts) != len(want) {
		t.Fatalf("expected %d calls %v, got %d: %v", len(want), want, len(scripts), scripts)
	}
	for i, s := range want {
		if scripts[i] != s {
			t.Errorf("call[%d]: want %q, got %q", i, s, scripts[i])
		}
	}
}

// TestDependencyFailureProducesDependencyError verifies that a failure inside a
// dependency is surfaced as a *DependencyError with the full command chain.
func TestDependencyFailureProducesDependencyError(t *testing.T) {
	cfg := newTestCfg(t)

	dep := newCmd("dep", "dep-script")
	cfg.Commands["dep"] = dep

	deps := &config.Deps{}
	deps.Set("dep", config.Dep{Name: "dep"})

	main := newCmd("main", "main-script")
	main.Depends = deps

	r := newRecordingRunner()
	r.failOn(0, fmt.Errorf("dep failed"))
	ex := NewExecutor(cfg, r.run)

	err := ex.Execute(execCtx(main))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var depErr *DependencyError
	if !errors.As(err, &depErr) {
		t.Fatalf("expected *DependencyError, got %T: %v", err, err)
	}
	if len(depErr.Chain) < 2 {
		t.Fatalf("expected chain of at least 2, got %v", depErr.Chain)
	}
	if depErr.Chain[0] != "main" {
		t.Errorf("chain[0]: want %q, got %q", "main", depErr.Chain[0])
	}
	if depErr.Chain[1] != "dep" {
		t.Errorf("chain[1]: want %q, got %q", "dep", depErr.Chain[1])
	}
}

// TestSequentialScriptsRunInOrder verifies that multiple cmd scripts in a single
// command are executed in declaration order.
func TestSequentialScriptsRunInOrder(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("multi", "script-a", "script-b", "script-c")

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	scripts := r.allScripts()
	want := []string{"script-a", "script-b", "script-c"}
	if len(scripts) != len(want) {
		t.Fatalf("expected %d calls %v, got %d: %v", len(want), want, len(scripts), scripts)
	}
	for i, s := range want {
		if scripts[i] != s {
			t.Errorf("call[%d]: want %q, got %q", i, s, scripts[i])
		}
	}
}

// TestAfterScriptRunsOnMainFailure verifies that the 'after' script is invoked
// even when the main command script fails.
func TestAfterScriptRunsOnMainFailure(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("cmd", "main-script")
	cmd.After = "after-script"

	r := newRecordingRunner()
	r.failOn(0, fmt.Errorf("exit 1"))
	ex := NewExecutor(cfg, r.run)

	err := ex.Execute(execCtx(cmd))
	if err == nil {
		t.Fatal("expected error from main script, got nil")
	}

	scripts := r.allScripts()
	found := false
	for _, s := range scripts {
		if s == "after-script" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("after-script was not invoked; recorded calls: %v", scripts)
	}
}

// TestAfterScriptErrorDoesNotOverrideMainError verifies that when both the main
// script and the 'after' script fail, Execute() returns the main script's error.
func TestAfterScriptErrorDoesNotOverrideMainError(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("cmd", "main-script")
	cmd.After = "after-script"

	mainErr := fmt.Errorf("main script failed")

	r := newRecordingRunner()
	r.failOn(0, mainErr)                        // main script fails
	r.failOn(1, fmt.Errorf("after script failed")) // after script also fails (logged, not returned)
	ex := NewExecutor(cfg, r.run)

	err := ex.Execute(execCtx(cmd))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// mainErr must be reachable via the error chain; the after error must not replace it.
	if !errors.Is(err, mainErr) {
		t.Errorf("returned error %v must wrap mainErr; after-script error must not override it", err)
	}
}

// TestParallelDispatchesAllScripts verifies that parallel mode invokes all scripts
// (order-independent).
func TestParallelDispatchesAllScripts(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newParallelCmd("par", "script-x", "script-y", "script-z")

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	got := r.allScripts()
	sort.Strings(got)
	want := []string{"script-x", "script-y", "script-z"}
	// want is already sorted

	if len(got) != len(want) {
		t.Fatalf("expected %d scripts %v, got %d: %v", len(want), want, len(got), got)
	}
	for i, s := range want {
		if got[i] != s {
			t.Errorf("script[%d]: want %q, got %q", i, s, got[i])
		}
	}
}

// TestParallelFailurePropagates verifies that a failure in one parallel script
// causes Execute() to return an error.
func TestParallelFailurePropagates(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newParallelCmd("par", "script-a", "script-b")

	r := newRecordingRunner()
	// Fail whichever script is dispatched first (index 0 in concurrent race).
	// The goal is that at least one failure causes Execute() to return an error.
	r.failOn(0, fmt.Errorf("one parallel script failed"))
	ex := NewExecutor(cfg, r.run)

	err := ex.Execute(execCtx(cmd))
	if err == nil {
		t.Fatal("expected error from parallel failure, got nil")
	}
}

// TestChecksumEnvVarsPresentInRunnerInvocation verifies that LETS_CHECKSUM_*
// env vars would be set in the runner invocation when a checksum is defined.
// RecordingRunner captures the command at call time; getChecksumEnvMap derives
// the expected env keys.
func TestChecksumEnvVarsPresentInRunnerInvocation(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("build", "build-script")
	// Pre-populate ChecksumMap directly — ChecksumCalculator is a no-op when
	// ChecksumSources is empty, so the map survives through initCmd unchanged.
	cmd.ChecksumMap = map[string]string{
		checksum.DefaultChecksumKey: "abc123",
	}

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if r.callCount() != 1 {
		t.Fatalf("expected 1 runner call, got %d", r.callCount())
	}

	capturedCmd := r.calls[0].command
	if capturedCmd.ChecksumMap[checksum.DefaultChecksumKey] != "abc123" {
		t.Errorf("ChecksumMap missing at runner call time: %v", capturedCmd.ChecksumMap)
	}

	// Verify the env-var key that ShellRunner would derive from this checksum.
	envMap := getChecksumEnvMap(capturedCmd.ChecksumMap)
	if _, ok := envMap["LETS_CHECKSUM"]; !ok {
		t.Errorf("expected LETS_CHECKSUM in computed env vars; got %v", envMap)
	}
}

// TestChecksumPersistedAfterSuccess verifies that a successful Execute() persists
// the checksum to disk when PersistChecksum is true.
func TestChecksumPersistedAfterSuccess(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("build", "build-script")
	cmd.PersistChecksum = true
	cmd.ChecksumMap = map[string]string{
		checksum.DefaultChecksumKey: "abc123",
	}

	r := newRecordingRunner()
	ex := NewExecutor(cfg, r.run)

	if err := ex.Execute(execCtx(cmd)); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if !checksum.IsChecksumForCmdPersisted(cfg.ChecksumsDir, "build") {
		t.Error("expected checksum to be persisted to disk after successful execution")
	}
}

// TestChecksumNotPersistedAfterFailure verifies that a failed Execute() does NOT
// persist the checksum to disk.
func TestChecksumNotPersistedAfterFailure(t *testing.T) {
	cfg := newTestCfg(t)
	cmd := newCmd("build", "build-script")
	cmd.PersistChecksum = true
	cmd.ChecksumMap = map[string]string{
		checksum.DefaultChecksumKey: "abc123",
	}

	r := newRecordingRunner()
	r.failOn(0, fmt.Errorf("build failed"))
	ex := NewExecutor(cfg, r.run)

	err := ex.Execute(execCtx(cmd))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if checksum.IsChecksumForCmdPersisted(cfg.ChecksumsDir, "build") {
		t.Error("expected checksum NOT to be persisted to disk after failed execution")
	}
}
