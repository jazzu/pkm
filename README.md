# Pelaajakuvamasiina
Golla kirjoitettu liima Counter Strike: Global Offensiven kuunteluun ja CasparCG:n sekä vMixin ohjaamiseen halutun USB-liitäntäisen pelaajakameran valitsemiseksi.
# Tarkoituksenmukainen kytkentä
```
[Pelaaja #1] ───┐
[Pelaaja #2] ───┤
[Pelaaja #3] ───┼ USB ─→┃CasparCG serveri #1┠── SDI ─┐
[Pelaaja #4] ───┤               ↑                    │
[Pelaaja #5] ───┘               │                    │
                            AMCP/Telnet              │
                                │                    ↓
┃Observer┠─────── JSON ───→┃PKM kone┠── XML ──────→[vMix]
                                │                    ↑
                            AMCP/Telnet              │
[Pelaaja #6] ───┐               │                    │
[Pelaaja #7] ───┤               ↓                    │
[Pelaaja #8] ───┼ USB ─→┃CasparCG serveri #2┠── SDI ─┘
[Pelaaja #9] ───┤
[Pelaaja #10] ──┘
```
# TODO
- [ ] Tee parempi README.md
