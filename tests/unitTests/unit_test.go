package unitTests

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
	"yandexLyceumTheme3gRPC/pkg/homework7"
)

func TestRetry(t *testing.T) {

}

func TestTimeoutElapsedTime(t *testing.T) {
	operationTooLong := func() error {
		time.Sleep(5 * time.Second) // эмулируем долгую операцию
		return nil
	}

	start := time.Now()
	_ = homework7.Timeout(operationTooLong, 1000) // таймаут 1 секунда
	timeElapsed := time.Since(start)
	if math.Abs(float64((timeElapsed - 1000*time.Millisecond).Milliseconds())) > 100 { // difference from 1s must be <0.1s
		t.Fatalf("Task with timeout=1000 must have a timeout of 1000ms, but it doesn't!\n"+
			"Elapsed time is %v", timeElapsed)
	}

	start = time.Now()
	_ = homework7.Timeout(operationTooLong, 2000) // таймаут 2 секунды
	timeElapsed = time.Since(start)
	if math.Abs(float64((timeElapsed - 2000*time.Millisecond).Milliseconds())) > 100 { // difference from 2s must be <0.1s
		t.Fatalf("Task with timeout=2000 must have a timeout of 2000ms, but it doesn't!\n"+
			"Elapsed time is %v", timeElapsed)
	}
}

func TestTimeoutResult(t *testing.T) {
	operationTooLong := func() error {
		time.Sleep(2 * time.Second) // эмулируем долгую операцию
		return nil
	}

	err := homework7.Timeout(operationTooLong, 1000) // таймаут 1 секунда
	if err == nil {
		t.Fatalf("Task must return an error due to timeout (func 2s > timeout 1s), but it's nil")
	}

	operationBad := func() error {
		return errors.New("error!")
	}
	err = homework7.Timeout(operationBad, 1000)
	if err == nil {
		t.Fatalf("Task must have been returned an error because it's what the operation does, but it did not")
	}

	operationCool := func() error {
		return nil
	}
	err = homework7.Timeout(operationCool, 1000)
	if err != nil {
		t.Fatalf("Task wasn't supposed to return any error because the operation returns nil immediately")
	}
}

func TestDeadLetterQueue(t *testing.T) {
	messages := []string{"msg1", "msg2", "msg3"}
	exceptedDeadMessages := []string{"msg2"}

	dlq := homework7.NewDeadLineQueue()

	homework7.ProcessWithDLQ(messages, func(msg string) error {
		if msg == "msg2" {
			return errors.New("processing failed")
		}
		fmt.Printf("Processed: %s\n", msg)
		return nil
	}, dlq)

	receivedDeadMessages := dlq.GetMessages()

	for index, deadMessage := range receivedDeadMessages {
		if deadMessage != exceptedDeadMessages[index] {
			t.Fatalf("Wrong dead letter list:\n"+
				"expected: %v\n"+
				"received: %v\n"+
				"wrong item on index %d\n"+
				"expected message: %s\n"+
				"received message: %s",
				exceptedDeadMessages, receivedDeadMessages, index, exceptedDeadMessages[index], deadMessage)
		}
	}
}
