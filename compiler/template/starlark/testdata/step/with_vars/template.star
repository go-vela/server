def main(ctx):
    steps = [step(x, ctx["vars"]["pull_policy"], ctx["vars"]["commands"]) for x in ctx["vars"]["tags"]]

    pipeline = {
        'version': '1',
        'steps': steps,
    }

    return pipeline

def step(tag, pull_policy, commands):
    return {
        "name": "build %s" % tag,
        "image": "golang:%s" % tag,
        "pull": pull_policy,
        'commands': commands.values(),
    }
