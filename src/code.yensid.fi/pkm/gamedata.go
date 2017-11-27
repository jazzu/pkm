package pkm

// FIXME: Tämä tyyppimäärittely on täysin overkill tällä hetkellä. Jokin https://github.com/jmoiron/jsonq kaltainen
// palikka voisi palvella paremmin.
type (
	GameData struct {
		Provider         *Provider         `json:"provider"`
		Map              *Map              `json:"map"`
		Round            *Round            `json:"round"`
		PlayerID         *PlayerID         `json:"player"`
		PlayerState      *StatePlayer      `json:"player_state"`
		PlayerWeapons    *PlayerWeapons    `json:"player_weapons"`
		PlayerMatchStats *MatchStatsPlayer `json:"player_match_stats"`
	}

	Provider struct {
		Name      string `json:"name"`
		AppID     int    `json:"appid"`
		Version   int    `json:"version"`
		SteamID   string `json:"steamid"`
		Timestamp int64  `json:"timestamp"`
	}

	Map struct {
		Name   string         `json:"name"`
		Phase  string         `json:"phase"`
		Round  int            `json:"round"`
		TeamCT map[string]int `json:"team_ct"`
		TeamT  map[string]int `json:"team_t"`
	}

	Round struct {
		Phase   string `json:"phase"`
		Bomb    string `json:"bomb"`
		WinTeam string `json:"win_team"`
	}

	PlayerID struct {
		SteamID  string `json:"steamid"`
		Name     string `json:"name"`
		Activity string `json:"activity"`
		Team     string `json:"team"`
	}

	StatePlayer struct {
		Player     PlayerState `json:"player"`
		Previously struct {
			Player PlayerState `json:"player"`
		} `json:"previously"`
		Added struct {
			Player PlayerState `json:"player"`
		} `json:"added"`
	}

	WeaponsPlayer struct {
		Player     PlayerWeapons `json:"player"`
		Previously struct {
			Player PlayerWeapons `json:"player"`
		} `json:"previously"`
		Added struct {
			Player PlayerWeapons `json:"player"`
		} `json:"added"`
	}

	PlayerWeapons struct {
		Weapons struct {
			Weapon01 Weapon `json:"weapon_01"`
			Weapon02 Weapon `json:"weapon_02"`
			Weapon03 Weapon `json:"weapon_03"`
			Weapon04 Weapon `json:"weapon_04"`
			Weapon05 Weapon `json:"weapon_05"`
			Weapon06 Weapon `json:"weapon_06"`
		} `json:"weapons"`
	}

	MatchStatsPlayer struct {
		Player     PlayerMatchStats `json:"player"`
		Previously struct {
			Player PlayerMatchStats `json:"player"`
		} `json:"previously"`
		Added struct {
			Player PlayerMatchStats `json:"player"`
		} `json:"added"`
	}

	PlayerState struct {
		State State `json:"state"`
	}

	PlayerMatchStats struct {
		MatchStats MatchStats `json:"match_stats"`
	}

	State struct {
		Health      int  `json:"health"`
		Armor       int  `json:"armor"`
		Helmet      bool `json:"helmet"`
		Flashed     int  `json:"flashed"`
		Smoked      int  `json:"smoked"`
		Burning     int  `json:"burning"`
		Money       int  `json:"money"`
		RoundKills  int  `json:"round_kills"`
		RoundKillHS int  `json:"round_killhs"`
	}

	Weapon struct {
		Name        string `json:"name"`
		Paintkit    string `json:"paintkit"`
		Type        string `json:"type"`
		State       string `json:"state"`
		AmmoClip    int    `json:"ammo_clip"`
		AmmoClipMax int    `json:"ammo_clip_max"`
		AmmoReserve int    `json:"ammo_reserve"`
	}

	MatchStats struct {
		Kills   int `json:"kills"`
		Assists int `json:"assists"`
		Deaths  int `json:"deaths"`
		MVPS    int `json:"mvps"`
		Score   int `json:"score"`
	}
)
