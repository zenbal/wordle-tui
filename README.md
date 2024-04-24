## wordle-tui

A simple terminal-based word guessing game written in Go.

### Features

1. **Wordle Game**: Aims to provide a similar look and feel to the original game.
2. **Suggestions**: Get a suggested next guess based on the current state of the game. Implemented using a backtracking algorithm and a trie data structure.
3. **Hints**: Receive hints when you've made a suboptimal guess.

### Installation

#### From Source

To build the Wordle Game App from source, clone the repository and use `go build`.

```bash
git clone https://github.com/zenbal/wordle-tui.git
cd wordle-tui 
go build
```

#### Pre-built Binaries

Alternatively, you can download the latest binary build from the [Releases](https://github.com/zenbal/wordle-tui/releases) page.

### Usage

Once installed, simply run the executable to start playing. Press `?` to show the available shortcuts.

### Credits

- Original game by [Josh Wardle](https://www.powerlanguage.co.uk/)
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)

### Contributions

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

Enjoy playing! 🎉
