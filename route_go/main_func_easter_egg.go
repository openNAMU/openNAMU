package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	select_list := []string{
		"PWD0ZbR7AOY", // Shanghai Teahouse ~ Chinese Tea
		"HoU29ljOmTE", // Flawless Clothing of Celestials
		"PR2vUm-Ald8", // U.N. Owen Was Her
		"opZoEmsu_Lo", // Night of Nights
		"txZFFTusSvw", // Reach for the Moon ~ Immortal Smoke
		"Ixq9xL2tvRU", // Phantom Ensemble
		"-3IAx_r4Au0", // Entrusting This World to Idols ~ Idolatrize World
		"wObZkycA6sc", // Last Remote
		"hZxYLa97gDg", // Emotional Skyscraper ~ Cosmic Mind
		"hwn2kw4eFJM", // Border of Life
		"wX2t_8HOtiY", // Voyage 1969
		"tLQjcf45fKE", // Necrofantasia
		"7DvMRAMuMrU", // Where Is That Bustling Marketplace Now ~ Immemorial Marketeers
		// Remix by NyxTheShield
		"SXFP9HgWBYQ", // 세계는 귀엽게 만들어져 있다
		"YDrgO0Oj3Fg", // 죽취비상
		"wxWV_sUGPB0", // 디자이어 드라이브
		"uw0h2O7UaZ8", // 100번째 블랙 마켓
		"blE4lnfEWbU", // 일렉트릭 헤리티지
		// Remix by KR. Palto47
	}
	select_str := select_list[rand.Intn(len(select_list)-1)]

	fmt.Print("<iframe width=\"640\" height=\"360\" src=\"https://www.youtube.com/embed/" + select_str + "\" frameborder=\"0\" allowfullscreen></iframe>")
}
