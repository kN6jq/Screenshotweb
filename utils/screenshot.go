package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	MaxWidth  = 1920
	MinHeight = 1080
)

type PocWebAppData struct {
	Image string `json:"image"` //网站截图（大图）
}

func Screenshot(url string) (bool, string) {
	site := &PocWebAppData{}
	siteImageName := fmt.Sprintf(`%s.png`, NewMd5(url))
	if Exists("./tmp/" + siteImageName) {
		os.Remove("./tmp/" + siteImageName)
	}
	status := DoFullScreenshot(url, fmt.Sprintf("./tmp/%s", siteImageName))
	if status {
		site.Image = siteImageName
		return true, fmt.Sprintf("/tmp/%s.png", NewMd5(url))
	} else {
		return false, ""
	}
}
func NewMd5(str ...string) string {
	h := md5.New()
	for _, v := range str {
		h.Write([]byte(v))
	}
	return hex.EncodeToString(h.Sum(nil))
}

/*执行截图*/
func DoFullScreenshot(url, path string) bool {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.WindowSize(MaxWidth, MinHeight),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 创建chrome实例
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// 创建超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 缓存对象
	var buf []byte

	// 运行截屏
	if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
		return false
	}

	// 保存文件
	if err := ioutil.WriteFile(path, buf, 0644); err != nil {
		return false
	}

	return true
}

/*全屏截图*/
func fullScreenshot(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*res, err = page.CaptureScreenshot().WithQuality(quality).WithClip(&page.Viewport{
				X:      0,
				Y:      0,
				Width:  MaxWidth,
				Height: MinHeight,
				Scale:  1,
			}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}

}
