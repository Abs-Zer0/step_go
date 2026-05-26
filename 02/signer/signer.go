package main

import (
	//	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	input := make(chan interface{}, 5)
	output := make(chan interface{}, 5)

	for _, jobFunc := range jobs {
		wg.Add(1)
		go pipelineWorker(jobFunc, input, output, wg)

		input = output
		output = make(chan interface{}, 5)
	}

	wg.Wait()
}

func pipelineWorker(jobFunc job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	jobFunc(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		if number, ok := data.(int); ok {
			str := strconv.Itoa(number)

			first := crc32Calc(str)
			second := crc32Md5Calc(str)

			wg.Add(1)
			go singleHashOut(first, second, out, wg)
		}
	}

	wg.Wait()
}

func singleHashOut(first, second <-chan string, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	result := <-first + "~" + <-second
	out <- result
}

func crc32Calc(data string) chan string {
	output := make(chan string)
	go func(str string, out chan string) {
		out <- DataSignerCrc32(str)
	}(data, output)

	return output
}

var mutex = &sync.Mutex{}

func crc32Md5Calc(data string) chan string {
	output := make(chan string)
	go func(str string, out chan string) {
		mutex.Lock()
		md5Hash := DataSignerMd5(str)
		mutex.Unlock()

		out <- DataSignerCrc32(md5Hash)
	}(data, output)

	return output
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		if str, ok := data.(string); ok {
			size := 6

			chans := make([]chan string, 0, size)
			for i := 0; i < size; i++ {
				chans = append(chans, crc32Calc(strconv.Itoa(i)+str))
			}

			wg.Add(1)
			go multiHashOut(chans, out, wg)
		}
	}

	wg.Wait()
}

func multiHashOut(chans []chan string, out chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	size := len(chans)
	hashs := make([]string, 0, size)
	for i := 0; i < size; i++ {
		hashs = append(hashs, <-chans[i])
	}

	result := strings.Join(hashs, "")
	out <- result
}

func CombineResults(in, out chan interface{}) {
	allInputs := make([]string, 0)

	for data := range in {
		if str, ok := data.(string); ok {
			allInputs = append(allInputs, str)
		}
	}

	sort.Strings(allInputs)

	result := strings.Join(allInputs, "_")
	out <- result
}
