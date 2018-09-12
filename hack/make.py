#!/usr/bin/env python


# http://stackoverflow.com/a/14050282
def check_antipackage():
    from sys import version_info
    sys_version = version_info[:2]
    found = True
    if sys_version < (3, 0):
        # 'python 2'
        from pkgutil import find_loader
        found = find_loader('antipackage') is not None
    elif sys_version <= (3, 3):
        # 'python <= 3.3'
        from importlib import find_loader
        found = find_loader('antipackage') is not None
    else:
        # 'python >= 3.4'
        from importlib import util
        found = util.find_spec('antipackage') is not None
    if not found:
        print('Install missing package "antipackage"')
        print('Example: pip install git+https://github.com/ellisonbg/antipackage.git#egg=antipackage')
        from sys import exit
        exit(1)
check_antipackage()

# ref: https://github.com/ellisonbg/antipackage
import antipackage
from github.appscode.libbuild import libbuild, pydotenv

import os
import os.path
import subprocess
import sys
from os.path import expandvars, join, dirname

libbuild.REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/kubevault/csi-driver'
BUILD_METADATA = libbuild.metadata(libbuild.REPO_ROOT)
libbuild.BIN_MATRIX = {
    'csi-vault': {
        'type': 'go',
        'go_version': True,
        'use_cgo': False,
        'distro': {
            'alpine': ['amd64'],
            'linux': ['amd64']
        }
    }
}
libbuild.BUCKET_MATRIX = {
    'prod': 'gs://appscode-cdn',
    'dev': 'gs://appscode-dev'
}


def call(cmd, stdin=None, cwd=libbuild.REPO_ROOT):
    print(cmd)
    return subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def die(status):
    if status:
        sys.exit(status)


def check_output(cmd, stdin=None, cwd=libbuild.REPO_ROOT):
    print(cmd)
    return subprocess.check_output([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def version():
    # json.dump(BUILD_METADATA, sys.stdout, sort_keys=True, indent=2)
    for k in sorted(BUILD_METADATA):
        print(k + '=' + BUILD_METADATA[k])


def fmt():
    libbuild.ungroup_go_imports('*.go', 'cloud', 'cmds')
    die(call('goimports -w *.go driver cmds'))
    call('gofmt -s -w *.go driver cmds')


def vet():
    call('go vet *.go')
    call('go vet $(go list ./... | grep -v /vendor/)')


def lint():
    call('golint *.go')
    call('golint $(go list ./... | grep -v /vendor/)')


def gen_assets():
    pass


def gen():
    gen_assets()


def build_cmd(name):
    cfg = libbuild.BIN_MATRIX[name]
    if cfg['type'] == 'go':
        if 'distro' in cfg:
            for goos, archs in cfg['distro'].items():
                for goarch in archs:
                    libbuild.go_build(name, goos, goarch, main='*.go')
        else:
            libbuild.go_build(name, libbuild.GOHOSTOS, libbuild.GOHOSTARCH, main='*.go')


def build_cmds():
    gen()
    for name in libbuild.BIN_MATRIX:
        build_cmd(name)


def build(name=None):
    if name:
        cfg = libbuild.BIN_MATRIX[name]
        if cfg['type'] == 'go':
            gen()
            build_cmd(name)
    else:
        build_cmds()


def install():
    die(call('GOBIN={} {} install *.go'.format(libbuild.GOBIN, libbuild.GOC)))


def default():
    gen()
    fmt()
    die(call('GOBIN={} {} install .'.format(libbuild.GOBIN, libbuild.GOC)))


def test(type, *args):
    pydotenv.load_dotenv(join(libbuild.REPO_ROOT, 'hack/config/.env'))
    if type == 'unit':
        die(call(libbuild.GOC + ' test -v ./pkg/...'))
    elif type == 'e2e':
        die(call('ginkgo -r --v --progress --trace -- --v=3'))
    else:
        print '{test unit|e2e}'


if __name__ == "__main__":
    if len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        globals()[sys.argv[1]](*sys.argv[2:])
    else:
        default()
