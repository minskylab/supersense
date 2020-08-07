package supersense

// Person is a struct to identify a person who post any on supersense
type Person struct {
	Name        string
	Photo       string
	SourceOwner string
	Email       *string
	ProfileURL  *string
	Username    *string
}
