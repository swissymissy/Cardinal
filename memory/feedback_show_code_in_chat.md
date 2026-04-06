---
name: Show code in chat before editing files
description: User wants to see code shown in the chat first and edit files themselves, not have Claude edit files directly
type: feedback
---

Always show code in the chat message first. Do not use Write or Edit tools to modify files unless the user explicitly says "go ahead" or "make the change" or similar confirmation. The user wants to review the code first and make edits themselves.

**Why:** User prefers to review before applying, and has rejected multiple direct file edits.

**How to apply:** When asked to write/show code, output it as a code block in the chat. Only use Edit/Write tools after explicit user approval.
