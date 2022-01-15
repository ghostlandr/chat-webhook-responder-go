package tokens

var (
	slackTokens = []string{
		"70XqnEL12zOlA08Fo0lraciE",
		"ayYWtEzhfqh5GcXdEqrD3H3h",
		"7b5WbqiybRqPDRTm2e9GvTUL",
		"x56o3ZQzti2l7YEb7ntRu4gE",
		"NogyqLDNMuukzKqEZmh5Q5l2",
		"5ZdfgUhazo0aUv17yLKK7VFG",
	}
)

// IsAuthorizedToken checks a token passed in to make sure it's one of the ones
// that has been pre-authorized to access the application
func IsAuthorizedToken(token string) bool {
	for _, tok := range slackTokens {
		if tok == token {
			return true
		}
	}
	return false
}
