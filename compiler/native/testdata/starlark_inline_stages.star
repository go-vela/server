def main(ctx):
  stageNames = ["foo", "bar"]

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
          "image": "alpine",
          'commands': [
              "echo hello from %s" % word
          ]
        }
      ]
  }