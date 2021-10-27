def main(ctx):
  return {
    'version': '1',
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