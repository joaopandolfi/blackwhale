package business

import (
	"octopus.com/mi/models"
	"octopus.com/mi/services"
)
// AVALIAR PRA VER SE VALE A PENA

type HomeBusiness interface {
	businessBase
	Init(level int,service services.PatientService)
	GetPatients()(patients []models.Patient,found bool, err error)
}

// AsyncHomeBusinessBase - base
type AsyncHomeBusinessBase struct {
	Institution int
	Level int
	Specialty int
	Service services.PatientService
}

// asyncHomeBusinessClient - client business layer
type asyncHomeBusinessClient struct {
	AsyncHomeBusinessBase
}

// asyncHomeBusinessProfessional - professional business layer
type asyncHomeBusinessProfessional struct {
	AsyncHomeBusinessBase
}

// ====== BASE =====
func (cc *AsyncHomeBusinessBase) Init(level int,service services.PatientService){
	cc.Level = level
	cc.Service = service
}

func (cc *AsyncHomeBusinessBase) SetInstitution(institution int){
	cc.Institution = institution
}

func (cc *AsyncHomeBusinessBase) SetSpecialty(specialty int){
	cc.Specialty = specialty
}



// ===== HERANCA ======
