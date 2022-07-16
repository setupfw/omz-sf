plugins=()
while read -r line; do
    if [[ "$line" != '#'* ]]; then
        read -A list <<<"$line"
        plugins+=("${list[@]}")
    fi
done <~/.zshrc.pluglist

source ~/.oh-my-zsh/lib/clipboard.zsh
detect-clipboard && plugins+=(copybuffer copyfile copypath)

for it in rsync tmux docker yarn rustc; do
    [ -x "$(command -v "$it")" ] && plugins+=("$it")
done

[ -x "$(command -v pip)" ] || [ -x "$(command -v pip3)" ] && plugins+=(pip)
