_git_contexts()
{
  local cur;
  cur=${COMP_WORDS[COMP_CWORD]}
  COMPREPLY=( $(compgen -W "$(gitctx)" -- $cur) );
}

complete -F _git_contexts gitctx