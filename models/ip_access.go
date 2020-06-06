package models

// IPAccess represents a user ip access event with nonessential columns from
// the database model dropped
type IPAccess struct {
	*Geography
	IP        string `json:"ip"`
	Speed     int64  `json:"speed"`
	Timestamp int64  `json:"timestamp"`
}
