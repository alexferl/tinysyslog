package constants

const (
	MutatorText = "text"
	MutatorJSON = "json"
)

var Mutators = []string{MutatorText, MutatorJSON}

const (
	FilterRegex = "regex"
)

var Filters = []string{FilterRegex}

const (
	SinkConsole       = "console"
	SinkElasticsearch = "elasticsearch"
	SinkFilesystem    = "filesystem"
)

var Sinks = []string{SinkConsole, SinkElasticsearch, SinkFilesystem}

const (
	ConsoleStdOut = "stdout"
	ConsoleStdErr = "stderr"
)

var ConsoleOutputs = []string{ConsoleStdOut, ConsoleStdErr}
