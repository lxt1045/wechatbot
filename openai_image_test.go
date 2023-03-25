package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// log "github.com/sirupsen/logrus"

/*
接口使用:
https://platform.openai.com/docs/api-reference/images/create
模型介绍：
https://platform.openai.com/docs/models/overview
*/
func TestAA(t *testing.T) {
	t.Run("editImage", func(t *testing.T) {

		file := `./image2.png`
		file2 := `./image.png`

		bs, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		mask := []byte{}
		mask, err = ioutil.ReadFile(file2)
		if err != nil {
			t.Fatal(err)
		}
		urls, err := editImage2(bs, mask, 2, "1024x1024", "变成女生")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("urls:%+v", urls)
	})
}
func TestPng(t *testing.T) {
	t.Run("editImage", func(t *testing.T) {
		// file := `/Users/bytedance/Downloads/openai.2296579348.png`
		file := `/Users/bytedance/Documents/work_doc/11111.png`
		// file := `/Users/bytedance/Downloads/image_edit_original.png`
		err := savePng(file, `/Users/bytedance/Documents/work_doc/222.png`)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func Draw(src image.Image) (dstImg image.Image) {
	// jpeg 解码器返回的 image 对象（姑且称为对象）是只读的并不能在上面自由绘制，我们需要创建一个画布：
	width := 1024
	height := 1024
	dst := image.NewRGBA(image.Rect(0, 0, width, height))                           // 创建一块画布
	draw.Draw(dst, image.Rect(0, 0, width, height), src, image.Pt(0, 0), draw.Over) // 绘制第一幅图

	return dst
}

func ImgType(src image.Image) {
	if dst, ok := src.(*image.YCbCr); ok {
		log.Println(1)
		dst.Opaque()
	} else if _, ok := src.(*image.RGBA); ok {
		log.Println(2)
	} else if _, ok := src.(*image.NRGBA); ok {
		log.Println(3)
	} else {
	}
}

func savePng(in, out string) error {
	input, _ := os.Open(in)
	defer input.Close()
	img, str, err := image.Decode(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(img.Bounds(), ":", str)
	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()

	ImgType(img)
	AddOpacity(img)
	img = Draw(img)
	ImgType(img)

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	return nil
}

func AddOpacity(src image.Image) (img1 image.Image) {
	const width, height = 1024, 1024
	dst, ok := src.(*image.RGBA)
	if ok {
		log.Println(1)
		dst.Opaque()
	} else {
		return
	}
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rgb := dst.RGBA64At(x, y)
			var a uint8 = 0
			if x < 100 || x > 900 || y < 100 || y > 900 {
				a = 255
			}
			img.Set(x, y, color.NRGBA{
				R: uint8(rgb.R),
				G: uint8(rgb.G),
				B: uint8(rgb.B),
				A: a,
			})
		}
	}
	ImgType(img)

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return img
}

func TestWW(t *testing.T) {
	t.Run("editImage", func(t *testing.T) {
		file := `https://oaidalleapiprodscus.blob.core.windows.net/private/org-sUf6mdj53ewcqDn8rL39HauI/user-zGTpHjQlZYjhioNZZAVB63Fl/img-aEnaRzMyNn8uiHZftIRvLKCb.png?st=2023-03-24T13%3A06%3A08Z&se=2023-03-24T15%3A06%3A08Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-03-24T12%3A10%3A06Z&ske=2023-03-25T12%3A10%3A06Z&sks=b&skv=2021-08-06&sig=%2BFsZZtD9o8KxFlG%2BJJ031nj984EL5ULExdj9m/CI4Rs%3D`
		// bs, err := ioutil.ReadFile(file)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		urls, err := editImage(file, 1, "1024x1024", "红一点")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("urls:%+v", urls)
	})
}
