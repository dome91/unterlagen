package configuration

import (
	"github.com/spf13/viper"
)

const (
	None      AssistantProvider = "none"
	OpenAI    AssistantProvider = "openai"
	Ollama    AssistantProvider = "ollama"
	FixedSize ChunkerType       = "fixed"
	Recursive ChunkerType       = "recursive"
)

type AssistantProvider string
type ChunkerType string

type Configuration struct {
	Production bool
	Server     ServerConfiguration
	Assistant  AssistantConfiguration
	Data       DataConfiguration
}

type AssistantConfiguration struct {
	Provider AssistantProvider
	ApiKey   string
	Chunker  AssistantChunkerConfiguration
}

type AssistantChunkerConfiguration struct {
	Type         ChunkerType
	MaxChunkSize int
	ChunkOverlap int
}

type DataConfiguration struct {
	Directory string
}

type ServerConfiguration struct {
	Port       string
	BaseURL    string
	SessionKey string
}

type DatabaseConfiguration struct {
}

func Load() Configuration {
	viper.SetEnvPrefix("unterlagen")
	viper.AutomaticEnv()

	// Define default values
	setDefaults()

	config := Configuration{
		Production: viper.GetBool("production"),
		Server: ServerConfiguration{
			Port:       viper.GetString("server_port"),
			BaseURL:    viper.GetString("server_baseurl"),
			SessionKey: viper.GetString("server_session_key"),
		},
		Assistant: AssistantConfiguration{
			Provider: AssistantProvider(viper.GetString("assistant_provider")),
			ApiKey:   viper.GetString("assistant_api_key"),
			Chunker: AssistantChunkerConfiguration{
				Type:         ChunkerType(viper.GetString("assistant_chunker_type")),
				MaxChunkSize: viper.GetInt("assistant_chunker_max_chunk_size"),
				ChunkOverlap: viper.GetInt("assistant_chunker_chunk_overlap"), // Fixed key name
			},
		},
		Data: DataConfiguration{
			Directory: viper.GetString("data_directory"),
		},
	}

	if config.Server.SessionKey == "" {
		config.Server.SessionKey = "my-secret-key"
	}
	return config
}

func setDefaults() {
	viper.SetDefault("production", true)
	// Server defaults
	viper.SetDefault("server_port", "8080")
	viper.SetDefault("server_baseurl", "http://localhost:8080")
	viper.SetDefault("server_session_key", "") // Add this if you want a default

	// Assistant defaults
	viper.SetDefault("assistant_provider", string(None))
	viper.SetDefault("assistant_api_key", "")

	// Chunker defaults
	viper.SetDefault("assistant_chunker_type", string(FixedSize))
	viper.SetDefault("assistant_chunker_max_chunk_size", 100)
	viper.SetDefault("assistant_chunker_chunk_overlap", 20)

	// Data defaults
	viper.SetDefault("data_directory", "data")
}
