package model

// Data содержит данные для отправки
type Data struct {
	PeriodStart         string `validate:"required" json:"period_start"`
	PeriodEnd           string `validate:"required" json:"period_end"`
	PeriodKey           string `validate:"required" json:"period_key"`
	IndicatorToMoId     string `validate:"required,numeric" json:"indicator_to_mo_id"`
	IndicatorToMoFactId string `validate:"required,numeric" json:"indicator_to_mo_fact"`
	Value               string `validate:"required,numeric" json:"value"`
	FactTime            string `validate:"required" json:"fact_time"`
	IsPlan              string `validate:"required,numeric" json:"is_plan"`
	AuthUserID          string `validate:"required,numeric" json:"auth_user_id"`
	Comment             string `json:"comment"`
}

// Config содержит конфигурацию приложения
type Config struct {
	Host        string `validate:"required,hostname" yaml:"host"`
	Port        string `validate:"required,numeric" yaml:"port"`
	BearerToken string `yaml:"bearer_token"`
	Href        string `validate:"required,url" yaml:"href"`
	Buffer      int    `validate:"required,numeric" yaml:"buffer"`
}
