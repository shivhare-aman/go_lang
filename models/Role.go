package models

// Role enum definition using iota
type Role int

const (
	RoleAdmin Role = iota
	RoleUser
)

// String method to convert Role enum to string
func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "RoleAdmin"
	case RoleUser:
		return "RoleUser"
	}
	return ""
}

// ParseRole converts a string to a Role enum
func ParseRole(roleStr string) Role {

	if roleStr == "RoleAdmin" {
		return RoleAdmin
	}
	return RoleUser
}
