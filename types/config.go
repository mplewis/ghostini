package types

type Config struct {
	GhostSite       string `figyr:"required,description=The base URL of your Ghost website"`
	ContentKey      string `figyr:"required,description=The Content API key for your Ghost website"`
	Domains         string `figyr:"required,description=The domains for which to serve Gemini content"`
	GeminiCertsPath string `figyr:"optional,description=The path to the certificates and keys for your Gemini domains"`
	Host            string `figyr:"optional,description=The host on which to listen"`
	Port            int    `figyr:"default=1965,description=The port on which to listen"`
}
