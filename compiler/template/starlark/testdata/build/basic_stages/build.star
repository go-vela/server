def main(ctx):
  stageNames = ["foo", "bar", "star"]

  stages = {}

  for name in stageNames:
    stages[name] = stage(name)

  return {
      'version': '1',
      'stages': stages
  }

def stage(word):
  return {
      "steps": [
        {
          "name": "build_%s" % word,
          "image": "alpine:latest",
          'commands': [
              "echo hello from %s" % word
          ]
        }
      ]
  }
