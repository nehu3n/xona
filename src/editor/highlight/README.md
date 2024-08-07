<div align="center">

# Xona Highlighting ✨

_Xona provides syntax highlighting for various language - file extensions. Here we detail how the system is built and why it is a success in terms of extensibility._

</div>

## 📚 Architecture

Xona's architecture focuses on flexibility and extensibility to one's preferences. A simple and easy-to-use approach is used.

Syntax highlighting is based on themes, each configurable based on general grammar rules that apply to each language.

For example, the default theme for the editor is as follows, and is located at `root_path/xona/config/themes/default.toml`. (**Note:** _TOML is used for all configuration files_).

```toml
[rules.words]
color = "#8F00FF" # purple

[rules.functions]
color = "#0000FF" # blue

[rules.types.strings]
color = "FFA500" # orange

# ...
```

As you can see, the files follow rules that apply to languages. `[rules.words]` are the reserved words; `[rules.functions]` are the calls to functions and methods; `[rules.types.strings]` are, within the types, the strings; and so on.

This allows that by simply changing the `color` property for each rule, you can customize each theme in a fast and intuitive way. The colors use hexadecimal base. Each theme is interchangeable at any time.

To make one, it is as easy as creating a file in `root_path/xona/config/themes/` with the name of the theme and extension `.toml`. The complete guide can be found in the docs folder of this repository.
