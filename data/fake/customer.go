package fake


type Customer struct {
	Name string `fake:"{name}" json:"name"`
	Phone string `fake:"+{phone}" json:"phone"`
	Email string `fake:"{email}" json:"email"`
}