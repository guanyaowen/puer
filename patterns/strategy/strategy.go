package strategy

import "github.com/guanyaowen/puer/util/maths"

type UniqueStrategyKey string

// Operator 策略控制者
type Operator interface {
	// Register 注册策略
	Register(UniqueStrategyKey, Strategy)

	// Apply 按key执行策略
	Apply(key UniqueStrategyKey, params interface{}) interface{}
}

type Strategy interface {
	Apply(params interface{}) interface{}
}

// NewOperator 新增一个策略控制器
func NewOperator(strategySize int) Operator {
	return newBaseOperator(strategySize)
}

func newBaseOperator(size int) *baseOperator {
	return &baseOperator{
		make(map[UniqueStrategyKey]Strategy, maths.Max(size, 0)),
	}
}

type baseOperator struct {
	strategyMap map[UniqueStrategyKey]Strategy
}

func (b *baseOperator) Register(key UniqueStrategyKey, strategy Strategy) {
	b.strategyMap[key] = strategy
}

func (b *baseOperator) Apply(key UniqueStrategyKey, params interface{}) interface{} {
	if strategy, ok := b.strategyMap[key]; ok {
		return strategy.Apply(params)
	}
	return nil
}
