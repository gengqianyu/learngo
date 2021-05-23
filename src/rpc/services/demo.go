package services

import "errors"

type Demo struct {
}

type Args struct {
	A, B int
}

// Div rpc method example
func (d Demo) Div(args Args, result *float64) error {
	if args.B == 0 {
		return errors.New("divisor by zero")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}
