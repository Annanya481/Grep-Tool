# Grep-Tool
This is a custom-built implementation of a grep-like tool written in Go. It allows you to search for patterns within text files or strings using regular expressions, much like the Unix grep command but implemented from scratch to understand pattern matching at a deeper level.

## Features
Supports basic regular expressions including:
- Anchors: ^ for the beginning of a line and $ for the end.
- Character Classes: For example, [abc] matches any one of 'a', 'b', or 'c'.
- Negated Character Classes: [^abc] matches any character except 'a', 'b', or 'c'.
- Quantifiers:
    - '?' for zero or one occurrence.
    - '+' for one or more occurrences.
- Escape Sequences: Support for matching digits (\d) and word characters (\w).
- Wildcards: '.' matches any single character.
- Alternation: For example, (cat|dog) matches 'cat' or 'dog'.

## Installation
Clone the repository
   ```
   git clone https://github.com/Annanya481/Grep-Tool.git
   ```
## Usage
 Run the script using the following command
   ```
   echo -n "<line>" | ./script.sh -E "<pattern>"
   ```
For example:
   ```
   echo -n "foo" | ./script.sh -E "(foo|bar)"
   ```
## Explanation
For a detailed explanation of this tool, read the article: 

## Contributing
Feel free to open issues or submit pull requests if you would like to contribute to this project.

1. Fork the repository
2. Create a new branch for your feature (git checkout -b feature-branch)
3. Commit your changes (git commit -m 'Add some feature')
4. Push to the branch (git push origin feature-branch)
5. Open a pull request

## License
This project is licensed under the MIT License - see the LICENSE file for details.
