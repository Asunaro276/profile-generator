package model

type User struct {
	Gender string
	Name   struct {
		Title string
		First string
		Last  string
	}
	Location struct {
		Street struct {
			Number int
			Name   string
		}
		City        string
		State       string
		Country     string
		Postcode    string
		Coordinates struct {
			Latitude  string
			Longitude string
		}
	}
	Email string
	Login struct {
		UUID     string
		Username string
		Password string
		Salt     string
		MD5      string
		SHA1     string
		SHA256   string
	}
	Dob struct {
		Date string
		Age  int
	}
	Registered struct {
		Date string
		Age  int
	}
	Phone string
	Cell  string
	ID    struct {
		Name  string
		Value string
	}
	Picture struct {
		Large     string
		Medium    string
		Thumbnail string
	}
	NAT string
}
