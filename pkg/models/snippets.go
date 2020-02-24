package models

type Snippet struct {
	Name          string
	Category      SnippetCategory
	UsageID       string
	UsageName     string
	DescriptionRu string
	DescriptionEn string
	Code          string
}

type SnippetResource struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	UsageID     string `json:"usage_id"`
	UsageName   string `json:"usage_name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type SnippetListResource []SnippetResource

type SnippetCategory int

const (
	TagSnippet SnippetCategory = iota + 1
	AdapterSnippet
)

func (s Snippet) ToResource(locale string) SnippetResource {
	var description string
	switch locale {
	case "en":
		description = s.DescriptionEn
	case "ru":
		description = s.DescriptionRu
	default:
		description = s.DescriptionEn
	}

	resource := SnippetResource{
		Name:        s.Name,
		Category:    s.Category.String(),
		UsageID:     s.UsageID,
		UsageName:   s.UsageName,
		Description: description,
		Code:        s.Code,
	}
	return resource
}

func StringToSnippetCategory(s string) (SnippetCategory, bool) {
	switch s {
	case TagSnippet.String():
		return TagSnippet, true
	case AdapterSnippet.String():
		return AdapterSnippet, true
	}
	return 0, false
}

func (snippetCategory SnippetCategory) String() string {
	names := [...]string{
		"unknown",
		"tag",
		"adapter"}

	if snippetCategory < TagSnippet || snippetCategory > AdapterSnippet {
		return names[0]
	}
	return names[snippetCategory]
}
