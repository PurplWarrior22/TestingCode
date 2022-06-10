package models

import (
	"encoding/json"
	"fmt"
)

//this file holds shared models for requirements files
type Requirement struct {
	System        string `json:"system"`
	Validation_id string `json:"validation_id"`
	Description   string `json:"description"`
	Originator    string `json:"originator"`
	Start_Date    string `json:"start_date"`
	Stop_Date     string `json:"stop_date"`
	Nipf          string `json:"nipf"`
	//type is enum of Standing or Dynamic
	Type RequirementType `json:"type"`
}

const (
	Standing      RequirementType = "standing"
	Dynamic       RequirementType = "adhoc"
	Amplification RequirementType = "amplification"
)

type RequirementType string

func (rt *RequirementType) UnmarshalJSON(b []byte) error {
	var s string
	json.Unmarshal(b, &s)
	requirementType := RequirementType(s)
	switch requirementType {
	case Standing, Dynamic, Amplification:
		{
			*rt = requirementType
			return nil
		}
	}
	return fmt.Errorf("invalid Leave type, expected one of %s, %s, got %s", Standing, Dynamic, string(b))
}
