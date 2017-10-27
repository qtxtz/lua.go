package codegen

import . "luago/compiler/ast"
import . "luago/compiler/lexer"
import . "luago/vm"

// kind of operands
const (
	ARG_CONST  = 1 // const index
	ARG_REG    = 2 // register index
	ARG_UPVAL  = 4 // upvalue index
	ARG_RK     = ARG_REG | ARG_CONST
	ARG_RU     = ARG_REG | ARG_UPVAL
	ARG_RUK    = ARG_REG | ARG_UPVAL | ARG_CONST
)

// todo: rename to evalExp()?
func (self *codeGen) cgExp(node Exp, a, n int) {
	switch exp := node.(type) {
	case *NilExp:
		self.emitLoadNil(exp.Line, a, n)
	case *FalseExp:
		self.emitLoadBool(exp.Line, a, 0, 0)
	case *TrueExp:
		self.emitLoadBool(exp.Line, a, 1, 0)
	case *IntegerExp:
		self.emitLoadK(exp.Line, a, exp.Val)
	case *FloatExp:
		self.emitLoadK(exp.Line, a, exp.Val)
	case *StringExp:
		self.emitLoadK(exp.Line, a, exp.Str)
	case *VarargExp:
		self.emitVararg(exp.Line, a, n)
	case *ParensExp:
		self.cgExp(exp.Exp, a, 1)
	case *NameExp:
		self.cgNameExp(exp, a)
	case *TableConstructorExp:
		self.cgTableConstructorExp(exp, a)
	case *FuncDefExp:
		self.cgFuncDefExp(exp, a)
	case *FuncCallExp:
		self.cgFuncCallExp(exp, a, n)
	case *BracketsExp:
		self.cgBracketsExp(exp, a)
	case *ConcatExp:
		self.cgConcatExp(exp, a)
	case *UnopExp:
		self.cgUnopExp(exp, a)
	case *BinopExp:
		self.cgBinopExp(exp, a)
	}
}

func (self *codeGen) cgTableConstructorExp(node *TableConstructorExp, a int) {
	nArr := node.NArr
	nExps := len(node.KeyExps)
	multRet := nExps > 0 &&
		isVarargOrFuncCallExp(node.ValExps[nExps-1])

	self.emitNewTable(node.Line, a, nArr, nExps - nArr)

	for i, keyExp := range node.KeyExps {
		valExp := node.ValExps[i]

		if nArr > 0 {
			if idx, ok := keyExp.(int); ok {
				_a := self.allocReg()
				if i == nExps-1 && multRet {
					self.cgExp(valExp, _a, -1)
				} else {
					self.cgExp(valExp, _a, 1)
				}

				if idx%50 == 0 || idx == nArr { // LFIELDS_PER_FLUSH
					if idx%50 == 0 {
						self.freeRegs(50)
					} else {
						self.freeRegs(idx%50)
					}
					line := lastLineOfExp(valExp)
					if i == nExps-1 && multRet {
						self.emitSetList(line, a, 0, idx/50 + 1)
					} else {
						self.emitSetList(line, a, idx%50, idx/50 + 1)
					}
				}

				continue
			}
		}

		b := self.allocReg()
		self.cgExp(keyExp, b, 1)
		c := self.allocReg()
		self.cgExp(valExp, c, 1)
		self.freeRegs(2)

		line := lastLineOfExp(valExp)
		self.emitSetTable(line, a, b, c)
	}
}

// f[a] := function(args) body end
func (self *codeGen) cgFuncDefExp(node *FuncDefExp, a int) {
	bx := self.genSubProto(node)
	self.emitClosure(node.LastLine, a, bx)
}

// r[a] := f(args)
func (self *codeGen) cgFuncCallExp(node *FuncCallExp, a, n int) {
	nArgs := self.prepFuncCall(node, a)
	self.emitCall(node.Line, a, nArgs, n)
}

// return f(args)
func (self *codeGen) cgTailCallExp(node *FuncCallExp, a int) {
	nArgs := self.prepFuncCall(node, a)
	self.emitTailCall(node.Line, a, nArgs)
}

func (self *codeGen) prepFuncCall(node *FuncCallExp, a int) int {
	nArgs := len(node.Args)
	lastArgIsVarargOrFuncCall := false

	self.cgExp(node.PrefixExp, a, 1)
	if node.MethodName != "" {
		self.allocReg()
		idx := self.indexOfConstant(node.MethodName)
		self.emitSelf(node.Line, a, a, idx)
	}
	for i, arg := range node.Args {
		tmp := self.allocReg()
		if i == nArgs-1 && isVarargOrFuncCallExp(arg) {
			lastArgIsVarargOrFuncCall = true
			self.cgExp(arg, tmp, -1)
		} else {
			self.cgExp(arg, tmp, 1)
		}
	}
	self.freeRegs(nArgs)

	if lastArgIsVarargOrFuncCall {
		nArgs = -1
	}
	if node.MethodName != "" {
		self.freeReg()
		nArgs++
	}

	return nArgs
}

// r[a] := name
func (self *codeGen) cgNameExp(node *NameExp, a int) {
	if r := self.indexOfLocVar(node.Name); r >= 0 {
		self.emitMove(node.Line, a, r)
	} else if idx := self.indexOfUpval(node.Name); idx >= 0 {
		self.emitGetUpval(node.Line, a, idx)
	} else { // x => _ENV['x']
		bracketsExp := &BracketsExp{
			Line:      node.Line,
			PrefixExp: &NameExp{node.Line, "_ENV"},
			KeyExp:    &StringExp{node.Line, node.Name},
		}
		self.cgBracketsExp(bracketsExp, a)
	}
}

// r[a] := prefix[key]
func (self *codeGen) cgBracketsExp(node *BracketsExp, a int) {
	oldRegs := self.usedRegs()
	b, kindB := self.expToOpArg(node.PrefixExp, ARG_RU)
	c, _ := self.expToOpArg(node.KeyExp, ARG_RK)
	self.freeRegs(self.usedRegs() - oldRegs)

	if kindB == ARG_UPVAL {
		self.emitGetTabUp(node.Line, a, b, c)
	} else {
		self.emitGetTable(node.Line, a, b, c)
	}
}

// r[a] := op exp
func (self *codeGen) cgUnopExp(node *UnopExp, a int) {
	oldRegs := self.usedRegs()
	b, _ := self.expToOpArg(node.Exp, ARG_REG)
	self.emitUnaryOp(node.Line, node.Op, a, b)
	self.freeRegs(self.usedRegs() - oldRegs)
}

// r[a] := exp1 op exp2
func (self *codeGen) cgBinopExp(node *BinopExp, a int) {
	switch node.Op {
	case TOKEN_OP_AND, TOKEN_OP_OR:
		oldRegs := self.usedRegs()

		b, _ := self.expToOpArg(node.Exp1, ARG_REG)
		self.freeRegs(self.usedRegs() - oldRegs)
		if node.Op == TOKEN_OP_AND {
			self.emitTestSet(node.Line, a, b, 0)
		} else {
			self.emitTestSet(node.Line, a, b, 1)
		}
		pcOfJmp := self.emitJmp(node.Line, 0)

		b, _ = self.expToOpArg(node.Exp2, ARG_REG)
		self.freeRegs(self.usedRegs() - oldRegs)
		self.emitMove(node.Line, a, b)		
		self.fixSbx(pcOfJmp, self.pc()-pcOfJmp)
	default:
		oldRegs := self.usedRegs()
		b, _ := self.expToOpArg(node.Exp1, ARG_RK)
		c, _ := self.expToOpArg(node.Exp2, ARG_RK)
		self.emitBinaryOp(node.Line, node.Op, a, b, c)
		self.freeRegs(self.usedRegs() - oldRegs)
	}
}

// r[a] := exp1 .. exp2
func (self *codeGen) cgConcatExp(node *ConcatExp, a int) {
	for _, subExp := range node.Exps {
		a := self.allocReg()
		self.cgExp(subExp, a, 1)
	}

	c := self.usedRegs() - 1
	b := c - len(node.Exps) + 1
	self.freeRegs(c - b + 1)
	self.emit(node.Line, OP_CONCAT, a, b, c)
}

func (self *codeGen) expToOpArg(node Exp, argKinds int) (arg, argKind int) {
	if argKinds&ARG_CONST > 0 {
		switch x := node.(type) {
		case *NilExp:
			return self.indexOfConstant(nil), ARG_CONST
		case *FalseExp:
			return self.indexOfConstant(false), ARG_CONST
		case *TrueExp:
			return self.indexOfConstant(true), ARG_CONST
		case *IntegerExp:
			return self.indexOfConstant(x.Val), ARG_CONST
		case *FloatExp:
			return self.indexOfConstant(x.Val), ARG_CONST
		case *StringExp:
			return self.indexOfConstant(x.Str), ARG_CONST
		}
	}
	if argKinds&ARG_REG > 0 {
		if nameExp, ok := node.(*NameExp); ok {
			if r := self.indexOfLocVar(nameExp.Name); r >= 0 {
				return r, ARG_REG
			}
		}
	}
	if argKinds&ARG_UPVAL > 0 {
		if nameExp, ok := node.(*NameExp); ok {
			if idx := self.indexOfUpval(nameExp.Name); idx >= 0 {
				return idx, ARG_UPVAL
			}
		}
	}
	a := self.allocReg()
	self.cgExp(node, a, 1)
	return a, ARG_REG
}
