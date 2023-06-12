def main(ctx):
  return {
    'version': '1',
    'environment': {
      'star': 'test3',
      'bar': 'test4',
    },
    'steps': [
      {
        'name': 'build',
        'image': 'golang:latest',
        'commands': [
          'go build',
          'go test',
        ]
      },
    ],
}
