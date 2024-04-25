from sentence_transformers import SentenceTransformer, util
model = SentenceTransformer('multi-qa-MiniLM-L6-cos-v1')
tokenizer = model.tokenizer

query_embedding = model.encode('How big is London')
query_embedding_plus = tokenizer.encode_plus('How big is London', add_special_tokens=False, return_offsets_mapping=True)
# check if the two embeddings are the same
print("Embedding:", query_embedding)
print("Embedding Plus:", query_embedding_plus)
passage_embedding = model.encode(['London has 9,787,426 inhabitants at the 2011 census',
                                  'London is known for its finacial district'])

print("Similarity:", util.dot_score(query_embedding, passage_embedding))