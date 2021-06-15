package auth

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/bakape/meguca/config"
)

// List of historical ruler titles
var titles = [...]string{
	"Confessor",
	"Ecclesiastic",
	"Venetian",
	"Prithvi-vallabha",
	"Illustrious",
	"Sebastos",
	"Young",
	"Jagirdar",
	"Pir",
	"Guardian Immortal",
	"Pirani",
	"Satrap",
	"Short",
	"Legatus",
	"Assistant to the President & Deputy National Security Advisor",
	"Fidalgo",
	"Lawgiver",
	"Tribune",
	"Chieftain",
	"Fortunate",
	"Chansonnier",
	"Seer",
	"Humanist",
	"Assistant Professor",
	"Apostle",
	"Tirbodh",
	"God-Given",
	"Kind",
	"Dilochitès",
	"Earl",
	"Red",
	"Pursuivant",
	"She-Wolf of France",
	"Assistant in Virtue",
	"Tyrant",
	"Toqui",
	"Distinguished Professor",
	"Amban",
	"Aspet",
	"Grand prince",
	"Saint",
	"Leading Aircraftman",
	"Apostle",
	"Dathapatish",
	"Peaceful",
	"Iron",
	"Stadtholder",
	"Bishop",
	"Master of the Sacred Palace",
	"Bapu",
	"Bloodthirsty",
	"Augusta",
	"Chhatrapati",
	"Trierarch",
	"Old",
	"Ganden Tripa",
	"Bewitched",
	"Precious",
	"Hammer of the Scots",
	"Fair Sun",
	"Allamah",
	"Professor",
	"Moor",
	"Patriarch",
	"Lover of Elegance",
	"Consul",
	"Cardinal-King",
	"Quarreller",
	"Curly",
	"Concubinus",
	"Recipient from the Inner Chamber",
	"Yishi",
	"Abbot",
	"Nushi",
	"Philosopher",
	"Desired",
	"Baroness",
	"Leading Aircraftwoman",
	"Fat",
	"Sacristan",
	"Unready",
	"Deed-Doer",
	"Junior Technician",
	"God's Wife",
	"Commissioner of Baseball",
	"Soter",
	"Sluggard",
	"Terrible",
	"Shaman",
	"Ilarchès",
	"Ancient",
	"Able",
	"Imperator",
	"Hekatontarchès",
	"Hazarapatish",
	"Captain",
	"Nuncio",
	"Champion",
	"Stern Counselor",
	"All-fair",
	"Fidei defensor",
	"King of Arms",
	"Most Beautiful",
	"Beauty",
	"Goodman",
	"Restorer",
	"Battler",
	"Agha",
	"African",
	"Crusader",
	"Admiral of the Fleet",
	"Absolutist",
	"Impotent",
	"Amir al-Mu'minin",
	"Kanstresios",
	"Corporal",
	"Learned",
	"High priestess",
	"Centurion",
	"Elder",
	"Madman",
	"Silent",
	"God-Loving",
	"Nomarch",
	"Hojatoleslam",
	"Inconstant",
	"King",
	"Bruce",
	"Agonothetes",
	"Magnanimous",
	"Sailor King",
	"Maharao",
	"Sun King",
	"Bastard",
	"Saver of Europe",
	"Archdeacon",
	"Governor-General",
	"Roju",
	"Bad",
	"Ell-High",
	"Ayatollah",
	"Mad",
	"Corrupted",
	"Lecturer",
	"Victorious",
	"Brown",
	"Patroon",
	"Dikastes",
	"Pilgrim",
	"Astrologer",
	"Grand duke",
	"Treacherous",
	"Warlike",
	"Baivarapatish",
	"Trembling",
	"Nobilissimus",
	"Towel Attendant",
	"Despot",
	"Maharaja",
	"Thakurani",
	"Dàifu",
	"Righteous",
	"Herzog",
	"Jiàoshòu",
	"Herald",
	"Chakravartin",
	"Hegumenor Hegumenia",
	"Shifu",
	"Monsignor",
	"Bearded",
	"Kolakretai",
	"Foreign minister",
	"Independentist",
	"Furén",
	"Councillor Pensionary",
	"Upajjhaya",
	"Mullah",
	"Grand Master",
	"Sushri",
	"Field Marshal",
	"Master of the Horse",
	"Clubfoot",
	"Tsar",
	"Great",
	"Sergeant",
	"Albanian-slayer",
	"Tenzo",
	"Prytaneis",
	"Inexorable",
	"Sahib",
	"Basilissa",
	"Peshwa",
	"Pope",
	"Dame",
	"Kumar",
	"Sibyl",
	"Slacker",
	"Kind-Hearted",
	"Soldier-King",
	"Determined",
	"Lugal",
	"Lochagos",
	"Memorable",
	"President pro tempore",
	"Longshanks",
	"Maha-kshtrapa",
	"Vain",
	"Optio",
	"Brilliant",
	"Cardinal-nephew",
	"Builder King",
	"Thakore",
	"Yuvraj",
	"Temple boy",
	"Purple-Born",
	"Comes",
	"Good Mother",
	"Popular",
	"President",
	"Wicked",
	"Devil",
	"Buddha",
	"Sahibah",
	"Handsome Fairness",
	"Crowned",
	"Catholic",
	"Noble",
	"Arahant",
	"Abbess",
	"Exile",
	"Fellow",
	"Empress",
	"Good",
	"Lionheart",
	"Dimoirites",
	"Begum",
	"Ceremonious",
	"Unique",
	"Vicereine",
	"Quiet",
	"Deacon",
	"Immoral",
	"Dux",
	"Blessed",
	"White",
	"Orphan",
	"Asapatish",
	"God-Like One",
	"Protodeacon",
	"Baron",
	"Agoranomos",
	"Taxiarch",
	"Well-Beloved",
	"Well-Served",
	"Syntagmatarchis",
	"Liberator",
	"Fratricide",
	"Sunim",
	"Emperor-Sacristan",
	"Sacrifice",
	"Unchaste",
	"Palatine",
	"Blessed",
	"Vicar general",
	"Last",
	"Bloody",
	"Tough",
	"Affable",
	"Hellenotamiae",
	"Principal Lecturer",
	"Argbadh",
	"Maharana",
	"Lady",
	"Talented",
	"Lamb",
	"Spirited",
	"Duke",
	"Shofet",
	"Post Master General",
	"Zhuxi",
	"Paygan Salarapoo",
	"Mandarin",
	"Avenger",
	"Kumari",
	"Metropolitan Bishop",
	"Dung-Named",
	"Cardinal",
	"Slut",
	"Troubadour",
	"Lord High Constable",
	"Wench",
	"Adopted",
	"Selected Lady",
	"Rani",
	"Sluggard",
	"Countess",
	"Middle",
	"Shōgun",
	"Chairman",
	"Lord Great Chamberlain",
	"Wise",
	"Chaste",
	"Hidalgo",
	"Deacon",
	"Vajracharya",
	"Merry",
	"Drunkard",
	"Anax",
	"Unavoidable",
	"Emperor",
	"Mawlawi",
	"Consort",
	"Sorcerer",
	"Bolognian",
	"Conqueror",
	"Madam",
	"Duchess",
	"Priest Hate",
	"Outlaw",
	"Archimandrite",
	"Thunderbolt",
	"Bold",
	"Diplomat",
	"Imam",
	"Bald",
	"Handsome",
	"Grand Admiral",
	"Simple",
	"The Most Honourable",
	"Be-shitten",
	"Maharani",
	"Gentle",
	"Sir",
	"Harlot",
	"Posthumous",
	"Warrior",
	"Service Provider",
	"Little Impaler",
	"Nakharar",
	"Hardy",
	"Magister Militum",
	"Major archbishop",
	"Courteous",
	"Pastor",
	"Indolent",
	"Ephor",
	"Malikah",
	"Fighter",
	"Soldier",
	"Caesar",
	"Gong",
	"Strategos",
	"Chancellor",
	"Spahbod",
	"Shimu",
	"Priest",
	"Tanuter",
	"Archduchess",
	"Theoroi",
	"Magnificent",
	"Priest",
	"Sultan",
	"Usurper",
	"Fowler",
	"Proxenos",
	"Perfect Prince",
	"Praetor",
	"Caliph",
	"Esquire",
	"Lord",
	"Broom-plant",
	"Polemarch",
	"Allower",
	"Nawab",
	"Damned",
	"Exarch",
	"Inquisitor",
	"Holy",
	"Coiffure Attendant",
	"Eloquent",
	"Hammer",
	"Swami",
	"Viceroy",
	"Rabbi",
	"Hunchback",
	"Young King",
	"Hierodeacon",
	"Pacific",
	"Unlucky",
	"Oath-Taker",
	"Beloved",
	"Tremulous",
	"Khan",
	"Easy",
	"Archpriest",
	"Xiaozhang",
	"Gyoja",
	"Venerable",
	"Servant",
	"Taitai",
	"Theorodokoi",
	"Basileus",
	"Philosopher King",
	"Cabbage",
	"Farmer",
	"Constable Prince",
	"Hopeful",
	"Rash",
	"Prudent",
	"Savior",
	"Evangelist",
	"Enlightened",
	"Earl Marshal",
	"Akhoond",
	"Alderman",
	"Chanyu",
	"Bookish",
	"Senior Aircraftwoman",
	"Archiater",
	"Wrymouth",
	"High priest",
	"Amphipole",
	"Younger",
	"Proud",
	"Presbyter",
	"Archduke",
	"Savakabuddha",
	"Princeps",
	"Commissioner",
	"Hieromonk",
	"Sacristan",
	"Chorbishop",
	"Raja",
	"Maiden",
	"Just",
	"Prince",
	"Pious",
	"Moneybag",
	"Virgin Queen",
	"Tagmatarchis",
	"Air Marshal",
	"Tetrarch",
	"Taoiseach",
	"Hierophant",
	"Ill-Tempered",
	"Prime minister",
	"Law-Mender",
	"Elbow-High",
	"Nawabzadi",
	"Subedar",
	"Arhat",
	"Edifier",
	"Divine Adoratrice",
	"Xiaojie",
	"Khawaja",
	"Varma",
	"Sebaste",
	"Tall",
	"Holy Prince",
	"Redemptress",
	"Generalissimo",
	"German",
	"Swift",
	"Generous",
	"National Security Advisor",
	"Peacemaker",
	"Councillor",
	"Reverend",
	"Lame",
	"Propagator of Deportment",
	"Queen",
	"Gothi",
	"Weak",
	"Brave",
	"Impaler",
	"Dalai Lama",
	"Unready",
	"Professor Emeritus",
	"Tsaritsa",
	"Rightly Guided",
	"Gong Bao",
	"Headman",
	"Grand duchess",
	"Shah",
	"Haty-a",
	"Senior Lecturer",
	"Brash",
	"Spider",
	"Chairwoman",
	"Squire",
	"Fickle",
	"Blond",
	"Fair",
	"Mild",
	"Goodwife",
	"Executioner",
	"Theologian",
	"Weiyuán",
	"Lax",
	"Lands’ Advocate",
	"Red King",
	"Voivode",
	"Xry Hbt",
	"Nizam",
	"Glorious",
	"Stout",
	"Black Prince",
	"Ecumenical Patriarch",
	"Stammerer",
	"Liberal",
	"Reader",
	"Saint",
	"Countess",
	"Zamindar",
	"Brave",
	"Indolent",
	"Longhaired king",
	"Princess",
	"Archon",
	"Singer",
	"Chiliarch",
	"Aswaran Salar",
	"One-Eyed",
	"Xiansheng",
	"Populator",
	"Naib",
	"Tramp",
	"Redless",
	"Black",
	"Catholicos",
	"Grand Inquisitor",
	"Senior Aircraftman",
	"Mayor",
	"Sparapet",
	"Fan-bearer on the Right Side of the King",
	"Magister Officiorum",
	"Constant",
	"Cruel",
	"Candid",
	"Mouth",
	"Ambitious",
	"Chief",
	"Lonko",
	"Presiding Patriarch",
	"Aesymnetes",
	"Khazar",
	"Artist-King",
	"Fürstor Fürstin",
	"Dean",
	"En",
	"Caulker",
	"Recipient of Edicts",
	"Shigong",
	"Pharaoh",
	"Doctor",
	"Careless",
	"Corrector",
	"Oppressed",
	"Navarch",
	"Aggressor",
	"Hadrat",
	"Choregos",
	"Strong",
	"Upasaka",
	"Grim",
	"Karo",
	"Primate",
	"Governor",
	"Aircraftman",
	"Associate Professor",
	"Reformer",
	"Monk",
	"Fearless",
	"Bear",
	"Sri",
	"Zongshi",
	"Elder",
	"Redeless",
	"Bavarian",
	"Lisp and Lame",
	"Secretary of State",
	"Navigator",
	"Pale",
	"Martyr",
	"Apodektai",
	"Marzban",
	"Rector",
	"Imperatrice",
	"Somatophylax",
	"Servant in the Place of Truth",
	"Accursed",
	"Chief",
	"Arab",
	"Whore",
	"Hipparchus",
	"Diwan",
	"Lady of Treasure",
	"Big Nest",
	"Powerful",
	"Donor Doña",
	"Broad-shouldered",
	"Hunter",
	"Sebastokrator",
	"Count",
	"Great Elector",
	"Grand Pensionary",
	"Debonaire",
	"Invincible",
	"Lion",
	"Mirza",
	"Liberal",
	"City Manager",
	"Archbishop",
	"Hairy",
	"Valiant",
	"Rajmata",
	"Diakonissa",
	"Yisheng",
	"Crosseyed",
	"Lord Privy Seal",
	"Mighty",
	"Boneless",
	"Decurio",
	"Sakellarios",
	"Laoshi",
	"Humane",
	"Apostate",
	"Blind",
	"Capacidónio",
	"Dom",
	"Nawabzada",
	"Faqih",
	"Mahatma",
	"Yuvrani",
	"Ætheling",
	"Tóngzhi",
	"Starosta",
	"Captain",
	"Unsui",
	"The Right Honourable",
	"Desai",
	"Epihipparch",
	"Malik",
}

// Hash buffer and produce a cryptogrphically safe title and the source hash
// in hex format for displaying to users
func HashToTitle(buf []byte) (title string, hash string) {
	h := sha256.New()
	h.Write(buf)
	h.Write([]byte(config.Get().Salt))
	digest := h.Sum(nil)
	return titles[int(binary.LittleEndian.Uint64(digest)>>1)%len(titles)],
		hex.EncodeToString(digest)
}