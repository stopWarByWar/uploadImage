package getSnapshot

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GetSnapshot(urlTemplate string, from, to int) ([][]byte, error) {
	// Start Chrome
	// Remove the 2nd param if you don't need debug information logged
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var tasts []chromedp.Action
	var imagesBuf [][]byte
	// Run Tasks
	// List of actions to run in sequence (which also fills our image buffer)
	for i := 0; i < to-from; i++ {
		url := fmt.Sprintf(urlTemplate, i+from)
		var imageBuf []byte
		imagesBuf = append(imagesBuf, imageBuf)
		tasts = append(tasts, ScreenshotTasks(url, &imagesBuf[i]))
		fmt.Println("get snapshots", url)
	}

	if err := chromedp.Run(ctx, tasts...); err != nil {
		return nil, err
	}
	// Write our image to file
	return imagesBuf, nil
}

func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(95).Do(ctx)
			return err
		}),
	}
}
