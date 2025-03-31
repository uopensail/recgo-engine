# recgo-engine

![Version: 0.1.16](https://img.shields.io/badge/Version-0.1.16-informational?style=flat-square)

A Helm chart for recgo-engine

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| bisonlou | <bisonlou@gmail.com> |  |
| Yann-J | <yann.jouanique@gmail.com> |  |
| Nzeugaa | <jean.poutcheu@gmail.com> |  |

## TL;DR;

[recgo-engine](https://www.recgo-engine.io/) is an open-source remote configuration / activation flag service.

```console
$ helm repo add one-acre-fund https://one-acre-fund.github.io/oaf-public-charts
$ helm install my-release one-acre-fund/recgo-engine
```

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami/ | mongodb | ~13.6.1 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| autoscaling.enabled | bool | `false` |  |
| autoscaling.maxReplicas | int | `100` |  |
| autoscaling.minReplicas | int | `1` |  |
| autoscaling.targetCPUUtilizationPercentage | int | `80` |  |
| autoscaling.targetMemoryUtilizationPercentage | int | `80` |  |
| fullnameOverride | string | `""` |  |
| recgo-engine.extraEnvVars.API_HOST | string | `"https://api-recgo-engine.uopensail.com:443"` |  |
| recgo-engine.extraEnvVars.NODE_ENV | string | `"production"` |  |
| recgo-engine.jwtSecret | string | `"jwtSecretString"` |  |
| recgo-engine.persistence.accessModes[0] | string | `"ReadWriteMany"` |  |
| recgo-engine.persistence.enabled | bool | `true` |  |
| recgo-engine.persistence.storage | string | `"3Gi"` |  |
| recgo-engine.persistence.type | string | `"emptyDir"` |  |

| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"recgo-engine/recgo-engine"` |  |
| image.tag | string | `"latest"` |  |
| imagePullSecrets | list | `[]` |  |
| ingress.annotations."kubernetes.io/ingress.class" | string | `"nginx"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/cors-allow-headers" | string | `"Authorization,Referer,sec-ch-ua,sec-ch-ua-mobile,sec-ch-ua-platform,User-Agent,X-Organization,Content-Type"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/cors-allow-origin" | string | `"https://api-recgo-engine.uopensail.com"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/enable-cors" | string | `"true"` |  |
| ingress.annotations."nginx.ingress.kubernetes.io/force-ssl-redirect" | string | `"true"` |  |
| ingress.apiHostName | string | `"api-recgo-engine.uopensail.com"` |  |
| ingress.enabled | bool | `false` |  |
| ingress.name | string | `"recgo-engine-ingress"` |  |
| ingress.secretName | string | `"recgo-engine-tls"` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| port.apiHTTPPort | int | `8080` |  |
| port.promePort | int | `8082` |  |
| replicaCount | int | `1` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| tolerations | list | `[]` |  |
| volume.mountPath | string | `"/usr/local/src/app/packages/back-end/uploads"` |  |
| volume.name | string | `"uploads-persistent-storage"` |  |
