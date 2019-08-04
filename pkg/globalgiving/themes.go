package globalgiving

type Theme struct {
	ID   string
	Name string
}

func GetThemes() []Theme {
	return []Theme{
		{ID: "animals", Name: "Animals"},
		{ID: "children", Name: "Children"},
		{ID: "climate", Name: "Climate Change"},
		{ID: "democ", Name: "Democracy and Governance"},
		{ID: "disaster", Name: "Disaster Recovery"},
		{ID: "ecdev", Name: "Economic Development"},
		{ID: "edu", Name: "Education"},
		{ID: "env", Name: "Environment"},
		{ID: "finance", Name: "Microfinance"},
		{ID: "gender", Name: "Women and Girls"},
		{ID: "health", Name: "Health"},
		{ID: "human", Name: "Humanitarian Assistance"},
		{ID: "rights", Name: "Human Rights"},
		{ID: "sport", Name: "Sport"},
		{ID: "tech", Name: "Technology"},
		{ID: "hunger", Name: "Hunger"},
		{ID: "art", Name: "Arts and Culture"},
		{ID: "lgbtq", Name: "LGBTQAI+"},
	}
}
