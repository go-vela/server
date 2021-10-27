def main(ctx):
    return {
        'version': '1',
        'steps': [
            step('foo'),
            step('bar')
        ],
    }
    
def step(word):
    return {
        "name": "build_%s" % word,
        "image": "alpine:latest",
        'commands': [
            "echo %s" % word
        ]
    }