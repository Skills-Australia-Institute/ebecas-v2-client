package ebecasv2client

// CreatedBy represents the eBECAS user who created a resource.
type CreatedBy struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Pagination defines pagination parameters for list and search requests.
type Pagination struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

// SortDirection represents the direction used when sorting results.
type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

// FilterType represents how multiple filter conditions are combined.
type FilterType string

const (
	FilterTypeAnd FilterType = "and"
)

// Operator represents a comparison operator used in filter conditions.
type Operator string

const (
	OperatorEq Operator = "eq"
)

// StudentField represents a searchable or sortable student field.
type StudentField string

const (
	StudentFieldID            StudentField = "students-id"
	StudentFieldFirstName     StudentField = "students-first_name"
	StudentFieldLastName      StudentField = "students-last_name"
	StudentFieldStudentNumber StudentField = "students-student_number"
)

// StudentSort defines the field and direction used to sort student results.
type StudentSort struct {
	Field     StudentField  `json:"field"`
	Direction SortDirection `json:"direction"`
}

// StudentFilter defines how student search conditions are combined.
type StudentFilter struct {
	Type       FilterType         `json:"type"`
	Conditions []StudentCondition `json:"conditions"`
}

// StudentCondition defines a single filter condition for a student search.
type StudentCondition struct {
	Field    StudentField `json:"field"`
	Operator Operator     `json:"operator"`
	Value    string       `json:"value"`
}
