package utils

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	dpi      = float64(72)
	fontfile = "./font/simhei.ttf"
	spacing  = float64(1.5)
)

func Whitemark(imgpath, savefile string, netime time.Time) {

	addWhitemark(imgpath, savefile, netime)
}

// 加水印函数需要两个参数，文件路径以及文件名称。
func addWhitemark(imgpath string, savefile string, netime time.Time) {

	if Exists("./" + "out.png") {
		os.Remove("./" + "out.png")
	}

	imgorgin, _ := os.Open(imgpath)
	img, _ := png.Decode(imgorgin)
	defer imgorgin.Close()

	flag.Parse()
	//字体为黑体，字体需要提前下载
	fontBytes, err := ioutil.ReadFile(fontfile)
	checkError(err)
	f, err := freetype.ParseFont(fontBytes)
	checkError(err)

	fg, bg := image.Black, image.Transparent
	//获取原始图片尺寸，新建水印图片
	rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	maxxlat := rgba.Bounds().Dx()
	maxylat := rgba.Bounds().Dy()
	//设置水印字体大小，此处设置为图片高度的三十二分之一
	size := float64(maxylat / 32)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	//设置水印位置，位置为距离图片右下角，按比例缩放
	pt := freetype.Pt(maxxlat-int(maxxlat/4), maxylat-int(maxylat/4)+int(c.PointToFixed(size)>>6))
	_, err = c.DrawString(netime.Format("2006-01-02 15:04:05"), pt)
	checkError(err)
	pt.Y += c.PointToFixed(size * spacing)

	// 保存水印图片
	outFile, err := os.Create("./" + "out.png")
	checkError(err)
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	checkError(err)
	err = b.Flush()
	checkError(err)
	//读取水印图片
	wmb, _ := os.Open("./" + "out.png")
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()
	//把水印图片盖在原始图片上
	offset := image.Pt(0, 0)
	bou := img.Bounds()
	m := image.NewNRGBA(bou)

	draw.Draw(m, bou, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	folderPath := filepath.Join("./result/", netime.Format("200601021504"))
	exists := Exists(folderPath)
	if !exists {
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			// 必须分成两步：先创建文件夹、再修改权限
			os.Mkdir(folderPath, 0777) //0777也可以os.ModePerm
			os.Chmod(folderPath, 0777)
		}
	}
	//保存新的图片
	imgw, _ := os.Create(folderPath + "/" + savefile + ".jpg")
	jpeg.Encode(imgw, m, &jpeg.Options{100})
	defer imgw.Close()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
