package calculate

import (
	"fmt"
	"strconv"
)

// float64计算结果保留2位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
