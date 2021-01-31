package l10n

import (
	"github.com/loctools/go-l10n/loc"
	"github.com/loctools/go-l10n/locjson"
)

// https://atlas.limsi.fr/?tab=Hexagone
func LoadTranslations() *loc.Pool {
	lp := loc.NewPool("fr_FR")
	lp.Resources["ba_BAS"] = locjson.Load("l10n/ba_BAS.json")
	lp.Resources["br_BRE"] = locjson.Load("l10n/br_BRE.json")
	lp.Resources["fr_FR"] = locjson.Load("l10n/fr_FR.json")
	lp.Resources["gr_ALS"] = locjson.Load("l10n/gr_ALS.json")
	lp.Resources["gr_FLA"] = locjson.Load("l10n/gr_FLA.json")
	lp.Resources["gr_FRA"] = locjson.Load("l10n/gr_FRA.json")
	lp.Resources["oc_CRO"] = locjson.Load("l10n/oc_CRO.json")
	lp.Resources["oc_GAS"] = locjson.Load("l10n/oc_GAS.json")
	lp.Resources["oc_LAN"] = locjson.Load("l10n/oc_LAN.json")
	lp.Resources["oc_NOC"] = locjson.Load("l10n/oc_NOC.json")
	lp.Resources["oc_PRO"] = locjson.Load("l10n/oc_PRO.json")
	lp.Resources["ro_ANG"] = locjson.Load("l10n/ro_ANG.json")
	lp.Resources["ro_BOU"] = locjson.Load("l10n/ro_BOU.json")
	lp.Resources["ro_CAT"] = locjson.Load("l10n/ro_CAT.json")
	lp.Resources["ro_CEN"] = locjson.Load("l10n/ro_CEN.json")
	lp.Resources["ro_CHA"] = locjson.Load("l10n/ro_CHA.json")
	lp.Resources["ro_COR"] = locjson.Load("l10n/ro_COR.json")
	lp.Resources["ro_FRC"] = locjson.Load("l10n/ro_FRC.json")
	lp.Resources["ro_FRP"] = locjson.Load("l10n/ro_FRP.json")
	lp.Resources["ro_GAL"] = locjson.Load("l10n/ro_GAL.json")
	lp.Resources["ro_LIG"] = locjson.Load("l10n/ro_LIG.json")
	lp.Resources["ro_LOR"] = locjson.Load("l10n/ro_LOR.json")
	lp.Resources["ro_MAI"] = locjson.Load("l10n/ro_MAI.json")
	lp.Resources["ro_NOR"] = locjson.Load("l10n/ro_NOR.json")
	lp.Resources["ro_PIC"] = locjson.Load("l10n/ro_PIC.json")
	lp.Resources["ro_POI"] = locjson.Load("l10n/ro_POI.json")
	lp.Resources["ro_WAL"] = locjson.Load("l10n/ro_WAL.json")
	return lp
}
