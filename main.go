package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	data := getDataFromStdin()
	arr := dataToArray(string(data))
	flagType := getFlagType()
	var results = make(chan int, len(arr))
	k := 3
	createWorkerPool(arr, k, results, flagType)
	close(results)
	sum := 0
	for result := range results {
		sum += result
	}
	fmt.Println("Total:", sum)
}

func getFlagType() string {
	typeData := flag.String("type", "", "a type")
	flag.Parse()
	if *typeData != "url" && *typeData != "file" {
		fmt.Println("invalid flag: '--type url || file'")
		os.Exit(1)
	}
	return *typeData
}

func getDataFromStdin() []byte {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	return bytes
}

func dataToArray(data string) []string {
	return strings.Split(data, "\n")
}

func getUrlData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getFileData(path string) (string, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func worker(path string, results chan<- int, wg *sync.WaitGroup, flagType string) {
	var data string
	var err error
	switch flagType {
	case "url":
		data, err = getUrlData(path)
	case "file":
		data, err = getFileData(path)
	}
	if err != nil {
		fmt.Println("Error in", path+":", err)
		wg.Done()
		return
	}
	searchArr := searchWords("Go", data)
	wordLen := len(searchArr)
	fmt.Println("Count for", path+":", wordLen)
	results <- wordLen
	wg.Done()
}

func createWorkerPoolLimit(start int, finish int, list []string, results chan int, flagType string) {
	var wg sync.WaitGroup
	for i := start; i < finish; i++ {
		if len(list)-1 >= i && list[i] != "" {
			wg.Add(1)
			go worker(list[i], results, &wg, flagType)
		}
	}
	wg.Wait()
}

func createWorkerPool(list []string, k int, results chan int, flagType string) {
	listLen := len(list)
	if listLen > k {
		for i := 0; i < listLen; i += k {
			createWorkerPoolLimit(i, i+k, list, results, flagType)
		}
	} else {
		createWorkerPoolLimit(0, listLen, list, results, flagType)
	}
}

func searchWords(word string, data string) []string {
	reg := regexp.MustCompile(`\b` + word + `\b[’]?`) // в идеале `\bword\b(?![’])`, но поиск в go regexp не поддерживается
	arr := reg.FindAllString(data, -1)
	for i, item := range arr {
		if item == word+"’" {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}
