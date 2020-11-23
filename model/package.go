package model

//import (
//"encoding/json"
//)

type Package struct {
	Id                 int        `json:"id"`
	Version            string     `json:"version"`
	Submission         Submission `json:"submission"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Author             string     `json:"author"`
	Depends            string     `json:"depends"`
	Imports            string     `json:"imports"`
	Suggests           string     `json:"suggests"`
	SystemRequirements string     `json:"systemRequirements"`
	License            string     `json:"license"`
	Title              string     `json:"title"`
	Url                string     `json:"url"`
	Source             string     `json:"source"`
	Md5Sum             string     `json:"md5sum"`
	Active             bool       `json:"active"`
	Deleted            bool       `json:"deleted"`
}

//func (pkgs []Package) FormatJSON() ([]byte, error) {
//return json.Marshal(pkgs)
//}

type Submission struct {
	Id       int    `json:"id"`
	Changes  string `json:"changes"`
	Accepted bool   `json:"accepted"`
	Deleted  bool   `json:"deleted"`
}
