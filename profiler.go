package typhoon

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func fibo1(n int) int64 {
	if n == 0 || n == 1 {
		return int64(n)
	}
	time.Sleep(time.Millisecond)
	return int64(fibo1(n-1)) + int64(fibo1(n-2))
}

func fibo2(n int) int {
	fn := make(map[int]int)

	for i := 0; i <= n; i++ {
		var f int
		if i <= 2 {
			f = 1
		} else {
			f = fn[i-1] + fn[i-2]
		}
		fn[i] = f
	}
	time.Sleep(50 * time.Millisecond)
	return fn[n]
}

func N1(n int) bool {

	k := float64(n/2 + 1)
	for i := 2; i < int(k); i++ {
		if (n % i) == 0 {
			return false
		}
	}
	return true

}

func N2(n int) bool {
	for i := 2; i < n; i++ {
		if (n % i) == 0 {
			return false
		}
	}
	return true
}

func main2() {
	cpuFile, err := os.Create("/tmp/cpuProfile.out")
	if err != nil {
		fmt.Println(err)
		return
	}

	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	total := 0
	for i := 2; i < 100000; i++ {
		n := N1(i)
		if n {
			total++
		}
	}

	fmt.Println("Total primes:", total)

	total = 0
	for i := 2; i < 100000; i++ {
		n := N2(i)
		if n {
			total++
		}
	}

	fmt.Println("total primes:", total)
	for i := 1; i < 10; i++ {
		n := fibo1(i)
		fmt.Println(n, " ")

	}
	fmt.Println()
	for i := 1; i < 10; i++ {
		n := fibo2(i)
		fmt.Println(n, " ")
	}
	fmt.Println()
	runtime.GC()

}
