package guesslanguage

import (
	"./models"
	"strings"
	"testing"
)

func Test_getRuns(t *testing.T) {
	var (
		words []string
		runs  []string
	)

	words = strings.Split("This is a test of the language checker", " ")
	runs = getRuns(words)

	if !keyExists("Basic Latin", runs) {
		t.Errorf("Runs must contain 'Basic Latin'")
	}

	words = strings.Split("abcdééé", " ")
	runs = getRuns(words)
	if len(runs) != 2 || !keyExists("Basic Latin", runs) || !keyExists("Extended Latin", runs) {
		t.Errorf("Runs must contain both 'Basic Latin' and 'Extended Latin'")
	}

	s := "Сайлау нәтижесінде дауыстардың басым бөлігін ел премьер " +
		"министрі Виктор Янукович пен оның қарсыласы, оппозиция " +
		"жетекшісі Виктор Ющенко алды."
	words = strings.Split(s, " ")
	runs = getRuns(words)
	if !keyExists("Cyrillic", runs) {
		t.Errorf("Runs musts contain 'Cyrillic'")
	}
}

func Test_GetOrderedModel(t *testing.T) {
	var om []string

	om = models.GetOrderedModel("abc")
	if len(om) != 1 || !keyExists("abc", om) {
		t.Errorf("Model must be [abc]")
	}

	om = models.GetOrderedModel("abca")
	if len(om) != 2 || !keyExists("abc", om) || !keyExists("bca", om) {
		t.Errorf("Model must be [abc bca]")
	}

	om = models.GetOrderedModel("abcabdcab")
	if len(om) != 6 ||
		!keyExists("cab", om) ||
		!keyExists("abc", om) ||
		!keyExists("abd", om) ||
		!keyExists("bca", om) ||
		!keyExists("bdc", om) ||
		!keyExists("dca", om) {
		t.Errorf("Model must be [cab abc abd bca bdc dca]")
	}
}

func Test_GuessLanguage(t *testing.T) {
	var (
		err     error
		guessed string
		texts   map[string]string = map[string]string{
			unknownLanguage: "",
			"ar":            "ملايين الناخبين الأمريكيين يدلون بأصواتهم وسط إقبال قياسي على انتخابات هي الأشد تنافسا منذ عقود",
			"az":            "Daxil olan xəbərlərdə deyilir ki, 6 nəfər Bağdadın mərkəzində yerləşən Təhsil Nazirliyinin binası yaxınlığında baş vermiş partlayış zamanı həlak olub.",
			"bg":            "е готов да даде гаранции, че няма да прави ядрено оръжие, ако му се разреши мирна атомна програма",
			"cs":            "Francouzský ministr financí zmírnil výhrady vůči nízkým firemním daním v nových členských státech EU",
			"da":            "På denne side bringer vi billeder fra de mange forskellige forberedelser til arrangementet, efterhånden som vi får dem",
			"el":            "αναμένεται να σπάσουν παράδοση δεκαετιών και να συνεχίσουν την εκστρατεία τους ακόμη και τη μέρα των εκλογών",
			"en":            "This is a test of the language checker.",
			"eo":            "La akcento estas sur la antaŭlasta silabo.",
			"es":            "La respuesta de los acreedores a la oferta argentina para salir del default no ha sido muy positiv",
			"et":            "Ennetamaks reisil ebameeldivaid vahejuhtumeid vii end kurssi reisidokumentide ja viisade reeglitega ning muu praktilise informatsiooniga",
			"fi":            "on julkishallinnon verkkopalveluiden yhteinen osoite. Kansalaisten arkielämää helpottavaa tietoa on koottu eri aihealueisiin",
			"fr":            "Vérifions que le détecteur de langue fonctionne.",
			"hr":            "biće prilično izjednačena, sugerišu najnovije ankete. Oba kandidata tvrde da su sposobni da dobiju rat protiv terorizma",
			"hu":            "Hiába jön létre az önkéntes magyar haderő, hiába nem lesz többé bevonulás, változatlanul fennmarad a hadkötelezettség intézménye",
			"hy":            "հարաբերական",
			"ja":            "トヨタ自動車、フィリピンの植林活動で第三者認証取得　トヨタ自動車(株)（以下、トヨタ）は、2007年９月よりフィリピンのルソン島北部に位置するカガヤン州ペニャブランカ町",
			"kk":            "Сайлау нәтижесінде дауыстардың басым бөлігін ел премьер министрі Виктор Янукович пен оның қарсыласы, оппозиция жетекшісі Виктор Ющенко алды.",
			"ky":            "көрбөгөндөй элдик толкундоо болуп, Кокон шаарынын көчөлөрүндө бир нече миң киши нааразылык билдирди.",
			"mk":            "на јавното мислење покажуваат дека трката е толку тесна, што се очекува двајцата соперници да ја прекршат традицијата и да се појават и на самиот изборен ден.",
			"nb":            "Nominasjonskomiteen i Akershus KrF har skviset ut Einar Holstad fra stortingslisten. Ytre Enebakk-mannen har plass p Stortinget s lenge Valgerd Svarstad Haugland sitter i",
			"nl":            "Die kritiek was volgens hem bitter hard nodig, omdat Nederland binnen een paar jaar in een soort Belfast zou dreigen te veranderen",
			"pl":            "Sprawdźmy, czy odgadywacz języków pracuje",
			"pt_PT":         "Portugal é um país soberano unitário localizado no Sudoeste da Europa.",
			"ro":            "în acest sens aparţinînd Adunării Generale a organizaţiei, în ciuda faptului că mai multe dintre solicitările organizaţiei privind organizarea scrutinului nu au fost soluţionate",
			"ru":            "авай проверить  узнает ли наш угадатель русски язык",
			"sq":            "kaluan ditën e fundit të fushatës në shtetet kryesore për të siguruar sa më shumë votues.",
			"sv":            "Vi säger att Frälsningen är en gåva till alla, fritt och för intet.  Men som vi nämnt så finns det två villkor som måste",
			"tr":            "yakın tarihin en çekişmeli başkanlık seçiminde oy verme işlemi sürerken, katılımda rekor bekleniyor.",
			"uk":            "Американське суспільство, поділене суперечностями, збирається взяти активну участь у голосуванні",
			"uz":            "милиция ва уч солиқ идораси ходимлари яраланган. Шаҳарда хавфсизлик чоралари кучайтирилган.",
			"vi":            "Hai vấn đề khó chịu với màn hình thường gặp nhất khi bạn dùng laptop là vết trầy xước và điểm chết. Sau đây là vài cách xử lý chúng.",
			"zh":            "美国各州选民今天开始正式投票。据信，"}
		french string = "Vérifions que le détecteur de langue fonctionne."
	)

	for name, text := range texts {
		guessed, err = Guess(text)
		if err != nil {
			t.Errorf("Got error '%v' for language '%s'", err, name)
		}

		if guessed != name {
			t.Errorf("Expected '%s' got '%s'", name, guessed)
		}
	}

	if GuessId(french) != ianaMap["fr"] {
		t.Error("Language id must be %d", ianaMap["fr"])
	}

	if GuessName(french) != nameMap["fr"] {
		t.Error("Language name must be %s", nameMap["fr"])
	}
}
