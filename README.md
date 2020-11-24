go-template-http
================

go-template-http is a simple HTTP server written in Go for rendering
[Go templates](https://golang.org/pkg/text/template/).

It does one thing and does it well.

getenv
------

It supports rendering environment variables through a custom `getenv` function.

For example, the following Go template shows the `$PATH` variable:

	$ curl -d '$PATH is: {{getenv "PATH"}}' localhost:8080
	$PATH is: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

environ
-------

To get an overview of all available environment variables, the custom `environ`
function can be used:

	$ curl -d '{{range environ}}{{println .}}{{end}}' localhost:8080
	PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
	HOSTNAME=cf70703ecf6d
	HOME=/

Motivation
----------

The motivation for creating this tool was to have a way to share build secrets
with `docker build` without embedding the secret into the image.

It's a kindof tricky thing to do, but it's possible by calling the network, for
example:

	$ export BITBUCKET_USERNAME=...
	$ export BITBUCKET_PASSWORD=...
	$ docker network create build
	$ docker run --env=BITBUCKET_USERNAME --env=BITBUCKET_PASSWORD --network=build --name=secrets --rm -d aude/go-template-http
	$ cat << '_EOF_' > Dockerfile
	FROM python

	WORKDIR /app

	RUN printf '%s\n' \
		'machine bitbucket.org' \
		'login {{getenv "BITBUCKET_USERNAME"}}' \
		'password {{getenv "BITBUCKET_PASSWORD"}}' \
		| curl --fail --data-binary @- http://secrets:8080/ > ~/.netrc \
		&& pip install git+https://bitbucket.org/your-company/private-package.git \
		&& rm ~/.netrc

	# extra check, to prove that ~/.netrc is not embedded in Docker image
	RUN test ! -f ~/.netrc

	COPY . .

	CMD ["python", "main.py"]
	_EOF_
	$ docker build --network=build .
	...
	$ docker stop secrets
	$ docker network rm build

