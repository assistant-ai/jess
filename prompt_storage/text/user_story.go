package text

// go:embed user_story.md
var staticPrompt string

var UserStoryPrompt = "You need to use provided markdown instruction for generating user story:" + staticPrompt
