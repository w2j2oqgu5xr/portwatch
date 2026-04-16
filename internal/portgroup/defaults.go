package portgroup

// DefaultGroups contains well-known port groupings.
var DefaultGroups = []struct {
	Name  string
	Ports []int
}{
	{"web", []int{80, 443, 8080, 8443}},
	{"database", []int{3306, 5432, 1433, 27017, 6379}},
	{"mail", []int{25, 465, 587, 110, 143, 993, 995}},
	{"ssh", []int{22}},
	{"dns", []int{53}},
	{"ftp", []int{20, 21}},
}

// LoadDefaults registers all DefaultGroups into the given Registry.
// Groups that fail to register (e.g. duplicates) are silently skipped.
func LoadDefaults(r *Registry) {
	for _, d := range DefaultGroups {
		_ = r.Add(d.Name, d.Ports)
	}
}
