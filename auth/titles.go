package auth

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/bakape/meguca/config"
)

// List of historical ruler titles
var titles = [...]string{
	// Historical
	"Abbess",
	"Able",
	"Absolutist",
	"Accursed",
	"Admiral of the Fleet",
	"Adopted",
	"Affable",
	"African",
	"Aggressor",
	"Air Marshal",
	"Aircraftman",
	"Albanian-slayer",
	"All-fair",
	"Allower",
	"Ambitious",
	"Ancient",
	"Apostate",
	"Apostle",
	"Arab",
	"Archbishop",
	"Archduchess",
	"Archduke",
	"Archon",
	"Archpriest",
	"Artist-King",
	"Assistant Professor",
	"Assistant in Virtue",
	"Assistant to the President & Deputy National Security Advisor",
	"Associate Professor",
	"Astrologer",
	"Avenger",
	"Bad",
	"Bald",
	"Baron",
	"Baroness",
	"Bastard",
	"Battler",
	"Bavarian",
	"Bear",
	"Bearded",
	"Beauty",
	"Beloved",
	"Bewitched",
	"Bishop",
	"Black",
	"Black Prince",
	"Blessed",
	"Blind",
	"Blond",
	"Bloodthirsty",
	"Bloody",
	"Bold",
	"Boneless",
	"Bookish",
	"Brash",
	"Brave",
	"Brilliant",
	"Broad-shouldered",
	"Brown",
	"Bruce",
	"Buddha",
	"Builder King",
	"Cabbage",
	"Caesar",
	"Candid",
	"Captain",
	"Cardinal",
	"Careless",
	"Catholic",
	"Ceremonious",
	"Chairman",
	"Chairwoman",
	"Champion",
	"Chancellor",
	"Chaste",
	"Chief",
	"Chieftain",
	"City Manager",
	"Clubfoot",
	"Commissioner",
	"Commissioner of Baseball",
	"Confessor",
	"Conqueror",
	"Consort",
	"Constable Prince",
	"Constant",
	"Corporal",
	"Corrector",
	"Corrupted",
	"Councillor",
	"Councillor Pensionary",
	"Count",
	"Countess",
	"Courteous",
	"Crosseyed",
	"Crowned",
	"Cruel",
	"Crusader",
	"Curly",
	"Dame",
	"Damned",
	"Dean",
	"Desired",
	"Despot",
	"Determined",
	"Devil",
	"Diplomat",
	"Distinguished Professor",
	"Doctor",
	"Drunkard",
	"Duchess",
	"Duke",
	"Dung-Named",
	"Earl",
	"Earl Marshal",
	"Easy",
	"Elbow-High",
	"Elder",
	"Eloquent",
	"Emperor",
	"Empress",
	"Enlightened",
	"Evangelist",
	"Executioner",
	"Exile",
	"Fair",
	"Farmer",
	"Fat",
	"Fearless",
	"Fellow",
	"Fickle",
	"Field Marshal",
	"Fighter",
	"Foreign minister",
	"Fortunate",
	"Fratricide",
	"Generalissimo",
	"Generous",
	"Gentle",
	"German",
	"Glorious",
	"God's Wife",
	"God-Given",
	"God-Like One",
	"God-Loving",
	"Good",
	"Good Mother",
	"Goodman",
	"Goodwife",
	"Governor",
	"Governor-General",
	"Grand Admiral",
	"Grand Inquisitor",
	"Grand Master",
	"Grand Pensionary",
	"Grand duchess",
	"Grand duke",
	"Grand prince",
	"Great",
	"Great Elector",
	"Grim",
	"Guardian Immortal",
	"Hairy",
	"Hammer",
	"Hammer of the Scots",
	"Handsome",
	"Hardy",
	"Harlot",
	"Headman",
	"Herald",
	"High priest",
	"High priestess",
	"Holy",
	"Holy Prince",
	"Hopeful",
	"Humane",
	"Humanist",
	"Hunchback",
	"Hunter",
	"Ill-Tempered",
	"Illustrious",
	"Immoral",
	"Impaler",
	"Imperator",
	"Imperatrice",
	"Impotent",
	"Inconstant",
	"Indolent",
	"Inexorable",
	"Inquisitor",
	"Invincible",
	"Iron",
	"Junior Technician",
	"Just",
	"Khan",
	"Kind",
	"Kind-Hearted",
	"King",
	"King of Arms",
	"Lady",
	"Lady of Treasure",
	"Lamb",
	"Lame",
	"Last",
	"Law-Mender",
	"Lawgiver",
	"Lax",
	"Leading Aircraftman",
	"Leading Aircraftwoman",
	"Learned",
	"Lecturer",
	"Liberal",
	"Liberator",
	"Lion",
	"Lionheart",
	"Lisp and Lame",
	"Little Impaler",
	"Longshanks",
	"Lord",
	"Lover of Elegance",
	"Mad",
	"Madam",
	"Madman",
	"Magnanimous",
	"Magnificent",
	"Maiden",
	"Major archbishop",
	"Mandarin",
	"Martyr",
	"Master of the Horse",
	"Master of the Sacred Palace",
	"Mayor",
	"Memorable",
	"Merry",
	"Metropolitan Bishop",
	"Mighty",
	"Mild",
	"Moneybag",
	"Monk",
	"Most Beautiful",
	"National Security Advisor",
	"Navigator",
	"Noble",
	"Oath-Taker",
	"Old",
	"One-Eyed",
	"Oppressed",
	"Orphan",
	"Outlaw",
	"Pacific",
	"Pale",
	"Pastor",
	"Patriarch",
	"Peaceful",
	"Peacemaker",
	"Perfect Prince",
	"Pharaoh",
	"Philosopher",
	"Philosopher King",
	"Pilgrim",
	"Pious",
	"Pope",
	"Popular",
	"Populator",
	"Posthumous",
	"Powerful",
	"Precious",
	"President",
	"Presiding Patriarch",
	"Priest",
	"Priest Hate",
	"Primate",
	"Prime minister",
	"Prince",
	"Princess",
	"Principal Lecturer",
	"Professor",
	"Professor Emeritus",
	"Propagator of Deportment",
	"Proud",
	"Prudent",
	"Purple-Born",
	"Quarreller",
	"Queen",
	"Quiet",
	"Rash",
	"Reader",
	"Rector",
	"Red",
	"Red King",
	"Redeless",
	"Reformer",
	"Restorer",
	"Righteous",
	"Rightly Guided",
	"Sacrifice",
	"Sailor King",
	"Saint",
	"Saver of Europe",
	"Savior",
	"Secretary of State",
	"Seer",
	"Selected Lady",
	"Senior Aircraftman",
	"Senior Aircraftwoman",
	"Senior Lecturer",
	"Sergeant",
	"Servant",
	"Service Provider",
	"Shaman",
	"She-Wolf of France",
	"Short",
	"Shōgun",
	"Silent",
	"Simple",
	"Singer",
	"Sir",
	"Slacker",
	"Sluggard",
	"Slut",
	"Soldier",
	"Soldier-King",
	"Sorcerer",
	"Spider",
	"Spirited",
	"Squire",
	"Stern Counselor",
	"Stout",
	"Strategos",
	"Strong",
	"Sultan",
	"Sun King",
	"Swift",
	"Talented",
	"Tall",
	"Temple boy",
	"Terrible",
	"The Most Honourable",
	"Theologian",
	"Thunderbolt",
	"Tough",
	"Towel Attendant",
	"Tramp",
	"Treacherous",
	"Trembling",
	"Tremulous",
	"Troubadour",
	"Tsar",
	"Tsaritsa",
	"Tyrant",
	"Unavoidable",
	"Unchaste",
	"Unique",
	"Unlucky",
	"Unready",
	"Usurper",
	"Vain",
	"Valiant",
	"Venerable",
	"Venetian",
	"Viceroy",
	"Victorious",
	"Virgin Queen",
	"Warlike",
	"Warrior",
	"Weak",
	"Well-Beloved",
	"Well-Served",
	"Wench",
	"White",
	"Whore",
	"Wicked",
	"Wise",
	"Wrymouth",
	"Young",
	"Young King",
	"Younger",

	// Fate
	"Saber",
	"Lancer",
	"Archer",
	"Rider",
	"Caster",
	"Assassin",
	"Berserker",
	"Alter Ego",
	"Foreigner",
	"Shielder",
	"Gunner",
	"Faker",
	"Gatekeeper",
	"Watcher",
	"Voyager",
	"Temptress",
	"Boxer",
	"Butler",
	"Gorgeous",
	"Idol",
	"Manslayer",
	"Monster",
	"Beast",
	"Noble Phantasm",

	// Shadows
	"Shadow",
	"Living Doll",
	"Star Bearer",
	"Grandfather",
	"Clinger",
	"Ghost",

	// Puppets
	"Sword-Plundering Nemesis",
	"Edgeless Blade",
	"Austere Snow Crow",
	"Enigmatic Gale",
	"Phantom Bird",
	"Sky Hunter",
	"Sharp-Eyed Impaler",
	"Frozen Wonder",
	"Final Punishment",
	"Night's Lament",
	"Eradicator of Life",
	"Screaming Phoenix Killer",
	"Honorable Sage",
	"Skeleton Who Scorns the Heavens",
	"Bones of Creation",
	"Vicious",
	"Charming Huntress",
	"Withered Life",
	"Immortal Iron Flute",
	"Wandering Shaman Bard",
	"Song That Dooms Evil",
	"Listening Fang",
	"Roaring Maniac",
	"Hunting Fox",
	"Scorpion Necklace",
	"Princess of Cruelty",
	"Plague of Moths and Locusts",
	"Shockwave",
	"Absolute Emptiness",
	"Seven Blasphemous Deaths",
	"Death King",
	"Elder Sun Lord",
	"Destiny",
	"Mantra",
	"Derision",
	"Heavenly Crafter",
	"Paragon of a Hundred Battles",
	"Beyond Indistinct",
	"Crescent-Bladed Mantis",
	"Ghost Seizing Heavenly Crafter",
	"Shining Reign",
	"White Lotus",

	// Puella Magi Witches
	"Walpurgisnacht",
	"Brazilian Night",

	// GuP Tank Roles
	"Tank Commander",
	"Tank Radio Operator",
	"Tank Gunner",
	"Tank Loader",
	"Tank Driver",

	// Stereotypes and Tropes, Misc.
	"Tsundere",
	"Yandere",
	"Kuudere",
	"Dandere",
	"Yangire",
	"Tsunshun",
	"Chuunibyou",
	"Baka",
	"MC",
	"Loser Protagonist",
	"Harem Protagonist",
	"Shōnen Protagonist",
	"Childhood Friend",
	"Student Council President",
	"Student Council Vice-President",
	"Student Council Secretary",
	"Student Council Treasurer",
	"Class Representative",
	"Club President",
	"Sensei",
	"School Nurse",
	"Transfer Student",
	"Senpai",
	"Onee-sama",
	"Siscon",
	"Aniki",
	"Brocon",
	"Loli",
	"Lolicon",
	"Shota",
	"Shotacon",
	"Oba-san",
	"OL",
	"Salaryman",
	"Tomboy",
	"Gyaru",
	"Bishoujo",
	"Bishōnen",
	"Yamato Nadeshiko",
	"Otokonoko",
	"Meganekko",
	"Mahou Shoujo",
	"Majokko",
	"Nekomimi",
	"Miko",
	"Ojou-sama",
	"Hime-sama",
	"Meido",
	"Elegant Gothic Lolita",
	"Cute Gothic Lolita",
	"Classic Lolita",
	"Sweet Lolita",
	"Kuro Lolita",
	"Shiro Lolita",
	"Military Lolita",
	"Pervert",
	"Panty Thief",
	"Hikikomori",
	"Otaku",
	"Weeaboo",
	"Fujoshi",
	"Densha Otaku",
	"Cosplay Otaku",
}

// Hash buffer and produce a cryptogrphically safe title and the source hash
// in hex format for displaying to users
func HashToTitle(buf []byte) (title string,
	hash string) {
	h := sha256.New()
	h.Write(buf)
	h.Write([]byte(config.Get().Salt))
	digest := h.Sum(nil)
	return titles[int(binary.LittleEndian.Uint64(digest)>>1)%len(titles)],

		hex.EncodeToString(digest)
}