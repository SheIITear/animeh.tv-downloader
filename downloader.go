package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Hentai struct {
	Success bool   `json:"success"`
	Titulo  string `json:"titulo"`
	Tipo    string `json:"tipo"`
	HTML    string `json:"html"`
}

func UnmarshalHentai(data []byte) (Hentai, error) {
	var r Hentai
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Hentai) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func DownloadHentai(filepath string, url string) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {

	query := os.Args[1]
	splitted := strings.Split(query, "/")
	name := splitted[4]
	fmt.Println(name)

	body := strings.NewReader(`seccion=reproductor&data=parsetoken%3D76%26token%3D593256685a6d59795a4751304d324d324e7a63314d32466a4d444d354d6a4e695a574d304d6d49784e7a553d%26id%3D176%26length%3D2%26capitulo%3D2%26type%3Dreproductor&nombre=` + name + `&servidor=1`)
	req, err := http.NewRequest("POST", "https://animeh.tv/reproductor.php", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "animeh.tv")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://animeh.tv")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", query)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "__cfduid=df91acc5901c95ba882de8f02f07ab3611590985872; PHPSESSID=ceaff2dd43c67753ac03923bec42b175; _ga=GA1.2.6980347.1590985884; _gid=GA1.2.839479332.1590985884; player_42=558; player_80=5; player_78=0; _gat_gtag_UA_111280423_2=1")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	kk, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	data, err := UnmarshalHentai(kk)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Success:", data.Success)

	// fmt.Println(data.HTML)

	z := html.NewTokenizer(strings.NewReader(data.HTML))
	var videourl string

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			fmt.Println("Downloading...")
			break
		}

		switch {
		case tt == html.SelfClosingTagToken:
			t := z.Token()

			isAnchor := t.Data == "source"

			if isAnchor {
				fmt.Println("We found a video!")

				for _, a := range t.Attr {

					if a.Key == "src" {
						fmt.Println("Found src:", a.Val)
						videourl = a.Val
						break
					}
				}
			}
		}
	}

	err = DownloadHentai(name+".mp4", videourl)

	if err != nil {
		panic(err)
	}

	fmt.Println("Downloaded: " + videourl)
}
