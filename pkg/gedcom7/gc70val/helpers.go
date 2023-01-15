package gc70val

// IsValidTag returns true if the tag is a valid standard or extended tag
func (g *Specs) IsValidTag(tag string) bool {
	if _, ok := g.Tags[tag]; ok {
		return true
	}
	if _, ok := g.Tags["record-"+tag]; ok {
		return true
	}
	if tag[0] == '_' {
		return true
	}
	if v, ok := g.Tags["FAM-"+tag]; ok {
		if v.Tag == tag {
			return true
		}
	}
	if v, ok := g.Tags["INDI-"+tag]; ok {
		if v.Tag == tag {
			return true
		}
	}

	// weird special cases
	if tag == "STAT" {
		return true
	}

	if _, ok := g.depTags[tag]; ok {
		return true
	}

	return false
}
