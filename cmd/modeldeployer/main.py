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
        self.tokenizer = AutoTokenizer.from_pretrained(f"{self.org_name}/{self.repo_name}")
        self.model = SentenceTransformer(f"{self.org_name}/{self.repo_name}")
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
        lines = code.split("\n")
        chunks = []
        current_chunk = []
        current_line = 1
        current_token_count = 0
        chunk_start_line = 1
        chunk_start_column = 0
        for line in lines:
            line = line + '\n'  # Add newline character back in
            tokens_data = self.tokenizer.encode_plus(line, add_special_tokens=False, return_offsets_mapping=True)
            tokens = tokens_data['input_ids']
            offsets = tokens_data['offset_mapping']

            # If the line is empty or only contains whitespace, skip it
            if not tokens:
                current_line += 1
                continue

            # If the line is too long, split it into multiple sublines
            if len(tokens) > max_token_length:
                sublines = []
                for i in range(0, len(tokens), max_token_length):
                    subline_tokens = tokens[i:i + max_token_length]
                    subline_offsets = offsets[i:i + max_token_length]
                    sublines.append((subline_tokens, subline_offsets))
            else:
                sublines = [(tokens, offsets)]

            for subline_tokens, subline_offsets in sublines:
                if current_token_count + len(subline_tokens) <= max_token_length:
                    current_chunk.extend(subline_tokens)
                    current_token_count += len(subline_tokens)
                    chunk_end_column = subline_offsets[-1][1]  # End column for last token in chunk
                else:
                    # Finish current chunk
                    chunk_code = self.tokenizer.decode(current_chunk)
                    chunks.append({
                        "chunk_id": len(chunks),
                        "file_hash": hash,  # Add hash to chunk for tracking purposes
                        "code": chunk_code,
                        "start_line": chunk_start_line,
                        "end_line": current_line,
                        "start_column": chunk_start_column,
                        "end_column": chunk_end_column
                    })
                    # Start new chunk
                    current_chunk = subline_tokens
                    current_token_count = len(subline_tokens)
                    chunk_start_line = current_line
                    chunk_start_column = subline_offsets[0][0]  # Start column for first token in new chunk
                    chunk_end_column = subline_offsets[-1][1]  # End column for last token in new chunk

            current_line += 1

        # Finish last chunk
        chunk_code = self.tokenizer.decode(current_chunk)
        chunks.append({
            "chunk_id": len(chunks),
            "file_hash": hash,  # Add hash to chunk for tracking purposes
            "code": chunk_code,
            "start_line": chunk_start_line,
            "end_line": current_line - 1,
            "start_column": chunk_start_column,
            "end_column": chunk_end_column
        })

        return chunks

    def remove_prefix(self, input_string, prefix):
        if input_string.startswith(prefix):
            return input_string[len(prefix):]
        return input_string

    def remove_suffix(self, input_string, suffix):
        if input_string.endswith(suffix):
            return input_string[:-len(suffix)]
        return input_string


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
