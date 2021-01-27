package edit

import (
	"testing"

	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/store"
)

func TestCommandHistory(t *testing.T) {
	f := setup(storeOp(func(s store.Store) {
		s.AddCmd("echo 1")
		s.AddCmd("echo 2")
	}))
	defer f.Cleanup()

	// TODO(xiaq): Test session history too.
	evals(f.Evaler, `@cmds = (edit:command-history)`)
	testGlobal(t, f.Evaler,
		"cmds",
		vals.MakeList(
			vals.MakeMap("id", "1", "cmd", "echo 1"),
			vals.MakeMap("id", "2", "cmd", "echo 2")))
}

func TestInsertLastWord(t *testing.T) {
	f := setup(storeOp(func(s store.Store) {
		s.AddCmd("echo foo bar")
	}))
	defer f.Cleanup()

	evals(f.Evaler, "edit:insert-last-word")
	wantBuf := cli.CodeBuffer{Content: "bar", Dot: 3}
	if buf := cli.GetCodeBuffer(f.Editor.app); buf != wantBuf {
		t.Errorf("buf = %v, want %v", buf, wantBuf)
	}
}
