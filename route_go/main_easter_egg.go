package main

import (
	"C"
	"math/rand"
)

//export Do
func Do() *C.char {
	select_list := []string{
		"PWD0ZbR7AOY",
		"HoU29ljOmTE",
		"PR2vUm-Ald8",
		"opZoEmsu_Lo",
		"txZFFTusSvw",
		"Ixq9xL2tvRU",
		"-3IAx_r4Au0",
		"wObZkycA6sc",
		"hZxYLa97gDg",
		"hwn2kw4eFJM",
		"wX2t_8HOtiY",
		"tLQjcf45fKE",
		// Remix by NyxTheShield
	}
	select_str := select_list[rand.Intn(len(select_list)-1)]

	return C.CString("<iframe width=\"640\" height=\"360\" src=\"https://www.youtube.com/embed/" + select_str + "\" frameborder=\"0\" allowfullscreen></iframe>")
}

func main() {

}
