package gc70val

/*
https://gedcom.io/specifications/FamilySearchGEDCOMv7.pdf
p. 27

DateValue = date / DatePeriod / dateRange / dateApprox / ""
DateExact = day D month D year ; in Gregorian calendar
DatePeriod = %s"FROM" D date [D %s"TO" D date]
 / %s"TO" D date
 / ""
date = [calendar D] [[day D] month D] year [D epoch]
dateRange = %s"BET" D date D %s"AND" D date
 / %s"AFT" D date
 / %s"BEF" D date
dateApprox = (%s"ABT" / %s"CAL" / %s"EST") D date
dateRestrict = %s"FROM" / %s"TO" / %s"BET" / %s"AND" / %s"BEF"
 / %s"AFT" / %s"ABT" / %s"CAL" / %s"EST" / %s"BCE"
calendar = %s"GREGORIAN" / %s"JULIAN" / %s"FRENCH_R" / %s"HEBREW"
 / extTag
day = Integer
year = Integer
month = stdTag / extTag ; constrained by calendar
epoch = %s"BCE" / extTag ; constrained by calendar
*/

type calDef struct {
	Lang          string   `yaml:"lang"`
	Type          string   `yaml:"type"`
	URI           string   `yaml:"uri"`
	Cal           string   `yaml:"standard tag"`
	Specification []string `yaml:"specification"`
	Months        []string `yaml:"months"`
}

func loadCal(in []byte) (calDef, error) {
	cm := calDef{}
	if err := deserializeYAML(in, &cm); err != nil {
		return cm, err
	}

	cm.Cal = extractFullTag(cm.URI)
	for i, v := range cm.Months {
		cm.Months[i] = extractFullTag(v)
	}

	return cm, nil
}
