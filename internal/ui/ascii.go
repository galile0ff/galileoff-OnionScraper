package ui

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Art struct {
	Name  string
	Color string
	Data  string
}

var arts = []Art{
	{"banner1", ColorCyan, banner1},
	{"banner2", ColorRed, banner2},
	{"banner3", ColorGreen, banner3},
	{"banner4", ColorYellow, banner4},
	{"banner5", ColorPurple, banner5},
	{"banner6", ColorBlue, banner6},
	{"banner7", ColorWhite, banner7},
	{"banner8", ColorGreen, banner8},
	{"banner9", ColorRed, banner9},
	{"banner10", ColorCyan, banner10},
	{"banner11", ColorYellow, banner11},
	{"banner12", ColorPurple, banner12},
	{"banner13", ColorBlue, banner13},
	{"banner14", ColorWhite, banner14},
	{"banner15", ColorRed, banner15},
}

func RandomArt() Art {
	rand.Seed(time.Now().UnixNano())
	return arts[rand.Intn(len(arts))]
}

func PrintRandomBanner() {
	art := RandomArt()
	fmt.Println(art.Color + art.Data + ColorReset)
}

func GetArt(name string) (Art, error) {
	for _, a := range arts {
		if a.Name == name {
			return a, nil
		}
	}
	return RandomArt(), errors.New("ascii bulunamadı")
}

const banner1 = `
                .---.
          _/__~0_\_
     .---'   .---.   '---.           KOMUTAN LOGAR...
    / .---. /_____\ .---. \         ----------------------------------
   | |   |  | ^ ^ |  |   | |         "Şifreli bir cisim yaklaşıyor!"
   | |   |  |  o  |  |   | |         
   | '---'  '-----'  '---' |         "Bu bir uçak mı? Hayır..."
    \                     /          "Bu bir kuş mu? Hayır..."
     \                   /           "Bu bir EXIT NODE!"
      \                 /
       '---------------'             (G.O.R.A - Onion Defense System)
           \       /
            \     / 
             '---'
`
const banner2 = `
    _______________________
   /                       \
  |   W A N T E D          |
  |                        |          AZİZİM, BURASI VAHŞİ AĞ...
  |      _.-------._       |         ----------------------------
  |    .'    ___    '.     |          "- Senin o statik IP..."
  |   /     (o_o)     \    |          "- Kanun tanımaz!"
  |  |       (_)       |   |          "- Ama burada şerif benim."
  |   \     __|__     /    |          
  |    '.___|___|___.'     |          (Ödül: 1000 Onion Coin)
  |          | |           |
  |__________| |___________|
`
const banner3 = `
             /^\
            /   \
           /_____\
          |       |                  BİZ HACKER DEĞİLİZ!
        __|_______|__               --------------------
       /             \               "Bak şu an veritabanı boş..."
      (  _  .---.  _  )              "Hokus Pokus!"
       \  \|( . )|/  /               "Şimdi dolu!"
        \  |  |  |  /
         \ '-----' /                 (Sihirbazlık değil, Scraping)
          '-------'
`
const banner4 = `
       .-------.
      /   Rx    \
     /___________\                   USER-AGENT YETMEZ...
    |             |                 ---------------------
    |  __________ |                  "Bunun <head>'i ağrıyor,"
    | |  ______  ||                  "<body>'sinde sızı var."
    | | | (o)  | ||
    | | |______| ||                  "Sana bi Tor devresi yazıyorum,"
    | |__________| |                  "Günde 3 kere IP değiştireceksin."
    |_____________|
`
const banner5 = `
      __
     /  \      ___       vVv
    | oo |    (o o)     (O O)         VALLA GOL OLUR!
    |__u_|    ( v )     ( w )        -----------------
    /|  |\    /| |\     /| |\         "Bak bu portu görüyor musun?"
   / |  | \  / | | \   / | | \        "Bu port boşuna açık değil."
     /  \      / \       / \          "Ben o veriye golümü atarım!"
    /____\    /___\     /___\
`
const banner6 = `
     .-------------------.
    /| .---------------. |\
   | | |T O R   A G I| | |          PEKİ NSA DE BİZİ GÖRECEK Mİ?
   | | |_______________| | |        -----------------------------
   | | _________________ | |         "- Şifreli bu evladım..."
   | |/___(o)_____(o)___\| |         "- Tünel kazıyoruz tünel!"
   |_______________________|         "- Müren Bey görmez ama..."
   /_______________________\         "- Biz her şeyi görürüz."
  /_________________________\
`
const banner7 = `
         .---.
        /  ^  \
       /  / \  \                    ROBOT DEĞİLİM BEN!
      /  /   \  \                  --------------------
     /   '---'   \                  "Senin o 'Captcha' dediğin..."
    |    (o_o)    |                 "Benim için çerez parası."
     \    (_)    /
      \  __|__  /                   "Kurtar beni 403'ten!"
       \_______/
`
const banner8 = `
        _________
       / _______ \
      / /       \ \                 LİSANSLI MI BU ABİ?
     | |   DVD   | |               --------------------
     | |  [TOR]  | |                "- Orjinal veri abicim."
     | |         | |                "- Sinemada çekim değil,"
      \ \_______/ /                 "- Direk serverdan çekim!"
       \_________/
`
const banner9 = `
          _  _
        (  )(  )
      (  (TOR)   )                  VERİ NERDE?
     (____________)                --------------
           //                       "- Bulutta."
          //                        "- Bulut nerde?"
       _ //_                        "- Yağmur oldu, aktı..."
      ( o_o )                       "- E veri?"
       (___)                        "- Islandı abicim."
`
const banner10 = `
      /¯¯¯¯¯¯¯¯¯\ 
     |   [REC]   |                  TAMAM CANIM, IŞIK VERİN!
     |   .---.   |                 ------------------------
     |  / ( ) \  |                  "Çekiyorum veriyi..."
     | |   |   | |                  "Çek, çek, çek..."
     |  \ (_) /  |                  "Hah, tam <meta> etiketinden al!"
     |   '---'   |
      \_________/
`
const banner11 = `
       .------.
      | |¯¯|¯| |                    KARIŞIK YAPTIM
      | |__| | |                   ---------------
      | .--. | |                    "Onion linkleri..."
      | |OO| | |                    "Proxy zincirleri..."
      | '==' | |                    "Ortaya karışık yaptım."
       '------'
`
const banner12 = `
         _
        {_}
        | |                         BU NEYMİŞ LEMİ?
        | |                        ------------------
      .-| |-.                       "- Onion Elixir efendim."
     /  | |  \                      "- İçince ne oluyor?"
    |   TOR   |                     "- IP'niz kayboluyor."
    |         |                     "- Sevdim bunu."
     \_______/
`
const banner13 = `
       .---.
      ([o_o])                       İNSANIM BEN İNSAN!
      /|___|\                      --------------------
     //| H |\\                      "Bana ff. deme..."
       |___|                        "Bana 404 de!"
      /|   |\                       "Bana 502 de!"
     /_|   |_\
`
const banner14 = `
     .           .
      \  .-.  /                     
     .- ( o ) -.                   ------------------------
      /  '-'  \                     "- Ateşi buldum!"
     '    |    '                    "- Bırak ateşi..."
          |                         "- Veriyi buldun mu?"
       .--^--.                      "- Veri taşa yazılı şefim."
      /_______\
`
const banner15 = `
       (  ) 
        )(                          DOKUNMA O PORTA!
       (__)                        -------------------
      (____)                        "Yanarsın..."
     (______)                       "IP ban yersin..."
    (________)                      "Biz burda sanat yapıyoruz!"
     \ .--. /
      \|o |/
       |__|
`
