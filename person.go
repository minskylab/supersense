package supersense

// Person is a struct to identify a person who post any on supersense
type Person struct {
	Name       string `json:"name"`
	Photo      string `json:"photo"`
	Owner      string `json:"owner"`
	Email      *string `json:"email"`
	ProfileURL *string `json:"profileURL"`
	Username   *string `json:"username"`
}
