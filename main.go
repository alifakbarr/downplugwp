package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	// loop page 2-10
	for i := 2; i < 10; i++ {
		url := fmt.Sprintf("https://wordpress.org/plugins/browse/popular/page/%d/", i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// handle err
		}
		req.Header.Set("User-Agent", "NikkoBot")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}

		defer resp.Body.Close()

		var re = regexp.MustCompile(`(?m)<h3 class="entry-title"><a href="(.*)" rel="bookmark">.*</a></h3>`)

		resp_body, _ := ioutil.ReadAll(resp.Body)
		str_rp := string(resp_body)

		for _, match := range re.FindAllStringSubmatch(str_rp, -1) {
			// fmt.Println(match[1])
			urls := match[1]
			req, err := http.NewRequest("GET", urls, nil)
			if err != nil {
				// handle err
			}
			req.Header.Set("User-Agent", "NikkoBot")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
			}
			resp_page, _ := ioutil.ReadAll(resp.Body)
			str_rp_page := string(resp_page)

			var re_download = regexp.MustCompile(`(?m)<a class="plugin-download button download-button button-large" href="(.*)">Download</a>`)
			for _, match := range re_download.FindAllStringSubmatch(str_rp_page, -1) {
				url_downloads := match[1]
				name_file := strings.Replace(url_downloads, "https://downloads.wordpress.org/plugin/", "", -1)
				req, err := http.NewRequest("GET", url_downloads, nil)
				if err != nil {
					// handle err
				}
				req.Header.Set("User-Agent", "NikkoBot")

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					// handle err
				}
				defer resp.Body.Close()
				resp_downloads, _ := ioutil.ReadAll(resp.Body)
				ioutil.WriteFile("plugins/"+name_file, resp_downloads, 0644)
				fmt.Println("Downloaded: " + name_file)
				//exec command
				cmd := exec.Command("unzip", "plugins/"+name_file, "-d", "extractedplugins/")
				cmd.Run()

			}
		}
	}

}
