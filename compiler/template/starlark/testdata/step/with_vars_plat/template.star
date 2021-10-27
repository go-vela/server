def main(ctx):
    return {
        'version': '1',
        'steps': [
            step(ctx["vela"]["repo"]["full_name"]),
        ],
    }

def step(full_name):
    return {
        "name": "echo %s" % full_name,
        "image": "alpine:latest",
        'commands': [
            "echo %s" % full_name
        ]
    }
