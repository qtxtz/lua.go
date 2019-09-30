package state

import (
	"math"

	. "github.com/zxh0/lua.go/api"
	"github.com/zxh0/lua.go/number"
)

type operator struct {
	metamethod  string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd = func(a, b int64) int64 { return a + b }
	fadd = func(a, b float64) float64 { return a + b }
	isub = func(a, b int64) int64 { return a - b }
	fsub = func(a, b float64) float64 { return a - b }
	imul = func(a, b int64) int64 { return a * b }
	fmul = func(a, b float64) float64 { return a * b }
	div  = func(a, b float64) float64 { return a / b }
	iunm = func(a, _ int64) int64 { return -a }
	funm = func(a, _ float64) float64 { return -a }
	band = func(a, b int64) int64 { return a & b }
	bor  = func(a, b int64) int64 { return a | b }
	bxor = func(a, b int64) int64 { return a ^ b }
	bnot = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	operator{"__add", iadd, fadd},
	operator{"__sub", isub, fsub},
	operator{"__mul", imul, fmul},
	operator{"__mod", number.IMod, number.FMod},
	operator{"__pow", nil, math.Pow},
	operator{"__div", nil, div},
	operator{"__idiv", number.IFloorDiv, number.FFloorDiv},
	operator{"__band", band, nil},
	operator{"__bor", bor, nil},
	operator{"__bxor", bxor, nil},
	operator{"__shl", number.ShiftLeft, nil},
	operator{"__shr", number.ShiftRight, nil},
	operator{"__unm", iunm, funm},
	operator{"__bnot", bnot, nil},
}

// [-(2|1), +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_arith
func (state *luaState) Arith(op ArithOp) {
	var a, b luaValue // operands
	b = state.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = state.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		state.stack.push(result)
		return
	}

	mm := operator.metamethod
	if result, ok := callMetamethod(a, b, mm, state); ok {
		state.stack.push(result)
		return
	}

	if operator.floatFunc == nil {
		panic("number has no integer representation")
	}

	var typeName string
	if _, ok := convertToFloat(a); !ok {
		typeName = state.TypeName(typeOf(a))
	} else {
		typeName = state.TypeName(typeOf(b))
	}
	panic("attempt to perform arithmetic on a " + typeName + " value")
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil { // bitwise
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else { // arith
		if op.integerFunc != nil { // add,sub,mul,mod,idiv,unm
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}
