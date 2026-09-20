package career

// SkillDefinition is one entry of the skill list (pp. 304-314), which the
// book reprints "for the convenience of players and Referees creating
// characters" from the Clement Sector Core Rulebook.
//
// The engine needs it for the results that do not name a skill: "gain a
// level in any skill of your choice", "gain any skill at level 1 which you
// do not already have". Without the list those are a choice with no
// options.
type SkillDefinition struct {
	Name string

	// Specialties is what the page prints beneath the skill, and is empty
	// for the skills that take none. A result naming a specialty of "Any"
	// chooses from here.
	Specialties []string

	// Open marks the three skills whose specialties the book says are
	// examples rather than a set: Science ("by no means a full list of the
	// possible sciences"), Trade ("by no means a full list of the possible
	// trades") and Language, whose specialties are "any language which can
	// be chosen as a primary language".
	//
	// A specialty a table names for one of these is legitimate whether or
	// not it is printed here.
	Open bool
}

// Skills is the list, in the order the book prints it -- which is
// alphabetical apart from Art, which p. 305 starts after p. 304 has ended
// with Athletics.
//
// Names are the list's, with two exceptions the book's own tables make:
// Administration is written "Admin" and Jack of All Trades is not
// abbreviated. Specialties are the form the tables use where the two
// differ; ERRATA E-35 has the table. Every skill any transcribed table
// names is here, and TestEveryTranscribedSkillIsOnTheList says so.
func Skills() []SkillDefinition {
	return []SkillDefinition{
		{Name: "Admin"},
		{Name: "Advocate", Specialties: []string{"Legal", "Oratory", "Politics"}},
		{Name: "Animals", Specialties: []string{
			"Farming", "Herding", "Riding", "Training", "Veterinary",
		}},
		{Name: "Athletics", Specialties: []string{"Climbing", "Running", "Swimming"}},
		{Name: "Art", Specialties: []string{
			"Acting", "Dance", "Holography", "Instrument", "Painting",
			"Sculpting", "Writing",
		}},
		{Name: "Astrogation"},
		{Name: "Broker"},
		{Name: "Carouse"},
		{Name: "Chef"},
		{Name: "Deception", Specialties: []string{
			"Disguise", "Forgery", "Intrusion", "Lie", "Pickpocket",
		}},
		{Name: "Diplomat"},
		{Name: "Discipline"},
		{Name: "Draw"},
		{Name: "Drive", Specialties: []string{
			"Hover", "Mole", "Tracked", "Wagon", "Walker", "Wheeled",
		}},
		{Name: "Electronics", Specialties: []string{
			"Comms", "Computers", "Electrical Repair", "Remote Operations",
			"Robotics", "Sensors",
		}},
		{Name: "Engineer", Specialties: []string{
			"M-Drive", "Z-Drive", "Life Support", "Power",
		}},
		{Name: "Etiquette"},
		{Name: "Explosives"},
		{Name: "Flyer", Specialties: []string{
			"Glide", "Grav", "Grav Harness", "Reentry Kit", "Rotor", "Wing",
		}},
		{Name: "Gambler"},
		{Name: "Gun Combat", Specialties: []string{
			"Slug Rifle", "Slug Pistol", "Shotgun", "Energy Rifle",
			"Energy Pistol", "Sonic",
		}},
		{Name: "Gunner", Specialties: []string{
			"Capital Weapons", "Ortillery", "Screens", "Turrets",
		}},
		{Name: "Heavy Weapons", Specialties: []string{
			"Large Slug Weapons", "Launchers", "Man-Portable Energy",
			"Field Artillery", "Vehicle Mounted Weapons", "Armor Mounted",
		}},
		{Name: "Instruction"},
		{Name: "Interrogation", Specialties: []string{"Questioning", "Torture"}},
		{Name: "Investigate"},
		{Name: "Jack of All Trades"},
		{Name: "Language", Open: true},
		{Name: "Leadership"},
		{Name: "Mechanic"},
		{Name: "Medic", Specialties: []string{
			"Alteration", "Altrants", "Cryogenics", "Cybernetics", "Diagnosis",
			"First Aid", "Surgery", "Uplifts",
		}},
		{Name: "Melee", Specialties: []string{
			"Added Weapon", "Blade", "Bludgeon", "Natural Weapons",
			"Thrown Weapons", "Unarmed Combat",
		}},
		{Name: "Navigation"},
		{Name: "Persuade"},
		{Name: "Pilot", Specialties: []string{"Small Craft", "Spacecraft"}},
		{Name: "Recon"},
		{Name: "Science", Open: true, Specialties: []string{
			"Altrant Psychology", "Uplift Psychology", "Physics", "Chemistry",
			"Cybernetics", "Biology", "Archeology", "Economics", "History",
			"Philosophy", "Psychology",
		}},
		{Name: "Seafarer", Specialties: []string{
			"Canoe", "Motorboats", "Ocean Ships", "Sail", "Submarine",
		}},
		{Name: "Stealth"},
		{Name: "Streetwise"},
		{Name: "Suit", Specialties: []string{
			"Battle Armor", "Hostile Environment", "Vacc Suit",
		}},
		{Name: "Survival", Specialties: []string{
			"Barren", "Cold", "Desert", "Forest", "Freefall", "Heat",
			"High Gravity", "High Pressure", "Jungle", "Low Gravity",
			"Low Pressure", "Mountains", "Plains", "Ocean", "Swamp",
		}},
		{Name: "Tactics", Specialties: []string{"Military", "Naval", "Sport"}},
		{Name: "Trade", Open: true, Specialties: []string{
			"Bartender", "Blacksmith", "Carpenter", "Cybernetics", "Gardener",
			"Glass Blower", "Hydroponics", "Lumberjack", "Naval Architect",
			"Prospector", "Space Construction",
		}},
	}
}

// SkillByName finds a skill on the list.
func SkillByName(name string) (SkillDefinition, bool) {
	for _, held := range Skills() {
		if held.Name == name {
			return held, true
		}
	}

	return SkillDefinition{}, false
}
