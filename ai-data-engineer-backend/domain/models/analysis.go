package models

type AnalysisResult struct {
	UserId         string       `json:"user_id"`
	AnalysisResult *LLMResponse `json:"analysis_result"`
	Status         string       `json:"status"`
}
