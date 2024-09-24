NOLINT_RE = '^-- NOLINT(: .*)?$'
NOTEST_RE = '^-- NOTEST(: .*)?$'

GO_CI_PATH_EXCLUDE_RE = [
    '^_tests/.*',
]

disabled_master_tests = [
    'commit_chain_length',
    'failed_parent_commits',
]

tests_configs = [
    {
        'id': 'golang-tests',
        'paths_re': ['.*'],
        'paths_re_exc': GO_CI_PATH_EXCLUDE_RE,
        'messages_re_exc': [NOTEST_RE],
        'tests': [
            {
                'name': 'GO Test',
                'JOB_NAME': 'webdev/golib/golang-tests',
                'COMMANDLINE': 'python3 $BUILDTOOLS_PATH/projects/mairon/go_test.py'
            }
        ]
    },
    {
        'id': 'golang-linting',
        'paths_re': ['.*'],
        'paths_re_exc': GO_CI_PATH_EXCLUDE_RE,
        'messages_re_exc': [NOLINT_RE],
        'tests': [
            {
                'name': 'GO Lint',
                'JOB_NAME': 'webdev/golib/golang-linting',
                'COMMANDLINE': 'python3 $BUILDTOOLS_PATH/projects/mairon/go_lint.py'
            }
        ]
    },
]
