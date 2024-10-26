package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	for range 5 {
		go func() {
			for {
				time.Sleep(50 * time.Millisecond)
				start := time.Now()
				resp, err := http.Get(`https://ru.wikipedia.org/wiki/%D0%92%D0%B5%D0%BB%D0%B8%D0%BA%D0%B8%D0%B5_%D1%80%D0%B0%D0%B2%D0%BD%D0%B8%D0%BD%D1%8B`)
				if err != nil {
					panic(err)
				}

				if _, err := io.Copy(io.Discard, resp.Body); err != nil {
					panic(err)
				}

				fmt.Printf("request, time: %s\n", time.Since(start))
				resp.Body.Close()
			}
		}()
	}

	<-make(chan struct{})
}
