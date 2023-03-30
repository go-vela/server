def main(ctx):
  step_list = [
      struct(name="foo"),
      struct(name="bar"),
      struct(name="star")
  ]

  steps = []

  for step in step_list:
    steps.append(build_step(step))

  return {
      'version': '1',
      'steps': steps
  }

def build_step(step):
  return {
      "name": "build_%s" % step.name,
      "image": "alpine:latest",
      'commands': [
          "echo %s" % step.name
      ]
  }
