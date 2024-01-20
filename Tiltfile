# -*- mode: Starlark -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

# Records the current time, then kicks off a server update.
# Normally, you would let Tilt do deploys automatically, but this
# shows you how to set up a custom workflow that measures it.
#local_resource(
#    'deploy',
#    'python3 record-start-time.py',
#)

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/activly ./cmd/main.go'

# Local deployment
local_resource(
  'activly-compile',
  compile_cmd,
  deps=['./cmd/main.go',],
  #resource_deps=['deploy'],
)

docker_build_with_restart(
  'activly-image',
  '.',
  entrypoint=['/app/build/activly'],
  dockerfile='deployment/Dockerfile',
  only=[
    './build',
    #'./web',
  ],
  #disable_push=True,
  live_update=[
    sync('./build', '/app/build'),
    #sync('./web', '/app/web'),
  ],
)

#allow_k8s_contexts('kind-kind')

if k8s_context() != 'kind-kind':
    fail('Requires Kind to be the active k8s context')

k8s_yaml('deployment/kubernetes.yaml')

k8s_resource(
  'activly',
  port_forwards=8080,
  #resource_deps=['deploy', 'activly-compile']
  resource_deps=['activly-compile']
)
