package career

import "slices"

// Tag is a class of career. Twelve results modify or direct an enlistment
// by class rather than by name -- "-4 DM to enlist in any government
// related career", "-6 DM to any enlistment roll for a career involving
// violence", "enlist automatically in a business, military, corporate or
// colonist career" -- and the book defines none of the words it uses.
//
// The assignment is therefore a reading, one per career, argued from the
// career descriptions of pp. 132-145. ERRATA E-37 carries the whole table
// and the reason for each. Nothing else in the engine depends on it: a tag
// is read only when a result asks about one.
type Tag string

// The classes the results name. There are no others, because a class this
// engine invented would modify a throw the book does not modify.
const (
	// TagMilitary is service in a military organization, which is what the
	// five careers carrying it describe themselves as.
	TagMilitary Tag = "military"

	// TagGovernment is service to a government, military or civil.
	TagGovernment Tag = "government"

	// TagAcademic is a career whose work is knowledge and its teaching.
	TagAcademic Tag = "academic"

	// TagCriminal is a career the book describes as operating outside lawful
	// authority, plus the one a conviction puts a character into.
	TagCriminal Tag = "criminal"

	// TagViolent is a career whose ordinary work includes violence, as
	// against one where violence is a hazard that may occur.
	TagViolent Tag = "violent"

	// TagBusiness, TagCorporate and TagColonist are the three the one result that
	// names them distinguishes: "enlist automatically in a business,
	// military, corporate or colonist career". Corporate is service to a
	// corporation; Business is commerce on one's own account.
	TagBusiness  Tag = "business"
	TagCorporate Tag = "corporate"
	TagColonist  Tag = "colonist"
)

// HasTag reports whether a career is of a class.
func (c Career) HasTag(tag Tag) bool {
	return slices.Contains(c.Tags, tag)
}
