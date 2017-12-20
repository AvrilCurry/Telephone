package config

var telephoneID = 1

func GetID() int {
	return telephoneID
}

func SetID() {
	telephoneID = telephoneID + 1
}
