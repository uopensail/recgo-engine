package program

import (
	"fmt"
	"sync/atomic"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/uopensail/ulib/sample"
)

// Program represents a lazily compiled expr program.
// The program is compiled at most once successfully and then reused.
// Compilation may be attempted multiple times until the first success.
type Program struct {
	// expression is the original expr string
	expression string

	// options are compile-time expr options (copied on construction)
	options []expr.Option

	// program holds the successfully compiled vm.Program.
	// Once published, it is read-only and safe for concurrent use.
	program atomic.Pointer[vm.Program]
}

// NewProgram creates a Program with the given expression and compile options.
// Compilation is deferred until the first successful Eval call.
func NewProgram(expression string, opts ...expr.Option) (*Program, error) {
	return &Program{
		expression: expression,
		options:    opts,
	}, nil
}

// Eval evaluates the expression against one or more feature sets.
// The expression will be compiled on the first successful call and reused
// by all subsequent calls.
//
// It is safe for concurrent use.
func (prog *Program) Eval(datas ...sample.Features) (any, error) {
	// Build the runtime environment from input features.
	// Capacity is estimated to reduce rehashing.
	env := make(map[string]any, 16)

	// Feature-to-env transformer
	applyFeature := func(key string, feature sample.Feature) error {
		switch feature.Type() {
		case sample.Int64Type:
			env[key] = feature.GetInt64Unsafe()
		case sample.Int64sType:
			env[key] = feature.GetInt64sUnsafe()
		case sample.Float32Type:
			env[key] = float64(feature.GetFloat32Unsafe())
		case sample.Float32sType:
			src := feature.GetFloat32sUnsafe()
			dst := make([]float64, len(src))
			for i, v := range src {
				dst[i] = float64(v)
			}
			env[key] = dst
		case sample.StringType:
			env[key] = feature.GetStringUnsafe()
		case sample.StringsType:
			env[key] = feature.GetStringsUnsafe()
		default:
			return fmt.Errorf("key:%s unsupported data type:%d", key, feature.Type())
		}
		return nil
	}

	// Merge all feature sets into one environment
	for _, data := range datas {
		if err := data.ForEach(applyFeature); err != nil {
			return nil, err
		}
	}

	// Fast path: already compiled
	if p := prog.program.Load(); p != nil {
		return expr.Run(p, env)
	}

	// Slow path: attempt to compile using current environment as schema
	opts := make([]expr.Option, 0, len(prog.options)+1)
	opts = append(opts, expr.Env(env))
	opts = append(opts, prog.options...)

	compiled, err := expr.Compile(prog.expression, opts...)
	if err != nil {
		return nil, err
	}

	// Publish the program atomically.
	// Only one goroutine will succeed; others discard their result.
	if prog.program.CompareAndSwap(nil, compiled) {
		return expr.Run(compiled, env)
	}

	// Another goroutine has already published the program.
	return expr.Run(prog.program.Load(), env)
}
