package gpt

const API_URL = "https://api.openai.com/v1/chat/completions"

type GPTModel string

const (
	ModelGPT4      GPTModel = "gpt-4"
	ModelGPT3Turbo GPTModel = "gpt-3.5-turbo"
)
