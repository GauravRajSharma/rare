package expressions

import (
	"strconv"
	"strings"
)

func stringComparator(equation func(string, string) string) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) < 2 {
			return stageError(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			val := args[0](context)
			for i := 1; i < len(args); i++ {
				val = equation(val, args[i](context))
			}

			return val
		})
	})
}

// Checks equality, and returns truthy if equals, and empty if not
func arithmaticEqualityHelper(test func(float64, float64) bool) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) != 2 {
			return stageError(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			left, err := strconv.ParseFloat(args[0](context), 64)
			if err != nil {
				return ErrorType
			}
			right, err := strconv.ParseFloat(args[1](context), 64)
			if err != nil {
				return ErrorType
			}

			if test(left, right) {
				return TruthyVal
			}
			return FalsyVal
		})
	})
}

func kfNot(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if Truthy(args[0](context)) {
			return FalsyVal
		}
		return TruthyVal
	})
}

// {and a b c ...}
func kfAnd(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			if arg(context) == FalsyVal {
				return FalsyVal
			}
		}
		return TruthyVal
	})
}

// {or a b c ...}
func kfOr(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			if arg(context) != FalsyVal {
				return TruthyVal
			}
		}
		return FalsyVal
	})
}

// {like string contains}
func kfLike(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.Contains(val, contains) {
			return val
		}
		return FalsyVal
	})
}
