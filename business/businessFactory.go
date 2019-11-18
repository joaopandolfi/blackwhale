package business

import (
	"octopus.com/mi/models"
	"octopus.com/mi/services"
)

// NewPropertyAsyncConsultBusiness - create business correct to user type
// @factory
func NewPropertyAsyncConsultBusiness(level int,service services.AsyncConsultService)(b AsyncConsultBusiness){
	switch level {
		case models.USER_CLIENT:
			c := &asyncConsultBusinessClient{}
			c.Particularity  = c
			b = c
		break;
		case models.USER_ESPECIALIST:
		case models.USER_DOCTOR:
		default: // Olhar isso daqui com calma
			c := &asyncConsultBusinessProfessional{}
			c.Particularity  = c
			b = c
	}

	b.Init(level,service)
	return
}
