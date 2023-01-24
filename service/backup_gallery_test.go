package service

import (
	"chiller/http_client"
	"sync"
	"testing"
)

func TestChannels(t *testing.T) {

	messages := make(chan int, 16)

	var wg sync.WaitGroup
	wg.Add(16)
	for i := 0; i < 16; i++ {
		go func(i int) {
			defer wg.Done()
			messages <- i
		}(i)
	}
	wg.Wait()
	close(messages)
	for v := range messages {
		t.Log(v)
	}
}

// poc
func TestBackupGallery(t *testing.T) {

	auth, _ := http_client.GetBearerToken()

	BackupGalleryService(auth)
}
