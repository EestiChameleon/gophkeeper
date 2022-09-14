package testdb

// Run method imitates the DB connection by returning TestVault.
func Run() (*TestVault, error) {
	return &TestVault{}, nil
}
