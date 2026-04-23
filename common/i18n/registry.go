package i18n

type Lang string

const (
	En Lang = "en"
	Fr Lang = "fr"
	Ar Lang = "ar"
)

const DefaultLang = En

var PrefixedLangs = []Lang{Fr, Ar}
var AllLangs = []Lang{En, Fr, Ar}

type Meta struct {
	Code      Lang
	HTMLLang  string
	Dir       string
	Native    string
	URLPrefix string
}

var metas = map[Lang]Meta{
	En: {Code: En, HTMLLang: "en", Dir: "ltr", Native: "English", URLPrefix: ""},
	Fr: {Code: Fr, HTMLLang: "fr", Dir: "ltr", Native: "Français", URLPrefix: "/fr"},
	Ar: {Code: Ar, HTMLLang: "ar", Dir: "rtl", Native: "العربية", URLPrefix: "/ar"},
}

func MetaFor(l Lang) Meta {
	if m, ok := metas[l]; ok {
		return m
	}
	return metas[DefaultLang]
}

func AllMetas() []Meta {
	out := make([]Meta, 0, len(AllLangs))
	for _, l := range AllLangs {
		out = append(out, metas[l])
	}
	return out
}
