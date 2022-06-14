package entity

type User struct {
	// FirstName of the User
	//
	// required: true
	// example: `Jared`
	Firstname string `json:"firstname" bson:"firstname"`

	// LastName of the User
	//
	// required: true
	// example: `Jackson`
	Lastname string `json:"lastname" bson:"lastname"`

	// Username of the User
	//
	// required: true
	// example: `JaredsonUsername`
	Username string `json:"username" bson:"username"`
}
