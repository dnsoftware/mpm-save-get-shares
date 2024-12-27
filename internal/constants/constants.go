package constants

const (
	ProjectRootAnchorFile = ".env"
	AppLogFile            = "app.log"
	TestLogFile           = "test.log"
)

// Работа с шарами
const (
	KafkaSharesGroup              = "sharesGroup"
	KafkaSharesTopic              = "shares"
	KafkaSharesAutocommitInterval = 5
)

const WorkerSeparator = "." // символ разделитель имени воркера от имени кошелька
