#!/bin/bash
# Translators script for uWelcome

# Repeating usage instructions if no language code is provided

if [ -z "$1" ]; then
    echo "Run the command like this: \`$0 <your-language-code>\`"
    echo ""
    echo "> Example: \`$0 fr\`"
    exit 1
fi

# Checking dependencies

which go >/dev/null 2>&1 || (echo "You don't have \`go\` installed :(" && exit 1)
which gettext >/dev/null 2>&1 || (echo "You don't have \`gettext\` installed :(" && exit 1)
[ -f "$(go env GOPATH)/bin/xgotext" ] || (go install github.com/leonelquinteros/gotext/cli/xgotext@latest || (echo "An error occurred while installing \`xgotext\` :(" && exit 1))

# If the language already exists, update it

if [ -f "locales/$1/default.po" ]; then
    echo "Language $1 already exists. Updating..."
    ~/go/bin/xgotext -in . -out locales/temp || (echo "An error occurred while running \`xgotext\` :(" && rm -rf locales/temp && exit 1)
    msgmerge --update locales/"$1"/default.po locales/temp/default.pot || (echo "An error occurred while running \`msgmerge\` :(" && rm -rf locales/temp && exit 1)
    rm -rf locales/temp
    rm -f locales/"$1"/default.po~
    echo "Translations for $1 updated!"
    exit 0
fi

# If the language does not exist, create it

echo "Creating new language $1..."
mkdir -p locales/"$1"
~/go/bin/xgotext -in . -out locales/temp || (echo "An error occurred while running \`xgotext\` :(" && rm -rf locales/temp && exit 1)
cp locales/temp/default.pot locales/"$1"/default.po
rm -rf locales/temp
echo "Translations for $1 generated. You can edit them in locales/$1/default.po"
