package shell

import (
	"bytes"
	"testing"

	"github.com/elves/elvish/pkg/eval"
	"github.com/elves/elvish/pkg/eval/vals"
	"github.com/elves/elvish/pkg/eval/vars"
)

func TestExtractExports(t *testing.T) {
	ns := eval.Ns{
		exportsVarName: vars.NewReadOnly(vals.EmptyMap.Assoc("a", "lorem")),
	}
	extractExports(ns, &bytes.Buffer{})
	if ns.HasName(exportsVarName) {
		t.Errorf("$%s not removed", exportsVarName)
	}
	if ns["a"].Get() != "lorem" {
		t.Errorf("$a not extracted from exports")
	}
}

func TestExtractExports_IgnoreNonMapExports(t *testing.T) {
	ns := eval.Ns{
		exportsVarName: vars.NewReadOnly("x"),
	}
	var errBuf bytes.Buffer
	extractExports(ns, &errBuf)
	if errBuf.Len() == 0 {
		t.Errorf("No error written with non-map exports")
	}
}

func TestExtractExports_IgnoreNonStringKeys(t *testing.T) {
	ns := eval.Ns{
		exportsVarName: vars.NewReadOnly(vals.EmptyMap.Assoc(vals.EmptyList, "lorem")),
	}
	var errBuf bytes.Buffer
	extractExports(ns, &errBuf)
	if errBuf.Len() == 0 {
		t.Errorf("No error written with non-string key")
	}
}

func TestExtractExports_DoesNotOverwrite(t *testing.T) {
	ns := eval.Ns{
		"a":            vars.NewReadOnly("lorem"),
		exportsVarName: vars.NewReadOnly(vals.EmptyMap.Assoc("a", "ipsum")),
	}
	var errBuf bytes.Buffer
	extractExports(ns, &errBuf)
	if ns["a"].Get() != "lorem" {
		t.Errorf("Existing variable overwritten")
	}
	if errBuf.Len() == 0 {
		t.Errorf("No error written with name conflict")
	}
}
