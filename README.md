# gcrproxy

Is an HTTP proxy that is intended to run on Google Cloud Run. 


## Why Google Cloud Run?

The main purpose of gcrproxy is to bypass government censorship on social media websites. The most common methods used to bypass these restrictions are to use the [TOR browser](https://www.torproject.org/) and to use paid VPN services but both methods have their own disadvantages. For example, it's very annoying to use social media websites through the TOR network because TOR exit nodes are restricted, constantly asked for captchas, or completely banned by social media providers. And since VPN services have their IP blocks known to the public, it's possible for governments to restrict their usage. Thus, this project comes with a different approach that is, to use google cloud run as a proxy. Since google cloud network is used by many companies, restricting access to these IPs would cause harm to those companies who pay taxes, would cause some important services to crash, and therefore, it would harm the economy and the government itself.


**TLDR:** This project assumes the google cloud network is too important to be restricted, so it uses it as a proxy.


## Setup

*Note: You need to have your gcloud SDK installed and authorized before following these steps. [Here](https://cloud.google.com/sdk/docs/quickstart-debian-ubuntu) is a quick start guide for Debian/Ubuntu.*

First clone the repository:
```sh
git clone https://github.com/paddlesteamer/gcrproxy.git
cd gcrproxy/
```

Now build the docker image. Note that `$PROJECT_ID` is your GCR project id:
```sh
docker image build -t gcr.io/$PROJECT_ID/gcrproxy .
```

Then push the image to google cloud:
```sh
docker push gcr.io/$PROJECT_ID/gcrproxy
```

And deploy:
```sh
gcloud run deploy --image gcr.io/$PROJECT_ID/gcrproxy --platform managed --region $REGION
```

You can choose the region according to this [document](https://cloud.google.com/run/docs/locations). The deploy command will return you an URL, note it to somewhere. Now, let's build and run the tunnel application:

```sh
go build ./cmd/tunnel
PROXY=$GCR_URL go run ./cmd/tunnel
```

And done. All you need to do is to configure `127.0.0.1:1080` as an HTTP proxy in your browser.

