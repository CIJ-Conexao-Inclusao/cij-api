package enum

type VacancyRequirementType string

const (
	Desirable  VacancyRequirementType = "desirable"
	Obligatory VacancyRequirementType = "obligatory"
)

type VacancyContractType string

const (
	CLT     VacancyContractType = "clt"
	PJ      VacancyContractType = "pj"
	Trainee VacancyContractType = "trainee"
)
