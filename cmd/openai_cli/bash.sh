alias aa='cd /path/to/go_tools/cmd/openai_cli; ./openai_cli'

function aaa() {
    f_path=$(readlink -f $1)
    cd /path/to/go_tools/cmd/openai_cli
    cat $f_path | ./openai_cli -a -n $1
    cd -
}
