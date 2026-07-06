// Additional providers - 22 more providers to reach 100+

pub struct AdditionalProviders;

impl AdditionalProviders {
    pub fn get_all_providers() -> Vec<Provider> {
        vec![
            // 22 New providers
            Provider { id: "expanse".into(), name: "Expanse Studios".into(), category: "slots".into(), games: 50 },
            Provider { id: "solidicon".into(), name: "Solidicon".into(), category: "slots".into(), games: 40 },
            Provider { id: "print".into(), name: "Print Studios".into(), category: "slots".into(), games: 35 },
            Provider { id: "fantasma".into(), name: "Fantasma Games".into(), category: "slots".into(), games: 30 },
            Provider { id: "truelab".into(), name: "TrueLab".into(), category: "slots".into(), games: 25 },
            Provider { id: "peter_sons".into(), name: "Peter & Sons".into(), category: "slots".into(), games: 30 },
            Provider { id: "apparat".into(), name: "Apparat".into(), category: "slots".into(), games: 25 },
            Provider { id: "bulletpot".into(), name: "Bulletpot".into(), category: "slots".into(), games: 20 },
            Provider { id: "3oaks".into(), name: "3 Oaks Gaming".into(), category: "slots".into(), games: 45 },
            Provider { id: "boomerang".into(), name: "Boomerang".into(), category: "slots".into(), games: 40 },
            Provider { id: "kalamba".into(), name: "Kalamba Games".into(), category: "slots".into(), games: 55 },
            Provider { id: "allforone".into(), name: "Allforone".into(), category: "slots".into(), games: 20 },
            Provider { name: "Bangaroo".into(), id: "bangaroo".into(), category: "slots".into(), games: 25 },
            Provider { id: "betsolutions".into(), name: "BetSolutions".into(), category: "slots".into(), games: 30 },
            Provider { id: "betermatic".into(), name: "Betermatic".into(), category: "slots".into(), games: 25 },
            Provider { id: "drl".into(), name: "Dr离婚".into(), category: "slots".into(), games: 35 },
            Provider { id: "eagaming".into(), name: "EA Gaming".into(), category: "live".into(), games: 40 },
            Provider { id: "lucky_streak".into(), name: "Lucky Streak".into(), category: "live".into(), games: 25 },
            Provider { id: "medialive".into(), name: "MediaLive".into(), category: "live".into(), games: 30 },
            Provider { id: "xprogaming".into(), name: "XPRO Gaming".into(), category: "live".into(), games: 25 },
            Provider { id: "globalgaming".into(), name: "Global Gaming".into(), category: "table".into(), games: 15 },
            Provider { id: "adsens".into(), name: "AdsPlay".into(), category: "slots".into(), games: 20 },
        ]
    }
}

pub struct Provider {
    pub id: String,
    pub name: String,
    pub category: String,
    pub games: i32,
}
