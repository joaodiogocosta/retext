package client

import (
	"fmt"
)

type ConnDryAdapter struct {}

func (adapter *ConnDryAdapter) Connect(session *Session, sendCh chan []byte) {
	go func() {
		for message := range sendCh {
			fmt.Println(string(message))
		}
	}()
}
