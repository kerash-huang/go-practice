package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	const check_path = "__PATH__"
	const check_string = "__CHECK_STR__"
	var input_filename string
	var output_filename string
	for isContinue := true; isContinue; {
		fmt.Print("Please enter filename: ")
		fmt.Scanln(&input_filename)
		if input_filename != "" {
			isContinue = false
		}
	}
	if input_filename == "q" {
		os.Exit(0)
	}
	var content string = readFileContent(input_filename)
	if content != "" {
		var arr = explode("\n", content)
		var wg sync.WaitGroup
		fmt.Println("Start checking")
		var exist_list []string
		for i := 0; i < len(arr); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				url := buildUrl(arr[i], check_path)
				response, success := curl_get(url)
				if success && response == check_string {
					exist_list = append(exist_list, arr[i])
				}
			}(i)
		}
		wg.Wait()
		if len(exist_list) > 0 {
			fmt.Println("You have " + fmt.Sprint(len(exist_list)) + " exist domain ..")
			fmt.Print("Enter file to save result(default: result.txt) : ")
			fmt.Scan(&output_filename)
			if output_filename == "" {
				output_filename = "result.txt"
			}
			var output_file, _err = os.OpenFile(output_filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
			if _err != nil {
				log.Fatal("Error open file: ", _err)
			}
			defer output_file.Close()
			for i := 0; i < len(exist_list); i++ {
				_, err := output_file.WriteString(exist_list[i] + "\n")
				if err != nil {
					log.Fatal("Error write file: ", err)
				}
			}
		} else {
			fmt.Println("No any exist domain ..")
		}
		fmt.Println("Complete checking task")
		var nothing string
		fmt.Scanln(&nothing)
	}
}

func buildUrl(url string, path string) string {
	url = checkAndFixHttps(url)
	url = strings.TrimRight(url, "/") + "/" + strings.TrimLeft(path, "/")
	return url
}

func curl_get(url string) (string, bool) {
	var response *http.Response
	var http_err error
	var httpInst = http.Client{
		Timeout: 2 * time.Second,
	}
	response, http_err = httpInst.Get(url)
	if http_err != nil {
		return "", false
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", false
	}
	var result, err = io.ReadAll(response.Body)
	if err != nil {
		return "", false
	}
	return string(result), true

}

func readFileContent(filename string) string {
	dat, err := os.ReadFile(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Fatal("File not exist: " + filename)
		return ""
	} else if err != nil {
		log.Fatal("Read error: ", err)
		return ""
	}
	return string(dat)
}

func explode(separator string, str string, limit ...int) []string {
	sep_limit := 0
	if len(limit) > 0 {
		sep_limit = limit[0]
	}
	var full_split_result = strings.Split(str, separator)
	if sep_limit > 0 {
		result := full_split_result[0 : sep_limit-1]
		return result
	} else {
		return full_split_result
	}
}

func checkAndFixHttps(url string) string {
	if strings.HasPrefix(url, "https://") {
		return url
	} else if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "https://", 1)
	} else {
		return "https://" + url
	}
}
