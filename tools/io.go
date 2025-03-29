package tools

import (
	"io"
	"log"
	"sync"
)

func IOCopy(w io.Writer, r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := io.Copy(w, r)
	if err != nil {
		log.Println("io err:", err)
	}
}
