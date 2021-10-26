def main(ctx):
  stepNames = ["foo", "bar", "star"]

  steps = []

  for name in stepNames:
    steps.append(step(name))

  return {
      'version': '1',
      'steps': steps
  }

def step(word):
  return {
      "name": "build_%s" % word,
      "image": "alpine:latest",
      'commands': [
          "echo %s" % word
      ]
  }
