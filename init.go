package tzdb

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

func init() {
	log.SetPrefix("tzdb ")
	initAsyncJob()
	time.Now().Zone()
}

var zonefile = "/tmp/zoneinfo.zip"

func download() (err error) {
	u := "https://api.github.com/repos/tiechui1994/tzdb/releases/latest"
	request, _ := http.NewRequest("GET", u, nil)
	request.Header.Set("Accept", "application/vnd.github.v3+json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	var release struct {
		Assets []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
		} `json:"assets"`
	}

	err = json.NewDecoder(response.Body).Decode(&release)
	if err != nil {
		log.Println(err)
		return err
	}

	response, err = http.Get(release.Assets[0].BrowserDownloadUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	fd, err := os.Open(zonefile)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = io.Copy(fd, response.Body)
	if err != nil {
		os.Remove(zonefile)
		log.Println(err)
		return err
	}

	return nil
}

// update zoneinfo information
func initAsyncJob() {
	if _, err := os.Open(zonefile); err != nil && os.IsNotExist(err) {
		err = download()
		if err != nil {
			return
		}
	}

	syscall.Setenv("ZONEINFO", zonefile)

	now := time.Now()
	y, m, d := now.Date()
	after := time.Date(y, m, d, 0, 0, 0, 0, time.Local).Sub(now)
	if after < 0 {
		y, m, d = now.Add(24 * time.Hour).Date()
		after = time.Date(y, m, d, 0, 0, 0, 0, time.Local).Sub(now)
	}

	go func() {
		time.Sleep(after)

		ticker := time.NewTicker(time.Hour * 24)
		for {
			select {
			case <-ticker.C:
				go download()
			}
		}
	}()
}
