package business

type businessBase interface {
	SetInstitution(institution int)
	SetSpecialty(specialty int)
	GetSpecialty() int
}
