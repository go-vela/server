def main(ctx):
  image = "alpine"

  return {
      'version': '1',
      'steps': [
        {
            "name": "foo",
            "image": image,
            "parameters": {
                "registry": "foo"
            }
        }
      ]
  }
