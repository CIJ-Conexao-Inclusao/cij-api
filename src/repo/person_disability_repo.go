package repo

import (
	"cij_api/src/model"
	"cij_api/src/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PersonDisabilityRepo interface {
	BaseRepoMethods

	GetPersonDisabilities(personId int) ([]model.PersonDisability, utils.Error)
	GetDisabilityById(disabilityId int) (model.Disability, utils.Error)
	UpsertPersonDisability(personDisability model.PersonDisability, tx *gorm.DB) utils.Error
	ClearPersonDisability(personId int, tx *gorm.DB) utils.Error

	// reports
	CountDisability() (model.DisabilityTotals, utils.Error)
	CountDisabilityByNeighborhood(neighborhood string) (model.DisabilityTotalsByNeighborhood, utils.Error)
}

type personDisabilityRepo struct {
	BaseRepo
	db *gorm.DB
}

func NewPersonDisabilityRepo(db *gorm.DB) PersonDisabilityRepo {
	repo := &personDisabilityRepo{
		db: db,
	}

	repo.SetRepo(repo.db)

	return repo
}

func personDisabilityRepoError(message string, code string) utils.Error {
	errorCode := utils.NewErrorCode(utils.DatabaseErrorCode, utils.DisabilityErrorType, code)

	return utils.NewError(message, errorCode)
}

func (n *personDisabilityRepo) GetPersonDisabilities(personId int) ([]model.PersonDisability, utils.Error) {
	var disabilities []model.PersonDisability

	err := n.db.Model(model.PersonDisability{}).Preload("Disability").Where("person_id = ?", personId).Find(&disabilities).Error
	if err != nil {
		return disabilities, personDisabilityRepoError("failed to get the person disabilities", "01")
	}

	return disabilities, utils.Error{}
}

func (n *personDisabilityRepo) GetDisabilityById(disabilityId int) (model.Disability, utils.Error) {
	var disability model.Disability

	err := n.db.Model(model.Disability{}).Where("id = ?", disabilityId).Find(&disability).Error
	if err != nil {
		return disability, personDisabilityRepoError("failed to get the disability", "02")
	}

	return disability, utils.Error{}
}

func (n *personDisabilityRepo) UpsertPersonDisability(personDisability model.PersonDisability, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	err := databaseConn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "person_id"}, {Name: "disability_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"acquired"}),
	}).Create(&personDisability).Error

	if err != nil {
		return personDisabilityRepoError("failed to upsert the person disability", "03")
	}

	return utils.Error{}
}

func (n *personDisabilityRepo) ClearPersonDisability(personId int, tx *gorm.DB) utils.Error {
	databaseConn := n.db

	if tx != nil {
		databaseConn = tx
	}

	if err := databaseConn.Where("person_id = ?", personId).Delete(&model.PersonDisability{}).Error; err != nil {
		return personDisabilityRepoError("failed to clear the person disability", "04")
	}

	return utils.Error{}
}

func (n *personDisabilityRepo) CountDisability() (model.DisabilityTotals, utils.Error) {
	var result []struct {
		Category string
		Total    int
	}

	query := `
		SELECT d.category, COUNT(*) AS total
		FROM person_disabilities pd
		JOIN disabilities d ON pd.disability_id = d.id
		GROUP BY d.category;
	`

	if err := n.db.Raw(query).Scan(&result).Error; err != nil {
		return model.DisabilityTotals{}, personDisabilityRepoError("failed to count the disabilities", "05")
	}

	totals := model.DisabilityTotals{}
	for _, row := range result {
		switch row.Category {
		case "Visual":
			totals.Visual = row.Total
		case "Hearing":
			totals.Hearing = row.Total
		case "Physical":
			totals.Physical = row.Total
		case "Intellectual":
			totals.Intellectual = row.Total
		case "Psychosocial":
			totals.Psychosocial = row.Total
		}
	}

	return totals, utils.Error{}
}

func (n *personDisabilityRepo) CountDisabilityByNeighborhood(neighborhood string) (model.DisabilityTotalsByNeighborhood, utils.Error) {
	var result []struct {
		Category string
		Total    int
	}

	query := `
		SELECT d.category, COUNT(*) AS total
		FROM person_disabilities pd
		JOIN disabilities d ON pd.disability_id = d.id
		JOIN people p ON pd.person_id = p.id
		JOIN addresses a ON p.address_id = a.id
		WHERE REPLACE(LOWER(NormalizeText(a.neighborhood)), ' ', '') = REPLACE(LOWER(NormalizeText(?)), ' ', '')
		GROUP BY d.category;
	`

	if err := n.db.Raw(query, neighborhood).Scan(&result).Error; err != nil {
		return model.DisabilityTotalsByNeighborhood{}, personDisabilityRepoError("failed to count the disabilities by neighborhood", "06")
	}

	totals := model.DisabilityTotalsByNeighborhood{}
	for _, row := range result {
		switch row.Category {
		case "Visual":
			totals.Visual = row.Total
		case "Hearing":
			totals.Hearing = row.Total
		case "Physical":
			totals.Physical = row.Total
		case "Intellectual":
			totals.Intellectual = row.Total
		case "Psychosocial":
			totals.Psychosocial = row.Total
		}
	}

	return totals, utils.Error{}
}
