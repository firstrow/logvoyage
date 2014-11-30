package common

var cachedUsers = make(map[string]*User)

// Search user and store record in memory
func FindCachedUser(email string) (*User, error) {
	if u, ok := cachedUsers[email]; ok {
		return u, nil
	} else {
		user, err := FindUserByEmail(email)
		if err != nil {
			return nil, err
		}
		if user != nil {
			cachedUsers[email] = user
			return cachedUsers[email], nil
		}
	}
	return nil, nil
}
