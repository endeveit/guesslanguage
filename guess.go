// Package guesslanguage provides a simple way to detect language of unicode string.
//
// Source code and project home:
// https://github.com/endeveit/guesslanguage
//
package guesslanguage

import (
	"errors"
	"github.com/endeveit/guesslanguage/models"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	regexWords      *regexp.Regexp
	maxLength       int      = 4096
	maxDistance     int      = maxLength * 300
	minLength       int      = 20
	maxGrams        int      = 300
	unknownLanguage string   = "UNKNOWN"
	codesBasicLatin []string = []string{
		"ceb",
		"en",
		"eu",
		"ha",
		"haw",
		"id",
		"la",
		"nr",
		"nso",
		"so",
		"ss",
		"st",
		"sw",
		"tlh",
		"tn",
		"ts",
		"xh",
		"zu"}
	codesExtendedLatin []string = []string{
		"af",
		"az",
		"ca",
		"cs",
		"cy",
		"da",
		"de",
		"eo",
		"es",
		"et",
		"fi",
		"fr",
		"hr",
		"hu",
		"is",
		"it",
		"lt",
		"lv",
		"nb",
		"nl",
		"pl",
		"pt",
		"ro",
		"sk",
		"sl",
		"sq",
		"sv",
		"tl",
		"tr",
		"ve",
		"vi"}
	codesAllLatin   []string = append(codesBasicLatin, codesExtendedLatin...)
	codesCyrillic   []string = []string{"bg", "kk", "ky", "mk", "mn", "ru", "sr", "uk", "uz"}
	codesArabic     []string = []string{"ar", "fa", "ps", "ur"}
	codesDevanagari []string = []string{"hi", "ne"}
	codesPt         []string = []string{"pt_BR", "pt_PT"}
	// NOTE mn appears twice, once for mongolian script and once for codesCyrillic
	singletons map[string]string = map[string]string{
		"Armenian":  "hy",
		"Hebrew":    "he",
		"Bengali":   "bn",
		"Gurmukhi":  "pa",
		"Greek":     "el",
		"Gujarati":  "gu",
		"Oriya":     "or",
		"Tamil":     "ta",
		"Telugu":    "te",
		"Kannada":   "kn",
		"Malayalam": "ml",
		"Sinhala":   "si",
		"Thai":      "th",
		"Lao":       "lo",
		"Tibetan":   "bo",
		"Burmese":   "my",
		"Georgian":  "ka",
		"Mongolian": "mn-Mong",
		"Khmer":     "km"}
	nameMap map[string]string = map[string]string{
		"ab":    "Abkhazian",
		"af":    "Afrikaans",
		"ar":    "Arabic",
		"az":    "Azeri",
		"be":    "Byelorussian",
		"bg":    "Bulgarian",
		"bn":    "Bengali",
		"bo":    "Tibetan",
		"br":    "Breton",
		"ca":    "Catalan",
		"ceb":   "Cebuano",
		"cs":    "Czech",
		"cy":    "Welsh",
		"da":    "Danish",
		"de":    "German",
		"el":    "Greek",
		"en":    "English",
		"eo":    "Esperanto",
		"es":    "Spanish",
		"et":    "Estonian",
		"eu":    "Basque",
		"fa":    "Farsi",
		"fi":    "Finnish",
		"fo":    "Faroese",
		"fr":    "French",
		"fy":    "Frisian",
		"gd":    "Scots Gaelic",
		"gl":    "Galician",
		"gu":    "Gujarati",
		"ha":    "Hausa",
		"haw":   "Hawaiian",
		"he":    "Hebrew",
		"hi":    "Hindi",
		"hr":    "Croatian",
		"hu":    "Hungarian",
		"hy":    "Armenian",
		"id":    "Indonesian",
		"is":    "Icelandic",
		"it":    "Italian",
		"ja":    "Japanese",
		"ka":    "Georgian",
		"kk":    "Kazakh",
		"km":    "Cambodian",
		"ko":    "Korean",
		"ku":    "Kurdish",
		"ky":    "Kyrgyz",
		"la":    "Latin",
		"lt":    "Lithuanian",
		"lv":    "Latvian",
		"mg":    "Malagasy",
		"mk":    "Macedonian",
		"ml":    "Malayalam",
		"mn":    "Mongolian",
		"mr":    "Marathi",
		"ms":    "Malay",
		"nd":    "Ndebele",
		"ne":    "Nepali",
		"nl":    "Dutch",
		"nn":    "Nynorsk",
		"no":    "Norwegian",
		"nso":   "Sepedi",
		"pa":    "Punjabi",
		"pl":    "Polish",
		"ps":    "Pashto",
		"pt":    "Portuguese",
		"pt_PT": "Portuguese (Portugal)",
		"pt_BR": "Portuguese (Brazil)",
		"ro":    "Romanian",
		"ru":    "Russian",
		"sa":    "Sanskrit",
		"sh":    "Serbo-Croatian",
		"sk":    "Slovak",
		"sl":    "Slovene",
		"so":    "Somali",
		"sq":    "Albanian",
		"sr":    "Serbian",
		"sv":    "Swedish",
		"sw":    "Swahili",
		"ta":    "Tamil",
		"te":    "Telugu",
		"th":    "Thai",
		"tl":    "Tagalog",
		"tlh":   "Klingon",
		"tn":    "Setswana",
		"tr":    "Turkish",
		"ts":    "Tsonga",
		"tw":    "Twi",
		"uk":    "Ukrainian",
		"ur":    "Urdu",
		"uz":    "Uzbek",
		"ve":    "Venda",
		"vi":    "Vietnamese",
		"xh":    "Xhosa",
		"zh":    "Chinese",
		"zh_TW": "Traditional Chinese (Taiwan)",
		"zu":    "Zulu"}
	ianaMap map[string]int = map[string]int{
		"ab":    12026,
		"af":    40,
		"ar":    26020,
		"az":    26030,
		"be":    11890,
		"bg":    26050,
		"bn":    26040,
		"bo":    26601,
		"br":    1361,
		"ca":    3,
		"ceb":   26060,
		"cs":    26080,
		"cy":    26560,
		"da":    26090,
		"de":    26160,
		"el":    26165,
		"en":    26110,
		"eo":    11933,
		"es":    26460,
		"et":    26120,
		"eu":    1232,
		"fa":    26130,
		"fi":    26140,
		"fo":    11817,
		"fr":    26150,
		"fy":    1353,
		"gd":    65555,
		"gl":    1252,
		"gu":    26599,
		"ha":    26170,
		"haw":   26180,
		"he":    26592,
		"hi":    26190,
		"hr":    26070,
		"hu":    26200,
		"hy":    26597,
		"id":    26220,
		"is":    26210,
		"it":    26230,
		"ja":    26235,
		"ka":    26600,
		"kk":    26240,
		"km":    1222,
		"ko":    26255,
		"ku":    11815,
		"ky":    26260,
		"la":    26280,
		"lt":    26300,
		"lv":    26290,
		"mg":    1362,
		"mk":    26310,
		"ml":    26598,
		"mn":    26320,
		"mr":    1201,
		"ms":    1147,
		"ne":    26330,
		"nl":    26100,
		"nn":    172,
		"no":    26340,
		"pa":    65550,
		"pl":    26380,
		"ps":    26350,
		"pt":    26390,
		"ro":    26400,
		"ru":    26410,
		"sa":    1500,
		"sh":    1399,
		"sk":    26430,
		"sl":    26440,
		"so":    26450,
		"sq":    26010,
		"sr":    26420,
		"sv":    26480,
		"sw":    26470,
		"ta":    26595,
		"te":    26596,
		"th":    26594,
		"tl":    26490,
		"tlh":   26250,
		"tn":    65578,
		"tr":    26500,
		"tw":    1499,
		"uk":    26520,
		"ur":    26530,
		"uz":    26540,
		"vi":    26550,
		"zh":    26065,
		"zh_TW": 22}
)

func init() {
	regexWords, _ = regexp.Compile(`(?:[^\d\s_-]|['])+`)
}

// Return the ISO 639-1 language code.
func Guess(text string) (result string, err error) {
	if !utf8.ValidString(text) {
		return result, errors.New("Input string contains invalid UTF-8-encoded runes")
	}

	runed := []rune(text)

	if len(runed) > maxLength {
		runed = runed[:maxLength]
	}

	words := regexWords.FindAllString(strings.Replace(string(runed), "’", "'", -1), -1)

	return guessLanguage(words, getRuns(words)), nil
}

// Return the language IANA ID.
func GuessId(text string) int {
	code, err := Guess(text)

	if err != nil || code == unknownLanguage {
		return 0
	}

	return ianaMap[code]
}

// Return the language name (in English).
func GuessName(text string) string {
	code, err := Guess(text)

	if err != nil || code == unknownLanguage {
		return unknownLanguage
	}

	return nameMap[code]
}

// Identify the language.
func guessLanguage(words []string, scripts []string) string {
	if keyExists("Hangul Syllables", scripts) ||
		keyExists("Hangul Jamo", scripts) ||
		keyExists("Hangul Compatibility Jamo", scripts) ||
		keyExists("Hangul", scripts) {
		return "ko"
	}

	if keyExists("Greek and Coptic", scripts) {
		return "el"
	}

	if keyExists("Kana", scripts) {
		return "ja"
	}

	if keyExists("CJK Unified Ideographs", scripts) ||
		keyExists("Bopomofo", scripts) ||
		keyExists("Bopomofo Extended", scripts) ||
		keyExists("KangXi Radicals", scripts) {
		return "zh"
	}

	if keyExists("Cyrillic", scripts) {
		return getFromModel(words, codesCyrillic)
	}

	if keyExists("Arabic", scripts) ||
		keyExists("Arabic Presentation Forms-A", scripts) ||
		keyExists("Arabic Presentation Forms-B", scripts) {
		return getFromModel(words, codesArabic)
	}

	if keyExists("Devanagari", scripts) {
		return getFromModel(words, codesDevanagari)
	}

	// Try languages with unique scripts
	for blockName, langName := range singletons {
		if keyExists(blockName, scripts) {
			return langName
		}
	}

	if keyExists("Extended Latin", scripts) {
		latinLang := getFromModel(words, codesExtendedLatin)
		if latinLang == "pt" {
			return getFromModel(words, codesPt)
		} else {
			return latinLang
		}
	}

	if keyExists("Basic Latin", scripts) {
		return getFromModel(words, codesAllLatin)
	}

	return unknownLanguage
}

// Count the number of characters in each character block
func getRuns(words []string) (relevantRuns []string) {
	var (
		runTypes     map[string]int = make(map[string]int)
		nbTotalChars int            = 0
		charBlock    string
		percentage   int
	)

	for _, word := range words {
		for _, char := range word {
			charBlock = getBlock(char)

			if _, ok := runTypes[charBlock]; ok {
				runTypes[charBlock]++
			} else {
				runTypes[charBlock] = 1
			}

			nbTotalChars++
		}
	}

	// return run types that used for 40% or more of the string
	// return Basic Latin if found more than 15%
	for key, value := range runTypes {
		percentage = value * 100

		if percentage >= 40 || percentage >= 15 && key == "Basic Latin" {
			relevantRuns = append(relevantRuns, key)
		}
	}

	return relevantRuns
}

// Check if key exists
func keyExists(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check words against known models
func getFromModel(words []string, languages []string) (result string) {
	sample := strings.Join(words, " ")

	if len([]rune(sample)) < minLength {
		return unknownLanguage
	}

	var (
		scores  map[string]int = make(map[string]int, len(languages))
		model   []string       = models.GetOrderedModel(sample)
		minimal int            = maxDistance
	)

	for _, lang := range languages {
		scores[lang] = getDistance(model, models.GetModels()[strings.ToLower(lang)])
	}

	for lang, score := range scores {
		if score < minimal {
			minimal = score
			result = lang
		}
	}

	return result
}

// Calculate the distance to the known model.
func getDistance(model []string, known_model map[string]int) int {
	var (
		data []string
		dist int
		sub  int
	)

	if len(model) > maxGrams {
		data = model[:maxGrams]
	} else {
		data = model
	}

	for i, value := range data {
		if _, ok := known_model[value]; ok {
			sub = i - known_model[value]
			if sub < 0 {
				sub = -sub
			}

			dist += sub
		} else {
			dist += maxGrams
		}
	}

	return dist
}
