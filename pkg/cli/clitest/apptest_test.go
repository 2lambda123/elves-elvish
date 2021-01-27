package clitest

import (
	"testing"

	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/cli/term"
)

func TestFixture(t *testing.T) {
	f := Setup(
		WithSpec(func(spec *cli.AppSpec) {
			spec.CodeAreaState.Buffer = cli.CodeBuffer{Content: "test", Dot: 4}
		}),
		WithTTY(func(tty TTYCtrl) {
			tty.SetSize(20, 30) // h = 20, w = 30
		}),
	)
	defer f.Stop()

	// Verify that the functions passed to Setup have taken effect.
	if cli.GetCodeBuffer(f.App).Content != "test" {
		t.Errorf("WithSpec did not work")
	}

	buf := f.MakeBuffer()
	// Verify that the WithTTY function has taken effect.
	if buf.Width != 30 {
		t.Errorf("WithTTY did not work")
	}

	f.TestTTY(t, "test", term.DotHere)

	f.App.Notify("something")
	f.TestTTYNotes(t, "something")

	f.App.CommitCode()
	if code, err := f.Wait(); code != "test" || err != nil {
		t.Errorf("Wait returned %q, %v", code, err)
	}
}
