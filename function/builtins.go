package function

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func add[T float32 | int64](a, b T) T {
	return a + b
}

func sub[T float32 | int64](a, b T) T {
	return a - b
}

func mul[T float32 | int64](a, b T) T {
	return a * b
}

func div[T float32 | int64](a, b T) T {
	return a / b
}

func abs[T float32 | int64](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func mod(a, b int64) int64 {
	return a % b
}

func round(val float32) int64 {
	return int64(math.Round(float64(val)))
}

func floor(val float32) int64 {
	return int64(math.Floor(float64(val)))
}

func ceil(val float32) int64 {
	return int64(math.Ceil(float64(val)))
}

func exp(val float32) float32 {
	return float32(math.Exp(float64(val)))
}

func log(val float32) float32 {
	return float32(math.Log(float64(val)))
}

func log10(val float32) float32 {
	return float32(math.Log10(float64(val)))
}

func log2(val float32) float32 {
	return float32(math.Log2(float64(val)))
}

func sqrt(val float32) float32 {
	return float32(math.Sqrt(float64(val)))
}

func sin(val float32) float32 {
	return float32(math.Sin(float64(val)))
}

func asin(val float32) float32 {
	return float32(math.Asin(float64(val)))
}

func asinh(val float32) float32 {
	return float32(math.Asinh(float64(val)))
}

func sinh(val float32) float32 {
	return float32(math.Sinh(float64(val)))
}

func cos(val float32) float32 {
	return float32(math.Cos(float64(val)))
}

func acos(val float32) float32 {
	return float32(math.Acos(float64(val)))
}

func cosh(val float32) float32 {
	return float32(math.Cosh(float64(val)))
}

func acosh(val float32) float32 {
	return float32(math.Acosh(float64(val)))
}

func tan(val float32) float32 {
	return float32(math.Tan(float64(val)))
}

func tanh(val float32) float32 {
	return float32(math.Tanh(float64(val)))
}

func atan(val float32) float32 {
	return float32(math.Atan(float64(val)))
}

func atanh(val float32) float32 {
	return float32(math.Atanh(float64(val)))
}

func sigmoid(val float32) float32 {
	return float32(1.0 / (1.0 + math.Exp(-float64(val))))
}

func pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func reverse(val string) string {
	runes := []rune(val)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func trim(val string) string {
	return strings.TrimSpace(val)
}

func substr(val string, start, len int) string {
	return val[start : start+len-1]
}

func upper(val string) string {
	return strings.ToUpper(val)
}

func lower(val string) string {
	return strings.ToLower(val)
}

func concat(args ...string) string {
	var str strings.Builder
	for i := 0; i < len(args); i++ {
		str.WriteString(args[i])
	}
	return str.String()
}

func from_unixtime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func unix_timestamp(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return t.Unix()
}

func date_add(date string, n int64) string {
	t, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return t.Add(time.Duration(n*86400) * time.Second).Format("2006-01-02")
}

func date_sub(date string, n int64) string {
	t, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return t.Add(time.Duration(n*(-86400)) * time.Second).Format("2006-01-02")
}

func date_diff(a, b string) int64 {
	first, _ := time.ParseInLocation("2006-01-02", a, time.Local)
	second, _ := time.ParseInLocation("2006-01-02", b, time.Local)
	return int64(math.Floor(first.Sub(second).Seconds() / 86400.0))
}

func day() string {
	return fmt.Sprintf("%d", time.Now().Day())
}

func date() string {
	return time.Now().Format("2006-01-02")
}

func year() string {
	return fmt.Sprintf("%d", time.Now().Year())
}

func month() string {
	return fmt.Sprintf("%d", time.Now().Month())
}

func hour() string {
	return fmt.Sprintf("%d", time.Now().Hour())
}

func minute() string {
	return fmt.Sprintf("%d", time.Now().Minute())
}

func second() string {
	return fmt.Sprintf("%d", time.Now().Second())
}

func now() int64 {
	return time.Now().Unix()
}

func cast2str(val interface{}) string {
	return cast.ToString(val)
}

func cast2int64(val interface{}) int64 {
	return cast.ToInt64(val)
}

func cast2float32(val interface{}) float32 {
	return cast.ToFloat32(val)
}

var Functions = map[string]reflect.Value{
	"addi":           reflect.ValueOf(add[int64]),
	"addf":           reflect.ValueOf(add[float32]),
	"subi":           reflect.ValueOf(sub[int64]),
	"subf":           reflect.ValueOf(sub[float32]),
	"muli":           reflect.ValueOf(mul[int64]),
	"mulf":           reflect.ValueOf(mul[float32]),
	"divi":           reflect.ValueOf(div[int64]),
	"divf":           reflect.ValueOf(div[float32]),
	"mod":            reflect.ValueOf(mod),
	"pow":            reflect.ValueOf(pow),
	"round":          reflect.ValueOf(round),
	"ceil":           reflect.ValueOf(ceil),
	"floor":          reflect.ValueOf(floor),
	"log":            reflect.ValueOf(log),
	"log2":           reflect.ValueOf(log2),
	"log10":          reflect.ValueOf(log10),
	"exp":            reflect.ValueOf(exp),
	"sqrt":           reflect.ValueOf(sqrt),
	"absi":           reflect.ValueOf(abs[int64]),
	"absf":           reflect.ValueOf(abs[float32]),
	"asinh":          reflect.ValueOf(asinh),
	"sinh":           reflect.ValueOf(sinh),
	"asin":           reflect.ValueOf(asin),
	"sin":            reflect.ValueOf(sin),
	"acosh":          reflect.ValueOf(acosh),
	"acos":           reflect.ValueOf(acos),
	"cosh":           reflect.ValueOf(cosh),
	"cos":            reflect.ValueOf(cos),
	"atanh":          reflect.ValueOf(atanh),
	"tanh":           reflect.ValueOf(tanh),
	"atan":           reflect.ValueOf(atan),
	"tan":            reflect.ValueOf(tan),
	"sigmoid":        reflect.ValueOf(sigmoid),
	"year":           reflect.ValueOf(year),
	"month":          reflect.ValueOf(month),
	"day":            reflect.ValueOf(day),
	"hour":           reflect.ValueOf(hour),
	"minute":         reflect.ValueOf(minute),
	"second":         reflect.ValueOf(second),
	"now":            reflect.ValueOf(now),
	"date":           reflect.ValueOf(date),
	"from_unixtime":  reflect.ValueOf(from_unixtime),
	"unix_timestamp": reflect.ValueOf(unix_timestamp),
	"date_add":       reflect.ValueOf(date_add),
	"date_sub":       reflect.ValueOf(date_sub),
	"date_diff":      reflect.ValueOf(date_diff),
	"reverse":        reflect.ValueOf(reverse),
	"upper":          reflect.ValueOf(upper),
	"lower":          reflect.ValueOf(lower),
	"substr":         reflect.ValueOf(substr),
	"trim":           reflect.ValueOf(trim),
	"concat":         reflect.ValueOf(concat),
	"cast2str":       reflect.ValueOf(cast2str),
	"cast2int64":     reflect.ValueOf(cast2int64),
	"cast2float32":   reflect.ValueOf(cast2float32),
}
