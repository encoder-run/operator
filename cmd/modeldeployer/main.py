import torch
from typing import Dict, List
from kserve import Model, ModelServer
from transformers import AutoTokenizer
from sentence_transformers import SentenceTransformer
import argparse

device = "cuda"  # for GPU usage or "cpu" for CPU usage

class CustomModel(Model):
    def __init__(self, name: str, org_name: str, repo_name: str, max_sequence_length: int):
        super().__init__(name)
        self.name = name
        self.org_name = org_name
        self.repo_name = repo_name
        self.max_sequence_length = max_sequence_length
        self.device = torch.device(device if torch.cuda.is_available() else "cpu")
        print("Using device:", self.device)
        self.ready = False

    def load(self):
        self.tokenizer = AutoTokenizer.from_pretrained(f"{self.org_name}/{self.repo_name}", trust_remote_code=True)
        self.model = SentenceTransformer(f"{self.org_name}/{self.repo_name}", trust_remote_code=True)
        self.model.max_seq_length = self.max_sequence_length
        self.ready = True

    def predict(self, payload: Dict, headers: Dict) -> Dict:
        inputs = payload["instances"]
        return self.process_multiple_files(inputs)

    def process_multiple_files(self, file_inputs: List[Dict]) -> Dict:
        all_chunks = []
        file_chunk_map = {}

        # Accumulate chunks from all files
        for file_input in file_inputs:
            file_path = file_input["file_path"]
            source_code = file_input["code"]
            hash = file_input["file_hash"]
            chunks = self.chunk_code(source_code, hash, 500)
            all_chunks.extend(chunks)
            file_chunk_map[file_path] = (len(all_chunks) - len(chunks), len(all_chunks))

        # Encode all chunks at once
        codes = [chunk["code"] for chunk in all_chunks]
        code_embs = self.model.encode(codes, convert_to_tensor=True)

        # Distribute embeddings back to respective files
        results = {}
        for file_path, (start, end) in file_chunk_map.items():
            file_chunks = all_chunks[start:end]
            for chunk, code_emb in zip(file_chunks, code_embs[start:end]):
                chunk["embedding"] = code_emb.tolist()
            results[file_path] = {"embeddings": file_chunks}

        return {"results": results}

    def chunk_code(self, code, hash, max_token_length):
        # Encode the entire code at once, ignoring special tokens
        tokens_data = self.tokenizer.encode_plus(code, add_special_tokens=False, return_offsets_mapping=True)
        tokens = tokens_data['input_ids']
        offsets = tokens_data['offset_mapping']

        chunks = []
        current_chunk_start_index = 0
        total_tokens = len(tokens)
        current_token_index = 0

        # Distribute tokens to each chunk as evenly as possible with some padding
        tokens_per_chunk = max_token_length
        while current_token_index < total_tokens:
            next_token_index = min(current_token_index + tokens_per_chunk, total_tokens)

            # Adjust the end index to not cut in the middle of a word
            while next_token_index < total_tokens and offsets[next_token_index - 1][1] != offsets[next_token_index][0]:
                next_token_index += 1

            # Ensure not to exceed total tokens while correcting
            next_token_index = min(next_token_index, total_tokens)

            chunk_tokens = tokens[current_chunk_start_index:next_token_index]
            start_index = offsets[current_chunk_start_index][0]
            end_index = offsets[next_token_index - 1][1] if next_token_index > 0 else 0

            # Manually reconstruct the text to ensure spaces are correctly included
            chunk_code = code[start_index:end_index]

            # Append the current chunk to the chunks list
            chunks.append({
                "chunk_id": len(chunks),
                "code": chunk_code,
                "file_hash": hash,
                "start_index": start_index,
                "end_index": end_index
            })

            # Update for the next chunk
            current_token_index = next_token_index
            current_chunk_start_index = next_token_index

        return chunks






if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Initialize the CustomModel with command line parameters.")
    parser.add_argument("--org_name", type=str, required=True, help="Organization name")
    parser.add_argument("--repo_name", type=str, required=True, help="Repository name")
    parser.add_argument("--max_sequence_length", type=int, default=512,
                        help="Maximum input sequence length of the model")

    args = parser.parse_args()  # Parse the arguments from the command line

    # Create an instance of the model using the parsed arguments
    model = CustomModel(name="custom-model", org_name=args.org_name, repo_name=args.repo_name, max_sequence_length=args.max_sequence_length)
    model.load()
    ModelServer().start([model])
