package business

import (
	"octopus.com/mi/models"
	"octopus.com/mi/services"
	"octopus.com/mi/utils"
)

type AsyncConsultBusiness interface {
	businessBase
	Init(level int,service services.AsyncConsultService)
	OpenedConsults() (consults []models.AsyncConsult,found bool, err error)
	RespondedConsults() (consults []models.AsyncConsult,found bool, err error)
	AllConsults() (consults []models.AsyncConsult,found bool, err error)
	GetConsult(idConsult int) (consult models.AsyncConsult,found bool, err error)
	CanResponse() bool
}

type AsyncConsultParticularity interface {
	GetAllByStatus(status int) (consults []models.AsyncConsult,found bool, err error)
}

// AsyncConsultBusinessBase - base
type AsyncConsultBusinessBase struct {
	Institution int
	Level int
	Specialty int
	Service services.AsyncConsultService
	Particularity AsyncConsultParticularity
}

// asyncConsultBusinessClient - client business layer
type asyncConsultBusinessClient struct {
	AsyncConsultBusinessBase
}

// asyncConsultBusinessProfessional - professional business layer
type asyncConsultBusinessProfessional struct {
	AsyncConsultBusinessBase
}


// ====== BASE =====
func (cc *AsyncConsultBusinessBase) Init(level int,service services.AsyncConsultService){
	cc.Level = level
	cc.Service = service
}

func (cc *AsyncConsultBusinessBase) SetInstitution(institution int){
	cc.Institution = institution
}

func (cc *AsyncConsultBusinessBase) SetSpecialty(specialty int){
	cc.Specialty = specialty
}

func (cc AsyncConsultBusinessBase) Get(idConsult int) ( models.AsyncConsult, bool, error){
	return cc.Service.Get(idConsult,cc.Institution)
}

func (cc AsyncConsultBusinessBase) OpenedConsults() (consults []models.AsyncConsult,found bool, err error){
	return cc.Particularity.GetAllByStatus(services.CONSULT_STATUS_OPEN)
}

func (cc AsyncConsultBusinessBase) RespondedConsults() (consults []models.AsyncConsult,found bool, err error){
	return cc.Particularity.GetAllByStatus(services.CONSULT_STATUS_RESPONDED)
}

func (cc AsyncConsultBusinessBase) AllConsults() (consults []models.AsyncConsult,found bool, err error){
	return cc.Particularity.GetAllByStatus(services.CONSULT_STATUS_ALL)
}

func (cc AsyncConsultBusinessBase) CanResponse() bool{
	return false
}

func (cc AsyncConsultBusinessBase) GetSpecialty() int{
	return cc.Specialty
}

func (cc AsyncConsultBusinessBase) GetConsult(idConsult int) (consult models.AsyncConsult,found bool, err error){
	return cc.Service.Get(idConsult,cc.Institution)
}

// getAllByStatus - Need be overrided
// @override
//func (cc AsyncConsultBusinessBase) getAllByStatus(status int) (consults []models.AsyncConsult,found bool, err error){
//	return
//}

// ===== Client ======
// getAllByStatus -
// @override
func (cc asyncConsultBusinessClient) GetAllByStatus(status int) (consults []models.AsyncConsult,found bool, err error){
	consults, found, err = cc.Service.GelAllByInstitution(cc.Institution,cc.Specialty,status)
	if err != nil {
		utils.Error("[AsyncConsultBusinessClient][Consult] - Error on get consults (status,error)",status,err.Error())
	}
	return
}


// ===== Business ======

// getAllByStatus -
// @override
func (cc asyncConsultBusinessProfessional) GetAllByStatus(status int) (consults []models.AsyncConsult,found bool, err error){
	consults, found, err = cc.Service.GelAllBySpecialty(cc.Specialty,status)
	if err != nil {
		utils.Error("[AsyncConsultBusinessProfessional][Consult] - Error on get consults (status,error)",status,err.Error())
	}
	return
}

func (cc asyncConsultBusinessProfessional) CanResponse() bool{
	return true
}
