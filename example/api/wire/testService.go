package wire

type TestService struct {
	TestName string
}

func NewTestService(testName string) *TestService {
	return &TestService{TestName: testName}
}

func (this *TestService) Name() string {
	return "test11"
}

func (this *TestService) Name2() string {
	return "test22"
}
