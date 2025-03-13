package url

type AllowOption func(*allowChecks)

func WithSchemeCheck(check AllowedSchemeCheck) AllowOption {
	return func(checks *allowChecks) {
		checks.isAllowedScheme = check
	}
}

func WithUserCheck(check AllowedUserCheck) AllowOption {
	return func(checks *allowChecks) {
		checks.isAllowedUser = check
	}
}

func WithHostCheck(check AllowedHostCheck) AllowOption {
	return func(checks *allowChecks) {
		checks.isAllowedHost = check
	}
}

func optionsToChecks(options []AllowOption) *allowChecks {
	checks := defaultAllowChecks()

	for _, o := range options {
		o(checks)
	}

	return checks
}
