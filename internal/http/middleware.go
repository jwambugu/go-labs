package httpapi

type Validator interface {
	Validate() error
}
