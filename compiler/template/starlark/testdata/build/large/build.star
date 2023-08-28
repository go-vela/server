######
## Setup the build matrix with the base versions a human will maintain.
######

DISTRO_WITH_VERSIONS = {
    # n.b. these reduce to DockerHub tags
    # https://hub.docker.com/_/python/tags?name=alpine
    # https://endoflife.date/alpine
    'alpine': [
        '3.17',  # EOL 22 Nov 2024
        '3.18'   # EOL 09 May 2025
    ],
    # https://hub.docker.com/_/python/tags?name=slim
    # https://endoflife.date/debian
    'debian': [
        'slim-bullseye',  # EOL 30 Jun 2026
        'slim-bookworm'   # EOL 10 Jun 2028
    ]
}
PYTHON_VERSIONS = [
    '3.8',
    '3.9',
    '3.10',
    '3.11'
]
POETRY_VERSIONS = [
    '1.6.1'
]

KANIKO_IMAGE = 'target/vela-kaniko:latest'
HADOLINT_IMAGE = 'hadolint/hadolint:v2.12.0-alpine'


## The base Docker container build step's config for push builds
def base():
    return {
        'image': KANIKO_IMAGE,
        'ruleset': {
            'event': 'push',
            'branch': 'main'
        },
        'pull': 'not_present',
        'secrets': [
            {
                'source': 'artifactory_password',
                'target': 'docker_password'
            }
        ]
    }


## The base Docker container plugin params for push builds
##
## These are parameters passed to Kaniko.
def base_params():
    return {
        'username': 'ibuildallthings',
        'registry': 'docker.example.com',
        'repo': 'docker.example.com/app/multibuild'
    }


## The step config for pull request builds
def pull_request():
    pr = base()
    pr['ruleset']['event'] = 'pull_request'
    pr['ruleset'].pop('branch')
    return pr


## The Kaniko params for pull request builds
def pull_request_params():
    prp = base_params()
    prp['dry_run'] = True
    return prp


## Define a linting stage that uses Hadolint inside of a Make task
##
## This keeps our Dockerfiles tidy and compliant with conventions
def stage_linting():
    return {
        'linting': {
            'steps': [{
                'name': 'check-docker',
                'image': HADOLINT_IMAGE,
                'pull': 'not_present',
                'commands': [
                    'time apk add --no-cache make',
                    'time make check-docker'
                ]
            }]
        }
    }


## Build stages comprised of a step for push and pull_request builds
def stage_build_tuple(distro, distro_version, python_version, poetry_version):
    pr = build_template("build", distro, distro_version, python_version, poetry_version, pull_request(), pull_request_params())
    base_step = build_template("publish", distro, distro_version, python_version, poetry_version, base(), base_params())
    combined = base_step | pr
    return combined


## Build a single stage for a build tuple, with its base step config and plugin parameters
def build_template(step_name, distro, distro_version, python_version, poetry_version, step_def_base, step_def_params):
    return {
        ('python_%s_%s_%s %s' % (python_version, distro, distro_version, step_def_base['ruleset']['event'])): {
                'steps': [step_def_base | {
                    'name': ('%s python-%s %s %s' % (step_name, python_version, distro, distro_version)),
                    'parameters': step_def_params | {
                        'dockerfile': ('python-%s.Dockerfile' % distro),
                        'build_args': [
                            'PYTHON_VERSION=%s' % python_version,
                            '%s_VERSION=%s' % (distro.upper(), distro_version),
                            'POETRY_VERSION=%s' % poetry_version
                        ],
                        'tags': [
                            '%s-%s-%s-%s' % (python_version, distro, distro_version, poetry_version),
                            '%s-%s-%s' % (python_version, distro, distro_version),
                            '%s-%s' % (python_version, distro)
                        ]
                    }
            }]
        }
    }


## Define a stage that uses the Slack template
def stage_slack_notify(needs):
    return {
        'slack': {
            'needs': needs,
            'steps': [{
                'name': 'slack',
                'template': {
                    'name': 'slack'
                }
            }]
        }
    }


## Builds the build matrix in the form of list of tuples from the constants defined at the top of the file
def build_matrix():
    BUILD_MATRIX = []
    for poetry_version in POETRY_VERSIONS:
        for python_version in PYTHON_VERSIONS:
            for distro in DISTRO_WITH_VERSIONS:
                for distro_version in DISTRO_WITH_VERSIONS[distro]:
                    BUILD_MATRIX.append((distro,
                                         distro_version,
                                         python_version,
                                         poetry_version))
    return BUILD_MATRIX


## Construct a secret
def secret(name, key, secret_type, engine='native'):
    return {'name': name, 'key': key, 'engine': engine, 'type': secret_type}


## Construct a template
def template(name, source, version=None, template_type='github'):
    real_source = '%s@%s' % (source, version) if version else source
    return {
        'name': name,
        'source': real_source,
        'type': template_type
    }

## The main method, the real deal.
##
## Vela actually calls this function, its return is what Vela uses.
def main(ctx):
    # Retrieve the org dynamically since we're using some org secrets
    vela_repo_org = ctx['vela']['repo']['org'] if 'vela' in ctx else "UNKNOWN-ORG"

    # Build the stages from the build matrix
    build_stages = {}
    for (distro, distro_version, python_version, poetry_version) in build_matrix():
        build_stages = build_stages | (stage_build_tuple(distro, distro_version, python_version, poetry_version))

    # assemble the stage list with the bookends of linting and notifications in place
    stages = stage_linting() | build_stages | stage_slack_notify(build_stages.keys())

    # Build the final output
    final = {
        'version': '1',
        'templates': [
            template(name='slack',
                     source='git.example.com/vela/vela-templates/slack/slack.yml')
        ],
        'stages': stages,
        'secrets': [
            secret('artifactory_password','platform/vela-secrets/artifactory_password_for_ibuildallthings', 'shared'),
            secret('slack_webhook', vela_repo_org + '/slack_webhook', 'org')
        ]
    }

    return final

