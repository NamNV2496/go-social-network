
#!/usr/bin/env bash

# exit immediately if any command return non-zero
set -e

: ${ROOT_PATH:=${PWD}}
: ${BUILD_PATH:=${PWD}/.build}
: ${GRPC_GO_BUILT_OUTPUT:="go-built-gprc-output"}
: ${DEFAULT_GRPC_LANGUAGES:=go}
: ${BUF_CACHE_DIR:="/tmp/buf"}

: ${GIT_HOST:="github.com"}
: ${GIT_ORG:="NamNV2496"}
: ${GIT_USER_NAME:="NamNV2496"}
: ${GIT_USER_EMAIL:="NamNV2496@gmail.com"}


GO_MODIFIED_PROTOS=()

main() {
    mkdir -p $BUILD_PATH
    echo "run command $@"

    eval $@ || return
}

git_config() {
    export GOPRIVATE="$GIT_HOST/$GIT_ORG/*"
    git config "user.name" "${GIT_USER_NAME}"
    git config "user.email" "${GIT_USER_EMAIL}"
}

install() {
    go-install buf github.com/bufbuild/buf/cmd/buf@v1.14.0
    go-install protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
    go-install protoc-gen-openapiv2 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
    go-install protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
    go-install protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.31.0
    go-install protoc-gen-validate github.com/envoyproxy/protoc-gen-validate@v1.0.2
}

go-install() {
    local bin="$1"
    local module="$2"
    command -v "$bin" >/dev/null 2>&1 && return
    go install "$module"
}


create_build_dir_with_nessary_files() {
    local proto_file=$1
    local build_path=$BUILD_PATH/go_generated

    cd $ROOT_PATH
    mkdir $BUILD_PATH
    # copy all protos file to build
    cp -r $ROOT_PATH/protos $BUILD_PATH
    # copy all base file into build folder
    # cp -r lib/. $BUILD_PATH/go_generated
    cp -r lib/. $BUILD_PATH
    echo $build_path
}


access_build_dir() {
    local build_dir=$1
    cd $build_dir
    # echo "% ${PWD/$ROOT_PATH/service}." >&2
    mkdir -p $build_dir
}

build_proto_func(){
    local proto_file_trigger=($1)
    local generate_proto_files=()
    if [ -n "$proto_file_trigger" ]; then
        if [ "$proto_file_trigger" == "all" ]; then
            echo "generate all services in protos"
            generate_proto_files=($(find protos -type f -name '*.proto' -exec dirname {} \; | sort -u | sed 's#protos/##'))

        else
            echo "run build protobuf file $proto_file_trigger"
            generate_proto_files=("${proto_file_trigger[@]}")
        fi
    else
        echo "find changed file to generate"
        generate_proto_files=($(get_file_change)) || return
    fi
    [ -z "$proto_file_trigger" ] && return

    # scan all changed proto files
    for proto in ${generate_proto_files[@]}; do
        echo "generate for $proto"
        eval generate_protoc_go $proto $lang || return
        # cdir $ROOT_PATH
    done

    echo "run build protobuf file done"

}

get_file_change() {
    local proto_file_changes=()
    proto_file_changes=("test", "test1")
    echo ${proto_file_changes[@]}
}

generate_protoc_go() {
    local proto=$1
    local lang=$2

    [ -z "$proto" ] && return
    local build_dir
    build_dir=$(create_build_dir_with_nessary_files $build_dir) || return
    access_build_dir $build_dir

    # backward 1 folder
    cd ..
    # run generate command
    # buf generate --template ${BUILD_PATH}/buf.gen.go.yaml || return

    cd $build_dir/protos/$proto
    echo pwd ${PWD}
    echo proto $proto
    local mod=$GIT_HOST/$GIT_ORG/$GRPC_GO_BUILT_OUTPUT/golang/$proto
    echo mod ${mod}
    version=$(basename "$proto")
    echo "$version"
    # mkdir -p "${BUILD_PATH}/go_generated/$proto"
    GO_MODULE=$mod envsubst '$GO_MODULE' < ${BUILD_PATH}/go.mod.tmpl > ./go.mod

    go mod tidy || return
    go test ./... || return
    go test -c || return

    # GO_MODIFIED_PROTOS+=("$proto:$version")

    echo built ${proto} done
}


# start main
main $@