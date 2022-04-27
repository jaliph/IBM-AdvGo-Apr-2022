package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	// common channel to get all the file digits
	ch := make(chan int)
	// odd only, even only
	oddCh := make(chan int)
	evenCh := make(chan int)

	chs := make([]chan bool, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		fmt.Println(os.Args[i])
		chs[i-1] = readFile(os.Args[i], ch)
	}

	// f1 := readFile("data2.dat", ch)
	// f2 := readFile("data2.dat", ch)
	f3 := splitter(ch, oddCh, evenCh)
	f4 := sum(oddCh)
	f5 := sum(evenCh)
	for _, ch := range chs {
		<-ch
	}
	// <-f1
	// <-f2
	close(ch)
	<-f3
	evenSum := <-f4
	oddSum := <-f5
	fmt.Println("OddSum is :", oddSum)
	fmt.Println("EvenSum is :", evenSum)
	writeSum(oddSum, evenSum)
	fmt.Println("All Done")

}

func writeSum(oddSum int, evenSum int) {
	// write the whole body at once
	s := "Even Total : " + strconv.Itoa(evenSum) + "\n" + "Odd Total : " + strconv.Itoa(oddSum)
	err := ioutil.WriteFile("actual-result.txt", []byte(s), 0644)
	if err != nil {
		panic(err)
	}
}

func sum(ch chan int) chan int {
	sumCh := make(chan int)
	go func() {
		var sum int
		for no := range ch {
			sum += no
		}
		sumCh <- sum
		close(sumCh)
	}()
	return sumCh
}

func splitter(ch chan int, o chan int, e chan int) chan bool {
	finish := make(chan bool)
	go func() {
		for ch != nil {
			select {
			case a, ok := <-ch:
				if ok {
					if a%2 == 0 {
						e <- a
					} else {
						o <- a
					}
				} else {
					// fmt.Println("Channel Closed")
					close(o)
					close(e)
					finish <- true
					ch = nil
				}
			}
		}
	}()
	return finish
}

func readFile(logfile string, splitterCh chan int) chan bool {
	finish := make(chan bool)
	go func() {
		f, err := os.OpenFile(logfile, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rd := bufio.NewReader(f)
		for {
			line, err := rd.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					finish <- true
					close(finish)
					break
				}

				panic(err)
			}
			number, err := strconv.Atoi(line[:len(line)-1])
			if err != nil {
				panic(err)
			}
			splitterCh <- number
		}
	}()
	return finish
}
