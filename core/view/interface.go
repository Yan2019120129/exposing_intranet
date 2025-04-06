package view

type AskType string

var (
	AskTypeGet     AskType = "GET"
	AskTypePost    AskType = "POST"
	AskTypePut     AskType = "PUT"
	AskTypeDelete  AskType = "DELETE"
	AskTypePatch   AskType = "PATCH"
	AskTypeHead    AskType = "HEAD"
	AskTypeOptions AskType = "OPTIONS"
)

type Body interface {
	GetOperates() []*Operate
}
