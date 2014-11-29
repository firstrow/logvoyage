package common

var cachedUsers = make(map[string]*User)

func FindCachedUser(email string) *User {
	if u, ok := cachedUsers[email]; ok {
		return u
	} else {
		user, _ := FindUserByEmail(email)
		if user != nil {
			cachedUsers[email] = user
			return cachedUsers[email]
		}
	}
	return nil
}
