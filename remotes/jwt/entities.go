package jwt

// Token -
type Token struct {
	ID          string                 `json:'id'`
	Permission  string                 `json:'permission'`
	Institution string                 `json:'institution'`
	Authorized  bool                   `json:'authorized'`
	Broker      map[string]interface{} `json:'broker'`
}
