package edit

import (
	"testing"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/ui"
)

func TestCommandMode(t *testing.T) {
	f := setup()
	defer f.Cleanup()

	evals(f.Evaler, `edit:insert:binding[Ctrl-'['] = $edit:command:start~`)
	feedInput(f.TTYCtrl, "echo")
	f.TTYCtrl.Inject(term.K('[', ui.Ctrl))
	f.TestTTY(t,
		"~> echo", Styles,
		"   vvvv", term.DotHere, "\n",
		" COMMAND ", Styles,
		"*********",
	)

	f.TTYCtrl.Inject(term.K('b'))
	f.TestTTY(t,
		"~> ", term.DotHere,
		"echo\n", Styles,
		"vvvv",
		" COMMAND ", Styles,
		"*********",
	)
}
