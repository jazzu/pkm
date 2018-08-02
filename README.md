# Pelaajakuvamasiina
Golla kirjoitettu liima Counter Strike: Global Offensiven kuunteluun ja OBS:n ohjaamiseen halutun USB-liitäntäisen pelaajakameran valitsemiseksi.

Kaksi videoserveriä lähettävät kuvaa NDI:llä vmixille tms kuvamikserille, alphakanavan kera siten että vain akviinen kamera näkyy, jommalta kummalta serveriltä. Nämä voidaan sitten aktivoida PIP:nä tai muuten halutusti.

# Tarkoituksenmukainen kytkentä
```
[Pelaaja #1] ───┐
[Pelaaja #2] ───┤
[Pelaaja #3] ───┼ USB ─→┃videoserveri A┠── NDI ─┐
[Pelaaja #4] ───┤               ↑                    │
[Pelaaja #5] ───┘               │                    │
                            websocket              │
                                │                    ↓
┃Observer┠─────── JSON ───→┃PKM kone┠ ←── XML ──────[vMix]
                                │                    ↑
                            websocket              │
[Pelaaja #6] ───┐               │                    │
[Pelaaja #7] ───┤               ↓                    │
[Pelaaja #8] ───┼ USB ─→┃videoserveri B┠── NDI ─┘
[Pelaaja #9] ───┤
[Pelaaja #10] ──┘
```

# Riippuvuudet ja asennus
1. OBS studio https://obsproject.com/
2. NDI-plugin https://github.com/Palakis/obs-ndi
3. NDI redist http://new.tk/NDIRedistV3
4. websocket-plugin https://github.com/Palakis/obs-websocket

Asenna ylläolevat molemmille videoservereille. Avaa/luo tyhjä OBS scene collection. Tarvitaan yksi scene, nimelään "Scene1". Tälle aktivoidaan NDI dedicated output, suositeltavaa nimetä ne esim serverA ja serverB

Kaikki kamarat tuodaan tähän sceneen, nimettynä "cam1"-"cam10". Kameroiden numerointi kannattaa aloittaa pelaajien selän takaa katsoen vasemmalta. Kaikki kamerat fit to screen, ja normaalitilassa kaikkien visibility pois päältä.

Websocket plugin on oletuksena portissa 4444, tällähetkellä ei autentikaatiotukea koska järjestelmä on tarkoitettu vain suljetussa verkossa ajettavaksi.

Observer-koneelle asennetaan GSI-asetustiedosto, peli tarvii pyöriä samassa verkossa tai palomuurissa pitää olla aukko peliverkosta PKM-koneen websocket-porttiin (oletus 1999)

Asetustiedostoon pelaajien steamID:t SteamID64-muodossa, järjestyksellä ei ole väliä.

# TODO
- [ ] Tee parempi README.md
- [ ] siirrä kovakoodatut asiat asetustiedostoon
- [ ] käyttöliittymä käsiohjaukseen ja tiimien valintaan
