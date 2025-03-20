package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

// Load an image from file
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Save an image to file
func saveImage(img image.Image, filename string, quality int) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	return err
}

// Resize an image
func resizeImage(img image.Image, width, height uint) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}


// Crop an image
func cropImage(img image.Image, width, height, x, y int) image.Image {
	return imaging.Crop(img, image.Rect(x, y, width, height))
}

// Rotate an image
func rotateImage(img image.Image, angle float64) image.Image {
	return imaging.Rotate(img, angle, color.Black)
}

// Flip an image horizontally
func flipImage(img image.Image) image.Image {
	return imaging.FlipH(img)
}

// Mirror an image (same as flip horizontally)
func mirrorImage(img image.Image) image.Image {
	return imaging.FlipH(img)
}

// Apply grayscale filter
func applyGrayscale(img image.Image) image.Image {
	return imaging.Grayscale(img)
}

// Apply sepia filter
func applySepia(img image.Image) image.Image {
	sepia := imaging.AdjustSaturation(img, -100)
	sepia = imaging.AdjustContrast(sepia, 20)
	sepia = imaging.AdjustBrightness(sepia, 10)
	return sepia
}

// Apply blur filter
func applyBlur(img image.Image, sigma float64) image.Image {
	return imaging.Blur(img, sigma)
}

// Add watermark to an image
type Position string

const (
	TopLeft     Position = "top-left"
	TopRight    Position = "top-right"
	BottomLeft  Position = "bottom-left"
	BottomRight Position = "bottom-right"
	Center      Position = "center"
)

func addWatermark(img image.Image, watermarkPath string, opacity float64, scale float64, position Position) (image.Image, error) {
	// Load watermark
	watermark, err := imaging.Open(watermarkPath)
	if err != nil {
		return nil, err
	}

	// Scale watermark based on original image width
	newWidth := int(float64(img.Bounds().Dx()) * scale)
	newHeight := int(float64(watermark.Bounds().Dy()) * (float64(newWidth) / float64(watermark.Bounds().Dx())))
	watermark = imaging.Resize(watermark, newWidth, newHeight, imaging.Lanczos)

	// Calculate position
	var offset image.Point
	switch position {
	case TopLeft:
		offset = image.Pt(10, 10)
	case TopRight:
		offset = image.Pt(img.Bounds().Dx()-newWidth-10, 10)
	case BottomLeft:
		offset = image.Pt(10, img.Bounds().Dy()-newHeight-10)
	case BottomRight:
		offset = image.Pt(img.Bounds().Dx()-newWidth-10, img.Bounds().Dy()-newHeight-10)
	case Center:
		offset = image.Pt(
			(img.Bounds().Dx()-newWidth)/2,
			(img.Bounds().Dy()-newHeight)/2,
		)
	}

	// Overlay watermark
	result := imaging.Overlay(img, watermark, offset, opacity)
	return result, nil
}


// Compress image (handled during saving)

// Convert image format
func convertFormat(img image.Image, filename string) error {
	return imaging.Save(img, filename)
}

func loadImageFromURL(url string) (image.Image, error) {
	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	// Decode image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}


type Resize struct {
	Width uint
	Height uint
}

type Crop struct {
	X int
	Y int
	Width int
	Height int
}

type Filter struct {
	Grayscale bool
	Blur float64
	Sepia bool
}

type Transformation struct {
	Resize Resize
	Crop Crop
	Rotate int
	Format string
	Filter Filter
	Quality	int
}

type ProcessImageOptions struct {
	Transformation Transformation
}

// Process the image (combining all operations)
func processImage(inputUrl string, options ProcessImageOptions) error {
	img, err := loadImageFromURL(inputUrl)
	if err != nil {
		return errors.New("load:" + err.Error())
	}

	if (options.Transformation.Resize.Width > 0) {
		if (options.Transformation.Resize.Height <= 0 || options.Transformation.Resize.Width <= 0) {
			return errors.New("invalid resize dimensions")
		}

		img = resizeImage(img, options.Transformation.Resize.Width, options.Transformation.Resize.Height)
	}

	if (options.Transformation.Crop.X > 0) {
		if (options.Transformation.Crop.Y <= 0 || options.Transformation.Crop.Width <= 0 || options.Transformation.Crop.Height <= 0) {
			return errors.New("invalid crop dimensions")
		}
		img = cropImage(img, options.Transformation.Crop.Width, options.Transformation.Crop.Height, options.Transformation.Crop.X, options.Transformation.Crop.Y)
	}


	if (options.Transformation.Rotate > 0) {
		img = rotateImage(img, float64(options.Transformation.Rotate))
	}


	if(options.Transformation.Filter.Blur > 0) {
		img = applyBlur(img, options.Transformation.Filter.Blur)
	}

	if(options.Transformation.Filter.Sepia) {
		img = applySepia(img)
	}

	if(options.Transformation.Filter.Grayscale) {
		img = applyGrayscale(img)
	}


	if(options.Transformation.Quality == 0){
		options.Transformation.Quality = 100
	}

	outputPath := fmt.Sprintf("output.%s", options.Transformation.Format)

	// Save (compress at 75% quality)
	err = saveImage(img, outputPath, options.Transformation.Quality)

	if err != nil {
		return errors.New("save:" + err.Error())
	}

	fmt.Println("Image processing complete:", outputPath)
	return nil
}

func main() {
	url := "https://imgs.search.brave.com/Jx8eAt3b3FFc7T8qGmK4AjpohtGB8b4pA1TYkGAkfTQ/rs:fit:860:0:0:0/g:ce/aHR0cHM6Ly90My5m/dGNkbi5uZXQvanBn/LzA5LzY1LzEzLzk0/LzM2MF9GXzk2NTEz/OTQwOV9JT21HVGVQ/Z2ZFVW44ek1YWnFw/YTlPRXRpRmJndkVC/TC5qcGc"
	err := processImage(url, ProcessImageOptions{
		Transformation: Transformation{
			Resize: Resize{
				Width: 800,
				Height: 600,
			},
			Filter: Filter{
				Grayscale: true,
			},
			Quality: 75,
			Format: "png",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
