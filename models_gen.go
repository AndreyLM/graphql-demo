// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql_demo

type NewUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type NewVideo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int    `json:"userId"`
	URL         string `json:"url"`
}
