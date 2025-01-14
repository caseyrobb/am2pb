# am2pb - Alertmanager to Pushbullet
[![Docker Repository on Quay](https://quay.io/repository/carobb/am2pb/status "Docker Repository on Quay")](https://quay.io/repository/carobb/am2pb)

A simple service that listens for alerts from Alertmanager via the [webhook receiver](https://prometheus.io/docs/alerting/configuration/#webhook_config),
and forwards them to [Pushbullet](http://pushbullet.com/).

This project was created to help me learn Go, do not use in production.

# Usage

Grab a token from your [Pushbullet account](https://www.pushbullet.com/#settings/account).

## Docker/Podman

Start the server:
```
$ podman run --rm \
    -e BEARER_TOKEN=<pushbullet api key> \
    -p 5000:5000 \
    quay.io/carobb/am2pb
```

Configure Alertmanager to use a webhook endpoint and ensure the URL points to the exposed service from the previous step:
```yaml
- name: "pb-receiver"
  webhook_configs:
  - url: http://127.0.0.1:5000/
```

## Kubernetes

Update [secret.yaml](https://github.com/caseyrobb/am2pb/blob/master/deploy/secret.yaml) with your Pushbullet token.  A sample `Deployment` and `Service` are included.

```
$ kubectl apply -k deploy -n <namespace>
```
