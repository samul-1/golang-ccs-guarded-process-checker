package main

/**
* A struct that represents a CCS process. The `kind` attribute contains the
* type of process (e.g. sum, restriction...)
 */
type Process struct {
	kind        string
	name        string
	left        *Process
	right       *Process
	restriction string
}

func NewProcess(kind, name string, left, right *Process, restriction string) *Process {
	return &Process{
		kind:        kind,
		name:        name,
		left:        left,
		right:       right,
		restriction: restriction,
	}
}

/**
* G(nil, X) = true
* G(x, X) = x not in X
* G(a.P, X) = G(P, emptyset)
* G(P\a, X) =  G(P, X)
* G(P[f], X) = G(P, X)
* G(P+Q X) = G(P, X) && G(Q, X)
* G(P|Q X) = G(P, X) && G(Q, X)
* G(rec x. P, X) = G(P, X U {x})

* P guarded iff G(P, emptyset) = true
**/

// Go doesn't have a set type, so we'll use a map[string]bool instead
func (p *Process) IsGuarded(X map[string]bool, resultChan chan<- bool) {
	switch p.kind {
	case "nil":
		// nil is always guarded
		resultChan <- true
	case "variable":
		// process variable is guarded iff it is not in X
		resultChan <- !X[p.name]
	case "prefix":
		// prefix is guarded iff its continuation is guarded
		go p.left.IsGuarded(map[string]bool{}, resultChan)
	case "sum", "composition":
		// sum and composition guarded iff both sides are guarded
		// spawn goroutines to check both sides
		leftChan := make(chan bool)
		rightChan := make(chan bool)
		go p.left.IsGuarded(X, leftChan)
		go p.right.IsGuarded(X, rightChan)

		// wait for both sides to finish and then combine results
		select { // we use select to avoid enforcing an order on the results
		case left := <-leftChan:
			resultChan <- (left && <-rightChan)
		case right := <-rightChan:
			resultChan <- (<-leftChan && right)
		}
	case "restriction", "relabeling":
		// Restriction and relabeling guarded iff P is guarded
		go p.left.IsGuarded(X, resultChan)
	case "recursion":
		// newX = X U {x}
		newX := make(map[string]bool)
		for k, v := range X {
			newX[k] = v // copy X into newX ...
		}
		newX[p.name] = true // ...and add x to newX

		// Recursion is guarded iff P is guarded with X U {x}
		go p.left.IsGuarded(newX, resultChan)
	default:
		panic("invalid process kind")
	}
}
