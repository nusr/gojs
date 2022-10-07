package parser

import (
	"testing"

	"github.com/nusr/gojs/scanner"
)

func TestFunction(t *testing.T) {

	source := `
	var a = 1

	console.log(1)

	function add(a, b) {
		return a + b
	}

	add(1, 3)

	class Base {
		a = 1
		method(c) {
			return c - this.a
		}
	}

	var b = new Base()

	b = [1, 2, 4]
	

	b[100] = 3

	b[0]

	b = function(a, b) {
		return a + b
	}

	a = class Test{
		c = 1
	}

	a = true && false
	a = 1 || 3

	a = {
		c: 1,
		b: true,
	}

	b = 1
	var i = 1;
	while(i < 10) {
		if (!i) {
			b += i
		}
	}

	`
	s := scanner.New(source)
	tokens := s.Scan()
	p := New(tokens)
	list := p.Parse()

	expects := []string{
		"var a=1;",
		"console.log(1);",
		"function add(a,b){return a+b;}",
		"add(1,3);",
		"class Base{a=1;method(c){return c-this.a;}}",
		"var b=new Base();",
		"b=[1,2,4];",
		"b[100]=3;",
		"b[0];",
		"b=function(a,b){return a+b;};",
		"a=class Test{c=1;};",
		"a=true&&false;",
		"a=1||3;",
		"a={c:1,b:true};",
		"b=1;",
		"var i=1;",
		"while(i<10){if(!i){b=b+i;}}",
	}
	for i, item := range list {
		if item.String() != expects[i] {
			t.Errorf("expect: %v,actual: %v", expects[i], item)
		}
	}
}
