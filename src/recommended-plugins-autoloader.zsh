source ~/.oh-my-zsh/lib/clipboard.zsh
detect-clipboard && plugins+=(copybuffer copyfile copypath)

for it in fzf rsync tmux docker yarn rustc; do
    [ -x "$(command -v "$it")" ] && plugins+=("$it")
done

[ -x "$(command -v pip)" ] || [ -x "$(command -v pip3)" ] && plugins+=(pip)
