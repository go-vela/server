def main(ctx):
  stepNames = ["foo", "bar"]

  steps = []

  for name in stepNames:
    steps.append(
          {
              "name": "build_%s" % name,
              "image": "alpine",
              'commands': [
                  "echo hello from %s" % name
              ]
          }
        )

  return {
      'version': '1',
      'steps': steps
  }