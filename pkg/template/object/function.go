package object

import (
	"errors"
	"fmt"
)

var _ Symbol = &SymbolFunction{}

// a fnValue evaluates Fn to return a value.
type SymbolFunction struct {
	Name        string
	Description string

	// required number of arguments
	ValidArgs ValidArgs

	// NOTE: should be called with the Call method. This way you get args validation for free
	Fn func(...Object) (Object, error)
}

// Type implements Symbol.
func (s *SymbolFunction) Type() symbolType {
	return SymbolTypeFunction
}

// Type implements Symbol.
func (s *SymbolFunction) SymbolName() string {
	return s.Name
}

func (s *SymbolFunction) String() string {
	return fmt.Sprintf("[%s %s: %s]", s.Type(), s.Name, s.Description)
}

// _symbol implements Symbol.
func (s *SymbolFunction) _symbol() {}

// does validation checks, then calls Fn
func (s *SymbolFunction) Call(objs ...Object) (Object, error) {
	if s.ValidArgs != nil {
		if err := s.ValidArgs(objs); err != nil {
			return nil, err
		}
	}

	return s.Fn(objs...)
}

type ValidArgs func(args []Object) error

// expects n args
func ExactArgs(n int) ValidArgs {
	return func(args []Object) error {
		length := len(args)
		if n != length {
			return fmt.Errorf("expected=%d args, got=%d", n, length)
		}
		return nil
	}
}

// expects at least n args, but allows for more
func MinArgs(n int) ValidArgs {
	return func(args []Object) error {
		length := len(args)
		if length < n {
			return fmt.Errorf("expected at least %d args, got=%d", n, length)
		}
		return nil
	}
}

// expects a number of args between min and max
func MinMaxArgs(min, max int) ValidArgs {
	return func(args []Object) error {
		length := len(args)
		if length < min || length > max {
			return fmt.Errorf("expected %d-%d args, got=%d", min, max, length)
		}
		return nil
	}
}

// allows for any number and types of args
func AnyArgs() ValidArgs {
	return func([]Object) error {
		return nil
	}
}

func ExpectTypesArgs(types ...ObjectType) ValidArgs {
	return func(args []Object) error {
		if len(args) != len(types) {
			return fmt.Errorf("expect %d args, got=%d", len(types), len(args))
		}

		errs := []error{}
		for idx, typ := range types {
			if typ != args[idx].ObjectType() {
				errs = append(errs, fmt.Errorf("expected type=%s, got=%s", typ, args[idx].ObjectType()))
			}
		}
		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}
}
