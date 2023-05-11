package translit

import (
	"testing"
)

func TestTransliterate(t *testing.T) {

	type testCase struct {
		name   string
		source string
		result string
	}

	//goland:noinspection SpellCheckingInspection

	testCases := []testCase{
		{name: "None", source: "",
			result: ""},
		{name: "SingleRune", source: "А",
			result: "a"},
		{name: "Latin", source: "English Message That Consists of Letters and 1000 Digits",
			result: "english-message-that-consists-of-letters-and-1000-digits"},
		{name: "StartsWithSpaces", source: "    сообщение с пробелами в начале",
			result: "soobshchenie-s-probelami-v-nachale"},
		{name: "Russian", source: "Ууух, я даже и не знаю",
			result: "uuuh-ya-dazhe-i-ne-znayu"},
		{name: "Blazers", source: "Пиджаки",
			result: "pidzhaki"},
		{name: "CellPhone", source: "Сотовый телефон",
			result: "sotovyj-telefon"},
		{name: "Wait", source: "Погоди-ка",
			result: "pogodi-ka"},
		{name: "Header-1", source: "Путин присвоил звания генералов 26 сотрудникам силовых ведомств",
			result: "putin-prisvoil-zvaniya-generalov-26-sotrudnikam-silovyh-vedomstv"},
		{name: "Header-2", source: "Попова предупредила о сохранении высоких рисков заразиться COVID-19",
			result: "popova-predupredila-o-sohranenii-vysokih-riskov-zarazitsya-covid-19"},
		{name: "Header-3", source: "Сотрудники посольства США обратились за вакциной «Спутник V»",
			result: "sotrudniki-posolstva-ssha-obratilis-za-vakcinoj-sputnik-v"},
		{name: "Header-4", source: "Почему аккумуляторы смартфонов взрываются и как защитить себя ...",
			result: "pochemu-akkumulyatory-smartfonov-vzryvayutsya-i-kak-zashchitit-sebya"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x := ToLatin(tc.source)
			if x != tc.result {
				t.Errorf("incorrect result: expected %s, got %s", tc.result, x)
			}
			if !Transliterated(x) {
				t.Errorf("incorrect result: %s is detected to be not transliterated", x)
			}
		})
	}
}
