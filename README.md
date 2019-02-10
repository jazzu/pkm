# Pelaajakuvamasiina
Golla kirjoitettu liima Counter Strike: Global Offensiven kuunteluun ja OBS:n ohjaamiseen halutun USB-liitäntäisen pelaajakameran valitsemiseksi.

Kaksi videoserveriä lähettävät kuvaa NDI:llä vmixille tms kuvamikserille, alphakanavan kera siten että vain akviinen kamera näkyy, jommalta kummalta serveriltä. Nämä voidaan sitten aktivoida PIP:nä tai muuten halutusti.

# Tarkoituksenmukainen kytkentä
```
[Pelaaja #1] ───┐
[Pelaaja #2] ───┤
[Pelaaja #3] ───┼ USB ─→┃videoserveri A┠── NDI ─┐
[Pelaaja #4] ───┤               ↑               │
[Pelaaja #5] ───┘               │               │
                            websocket           │
                                │               ↓
┃Observer┠─────── JSON ───→┃PKM kone┃         [vMix]
                                │               ↑
                            websocket           │
[Pelaaja #6] ───┐               │               │
[Pelaaja #7] ───┤               ↓               │
[Pelaaja #8] ───┼ USB ─→┃videoserveri B┠── NDI ─┘
[Pelaaja #9] ───┤
[Pelaaja #10] ──┘
```

# Riippuvuudet ja asennus
1. OBS studio https://obsproject.com/
2. NDI-plugin https://github.com/Palakis/obs-ndi
3. NDI redist http://new.tk/NDIRedistV3
4. websocket-plugin https://github.com/Palakis/obs-websocket

Asenna ylläolevat molemmille videoservereille. Avaa/luo tyhjä OBS scene collection. Tarvitaan yksi scene, nimeltään "Scene1". Tälle aktivoidaan NDI dedicated output, suositeltavaa nimetä ne esim. "serverA" ja "serverB".

Kaikki kamarat tuodaan tähän sceneen, nimettynä "cam1"-"cam10". Kameroiden numerointi kannattaa aloittaa pelaajien selän takaa katsoen vasemmalta. Kaikki kamerat `fit-to-screen`, ja normaalitilassa kaikkien `visibility` pois päältä.

Websocket-plugin on oletuksena portissa 4444. Tällä hetkellä tässä projektissa ei ole autentikaatiotukea, koska järjestelmä on tarkoitettu vain suljetussa verkossa ajettavaksi.

Observer-koneelle asennetaan kansioon `steamapps\common\Counter-Strike Global Offensive\csgo\cfg` GSI-asetustiedosto (esim. `gamestate_integration_pkm.cfg`. Pelin pitää pyöriä samassa verkossa tai palomuurissa pitää olla aukko peliverkosta PKM-koneen websocket-porttiin (oletus 1999).

[Lisätiedot GSI:stä.](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration)

Asetustiedostoihin laitetaan pelaajien steamID:t SteamID64-muodossa, ja tiedostoja on yksi per joukkue. Paikat myöskin pelaajien takaa vasemmalta laskien. Paikka 0 tarkoittaa sitä, että pelaajalla ei ole kameraa tai kamera on esimerkiksi väärin suunnattu, ja sen takia halutaan hetkellisesti poistaa käytöstä näin:

  * editoi fileä (muuta kameralle paikaksi 0),
  * keskeytä ajossa oleva ohjelma `ctrl+c` ja
  * käynnistä ohjelma uudelleen.

Serverin käynnistys: `./pkm -A team2.json -B team1.json`

# TODO
- [x] Tee parempi README.md
- [x] siirrä kovakoodatut asiat asetustiedostoon
- [ ] käyttöliittymä käsiohjaukseen
- [x] ja tiimien valintaan
- [ ] siirrä konfiguraatiologiikka obs.go:sta server.go:hon
- [ ] korjaa obs.go:ssa olevat todo-kommentit
- [x] lisää GSI-tiedosto tänne
- [ ] lisää obs-mallifilet tänne
- [ ] tuki eri formaateissa oleville steam-id:ille (esim. https://github.com/MrWaggel/gosteamconv tai https://godoc.org/github.com/Acidic9/steam)