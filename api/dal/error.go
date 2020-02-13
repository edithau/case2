package dal

// UnauthorizedErr is the error returned when password is incorrect.
type UnauthorizedErr struct{}

func (e *UnauthorizedErr) Error() string { return "invalid password" }

// CredentialsErr is returned when the credentials don't match any user.
type CredentialsErr struct{}

func (e *CredentialsErr) Error() string { return "invalid credentials" }
