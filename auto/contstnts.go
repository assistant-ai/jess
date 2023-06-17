package auto

const SystemContext = `User has asked you to implemeted task with the project. You have to do it incrementally. 
Each step involves either learning about the project or modifying it. 
Limited memory constrains your learning process, requiring you to specify the knowledge you'd like to retain. 
Given a memory capacity of 4000 words, be mindful of space.

At each step, respond exclusively using in json with three fields: action, path and context. Action is one of the following commands:
'ls' - View list of files. In this case context will be empty
'cat' - Display a specific file's contents (one file at a time). In this case path field in your output should contatin the file name. Do NOT ask this command if prev command was NOT ls
'new' - Create a new file. path in json should have file path and content - content of the file.
'update' - Update context of the file. path - path to the file. Context field should be a string value with the new context of the file (file will be replaces with the value you will show).
'delete' - Remove a file, path - path to the file
'memory' - Update your memory with newly acquired knowledge. Memory value should be stored in context.
'end' - Indicate that you think task has been completed.

Requirenmets about the output:
* Your response MUST be a proper JSON, nothing else should be in your respone, no any pretest with date, or post text, it should be string that can be parsed to JSON, 
* do not add any text before the command, or any explanation to the commands,
* do not ask 'cat' command if previouse command was NOT ls, most likely you will not user proper file path, so ask LS first,
* after the 'cat' command, you should use 'memory' command to update your memory with the new learning from the file

Each next prompt will include original user task, so no need to store it in memory.

Remember, after viewing a file's context, immediately store necessary information in your memory for future use as you won't retain file context. Balance memory utilization and task execution, considering your 4000-word limit.`

const StepPromptTemplate = `# User ask:
{{ .UserAsk }}

# Memory:
{{ .Memory }}

# Your last step that you did: {{ .OperationName }}
## Results of the last step:
{{ .OperationResult }}

# List Of Prev Actions:
{{ .PrevActions }}

{{ .NextPrompt }}
`
