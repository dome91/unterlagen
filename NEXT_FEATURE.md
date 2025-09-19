# Assistant
The next feature is the Assistant.

## Problem Scenario
The user has a lot of documents. Full-Text-Search only helps marginally due to missing understanding of context in the documents.
The user wants the possibility to "ask" questions about their documents and get direct answers plus links to check the relevant documents, instead searching through all of their documents.

## Solution
Implement an assistant that gets access to the documents and can answer questions about them.
The solution is split into two phases:
1. Implement a system that takes uploaded documents and performs the necessary pre-steps for Retrieval Augmented Generation: Chunk the document into meaningful pieces and generate embeddings for each chunk.
2. The assistant gets a question of the user. Generating the embedding of the question will lead to the chunks that are most likely to have the answer to the question. Collecting the chunks and the references to the documents, the assistant will formulate an answer and present it to the user.

The assistant will be its own page in the UI. It should start with a list of existing chats, in case the user wants to go back to an old chat or the user can create a new chat.
Use the already existing code where possible.
