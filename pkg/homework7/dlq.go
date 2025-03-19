package homework7

type DeadLineQueue struct {
	queue []string
}

func NewDeadLineQueue() *DeadLineQueue {
	return &DeadLineQueue{make([]string, 0)}
}

func (dlq *DeadLineQueue) Enqueue(str string) {
	dlq.queue = append(dlq.queue, str)
}

func (dlq *DeadLineQueue) GetMessages() []string {
	return dlq.queue
}

func ProcessWithDLQ(messages []string, operation func(string) error, dlq *DeadLineQueue) {
	for _, message := range messages {
		err := operation(message)
		if err != nil {
			dlq.Enqueue(message)
		}
	}
}
