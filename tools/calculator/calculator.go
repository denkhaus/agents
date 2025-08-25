package calculator

import (
	"context"
	"strings"

	"github.com/samber/do"
	"trpc.group/trpc-go/trpc-agent-go/tool"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// -----------------------------------------------------------------------------
// Constants & helpers ----------------------------------------------------------
// -----------------------------------------------------------------------------
const (
	ToolName = "calculator"
)

type Calculator interface {
}

// Constants for supported calculator operations.
const (
	opAdd      = "add"
	opSubtract = "subtract"
	opMultiply = "multiply"
	opDivide   = "divide"
)

// calculatorArgs holds the input for the calculator tool.
type calculatorArgs struct {
	Operation string  `json:"operation" description:"The operation: add, subtract, multiply, divide"`
	A         float64 `json:"a" description:"First number"`
	B         float64 `json:"b" description:"Second number"`
}

// calculatorResult holds the output for the calculator tool.
type calculatorResult struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
	Result    float64 `json:"result"`
}

// Calculator tool implementation.
// calculate performs the requested mathematical operation.
// It supports add, subtract, multiply, and divide operations.
func calculate(ctx context.Context, args calculatorArgs) (calculatorResult, error) {
	var result float64
	// Select operation based on input.
	switch strings.ToLower(args.Operation) {
	case opAdd:
		result = args.A + args.B
	case opSubtract:
		result = args.A - args.B
	case opMultiply:
		result = args.A * args.B
	case opDivide:
		if args.B != 0 {
			result = args.A / args.B
		}
	}
	return calculatorResult{
		Operation: args.Operation,
		A:         args.A,
		B:         args.B,
		Result:    result,
	}, nil
}

func NewTool(i *do.Injector) (tool.CallableTool, error) {
	// Create calculator tool for mathematical operations.
	calculatorTool := function.NewFunctionTool(
		calculate,
		function.WithName(ToolName),
		function.WithDescription(
			// Perform basic mathematical calculations (add, subtract, multiply, divide).
			"Perform basic mathematical calculations "+
				"(add, subtract, multiply, divide)",
		),
	)

	return calculatorTool, nil
}
