package networking

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestStablelity(t *testing.T) {
	t.Skip()
	client, err := Dial("https://eth.merkle.io")
	if err != nil {
		panic("init failed")
	}
	client2, err := Dial("https://ethereum.publicnode.com")
	if err != nil {
		panic(err)
	}
	solidClient := NewSolidEthClient(client2, []*StatefulEthClient{client2, client}...)

	for true {
		time.Sleep(1 * time.Millisecond)
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				block, err := solidClient.BlockNumber(context.Background())
				if err != nil {
					wg.Done()
					return
				}
				if block == 0 {
					panic("unexpected")
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
