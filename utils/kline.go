package utils

import (
	"golang.org/x/exp/rand"
	"math"
	"time"
)

// KlineAttrs K线图数据
type KlineAttrs struct {
	OpenPrice  float64 `json:"openPrice"`  //开盘价格
	HighPrice  float64 `json:"highPrice"`  //最高价格
	LowsPrice  float64 `json:"lowsPrice"`  //最低价格
	ClosePrice float64 `json:"closePrice"` //收盘价格
	Vol        float64 `json:"vol"`        //交易量
	Amount     float64 `json:"amount"`     //成交额
	CreatedAt  int64   `json:"createdAt"`  //开盘时间
}

// GenerateKline 生成K线图
func GenerateKline(startPrice, endPrice float64, startTime, endTime time.Time) []*KlineAttrs {
	var interval int64 = 60
	rows := int((endTime.Unix() - startTime.Unix()) / interval)
	volatility := float64(getIntegerDigits(int64((startPrice + endPrice) / 2)))
	//bridgePrice := generateBrownianBridge(startPrice, endPrice, rows, volatility)
	bridgePrice := sinusoidalChange(startPrice, endPrice, rows)
	return generateKLineData(startTime, interval, bridgePrice, rows, volatility)
}

// GetIntegerDigits 获取整数数字位
func getIntegerDigits(interval int64) int {
	count := 1
	for i := 0; i < 20; i++ {
		if interval/10 <= 1 {
			break
		}
		interval = interval / 10
		count *= 10
	}
	return count
}

// generateBrownianBridge 生成布朗桥价格路径，确保价格在 [min(start, end), max(start, end)] 范围内
func generateBrownianBridge(start, end float64, steps int, volatility float64) []float64 {
	// 确定价格范围
	minPrice := math.Min(start, end)
	maxPrice := math.Max(start, end)

	prices := make([]float64, steps+1)
	prices[0] = math.Max(math.Min(start, maxPrice), minPrice)   // 确保起始价格在范围内
	prices[steps] = math.Max(math.Min(end, maxPrice), minPrice) // 确保结束价格在范围内

	for i := 1; i < steps; i++ {
		// 时间点比例
		t := float64(i) / float64(steps)

		// 期望价格线性插值
		expected := start + t*(end-start)

		// 生成随机偏差，符合波动率
		// 这里使用标准正态分布乘以波动率和预定义的步长
		deviation := rand.NormFloat64() * GenerateRandomFloat64(volatility/10, volatility, 2) * (end - start) / float64(steps)

		// 当前价格 = 期望价格 + 偏差
		prices[i] = expected + deviation

		// 确保价格在 [minPrice, maxPrice] 范围内
		prices[i] = math.Max(math.Min(prices[i], maxPrice), minPrice)
	}

	// 为确保价格连续性，可以对价格进行平滑处理
	for i := 1; i < len(prices); i++ {
		// 简单的线性插值以确保价格不会出现剧烈跳变
		prices[i] = (prices[i-1] + prices[i]) / 2

		// 再次确保价格在 [minPrice, maxPrice] 范围内
		prices[i] = math.Max(math.Min(prices[i], maxPrice), minPrice)
	}

	// 确保最后价格严格等于 end
	if end > 0 {
		prices[steps] = end
	}

	return prices
}

// generateKLineData 根据价格路径生成 K 线数据
func generateKLineData(startTime time.Time, interval int64, prices []float64, rows int, volatility float64) []*KlineAttrs {
	var kLines []*KlineAttrs

	for i := 0; i < rows; i++ {
		open := prices[i]
		close := prices[i+1]

		// 计算最高价和最低价
		high := math.Max(open, close) + GenerateRandomFloat64(volatility/100, volatility/10, 2)
		low := math.Min(open, close) - GenerateRandomFloat64(volatility/100, volatility/10, 2)

		// 确保最低价不低于0
		if low < 0 {
			low = math.Min(open, close)
		}

		// 模拟成交量（可选）
		volume := rand.Float64()*1000 + 100 // 随机 100 到 1100

		// 添加到 K 线数据
		kLines = append(kLines, &KlineAttrs{
			CreatedAt:  startTime.Unix(),
			OpenPrice:  round(open, 2),
			ClosePrice: round(close, 2),
			LowsPrice:  round(low, 2),
			HighPrice:  round(high, 2),
			Vol:        round(volume, 2),
		})
		startTime = startTime.Add(time.Duration(interval) * time.Second)
	}

	return kLines
}

// round 辅助函数：四舍五入
func round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// GenerateRandomFloat64 生成指定范围内的随机float64数字，保留指定位小数
func GenerateRandomFloat64(min, max float64, decimalPlaces int) float64 {
	// 确保min小于max
	if min > max {
		min, max = max, min
	}

	// 生成随机float64
	random := min + rand.Float64()*(max-min)

	// 计算四舍五入的倍数
	roundMultiplier := math.Pow10(decimalPlaces)

	// 保留指定位小数
	return math.Round(random*roundMultiplier) / roundMultiplier
}

// sinusoidalChange 价格浮动
func sinusoidalChange(min, max float64, num int) []float64 {
	amplitude := (max - min) / 2 // 振幅
	offset := min + amplitude    // 基线

	vals := make([]float64, 0)
	for i := 0; i < num+1; i++ {
		// 生成正弦值，周期为1分钟
		progress := float64(i) / float64(60) * 2 * math.Pi
		value := offset + amplitude*math.Cos(progress) // Sinh Sin
		vals = append(vals, round(value, 2))
	}
	return vals
}
