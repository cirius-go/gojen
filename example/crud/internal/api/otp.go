package api

// Otp represents the otp API.
type Otp struct {
	svc OtpSvc
}

// NewOtp returns a new Otp instance.
func NewOtp(svc OtpSvc) *Otp {
	return &Otp{
		svc: svc,
	}
}
