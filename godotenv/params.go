package godotenv

type Params struct{}

type App interface {
	Load() string
}
