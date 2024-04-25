import { QueryInput, SearchResult } from "../../api/types.js";

// Dummy initial data
const searchResults: SearchResult[] = [
    {
        id: "embedding:key:1",
        hash: "hash1",
        chunkID: 1,
        content: '\nconst x = 5;\nlet y = 10;\nfunction multiply() {\n  return x * y;\n}\nmultiply();\ndivide();',
        path: 'example.js',
        owner: 'OpenAI',
        repo: 'GPT',
        startLine: 10,
        startIndex: 1,
        endIndex: 6,
        score: 0.9
    },
    {
        id: "embedding:key:2",
        hash: "hash2",
        chunkID: 2,
        content: 'import tensorflow as tf\n\nx = tf.constant(5)\ny = tf.constant(10)\nresult = tf.multiply(x, y)\nprint(result)',
        path: 'script.py',
        owner: 'Google',
        repo: 'TensorFlow',
        startLine: 1,
        startIndex: 1,
        endIndex: 6,
        score: 0.8
    },
    {
        id: "embedding:key:3",
        hash: "hash3",
        chunkID: 3,
        content: 'public class HelloWorld {\n  public static void main(String[] args) {\n    System.out.println("Hello, World!");\n  }\n}',
        path: 'HelloWorld.java',
        owner: 'Facebook',
        repo: 'React',
        startLine: 1,
        startIndex: 1,
        endIndex: 6,
        score: 0.7
    },
];

class SearchApi {
    search(query: QueryInput): SearchResult[] {
        return searchResults;
    }
}

export const searchApi = new SearchApi();