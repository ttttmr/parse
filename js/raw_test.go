package js

import (
	"io"
	"testing"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/test"
)

func TestRaw(t *testing.T) {
	var tests = []struct {
		js       string
		expected string
	}{
		// BlockStmt
		{"if (true) { x(1, 2, 3); };", "if (true) { x(1, 2, 3); }; "},

		// IfStmt
		{"if (true) { true; };", "if (true) { true; }; "},
		{"if (true) { true; } else { false; };", "if (true) { true; } else { false; }; "},
		{"if (true) { true; } else { if(true) { true; } else { false; } };", "if (true) { true; } else { if (true) { true; } else { false; }; }; "},

		// DoWhileStmt
		{"do { continue; } while (true);", "do { continue; } while (true); "},
		{"do { x = 1; } while (true);", "do { x = 1; } while (true); "},

		// WhileStmt
		{"while (true) { true; };", "while (true) { true; }; "},
		{"while (true) { x = 1; };", "while (true) { x = 1; }; "},

		// ForStmt
		{"for ( ; ; ) { true; };", "for ( ; ; ) { true; }; "},
		{"for (x = 1; ; ) { true; };", "for (x = 1; ; ) { true; }; "},
		{"for (x = 1; x < 2; ) { true; };", "for (x = 1; x < 2; ) { true; }; "},
		{"for (x = 1; x < 2; x++) { true; };", "for (x = 1; x < 2; x++) { true; }; "},
		{"for (x = 1; x < 2; x++) { x = 1; };", "for (x = 1; x < 2; x++) { x = 1; }; "},

		// ForInStmt
		{"for (var x in [1, 2]) { true; };", "for (var x in [1, 2]) { true; }; "},
		{"for (var x in [1, 2]) { x = 1; };", "for (var x in [1, 2]) { x = 1; }; "},

		// ForOfStmt
		{"for (const element of [1, 2]) { true; };", "for (const element of [1, 2]) { true; }; "},
		{"for (const element of [1, 2]) { x = 1; };", "for (const element of [1, 2]) { x = 1; }; "},

		// SwitchStmt
		{"switch (true) { case true: break; case false: false; };", "switch (true) { case true: break; case false: false; }; "},
		{"switch (true) { case true: x(); break; case false: x(); false; };", "switch (true) { case true: x(); break; case false: x(); false; }; "},
		{"switch (true) { default: false; };", "switch (true) { default: false; }; "},

		// BranchStmt
		{"for (i = 0; i < 3; i++) { continue; }; ", "for (i = 0; i < 3; i++) { continue; }; "},
		{"for (i = 0; i < 3; i++) { x = 1; }; ", "for (i = 0; i < 3; i++) { x = 1; }; "},

		// ReturnStmt
		{"return;", "return; "},
		{"return 1;", "return 1; "},

		// WithStmt
		{"with (true) { true; };", "with (true) { true; }; "},
		{"with (true) { x = 1; };", "with (true) { x = 1; }; "},

		// LabelledStmt
		{"loop: for (x = 0; x < 1; x++) { true; };", "loop: for (x = 0; x < 1; x++) { true; }; "},

		// ThrowStmt
		{"throw x;", "throw x; "},

		// TryStmt
		{"try { true; } catch(e) { };", "try { true; } catch(e) { }; "},
		{"try { true; } catch(e) { true; };", "try { true; } catch(e) { true; }; "},
		{"try { true; } catch(e) { x = 1; };", "try { true; } catch(e) { x = 1; }; "},

		// DebuggerStmt
		{"debugger;", "debugger; "},

		// Alias
		{"import * as name from 'module-name';", "import * as name from 'module-name'; "},

		// ImportStmt
		{"import defaultExport from 'module-name';", "import defaultExport from 'module-name'; "},
		{"import * as name from 'module-name';", "import * as name from 'module-name'; "},
		{"import { export1 } from 'module-name';", "import { export1 } from 'module-name'; "},
		{"import { export1 as alias1 } from 'module-name';", "import { export1 as alias1 } from 'module-name'; "},
		{"import { export1 , export2 } from 'module-name';", "import { export1 , export2 } from 'module-name'; "},
		{"import { foo , bar } from 'module-name/path/to/specific/un-exported/file';", "import { foo , bar } from 'module-name/path/to/specific/un-exported/file'; "},
		{"import defaultExport, * as name from 'module-name';", "import defaultExport , * as name from 'module-name'; "},
		{"import 'module-name';", "import 'module-name'; "},
		{"var promise = import('module-name');", "var promise = import('module-name'); "},

		// ExportStmt
		{"export { myFunction as default };", "export { myFunction as default }; "},
		{"export default k = 12;", "export default k = 12; "},

		// DirectivePrologueStmt
		{"'use strict';", "'use strict'; "},

		// BindingArray
		{"let [name = 5] = z;", "let [name = 5] = z; "},

		// BindingObject
		{"let {} = z;", "let { } = z; "},

		// BindingElement
		{"let [name = 5] = z;", "let [name = 5] = z; "},

		// VarDecl
		{"x = 1;", "x = 1; "},
		{"var x;", "var x; "},
		{"var x = 1;", "var x = 1; "},
		{"var x, y = [];", "var x, y = []; "},
		{"let x;", "let x; "},
		{"let x = 1;", "let x = 1; "},
		{"const x = 1;", "const x = 1; "},

		// Params
		{"function xyz(a, b) { };", "function xyz(a, b) { }; "},

		// FuncDecl
		{"function xyz(a, b) { };", "function xyz(a, b) { }; "},

		// MethodDecl
		{"class A { field; static get method () { }; };", "class A { field; static get method () { }; }; "},

		// FieldDefinition
		{"class A { field; };", "class A { field; }; "},
		{"class A { field = 5; };", "class A { field = 5; }; "},

		// ClassDecl
		{"class A { field; static get method () { }; };", "class A { field; static get method () { }; }; "},
		{"class B extends A { field; static get method () { }; };", "class B extends A { field; static get method () { }; }; "},

		// LiteralExpr
		{"'test';", "'test'; "},

		// ArrayExpr
		{"[1, 2, 3];", "[1, 2, 3]; "},

		// Property
		{`x = {x: "value"};`, `x = {x: "value"}; `},
		{`x = {"x": "value"};`, `x = {x: "value"}; `},

		// ObjectExpr
		{`x = {x: "value", y: "value"};`, `x = {x: "value", y: "value"}; `},

		// TemplateExpr
		{"x = `value`;", "x = `value`; "},
		{"x = `value${'hi'}`;", "x = `value${'hi'}`; "},

		// GroupExpr
		{"x = (1 + 1) / 1;", "x = (1 + 1) / 1; "},

		// IndexExpr
		{"x = y[1];", "x = y[1]; "},

		// DotExpr
		{"x = y.z;", "x = y.z; "},

		// NewTargetExpr
		{"x = new.target;", "x = new.target; "},

		// ImportMetaExpr
		{"x = import.meta;", "x = import.meta; "},

		// Args
		{"x(1, 2);", "x(1, 2); "},

		// NewExpr
		{"new x;", "new x; "},
		{"new x(1);", "new x(1); "},

		// CallExpr
		{"x();", "x(); "},

		// OptChainExpr
		{"x = y?.z;", "x = y?.z; "},

		// UnaryExpr
		{"x = 1 + 1;", "x = 1 + 1; "},

		// BinaryExpr
		{"a << b;", "a << b; "},

		// CondExpr
		{"a && b;", "a && b; "},
		{"a || b;", "a || b; "},

		// YieldExpr
		{"x = function* foo(x) { while (x < 2) { yield x; x++; }; };", "x = function* foo(x) { while (x < 2) { yield x; x++; }; }; "},

		// ArrowFunc
		{"(x) => { y(); };", "(x) => { y(); }; "},
		{"(x, y) => { z(); };", "(x, y) => { z(); }; "},
		{"async (x, y) => { z(); };", "async (x, y) => { z(); }; "},
	}
	for _, tt := range tests {
		t.Run(tt.js, func(t *testing.T) {
			ast, err := Parse(parse.NewInputString(tt.js))
			if err != io.EOF {
				test.Error(t, err)
			}
			test.String(t, ast.Raw(), tt.expected)
		})
	}
}

func TestRawRealWorldJS(t *testing.T) {
	js := `
	var _0x34d2=['log','415343ArZKCi','11ItVJMI','98599KfnlVw','139pQCPDx','526583DuLSJk','Hello\x20World!','5823JSTLxZ','543807ONUblA','2uewDkG','146389ygBdVV','2273BZpJsB'];(function(_0x6b2cbe,_0x3fd0f4){var _0xffa95b=_0x3d17;while(!![]){try{var _0x5239dc=-parseInt(_0xffa95b(0x132))+parseInt(_0xffa95b(0x12f))+-parseInt(_0xffa95b(0x131))*-parseInt(_0xffa95b(0x12c))+parseInt(_0xffa95b(0x135))*-parseInt(_0xffa95b(0x12e))+-parseInt(_0xffa95b(0x12d))+-parseInt(_0xffa95b(0x134))+-parseInt(_0xffa95b(0x133))*-parseInt(_0xffa95b(0x12b));if(_0x5239dc===_0x3fd0f4)break;else _0x6b2cbe['push'](_0x6b2cbe['shift']());}catch(_0x12a1ea){_0x6b2cbe['push'](_0x6b2cbe['shift']());}}}(_0x34d2,0x4d4a4));function hi(){var _0x2d5dfd=_0x3d17;console[_0x2d5dfd(0x12a)](_0x2d5dfd(0x130));}function _0x3d17(_0x59c992,_0x5be83e){_0x59c992=_0x59c992-0x12a;var _0x34d208=_0x34d2[_0x59c992];return _0x34d208;}hi();
	`

	ast, err := Parse(parse.NewInputString(js))
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// reparse to make sure JS is still valid
	_, err = Parse(parse.NewInputString(ast.Raw()))
	if err != nil && err == io.EOF {
		t.Error("Err: ", err)
	}
}
