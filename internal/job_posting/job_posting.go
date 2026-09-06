package jobposting


type Listing struct {
	Title       string   	`json:"title"`
	URL         string  	`json:"url"`
	CompanyName string  	`json:"company_name"`
	CompanyURL  string   	`json:"company_url"`
	Source      string   	`json:"source"`
	Category    string   	`json:"category"`
	ID          string   	`json:"id"`
	Active      bool        `json:"active"`
	Terms       []string 	`json:"terms"`
	Locations   []string 	`json:"locations"`
	Sponsorship string   	`json:"sponsorship"`
	IsVisible   bool     	`json:"is_visible"`
	Degrees     []string 	`json:"degrees"`
	DateUpdated int64   	`json:"date_updated"`		
	DatePosted	int64		`json:"date_posted"`
}

