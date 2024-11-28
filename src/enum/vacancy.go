package enum

type VacancyRequirementType string

const (
	Desirable  VacancyRequirementType = "desirable"
	Obligatory VacancyRequirementType = "obligatory"
)

func (v VacancyRequirementType) IsValid() bool {
	switch v {
	case Desirable, Obligatory:
		return true
	}
	return false
}

type VacancyContractType string

const (
	CLT     VacancyContractType = "clt"
	PJ      VacancyContractType = "pj"
	Trainee VacancyContractType = "trainee"
)

func (v VacancyContractType) IsValid() bool {
	switch v {
	case CLT, PJ, Trainee:
		return true
	}
	return false
}
