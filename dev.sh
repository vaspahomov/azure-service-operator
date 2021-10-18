#!/bin/sh

set -eu

GIT_ROOT=$(git rev-parse --show-toplevel)
TOOL_DEST=$GIT_ROOT/hack/tools

if [ ! -f "$TOOL_DEST/task" ]; then # check for local installation
    if [ ! -f "/usr/local/bin/task" ]; then # or devcontainer installation
        $GIT_ROOT/.devcontainer/install-dependencies.sh local # otherwise, install the tools
    fi
fi

# Setup envtest binaries
# NB: if you change this, .devcontainer/Dockerfile also likely needs updating
source <(setup-envtest use -i -p env 1.12.x) # this sets KUBEBUILDER_ASSETS
export PATH="$KUBEBUILDER_ASSETS:$TOOL_DEST:$PATH"

echo "Entering $SHELL with expanded PATH (use 'exit' to quit):"
echo "Try running 'task -l' to see possible commands."
$SHELL
