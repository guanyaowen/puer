package strategy

import (
	"fmt"
	"testing"
)

var (
	AdditionKey       UniqueStrategyKey = "Addition"
	MultiplicationKey UniqueStrategyKey = "Multiplication"
)

type Params struct {
	x int
	y int
}

type Addition struct{}

func (a Addition) Apply(params interface{}) interface{} {
	p, ok := params.(*Params)
	if ok {
		return p.x + p.y
	}
	return 0
}

type Multiplication struct{}

func (a Multiplication) Apply(params interface{}) interface{} {
	p, ok := params.(*Params)
	if ok {
		return p.x * p.y
	}
	return 0
}

func Test_baseOperator_Apply(t *testing.T) {
	operator := NewOperator(5)

	operator.Register(AdditionKey, Addition{})
	operator.Register(MultiplicationKey, Multiplication{})

	params := &Params{
		x: 10,
		y: 10,
	}

	fmt.Println(operator.Apply(AdditionKey, params))
	fmt.Println(operator.Apply(MultiplicationKey, params))
}
