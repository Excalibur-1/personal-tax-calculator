package main

import "fmt"

const (
	// 1    不超过36000元的	            3	0
	// 2	超过 36000 元至144000元的部分	10	2520
	// 3	超过 144000 元至300000元的部分	20	16920
	// 4	超过 300000 元至420000元的部分	25	31920
	// 5	超过 420000 元至660000元的部分	30	52920
	// 6	超过 660000 元至960000元的部分	35	85920
	// 7	超过 960000 元的部分	        45	181920
	tax1 = 36000 * 0.03
	tax2 = (144000 - 36000) * 0.1
	tax3 = (300000 - 144000) * 0.2
	tax4 = (420000 - 300000) * 0.25
	tax5 = (660000 - 420000) * 0.3
	tax6 = (960000 - 660000) * 0.35

	deduct                     float64 = 60000 // 固定扣除，每月5000的起征点
	SpecialAdditionalDeduction float64 = 18000 // 专项附加扣除（深圳房租固定最高18000）

	// 三险一金，个人缴纳的比例
	n1 float64 = 0.08
	n2 float64 = 0.02
	n3 float64 = 0.06
	// 失业险为固定值，根据实际情况填入
	n4 float64 = 7.08

	// 三险一金，公司缴纳的比例
	m1 float64 = 0.15
	m2 float64 = 0.062
	m3 float64 = 0.06
)

var (
	a      float64 = 18000 // 养老保险缴纳的基数
	b      float64 = 18000 // 医疗保险缴纳的基数
	c      float64 = 18000 // 公积金缴纳的基数
	salary float64 = 20000 // 每个月的实际工资
)

func main() {
	// 计算需要缴纳的个税和公司缴纳的公积金
	f1, f2 := test1(a, b, c, salary)
	fmt.Printf("需要缴纳的个税：%v￥，公司缴纳的公积金：%v￥\n", f1, f2)

	// 计算最佳人民币工资和发u工资的分配，默认计算范围(5000-你的每月工资额)
	test2(salary)
}

// 因为选择发u则相应的社保和公积金基数也会随着变化 所以，计算只考虑公积金部分，
// 如果把公司缴纳的社保部分也考虑进去的话，则人民币部分肯定是越大越好了
// 公积金缴纳比例可以按照实际情况修改上面定义的变量
func test2(salary float64) {
	// 初始化
	n1, n2 := test1(5000, 5000, 5000, 5000)
	// 最佳的人民币工资额
	var index, max float64
	for i := float64(5001); i < salary; i++ {
		// 三险一金缴纳基数默认取工资数，如果有不同的情况，可以在此加入三险一金的基数计算逻辑
		// 比如三险一金的缴纳基数是工资的一半，则: a, b, c := i/2, i/2, i/2
		a, b, c := i, i, i
		f1, f2 := test1(a, b, c, i)

		// 这个是少交的税
		m1 := f1 - n1

		// 这个是少拿的公积金
		m2 := f2 - n2

		m3 := m1 - m2
		// 如果少交的税大于少拿的公积金
		if m3 > max {
			// 记录最优值的下标
			index = i

			// 记录最优的值
			max = m3
		}

		// 修改n1,n2的值
		n1, n2 = f1, f2
	}
	fmt.Printf("最佳发人民币工资额：%v￥, 最佳发u工资额:%v￥\n", index, salary-index)
}

func test1(a, b, c, salary float64) (float64, float64) {
	// 三险一金
	num1 := a * n1
	num2 := b * n2
	num3 := c * n3
	num := (num1 + num2 + num3 + n4) * 12
	// 应纳税所得额
	res1 := (salary * 12) - num - deduct - SpecialAdditionalDeduction
	// 公司缴纳的公积金，这里只考虑公积金部分，如果要考虑社保部分，则肯定是基数越高越好
	res2 := (c * m3) * 12
	return test(res1), res2
}

func test(res float64) float64 {
	if res <= 0 {
		return 0
	}
	var num float64
	if res > 36000 {
		num = tax1 + (res-36000)*0.1
	} else if res > 144000 {
		num = tax1 + tax2 + ((res - 144000) * 0.2)
	} else if res > 300000 {
		num = tax1 + tax2 + tax3 + ((res - 300000) * 0.25)
	} else if res > 420000 {
		num = tax1 + tax2 + tax3 + tax4 + ((res - 420000) * 0.3)
	} else if res > 660000 {
		num = tax1 + tax2 + tax3 + tax4 + tax5 + ((res - 660000) * 0.35)
	} else if res > 960000 {
		num = tax1 + tax2 + tax3 + tax4 + tax5 + tax6 + ((res - 960000) * 0.45)
	} else {
		num = res * 0.03
	}
	return num
}
