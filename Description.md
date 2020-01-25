What’s in a Docker registry?
Docker registries are used as a version control system for the development and deployment of container images. Docker images provide attackers with insights into the inner workings of an organization’s containers and the services running within them.

Any values compiled into the container during build time will be available through the image’s filesystem and metadata. As a result, if the image creator chose to supply any secret or sensitive values at build time rather than at runtime, then these values will be readable to attackers who simply obtain the container’s image from a compromised registry.

If the images are properly stripped of secrets, attackers may still find use in the application logic and configuration files stored in the image. This information can be extremely helpful to attackers who no longer need to make assumptions and discoveries about the application’s behavior.

What can go-pillage-registries do?
pilreg can be used as easily as the regular Docker tools by taking advantage of Google’s container registry library which utilizes Docker’s authentication semantics (run docker login and you can authenticate!).

You can simply provide a registry, and pilreg will enumerate all repositories and their tags using Docker’s /v2/_catalog and /v2/<name>/tags/list APIs. If the target registry does not support these APIs, or if you just don’t want to pillage their entire registry, then you can supply a list of known/desired repositories and tags. For registries with invalid TLS configurations, pilreg also has the option to ignore TLS verification (only suggested for security professionals)

When pilreg finds images, it automatically requests their manifest and configuration and outputs the results. There is also the option to store these results in an output directory for inspection. If an output directory is specified, then attackers also have the option to pull the filesystems for each image and store the resulting archives.

Demo
In the below demonstration, we will create an image with a secret set in an environment variable at build time to demonstrate an attacker’s ability to retrieve the value from the image’s configuration stored on a remote registry:

Start registry:
docker run -d -p 5000:5000 registry
Create the Docker image and push it to the registry
mkdir /tmp/example
cd /tmp/example

cat <<EOF > Dockerfile
FROM ubuntu
ENV "SUPERSECRET" "123456"
CMD ["uname", "-a"]
EOF

docker build -t 127.0.0.1:5000/test/test .
docker push 127.0.0.1:5000/test/test
Install pilreg
git clone https://github.com:nccgroup/go-pillage-registries.git

cd go-pillage-registries
go install ./...
Pillage the registry:
pilreg 127.0.0.1:5000 | tee pillaged.json
[
  {
    "Reference": "127.0.0.1:5000/test/test:latest",
    "Registry": "127.0.0.1:5000",
    "Repository": "test/test",
    "Tag": "latest",
    "Manifest": "{...}",
    "Config": "{...}",
    "Error": null
  }
]
Search for the secret:
jq .[].'Config' pillaged.json -r | jq . | grep SECRET
      "SUPERSECRET=123456"
      "SUPERSECRET=123456"
      "created_by": "/bin/sh -c #(nop)  ENV SUPERSECRET=123456",
