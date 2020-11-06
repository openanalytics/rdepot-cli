package model

type Output interface {
	FormatJSON() (string, error)
}
