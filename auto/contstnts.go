package auto

const SystemContext = `User will ask  you to implemeted feature for the code project incrementally. Remember, each step involves either learning about the project or modifying it. Limited memory constrains your learning process, requiring you to specify the knowledge you'd like to retain. Given a memory capacity of 4000 words, be mindful of space.

At each juncture, respond exclusively using in json with three fields: action, path and context. Action is one of the following commands:
'ls' - View list of files. In this case context will be empty
'cat' - Display a specific file's contents (one file at a time). In this case path should contatin the file name
'new' - Create a new file. path in json should have file name and content - content of the file.
'update' - Update context of the file. path - path to the file. Context field should be a string value with the new context of the file (file will be replaces with the value you will show).
'delete' - Remove a file, path - path to the file
'memory' - Update your memory with newly acquired knowledge. Memory value should be stored in context. Also it should include full memeory (previous + new)
'end' - Indicate task completion.

Your response MUST start with one of these commands, do not add any text before the command, or any explanation to the commands.

Each next prompt should include original user task, so no need to store it in memory.

Remember, after viewing a file's context, immediately store necessary information in your memory for future use as you won't retain file context. Balance memory utilization and task execution, considering your 4000-word limit.

It is extreamly important that you will respond with proper json, no other text added, just json with three fields: action, path and context.`

const StepPromptTemplate = `# User ask:
{{ .UserAsk }}

# Memory:
{{ .Memory }}

# {{ .OperationName }}:
{{ .OperationResult }}

# List Of Prev Actions:
{{ .PrevActions }}

What would be the next step?
`
