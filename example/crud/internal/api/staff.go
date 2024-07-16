package api

// Staff represents the staff API.
type Staff struct {
	svc StaffSvc
}

// NewStaff returns a new Staff instance.
func NewStaff(svc StaffSvc) *Staff {
	return &Staff{
		svc: svc,
	}
}
