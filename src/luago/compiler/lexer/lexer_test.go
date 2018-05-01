package lexer

import "testing"
import "github.com/stretchr/testify/assert"

func TestNextToken(t *testing.T) {
	lexer := NewLexer("str", `;,()[]{}+-*^%%&|#`)
	assertNextTokenKind(t, lexer, TOKEN_SEP_SEMI)
	assertNextTokenKind(t, lexer, TOKEN_SEP_COMMA)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LPAREN)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RPAREN)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LBRACK)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RBRACK)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LCURLY)
	assertNextTokenKind(t, lexer, TOKEN_SEP_RCURLY)
	assertNextTokenKind(t, lexer, TOKEN_OP_ADD)
	assertNextTokenKind(t, lexer, TOKEN_OP_MINUS)
	assertNextTokenKind(t, lexer, TOKEN_OP_MUL)
	assertNextTokenKind(t, lexer, TOKEN_OP_POW)
	assertNextTokenKind(t, lexer, TOKEN_OP_MOD)
	assertNextTokenKind(t, lexer, TOKEN_OP_MOD)
	assertNextTokenKind(t, lexer, TOKEN_OP_BAND)
	assertNextTokenKind(t, lexer, TOKEN_OP_BOR)
	assertNextTokenKind(t, lexer, TOKEN_OP_LEN)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken2(t *testing.T) {
	lexer := NewLexer("str", `... .. . :: : // / ~= ~ == = << <= < >> >= >`)
	assertNextTokenKind(t, lexer, TOKEN_VARARG)
	assertNextTokenKind(t, lexer, TOKEN_OP_CONCAT)
	assertNextTokenKind(t, lexer, TOKEN_SEP_DOT)
	assertNextTokenKind(t, lexer, TOKEN_SEP_LABEL)
	assertNextTokenKind(t, lexer, TOKEN_SEP_COLON)
	assertNextTokenKind(t, lexer, TOKEN_OP_IDIV)
	assertNextTokenKind(t, lexer, TOKEN_OP_DIV)
	assertNextTokenKind(t, lexer, TOKEN_OP_NE)
	assertNextTokenKind(t, lexer, TOKEN_OP_WAVE)
	assertNextTokenKind(t, lexer, TOKEN_OP_EQ)
	assertNextTokenKind(t, lexer, TOKEN_OP_ASSIGN)
	assertNextTokenKind(t, lexer, TOKEN_OP_SHL)
	assertNextTokenKind(t, lexer, TOKEN_OP_LE)
	assertNextTokenKind(t, lexer, TOKEN_OP_LT)
	assertNextTokenKind(t, lexer, TOKEN_OP_SHR)
	assertNextTokenKind(t, lexer, TOKEN_OP_GE)
	assertNextTokenKind(t, lexer, TOKEN_OP_GT)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_keywords(t *testing.T) {
	keywords := `
	and       break     do        else      elseif    end
	false     for       function  goto      if        in
	local     nil       not       or        repeat    return
	then      true      until     while
    `
	lexer := NewLexer("str", keywords)
	assertNextTokenKind(t, lexer, TOKEN_OP_AND)
	assertNextTokenKind(t, lexer, TOKEN_KW_BREAK)
	assertNextTokenKind(t, lexer, TOKEN_KW_DO)
	assertNextTokenKind(t, lexer, TOKEN_KW_ELSE)
	assertNextTokenKind(t, lexer, TOKEN_KW_ELSEIF)
	assertNextTokenKind(t, lexer, TOKEN_KW_END)
	assertNextTokenKind(t, lexer, TOKEN_KW_FALSE)
	assertNextTokenKind(t, lexer, TOKEN_KW_FOR)
	assertNextTokenKind(t, lexer, TOKEN_KW_FUNCTION)
	assertNextTokenKind(t, lexer, TOKEN_KW_GOTO)
	assertNextTokenKind(t, lexer, TOKEN_KW_IF)
	assertNextTokenKind(t, lexer, TOKEN_KW_IN)
	assertNextTokenKind(t, lexer, TOKEN_KW_LOCAL)
	assertNextTokenKind(t, lexer, TOKEN_KW_NIL)
	assertNextTokenKind(t, lexer, TOKEN_OP_NOT)
	assertNextTokenKind(t, lexer, TOKEN_OP_OR)
	assertNextTokenKind(t, lexer, TOKEN_KW_REPEAT)
	assertNextTokenKind(t, lexer, TOKEN_KW_RETURN)
	assertNextTokenKind(t, lexer, TOKEN_KW_THEN)
	assertNextTokenKind(t, lexer, TOKEN_KW_TRUE)
	assertNextTokenKind(t, lexer, TOKEN_KW_UNTIL)
	assertNextTokenKind(t, lexer, TOKEN_KW_WHILE)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_identifiers(t *testing.T) {
	identifiers := `_ __ ___ a _HW_ hello_world HelloWorld HELLO_WORLD`
	lexer := NewLexer("str", identifiers)
	assertNextIdentifier(t, lexer, "_")
	assertNextIdentifier(t, lexer, "__")
	assertNextIdentifier(t, lexer, "___")
	assertNextIdentifier(t, lexer, "a")
	assertNextIdentifier(t, lexer, "_HW_")
	assertNextIdentifier(t, lexer, "hello_world")
	assertNextIdentifier(t, lexer, "HelloWorld")
	assertNextIdentifier(t, lexer, "HELLO_WORLD")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_numbers(t *testing.T) {
	numbers := `
	3   345   0xff   0xBEBADA
	3.0     3.1416     314.16e-2     0.31416E1     34e1
	0x0.1E  0xA23p-4   0X1.921FB54442D18P+1
	3.	.3	00001
	`
	lexer := NewLexer("str", numbers)
	assertNextNumber(t, lexer, "3")
	assertNextNumber(t, lexer, "345")
	assertNextNumber(t, lexer, "0xff")
	assertNextNumber(t, lexer, "0xBEBADA")
	assertNextNumber(t, lexer, "3.0")
	assertNextNumber(t, lexer, "3.1416")
	assertNextNumber(t, lexer, "314.16e-2")
	assertNextNumber(t, lexer, "0.31416E1")
	assertNextNumber(t, lexer, "34e1")
	assertNextNumber(t, lexer, "0x0.1E")
	assertNextNumber(t, lexer, "0xA23p-4")
	assertNextNumber(t, lexer, "0X1.921FB54442D18P+1")
	assertNextNumber(t, lexer, "3.")
	assertNextNumber(t, lexer, ".3")
	assertNextNumber(t, lexer, "00001")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_comments(t *testing.T) {
	lexer := NewLexer("str", `
	--
	--[[]]
	a -- short comment
	+ --[[ long comment ]] b --[===[ long
	comment
	]===] - c
	--`)
	assertNextIdentifier(t, lexer, "a")
	assertNextTokenKind(t, lexer, TOKEN_OP_ADD)
	assertNextIdentifier(t, lexer, "b")
	assertNextTokenKind(t, lexer, TOKEN_OP_MINUS)
	assertNextIdentifier(t, lexer, "c")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_strings(t *testing.T) {
	strs := `
	[[]] [[ long string ]]
	[=[
long string]=]
	[===[long\z
	string]===]
	'' '"' 'short string'
	"" "'" "short string"
	'\a\b\f\n\r\t\v\\\"\''
	'\8 \08 \64 \122 \x08 \x7a \x7A \u{6211} zzz'
	'foo \z  
	

	bar'
	`
	lexer := NewLexer("str", strs)
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, " long string ")
	assertNextString(t, lexer, "long string")
	assertNextString(t, lexer, "long\\z\n\tstring")
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, "\"")
	assertNextString(t, lexer, "short string")
	assertNextString(t, lexer, "")
	assertNextString(t, lexer, "'")
	assertNextString(t, lexer, "short string")
	assertNextString(t, lexer, "\a\b\f\n\r\t\v\\\"'")
	assertNextString(t, lexer, "\b \b @ z \b z z 我 zzz")
	assertNextString(t, lexer, "foo bar")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
	assert.Equal(t, lexer.line, 15)
}

func TestNextToken_strings2(t *testing.T) {
	strs := `'\\' --'foo' 
	'\
'`
	lexer := NewLexer("str", strs)
	assertNextString(t, lexer, "\\")
	assertNextString(t, lexer, "\n")
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestNextToken_whiteSpaces(t *testing.T) {
	strs := "\r\n \r\n \n\r \n \r \n \t\v\f"
	lexer := NewLexer("str", strs)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
	assert.Equal(t, lexer.line, 7)
}

func TestNextToken_hw(t *testing.T) {
	src := `print("Hello, World!")`
	lexer := NewLexer("str", src)

	assertNextIdentifier(t, lexer, "print")
	assertNextTokenKind(t, lexer, TOKEN_SEP_LPAREN)
	assertNextString(t, lexer, "Hello, World!")
	assertNextTokenKind(t, lexer, TOKEN_SEP_RPAREN)
	assertNextTokenKind(t, lexer, TOKEN_EOF)
}

func TestLookAhead(t *testing.T) {
	src := `print("Hello, World!")`
	lexer := NewLexer("str", src)

	assert.Equal(t, lexer.LookAhead(), TOKEN_IDENTIFIER)
	lexer.NextToken()
	assert.Equal(t, lexer.LookAhead(), TOKEN_SEP_LPAREN)
	lexer.NextToken()
	assert.Equal(t, lexer.LookAhead(), TOKEN_STRING)
	lexer.NextToken()
	assert.Equal(t, lexer.LookAhead(), TOKEN_SEP_RPAREN)
	lexer.NextToken()
	assert.Equal(t, lexer.LookAhead(), TOKEN_EOF)
}

func TestErrors(t *testing.T) {
	testError(t, "?", "src:1: unexpected symbol near '?'")
	testError(t, "[===", "src:1: invalid long string delimiter near '[='")
	testError(t, "[==[xx", "src:1: unfinished long string or comment")
	testError(t, "'abc\\defg", "src:1: unfinished string")
	testError(t, "'abc\\defg'", "src:1: invalid escape sequence near '\\d'")
	testError(t, "'\\256'", "src:1: decimal escape too large near '\\256'")
	testError(t, "'\\u{11FFFF}'", "src:1: UTF-8 value too large near '\\u{11FFFF}'")
	testError(t, "'\\'", "src:1: unfinished string")
}

func testError(t *testing.T, chunk, expectedErr string) {
	err := safeNextToken(NewLexer("src", chunk))
	assert.Equal(t, err, expectedErr)
}

func assertNextTokenKind(t *testing.T, lexer *Lexer, expectedKind int) {
	_, kind, _ := lexer.NextToken()
	assert.Equal(t, kind, expectedKind)
}

func assertNextIdentifier(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.Equal(t, kind, TOKEN_IDENTIFIER)
	assert.Equal(t, token, expectedToken)
}

func assertNextNumber(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.Equal(t, kind, TOKEN_NUMBER)
	assert.Equal(t, token, expectedToken)
}

func assertNextString(t *testing.T, lexer *Lexer, expectedToken string) {
	_, kind, token := lexer.NextToken()
	assert.Equal(t, kind, TOKEN_STRING)
	assert.Equal(t, token, expectedToken)
}

func safeNextToken(lexer *Lexer) (err string) {
	// catch error
	defer func() {
		if r := recover(); r != nil {
			err = r.(string)
		}
	}()
	_, _, err = lexer.NextToken()
	return
}
