package salesforce

type Description struct {
    Name    string              `json:"name"`
    Label   string              `json:"label"`
    Fields  []DescriptionField  `json:"fields"`
}
type DescriptionField struct {
    Name    string  `json:"name"`
    Label   string  `json:"label"`
    Type    string  `json:"type"`
}


