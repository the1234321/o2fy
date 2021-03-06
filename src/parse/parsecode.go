package parse

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"golang.design/x/clipboard"
)

func initCode() (*http.Client, string, error) {

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	HttpClient := &http.Client{Jar: jar, Timeout: time.Second * 10}
	resp1, err := HttpClient.Get("http://pso2s.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp1.Body.Close()

	var csrf *http.Cookie = nil

	for k, v := range resp1.Cookies() {
		if v.Name == "csrftoken" {
			csrf = resp1.Cookies()[k]
			break
		}
	}

	if csrf == nil {
		return nil, "", fmt.Errorf("csrf not found")
	}

	return HttpClient, csrf.Value, nil
}

func parseCode(HttpClient *http.Client, token string, img []byte) (string, error) {

	srcFile := bytes.NewReader(img)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", "OEMCaptchaImage.png")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://pso2s.com/upload/", body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-CSRFToken", token)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func getImgSize(img []byte) (*image.Config, error) {
	reader := bytes.NewReader(img)
	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, err
	}
	return &im, nil
}

func StartParse(result_string *binding.String) {

	HttpClient, token, err := initCode()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Println("????????????:")
	fmt.Println("??????????????????????????????, ??????????????????:21587709")

	// openbrowser("https://mcha.isao.net/profile_oem/OEMLogin.php?product_name=pso2&p_siteno=P00011")

	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	changedImg := clipboard.Watch(context.Background(), clipboard.FmtImage)

	for {
		img := <-changedImg
		{

			size, err := getImgSize(img)
			if err != nil {
				// fmt.Printf("????????????????????????:%v \n", err)
				continue
			}

			if size.Width != 200 || size.Height != 40 {
				fmt.Printf("???????????????????????????????????????\n")
				fmt.Println(size.Width, " ", size.Height)
				continue
			}

			fmt.Printf("???????????????????????????")

			result, err := parseCode(HttpClient, token, img)

			if err != nil {
				fmt.Printf("????????????????????????:%v \n", err)
				continue
			}
			fmt.Printf(", ???????????????:%s", result)
			if result != "????????????" {
				clipboard.Write(clipboard.FmtText, []byte(result))
				fmt.Printf(", ?????????????????????\n")

				if result_string != nil {
					(*result_string).Set(result)
				}
			}

		}
	}
}
