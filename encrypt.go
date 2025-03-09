package main

import (
	"crypto/aes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func readImageFile(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image file: %w", err)
	}
	defer file.Close()

	var img image.Image
	switch filepath.Ext(imagePath) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format")
	}

	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	return img, nil
}

func saveImageFile(img image.Image, imagePath string) error {
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("error creating image file: %w", err)
	}
	defer file.Close()

	switch filepath.Ext(imagePath) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, nil)
	case ".png":
		err = png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported image format")
	}

	if err != nil {
		return fmt.Errorf("error encoding image: %w", err)
	}

	return nil
}

func getRgbData(img image.Image) []byte {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	rgbData := make([]byte, 0, bounds.Dx()*bounds.Dy()*3)
	for i := 0; i < len(rgba.Pix); i += 4 {
		rgbData = append(rgbData, rgba.Pix[i], rgba.Pix[i+1], rgba.Pix[i+2])
	}
	return rgbData
}

func createImageFromRgbData(rgbData []byte, bounds image.Rectangle) image.Image {
	width := bounds.Dx()
	height := bounds.Dy()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			offset := (y*width + x) * 3
			r := rgbData[offset]
			g := rgbData[offset+1]
			b := rgbData[offset+2]
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func genRandomSeed() int64 {
	return rand.Int63()
}

func genRandomKey(seed int64) []byte {
	r := rand.New(rand.NewSource(seed))
	key := make([]byte, 16)
	for i := range key {
		key[i] = byte(r.Intn(256))
	}
	return key
}

func ecbEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddingLength := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	plaintext = append(plaintext, make([]byte, paddingLength)...)

	ciphertext := make([]byte, len(plaintext))

	for i := 0; i < len(plaintext); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], plaintext[i:i+aes.BlockSize])
	}

	ciphertext = ciphertext[:len(ciphertext)-paddingLength]

	return ciphertext, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run encrypt.go <image-file-path> [seed]")
		return
	}

	imagePath := os.Args[1]
	var seed int64
	var err error
	if len(os.Args) > 2 {
		seed, err = strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("Error parsing seed:", err)
			return
		}
	} else {
		seed = genRandomSeed()
		fmt.Println("Seed:", seed)
	}

	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer file.Close()

	img, err := readImageFile(imagePath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}

	key := genRandomKey(seed)

	pixels := getRgbData(img)
	encryptedPixels, err := ecbEncrypt(pixels, key)
	if err != nil {
		fmt.Println("Error encrypting pixels:", err)
		return
	}

	newImg := createImageFromRgbData(encryptedPixels, img.Bounds())

	newImagePath := "ecb_" + filepath.Base(imagePath)
	err = saveImageFile(newImg, newImagePath)
	if err != nil {
		fmt.Println("Error saving new image file:", err)
		return
	}
}
