# -*- mode: Starlark -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/activly ./cmd/main.go'

# Local deployment
local_resource(
  'activly-compile',
  compile_cmd,
  deps=['./cmd/main.go',],
)

docker_build_with_restart(
  'activly-backend',
  '.',
  entrypoint=['/app/activly'],
  dockerfile='deployments/Dockerfile',
  only=[
    './build',
    #'./web',
  ],
  live_update=[
    sync('./build/activly', '/app/activly'),
    #sync('./web', '/app/web'),
  ],
)

if k8s_context() != 'kind-kind':
    fail('Cannot access kind-kind Kubernetes context.  Is the Kind cluster accessible?')

k8s_yaml(['deployments/tilt/activly-backend.yaml', 'deployments/tilt/postgres.yaml', ])

k8s_resource(
  'activly-backend',
  port_forwards=8080,
  #resource_deps=['deploy', 'activly-compile']
  resource_deps=['activly-compile']
)

k8s_resource(
  'postgres',
  port_forwards=5432,
)
