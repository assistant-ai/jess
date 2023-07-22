package utils

import (
	"fmt"
	"github.com/agnivade/levenshtein"
	"strings"
)

var words = []string{
	"Finish",
	"Thanks",
	"Bye",
	"Goodbye",
	"Farewell",
	"Take care",
	"Later",
	"See you",
	"See ya",
	"Adios",
	"Bye-bye",
	"Catch you later",
	"Until next time",
	"Signing off",
	"Best wishes",
	"So long",
	"Ta-ta",
	"Cheerio",
	"Till we meet again",
	"Toodles",
	"Have a great day",
	"Have a good one",
	"Peace out",
	"Ciao",
	"Hasta la vista",
	"Tata",
	"Adieu",
	"Bye for now",
	"Goodnight",
	"Keep in touch",
	"Until we meet again",
	"May the Force be with you",
	"Sayonara",
	"Later gator",
	"Catch you on the flip side",
	"Be well",
	"Take it easy",
	"Keep on truckin'",
	"Be safe",
	"See you later, alligator",
	"Time to go",
	"Talk to you later",
	"Bon voyage",
	"Peace",
	"Gotta run",
	"Smell you later",
	"See you soon",
	"Hasta mañana",
	"Bon appétit",
	"Godspeed",
	"Take care of yourself",
	"Happy trails",
	"Bye bye for now",
	"Until we speak again",
	"Bon soir",
	"Until later",
	"Over and out",
	"Have a nice one",
	"Have a blessed day",
	"Stay in touch",
	"Good to see you",
	"Stay safe",
	"Later on",
	"G'bye",
	"Have a great one",
	"Have a good day",
	"Talk soon",
	"Take care now",
	"Peace and love",
	"Fond farewell",
	"Until next time, take care",
	"Good luck",
	"So long, farewell",
	"See you in a bit",
	"Catch you on the rebound",
	"Be good",
	"Adieu, my friend",
	"Have a wonderful day",
	"Enjoy your day",
	"God be with you",
	"Until next time, goodbye",
	"Till then",
	"See you around",
	"Goodbye, my friend",
	"Time to say goodbye",
	"Until we meet again, farewell",
	"Stay well",
	"May your day be filled with joy",
	"Take care, my friend",
	"G2G",
	"TTYL",
	"Peace and blessings",
	"Goodbye, take care",
	"Don't be a stranger",
	"See you later, buddy",
	"Until tomorrow",
	"Fare thee well",
	"exit",
}

func isSimilar(s1, s2 string, threshold int) bool {
	distance := levenshtein.ComputeDistance(s1, s2)
	return distance <= threshold
}

func IfAnswerInFinishingArray(finishingAnswer string) bool {
	for _, word := range words {
		word = strings.ToLower(word)
		finishingAnswer = strings.ToLower(finishingAnswer)
		if finishingAnswer == word {
			return true
		}
		if isSimilar(finishingAnswer, word, 3) {
			return true
		}
	}
	return false
}

func AskForConfirmation() bool {
	var response string

	PrintlnRed("Do you want to exit? ([Y]es/[N]o): ")
	PrintlnRed("Press [Enter] to continue dialog: ")
	fmt.Scanln(&response)

	response = strings.TrimSpace(strings.ToLower(response))

	switch response {
	case "yes", "y":
		return true
	case "no", "n", "":
		return false
	default:
		fmt.Println("Invalid input. Please respond with yes or no.")
		return AskForConfirmation()
	}
}
