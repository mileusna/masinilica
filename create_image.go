package main

import (
	"bytes"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
)

const (
	width      = 1200
	leftMargin = 50
)

type Line struct {
	X    int
	Y    int
	Text string
}

func CreateImage(text, fontSize, align string, darkMode bool) ([]byte, error) {

	backgroundColor := color.RGBA{255, 255, 255, 255}
	textColor := color.Black
	if darkMode {
		textColor = color.White
		backgroundColor = color.RGBA{36, 40, 53, 255}
	}

	textFontSize := 140
	topMargin := 140 + textFontSize/2

	switch strings.ToLower(fontSize) {
	case "small":
		textFontSize = 90
		topMargin = 90 + textFontSize/2
	case "large":
		textFontSize = 180
		topMargin = 180 + textFontSize/2
	}

	face, err := opentype.NewFace(masinilicaFont, &opentype.FaceOptions{
		Size:    float64(textFontSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	d := &font.Drawer{
		Face: face,
	}

	lines := parseLines(text, textFontSize, topMargin, align, d)
	height := len(lines)*textFontSize + topMargin

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	d.Dst = img
	d.Src = image.NewUniform(textColor)
	for _, l := range lines {
		d.Dot = fixed.P(l.X, l.Y)
		d.DrawString(l.Text)
	}

	// return image
	pngBuff := bytes.NewBuffer([]byte{})
	if err := png.Encode(pngBuff, img); err != nil {
		return nil, err
	}

	return pngBuff.Bytes(), nil
}

// parseLines returns lines slice coordinates where to write
func parseLines(text string, textFontSize int, topMargin int, align string, d *font.Drawer) []Line {
	lines := make([]Line, 0, 1)
	rawLines := strings.Split(text, "\n")
	currentY := topMargin
	for _, l := range rawLines {

		linesToWrite, widths := wrapTextByWidth(l, width-(2*leftMargin), d)
		for i, lineText := range linesToWrite {
			line := Line{
				Text: lineText,
				Y:    currentY,
			}
			switch align {
			case "center":
				line.X = (width - widths[i]) / 2
			case "right":
				line.X = width - widths[i] - leftMargin
			default:
				line.X = leftMargin
			}
			currentY = currentY + textFontSize
			lines = append(lines, line)
		}
	}
	return lines
}

// Prelamanje teksta po širini u pikselima
func wrapTextByWidth(text string, maxWidth int, d *font.Drawer) ([]string, []int) {
	words := strings.Fields(text)
	var lines []string
	var widths []int
	var currentLine string

	for _, word := range words {
		tryLine := word
		if currentLine != "" {
			tryLine = currentLine + " " + word
		}

		width := d.MeasureString(tryLine).Round()
		if width > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				widths = append(widths, d.MeasureString(currentLine).Round())
			}
			currentLine = word
		} else {
			currentLine = tryLine
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
		widths = append(widths, d.MeasureString(currentLine).Round())
	}

	return lines, widths
}
