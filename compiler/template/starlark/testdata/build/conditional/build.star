def main(ctx):
  stepNames = ["foo", "bar", "star"]

  steps = []

  for name in stepNames:
    if name == "foo" and ctx["vela"]["build"]["branch"] == "main":
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
