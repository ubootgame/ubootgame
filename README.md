# ubootgame

- Latest version on Github Pages: https://ubootgame.github.io/ubootgame/
- Blog on Github Pages: https://ubootgame.github.io/ubootgame/blog/

## Blogs

We use this [setup](https://t3st3ro.github.io/2022/11/02/self-contained-jekyll-with-docker.html) to avoid installing Jekyll locally.

### Initial setup

```bash
docker run --rm \
    --volume="${PWD}:/srv/jekyll" \
    --platform "linux/amd64" \
    -it jvconseil/jekyll-docker \
    sh -c "chown -R jekyll /usr/gem/ && jekyll new blog && bundle config set --local path 'vendor'"
```

### Run locally

```bash
docker-compose -f blog/docker-compose.yml up
```
