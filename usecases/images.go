package usecases

import (
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"
)

type ImageUsecase interface {
	InsertWatermark(srcImage io.Reader, fileName string) error
}

type imageUsecase struct {
	maxHeight      int
	maxWidth       int
	watermarkImage *image.NRGBA
}

func NewImageUsecase() ImageUsecase {
	wmImg, err := loadWatermark("./public/watermark.png")
	if err != nil {
		panic(err)
	}

	return imageUsecase{
		maxHeight:      720,
		maxWidth:       1280,
		watermarkImage: wmImg,
	}
}

func (u imageUsecase) InsertWatermark(srcImage io.Reader, fileName string) error {
	srcImg, _, err := image.Decode(srcImage)
	if err != nil {
		return err
	}

	// 画像をリサイズする
	sw := srcImg.Bounds().Dx()
	sh := srcImg.Bounds().Dy()

	resizedImage := &image.RGBA{}

	if sh >= sw {
		f := float64((sw * u.maxHeight))
		w := math.Round(f / float64(sh))
		resizedImage = image.NewRGBA(image.Rect(0, 0, int(w), u.maxHeight))
	} else {
		f := float64((sh * u.maxWidth))
		h := math.Round(f / float64(sw))
		resizedImage = image.NewRGBA(image.Rect(0, 0, u.maxWidth, int(h)))
	}

	draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	wmImg := u.watermarkImage
	x := resizedImage.Bounds().Dx() - wmImg.Bounds().Dx()
	y := resizedImage.Bounds().Dy() - wmImg.Bounds().Dy()

	dstImg := image.NewNRGBA(resizedImage.Bounds())
	draw.Draw(dstImg, resizedImage.Bounds(), resizedImage, image.Point{}, draw.Src)
	draw.Draw(dstImg, wmImg.Bounds().Add(image.Pt(x, y)), wmImg, image.Point{}, draw.Over)

	dstFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)

		return err
	}
	defer dstFile.Close()

	if err := jpeg.Encode(dstFile, dstImg, &jpeg.Options{Quality: 95}); err != nil {
		log.Println(err)

		return err
	}

	return nil

}

func loadWatermark(path string) (*image.NRGBA, error) {
	wmFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer wmFile.Close()

	src, _, err := image.Decode(wmFile)
	if err != nil {
		return nil, err
	}

	width := src.Bounds().Dx() / 2
	height := src.Bounds().Dy() / 2
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst, nil
}
