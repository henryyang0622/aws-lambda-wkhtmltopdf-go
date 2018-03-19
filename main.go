package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

type Data struct {
	HTMLBase64 string `json:"htmlbase64"`
	URL        string `json:"url"`
	Options    []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"options"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}
	var data Data
	fmt.Println(request.Body)
	body := strings.Replace(request.Body, "'", "\"", -1)
	fmt.Println("==>" + body)
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	opt := getoption(data)
	if len(data.URL) > 0 {
		url_to_pdf(data.URL, opt)

	} else {
		htmltopdf(data.HTMLBase64, opt)
	}
	ff, err := ioutil.ReadFile("/tmp/to.pdf")
	if err != nil {
		fmt.Println(err)
	}

	var output string = base64.StdEncoding.EncodeToString(ff)

	return events.APIGatewayProxyResponse{
		Body:       "{ \"pdfbase64\": \"" + output + "\"}",
		StatusCode: 200,
	}, nil
}
func getoption(data Data) (opt []string) {

	option := []string{"--encoding", "utf-8"}
	fmt.Println(len(data.Options))
	for i := 0; i < len(data.Options); i++ {
		fmt.Println("Options[", i, "] =", data.Options[i].Key)
		if (len(data.Options[i].Key)) > 0 {
			option = append(option, "--"+data.Options[i].Key)
		}
		if (len(data.Options[i].Value)) > 0 {
			option = append(option, data.Options[i].Value)
		}
	}
	return option
}
func url_to_pdf(url string, opt []string) {
	fmt.Println("Get HTML from url" + url)

	opt = append(opt, url)
	opt = append(opt, "/tmp/to.pdf")
	fmt.Println(opt)
	cmd := exec.Command("./wkhtmltopdf", opt...)
	//cmd := exec.Command("pwd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	cmd = exec.Command("ls", "-l", "/tmp")
	//cmd := exec.Command("ls")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func htmltopdf(base64data string, opt []string) {
	fmt.Printf("Get HTML base64")
	fmt.Println("-----------------------------" + base64data)

	decodeBytes, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(decodeBytes))

	tofile, _ := os.Create("/tmp/form.html") // 建立字檔
	tofile.WriteString(string(decodeBytes))
	opt = append(opt, "/tmp/form.html")
	opt = append(opt, "/tmp/to.pdf")
	fmt.Println(opt)
	cmd := exec.Command("./wkhtmltopdf", opt...)
	//cmd := exec.Command("pwd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	cmd = exec.Command("ls", "-l", "/tmp")
	//cmd := exec.Command("ls")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func main() {
	lambda.Start(Handler)
	//build code 參考 要由bash 裡去zip   ex: zip main.zip main wkhtmltopdf
	//https://github.com/aws/aws-lambda-go#for-developers-on-windows
}

func checkErr(err error) {

	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
}
