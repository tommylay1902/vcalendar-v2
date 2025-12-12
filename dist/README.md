---
language: en
license: apache-2.0
library_name: sentence-transformers
tags:
- sentence-transformers
- feature-extraction
- sentence-similarity
- transformers
- llama-cpp
- gguf-my-repo
datasets:
- s2orc
- flax-sentence-embeddings/stackexchange_xml
- ms_marco
- gooaq
- yahoo_answers_topics
- code_search_net
- search_qa
- eli5
- snli
- multi_nli
- wikihow
- natural_questions
- trivia_qa
- embedding-data/sentence-compression
- embedding-data/flickr30k-captions
- embedding-data/altlex
- embedding-data/simple-wiki
- embedding-data/QQP
- embedding-data/SPECTER
- embedding-data/PAQ_pairs
- embedding-data/WikiAnswers
pipeline_tag: sentence-similarity
base_model: sentence-transformers/all-MiniLM-L6-v2
---

# fartboner/all-MiniLM-L6-v2-Q8_0-GGUF
This model was converted to GGUF format from [`sentence-transformers/all-MiniLM-L6-v2`](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) using llama.cpp via the ggml.ai's [GGUF-my-repo](https://huggingface.co/spaces/ggml-org/gguf-my-repo) space.
Refer to the [original model card](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) for more details on the model.

## Use with llama.cpp
Install llama.cpp through brew (works on Mac and Linux)

```bash
brew install llama.cpp

```
Invoke the llama.cpp server or the CLI.

### CLI:
```bash
llama-cli --hf-repo fartboner/all-MiniLM-L6-v2-Q8_0-GGUF --hf-file all-minilm-l6-v2-q8_0.gguf -p "The meaning to life and the universe is"
```

### Server:
```bash
llama-server --hf-repo fartboner/all-MiniLM-L6-v2-Q8_0-GGUF --hf-file all-minilm-l6-v2-q8_0.gguf -c 2048
```

Note: You can also use this checkpoint directly through the [usage steps](https://github.com/ggerganov/llama.cpp?tab=readme-ov-file#usage) listed in the Llama.cpp repo as well.

Step 1: Clone llama.cpp from GitHub.
```
git clone https://github.com/ggerganov/llama.cpp
```

Step 2: Move into the llama.cpp folder and build it with `LLAMA_CURL=1` flag along with other hardware-specific flags (for ex: LLAMA_CUDA=1 for Nvidia GPUs on Linux).
```
cd llama.cpp && LLAMA_CURL=1 make
```

Step 3: Run inference through the main binary.
```
./llama-cli --hf-repo fartboner/all-MiniLM-L6-v2-Q8_0-GGUF --hf-file all-minilm-l6-v2-q8_0.gguf -p "The meaning to life and the universe is"
```
or 
```
./llama-server --hf-repo fartboner/all-MiniLM-L6-v2-Q8_0-GGUF --hf-file all-minilm-l6-v2-q8_0.gguf -c 2048
```
