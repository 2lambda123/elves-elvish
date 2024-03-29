// Code generated by "stringer -type=PrimaryType,RedirMode,ExprCtx -output=zstring.go"; DO NOT EDIT.

package parse

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BadPrimary-0]
	_ = x[Bareword-1]
	_ = x[SingleQuoted-2]
	_ = x[DoubleQuoted-3]
	_ = x[Variable-4]
	_ = x[Wildcard-5]
	_ = x[Tilde-6]
	_ = x[ExceptionCapture-7]
	_ = x[OutputCapture-8]
	_ = x[List-9]
	_ = x[Lambda-10]
	_ = x[Map-11]
	_ = x[Braced-12]
}

const _PrimaryType_name = "BadPrimaryBarewordSingleQuotedDoubleQuotedVariableWildcardTildeExceptionCaptureOutputCaptureListLambdaMapBraced"

var _PrimaryType_index = [...]uint8{0, 10, 18, 30, 42, 50, 58, 63, 79, 92, 96, 102, 105, 111}

func (i PrimaryType) String() string {
	if i < 0 || i >= PrimaryType(len(_PrimaryType_index)-1) {
		return "PrimaryType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PrimaryType_name[_PrimaryType_index[i]:_PrimaryType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BadRedirMode-0]
	_ = x[Read-1]
	_ = x[Write-2]
	_ = x[ReadWrite-3]
	_ = x[Append-4]
}

const _RedirMode_name = "BadRedirModeReadWriteReadWriteAppend"

var _RedirMode_index = [...]uint8{0, 12, 16, 21, 30, 36}

func (i RedirMode) String() string {
	if i < 0 || i >= RedirMode(len(_RedirMode_index)-1) {
		return "RedirMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RedirMode_name[_RedirMode_index[i]:_RedirMode_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NormalExpr-0]
	_ = x[CmdExpr-1]
	_ = x[LHSExpr-2]
	_ = x[BracedElemExpr-3]
	_ = x[strictExpr-4]
}

const _ExprCtx_name = "NormalExprCmdExprLHSExprBracedElemExprstrictExpr"

var _ExprCtx_index = [...]uint8{0, 10, 17, 24, 38, 48}

func (i ExprCtx) String() string {
	if i < 0 || i >= ExprCtx(len(_ExprCtx_index)-1) {
		return "ExprCtx(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ExprCtx_name[_ExprCtx_index[i]:_ExprCtx_index[i+1]]
}
