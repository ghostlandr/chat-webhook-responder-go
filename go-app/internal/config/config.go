package config

import "os"

var (
	DefinerClientID  = "11970859830.2903858978897"
	UdefinerClientID = "11970859830.2902453792341"
)

var DefinerClientSecret string = os.Getenv("DEFINERCLIENTSECRET")
var DefineSigningSecret string = os.Getenv("DEFINESIGNINGSECRET")
var UdefinerClientSecret string = os.Getenv("UDEFINERCLIENTSECRET")
var UdefineSigningSecret string = os.Getenv("UDEFINESIGNINGSECRET")
