package main
import "fmt"


type Parser struct {
	l *Lexer

	curToken  Token
	peekToken Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, setting curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Interface that defines the AST Processes and implements
// a String() method for pretty-printing
type AstProcess interface {
	String() string
}

type Identifier struct {
	Token Token
	Value string
}
func (i *Identifier) String() string {
	return i.Value
}

type Prefix struct {
	Token Token
	Name  *Identifier
	Right AstProcess
}
func (p *Prefix) String() string {
	return fmt.Sprintf("%s.%s", p.Name.String(), p.Right.String())
}

type Recursion struct {
	Token Token
	Name  *Identifier
	Body  AstProcess
}
func (r *Recursion) String() string {
	return fmt.Sprintf("rec %s. %s", r.Name.String(), r.Body.String())
}

type Composition struct {
	Token Token
	Left  AstProcess
	Right AstProcess
}
func (c *Composition) String() string {
	return fmt.Sprintf("(%s|%s)", c.Left.String(), c.Right.String())
}

type Summation struct {
	Token Token
	Left  AstProcess
	Right AstProcess
}
func (s *Summation) String() string {
	return fmt.Sprintf("(%s + %s)", s.Left.String(), s.Right.String())
}

type Restriction struct {
	Token Token
	Name  *Identifier
	Body  AstProcess
}
func (r *Restriction) String() string {
	return fmt.Sprintf("%s\\%s", r.Body.String(), r.Name.String())
}

type Relabelling struct {
	Token Token
	Name  *Identifier
	Body  AstProcess
}
func (r *Relabelling) String() string {
	return fmt.Sprintf("%s[%s]", r.Body.String(), r.Name.String())
}

type Nil struct {
	Token Token
}
func (r *Nil) String() string {
	return fmt.Sprintf("nil")
}

/* Actual parser */

func (p *Parser) ParseAstProcess() AstProcess {
    var left AstProcess
    canBePrefix := false // only idenfiers can be prefixes, not any process

    switch p.curToken.Type {
    case LPAREN:
        p.nextToken()
        // when encountering an open parenthesis, recursively parse the process
        // contained between the parentheses
        left = p.ParseAstProcess()

        // consume closed parenthesis
        if p.curToken.Type != RPAREN {
            panic("expected ')' after expression")
        }
        p.nextToken()
    case IDENT:
        canBePrefix = true // an identifier can be a prefix action
        left = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
        p.nextToken()
    case REC:
        // rec requires an identifier
        p.nextToken()
        if p.curToken.Type != IDENT {
            panic("expected identifier after 'rec'")
        }
        identifier := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
        
        p.nextToken()
        if p.curToken.Type != DOT {
            panic("expected '.' after identifier")
        }

        // parse the rec body
        p.nextToken()
        body := p.ParseAstProcess()
        return &Recursion{Name: identifier, Body: body}
    case NIL:
        p.nextToken()
        return &Nil{}
    default:
        // TODO review
        left = p.ParseAstProcess()
    }

    // after parsing the left process, if it's not one of the base cases, check if
    // it's inside of a sum or composition, or if it has a relabeling or restriction.
    // if we parsed an identifier, check if it's a prefix
    for {
        switch p.curToken.Type {
        case DOT:
            if canBePrefix {
              p.nextToken()
              right := p.ParseAstProcess()
              left = &Prefix{Name: left.(*Identifier), Right: right}
            } else {
              panic("prefixes must be identifiers")
            }
        case PLUS:
            p.nextToken()
            right := p.ParseAstProcess()
            left = &Summation{Left: left, Right: right}
        case PIPE:
            p.nextToken()
            right := p.ParseAstProcess()
            left = &Composition{Left: left, Right: right}
        case BACKSLASH:
            p.nextToken()
            if p.curToken.Type != IDENT {
                panic("expected identifier after '\\'")
            }
            identifier := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
            p.nextToken()
            left = &Restriction{Body: left, Name: identifier}
        case BRACKET_OPEN:
            p.nextToken()
            if p.curToken.Type != IDENT {
                panic("expected identifier after '['")
            }
            identifier := &Identifier{Token: p.curToken, Value: p.curToken.Literal}
            p.nextToken()
            if p.curToken.Type != BRACKET_CLOSE {
                panic("expected ']' after identifier")
            }
            p.nextToken()
            left = &Relabelling{Body: left, Name: identifier}
        default:
            return left
        }
    }
}

