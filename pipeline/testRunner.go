package pipeline

// RunTests runs a serivce's tests
func RunTests(srv *Service) error {
	return srv.execTests()
}
