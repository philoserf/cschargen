package career

// GraduateSchool is the school of pp. 96-100.
//
// It is the first institution whose prerequisite is not a characteristic:
// "the character must have achieved Success in an Undergraduate University
// or Military Academy" (p. 97). A second success there is a doctorate
// rather than a second master's (p. 99).
func GraduateSchool() Institution {
	return Institution{
		Name:           "Graduate School",
		Cite:           "pp. 96-100",
		Prerequisites:  []Check{{Characteristic: "EDU", Number: 10}},
		RequiresDegree: true,
		GraduateEDU:    12,
		Admission:      graduateAdmission(8),
		Success:        graduateSuccess(8),
		Honors:         graduateHonors(10),
		Skills:         undergraduateSkills(),
		Failure:        graduateFailure(),
		Events:         graduateEvents(),
		LifeEvents:     nil,
	}
}

// MedSchool is the school of pp. 100-103. It is named for the institution
// rather than for the enlistment modifier of the same name, which is why
// the constructor is not MedicalSchool.
func MedSchool() Institution {
	return Institution{
		Name:           "Medical School",
		Cite:           "pp. 100-103",
		Prerequisites:  []Check{{Characteristic: "EDU", Number: 8}},
		RequiresDegree: true,
		GraduateEDU:    12,
		Admission:      graduateAdmission(8),
		// ERRATA E-28: the box on p. 100 says EDU 9+ and the prose on
		// p. 101 says "roll 8 or higher". The prose wins.
		Success:    graduateSuccess(8),
		Honors:     graduateHonors(10),
		Skills:     []Effect{skill("Medic", "Any")},
		Failure:    medicalFailure(),
		Events:     medicalEvents(),
		LifeEvents: nil,
	}
}

// graduateFailure is the 2d6 table of p. 98.
func graduateFailure() MishapTable {
	return MishapTable{
		{
			Summary: "expelled for serious misconduct",
			Effects: []Effect{
				noAdmissionFor(2),
				chr("EDU", -2),
			},
		},
		{Summary: "the work proves beyond you", Effects: []Effect{chr("EDU", -1)}},
		{
			Summary: "funding that ends unexpectedly",
			Effects: []Effect{chr("EDU", -1), skill("Admin")},
		},
		{
			Summary: "the wrong professor offended, and an environment turned hostile",
			Effects: []Effect{chr("EDU", -1), relationshipAt(Rival, 1, -65)},
		},
		{
			Summary: "a crisis at home that outweighed the work",
			Effects: []Effect{loseRelative(), ratingOfRole("parent", false, 50)},
		},
		{
			Summary: "burnout, and a withdrawal before failure became unavoidable",
			Effects: []Effect{
				chr("EDU", -1),
				throwModifier("admission to any higher education", 2),
			},
		},
		{Summary: "graduate school was simply not the right course", Effects: nil},
		{
			Summary: "an accident that made progress impossible",
			Effects: []Effect{skill("Medic", "First Aid")},
		},
		{
			Summary: "competition with peers that turned toxic",
			Effects: []Effect{relationshipAt(Rival, 1, -65), skill("Persuade")},
		},
		{
			Summary: "a programme that collapsed through no fault of your own",
			Effects: []Effect{pickSkill("Admin", "Advocate", "Broker")},
		},
		{
			Summary: "a position too valuable to refuse, before completion",
			Effects: []Effect{
				relationshipAt(Contact, 1, 30),
				unimplemented("enlist automatically into any career"),
			},
		},
	}
}

// medicalFailure is the 2d6 table of pp. 101-102.
func medicalFailure() MishapTable {
	return MishapTable{
		{
			Summary: "an ethics violation, and immediate removal",
			Effects: []Effect{
				relationshipAt(Enemy, 1, -130),
				chr("EDU", -2),
				pick("choose what the removal taught you",
					opt("Deception", skill("Deception", "Forgery", "Lie")),
					opt("Streetwise", skill("Streetwise"))),
				noAdmissionFor(2),
			},
		},
		{
			Summary: "technical and diagnostic demands beyond your preparation",
			Effects: []Effect{chr("EDU", -1)},
		},
		{
			Summary: "a serious mistake during supervised care",
			Effects: []Effect{chr("EDU", -1), chr("CHA", -1), skill("Medic", "First Aid")},
		},
		{
			Summary: "stress and exposure to suffering, and a withdrawal",
			Effects: []Effect{chr("EDU", -1), chr("END", 1)},
		},
		{
			Summary: "tuition or sponsorship that ended",
			Effects: []Effect{chr("EDU", -1), pickSkill("Admin", "Broker")},
		},
		{
			Summary: "a crisis at home that outweighed the training",
			Effects: []Effect{loseRelative(), ratingOfRole("parent", false, 50)},
		},
		{
			Summary: "an irreparable disagreement with the faculty",
			Effects: []Effect{relationshipAt(Rival, 1, -60)},
		},
		{
			Summary: "the realisation that you lack the temperament for practice",
			Effects: []Effect{pick("choose what you left with",
				opt("Persuade", skill("Persuade")),
				opt("END", chr("END", 1)))},
		},
		{
			Summary: "licensing clearance denied, and training made impossible",
			Effects: []Effect{chr("EDU", -1), pickSkill("Admin", "Streetwise")},
		},
		{
			Summary: "a shift into laboratory work that ends the clinical training",
			Effects: []Effect{unimplemented(
				"automatic admission to Graduate School, or -- with a doctorate already -- " +
					"automatic enlistment in the Scientist career")},
		},
		{
			Summary: "a programme permanently shut down",
			Effects: []Effect{
				throwModifier("admission to a different medical school in two or more terms", 2),
				pickSkill("Admin", "Broker"),
				relationshipAt(Contact, 1, 50),
			},
		},
	}
}
