package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	f, err := os.Open("/usr/share/dict/american-english")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
		//fmt.Printf("line: %s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Num Words: %d\n", len(words))

	foundChannel := make(chan bool)

	shuffle(words)
	wordsCopy := shuffle(words)

	//load(wordsCopy, words, foundChannel)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go find(words, "*****", foundChannel, wg)
	for _, w := range wordsCopy {
		wg.Add(1)
		go find(words, w, foundChannel, wg)
	}
	wg.Add(1)
	go find(words, "&&&&&&", foundChannel, wg)

	go sentinal(wg, foundChannel)

	count := 0
	for found := range foundChannel {
		count++
		fmt.Printf("%v\n", found)
	}

	fmt.Printf("Count: %d\n", count)
}

func sentinal(group *sync.WaitGroup, ch chan bool) {
	group.Wait()
	close(ch)
}

//func load(toFind, cache []string, sentinal chan bool) {
//for _, w := range toFind {
//go find(cache, w, sentinal)
//}
//defer close(sentinal)
//}

func find(words []string, target string, sentinal chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	bottom := 0
	top := len(words)
	middle := top / 2

	attempts := 1

	for {
		attempts++

		if top == bottom {
			//fmt.Printf("%d iterations\n", attempts)
			sentinal <- false
			return
		}

		w := words[middle]
		if w == target {
			//fmt.Printf("%d iterations\n", attempts)
			sentinal <- true
			return
		} else if w < target {
			bottom = middle
		} else if w > target {
			top = middle
		}

		middle = ((top - bottom) / 2) + bottom
	}
}

func shuffle(words []string) []string {

	wordsCopy := make([]string, len(words), len(words))

	copy(wordsCopy, words)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(wordsCopy), func(i, j int) { wordsCopy[i], wordsCopy[j] = wordsCopy[j], wordsCopy[i] })

	return wordsCopy
}
