package settings

var (
	CurrentRunConf   RunConfig
	CurrentParseConf ParseConfig
)

//RunConfig is the configuration for the run
type RunConfig struct {
	URL          string
	RuleFilePath string
	Debug        bool
}

//ParseConfig is the configuration for the parse
type ParseConfig struct {
	LibDirPath   string
	RuleFilePath string
	Debug        bool
}
