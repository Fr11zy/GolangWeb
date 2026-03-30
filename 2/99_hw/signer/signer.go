package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	channels := make([]chan interface{}, len(jobs)+1)
	for i := range channels {
		channels[i] = make(chan interface{})
	}

	wg := &sync.WaitGroup{}

	for i, jobfunc := range jobs {
		wg.Add(1)
		go func(jf job, in, out chan interface{}) {
			defer wg.Done()
			jf(in, out)
			close(out)
		}(jobfunc, channels[i], channels[i+1])
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		md5dataHash := DataSignerMd5(strconv.Itoa(data.(int)))
		wg.Add(1)
		go func (data interface{},md5dataHash string)  {
			defer wg.Done()
			crc32data := make([]string,2)
			innerwg := &sync.WaitGroup{}
			innerwg.Add(2)
			go func ()  {
				defer innerwg.Done()
				crc32data[0] = DataSignerCrc32(strconv.Itoa(data.(int)))
			}()
			go func ()  {
				defer innerwg.Done()
				crc32data[1] = DataSignerCrc32(md5dataHash)
			}()
			innerwg.Wait()
			out <- strings.Join(crc32data, "~")
		}(data,md5dataHash)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func (data interface{}){
			defer wg.Done()
			results := make([]string,6)
			innerwg := &sync.WaitGroup{}
			for th := 0; th < 6; th++ {
				innerwg.Add(1)
				go func(th int) {
					defer innerwg.Done()
					s := strconv.Itoa(th)+data.(string)
					results[th] = DataSignerCrc32(s)
				}(th)	
			}
			innerwg.Wait()
			res := strings.Join(results, "")
			out <- res
		}(data)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var a []string
	for data := range in {
		a = append(a, data.(string))
	}
	sort.Strings(a)
	out <- strings.Join(a, "_")
}
