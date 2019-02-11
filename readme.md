openNAMU
====
[![Python 3.5 or later Required](https://img.shields.io/badge/python-3.5%20or%20higher-blue.svg)](https://python.org)
[![LICENSE](https://img.shields.io/badge/license-BSD%203--Clause-lightgrey.svg)](./LICENSE)

![](https://raw.githubusercontent.com/2du/openNAMU/master/.github/logo.png)

openNAMU is a Python-based wiki engine. You can use openNAMU by installing Python and its dependency modules, and you can modify the code yourself to create more specialized wikis.

 * [(README for Korean)](./readme-ko.md)

### Index
 * [Getting Started](#getting-started)
 * [Clone](#clone)
 * [Contribute](#contribute)
 * [License](#license)
 * [Authors](#authors)
 * [Etc.](#etc)

# Getting Started
openNAMU is based upon Python, and it requires a Python environment.

## Set-Ups
### Install Python
See [Python Installation Guide(KR)](https://github.com/404-sdok/how-to-python/blob/master/0.md).

### Download Releases
Download the [release version of openNAMU](https://github.com/2du/openNAMU/releases), and unzip the file. It is also possible to download releases by [cloning this repository](#Clone).

### Install Modules
Windows
```
pip install -r requirements.txt
```

Linux
```
pip3 install -r requirements.txt
```
## Launching Application
Windows
```
python app.py
```

Linux
```
python3 app.py
```

## Publishing Application

# Clone
You can clone this repository by entering the following command at the terminal (command prompt):
## Stable
 * `git clone -b stable https://github.com/2du/openNAMU.git`

## Beta
 * `git clone -b master https://github.com/2du/openNAMU.git`

# Contribute
openNAMU may have some untested bugs. Your use of openNAMU and bug discovery will help develop openNAMU.
[Create Issues](https://github.com/2du/openNAMU/issues/new)

openNAMU is open source project. Add new features and request pull requests. 
[Create Pull Requests](https://github.com/2du/openNAMU/compare)

# Lisence
openNAMU is protected by [BSD 3-Clause License](./LICNESE). Please refer to the documentation for details.

## External Projects
 * Quotes icon [Dave Gandy](http://www.flaticon.com/free-icon/quote-left_25672) CC 3.0 BY
 * Syntax highlighting [highlightjs](https://highlightjs.org/)
 * Numerical expression [MathJax](https://www.mathjax.org/)

# Authors
 * [Reference](https://github.com/2DU/openNAMU/graphs/contributors)

## Special Thanks
 * [Team Croatia](https://github.com/TeamCroatia)
 * Basix
 * Efrit
 * Other chat rooms

# Etc.
`set.json` is a configuration file that stores some local settings.
 * [filename].db = Database name

If you delete `set.json`, you can create a new one again.

[Test Server](http://namu.ml/)

Owner rights are granted to the first registor.
