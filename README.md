# ai-cli

This will be a small CLI app that will allow users to run local LLMs, or provide an API to interact with OpenAI, Claude, Gemini, etc. 

The goal is to replace Claude Code and similar tools with something that works locally. It may not be better than them, or do exactly the same thing, but I find that LLMs too often make changes to files. My idea is to keep it purely text based so a model can offer suggestions on work that's been completed, or a specific direction to take something, but to not actually implement the code. 

The release of Gemini 3 inspired this. I immediately found that no matter what I told it, it would still attempt to change my files, going as far as to delete them. This is annoying, and not how I like to do things.