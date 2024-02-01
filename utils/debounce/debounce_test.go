package debounce

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {

	tag := "teste"
	value := "TestValue"
	channel := Channel(tag)
	for i := 0; i < 2; i++ {
		channel <- fmt.Sprintf("%s.%d", value, i)
		Run(tag, time.Second*2, channel, func(payload interface{}) {
			valueInputed := <-channel
			assert.Equal(t, fmt.Sprintf("%s.1", value), valueInputed)
		})
	}

}
